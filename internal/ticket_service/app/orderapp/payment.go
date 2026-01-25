package orderapp

import (
	"context"
	"database/sql"
	"errors"
	"example_shop/common/db"
	model2 "example_shop/internal/model"
	kitexuser "example_shop/kitex_gen/user"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PayOrder 支付发起用例：将订单状态推进到 PAYING（简化版）。
//
// 该接口表达“我准备去支付了”，并生成/返回第三方支付流水号 pay_no（本项目用 MOCKPAY-* 代替）。
//
// 状态机约束：
// - 允许从 PENDING_PAY 进入 PAYING
// - PAYING 可重复调用（幂等），但不应生成新的 pay_no（复用旧的）
// - 已 ISSUED 的订单直接返回已完成状态
//
// 并发控制：
// - 通过对订单行加 FOR UPDATE，保证并发请求对同一订单的状态推进是串行的
func (s *Service) PayOrder(ctx context.Context, req *kitexuser.PayOrderReq) (*kitexuser.PayOrderResp, error) {
	if strings.TrimSpace(req.UserId) == "" || strings.TrimSpace(req.OrderId) == "" {
		return &kitexuser.PayOrderResp{BaseResp: baseResp(400, "user_id/order_id不能为空")}, nil
	}

	var payNo string
	var outStatus string
	err := db.MysqlDB.Transaction(func(tx *gorm.DB) error {
		var o model2.OrderInfo
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("order_id = ? AND user_id = ?", req.OrderId, req.UserId).First(&o).Error; err != nil {
			return err
		}
		if o.OrderStatus == "ISSUED" {
			outStatus = "ISSUED"
			payNo = o.PayNo.String
			return nil
		}
		if o.OrderStatus != "PENDING_PAY" && o.OrderStatus != "PAYING" {
			return fmt.Errorf("订单状态不允许支付")
		}
		if o.PayDeadline.Valid && time.Now().After(o.PayDeadline.Time) {
			return fmt.Errorf("支付超时")
		}

		if o.PayNo.Valid {
			payNo = o.PayNo.String
		} else {
			payNo = "MOCKPAY-" + uuid.New().String()
		}
		updates := map[string]interface{}{
			"order_status": "PAYING",
			"pay_channel":  nullString(req.PayChannel),
			"pay_no":       nullString(payNo),
		}
		if err := tx.Model(&model2.OrderInfo{}).Where("order_id = ?", o.ID).Updates(updates).Error; err != nil {
			return err
		}
		writeOrderAuditLog(tx, o.ID, "PAY_ORDER", req.UserId, o.OrderStatus, "PAYING", map[string]any{"pay_channel": req.PayChannel, "pay_no": payNo})
		outStatus = "PAYING"
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &kitexuser.PayOrderResp{BaseResp: baseResp(404, "订单不存在")}, nil
		}
		return &kitexuser.PayOrderResp{BaseResp: baseResp(400, err.Error())}, nil
	}
	resp := &kitexuser.PayOrderResp{BaseResp: baseResp(200, "success"), OrderStatus: outStatus, PayStatus: "PAYING"}
	if payNo != "" {
		resp.PayNo = &payNo
	}
	if outStatus == "ISSUED" {
		resp.PayStatus = "PAID"
	}
	return resp, nil
}

// ConfirmPay 支付回调确认用例：锁座占用变 SOLD，订单状态推进到 ISSUED。
//
// 该接口表达“第三方已支付成功，请出票”。
//
// 关键一致性点：
// - 必须校验订单当前状态为 PAYING，避免重复回调把已取消/已退款订单错误推进
// - 必须校验 pay_no 匹配，避免串单（把别的订单的支付回调应用到当前订单）
// - 必须在同一事务内完成：
//   1) 占用从 LOCKED -> SOLD
//   2) 订单从 PAYING -> ISSUED
//   否则会出现“订单已出票但座位仍是 LOCKED”的不一致
//
// 幂等性：
// - 如果订单已 ISSUED，直接返回成功（允许第三方重复通知）
func (s *Service) ConfirmPay(ctx context.Context, req *kitexuser.ConfirmPayReq) (*kitexuser.ConfirmPayResp, error) {
	if strings.TrimSpace(req.UserId) == "" || strings.TrimSpace(req.OrderId) == "" || strings.TrimSpace(req.PayNo) == "" {
		return &kitexuser.ConfirmPayResp{BaseResp: baseResp(400, "user_id/order_id/pay_no不能为空")}, nil
	}

	var outStatus string
	err := db.MysqlDB.Transaction(func(tx *gorm.DB) error {
		var o model2.OrderInfo
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("order_id = ? AND user_id = ?", req.OrderId, req.UserId).First(&o).Error; err != nil {
			return err
		}
		if o.OrderStatus == "ISSUED" {
			outStatus = "ISSUED"
			return nil
		}
		if o.OrderStatus != "PAYING" {
			return fmt.Errorf("订单状态不允许确认支付")
		}
		if !o.PayNo.Valid || o.PayNo.String != req.PayNo {
			return fmt.Errorf("pay_no不匹配")
		}
		if o.PayDeadline.Valid && time.Now().After(o.PayDeadline.Time) {
			return fmt.Errorf("支付已超时")
		}
		third := strings.TrimSpace(req.GetThirdPartyStatus())
		if third != "" && strings.ToUpper(third) != "SUCCESS" {
			return fmt.Errorf("第三方支付未成功")
		}

		var occs []model2.SeatSegmentOccupancy
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("order_id = ? AND status = ?", o.ID, "LOCKED").Find(&occs).Error; err != nil {
			return err
		}
		if len(occs) == 0 {
			return fmt.Errorf("未找到锁定座位")
		}
		now := time.Now()
		for _, oc := range occs {
			if oc.LockExpireTime.Valid && now.After(oc.LockExpireTime.Time) {
				return fmt.Errorf("锁座已过期")
			}
		}
		if err := tx.Model(&model2.SeatSegmentOccupancy{}).
			Where("order_id = ? AND status = ?", o.ID, "LOCKED").
			Updates(map[string]interface{}{"status": "SOLD", "lock_expire_time": sql.NullTime{Valid: false}}).Error; err != nil {
			return err
		}

		payTime := time.Now()
		if err := tx.Model(&model2.OrderInfo{}).Where("order_id = ?", o.ID).Updates(map[string]any{
			"order_status": "ISSUED",
			"pay_time":     nullTime(payTime),
		}).Error; err != nil {
			return err
		}
		writeOrderAuditLog(tx, o.ID, "PAY_CALLBACK", req.UserId, "PAYING", "ISSUED", map[string]any{"pay_no": req.PayNo, "third_party_status": third})
		outStatus = "ISSUED"
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &kitexuser.ConfirmPayResp{BaseResp: baseResp(404, "订单不存在")}, nil
		}
		return &kitexuser.ConfirmPayResp{BaseResp: baseResp(400, err.Error())}, nil
	}
	return &kitexuser.ConfirmPayResp{BaseResp: baseResp(200, "success"), OrderStatus: outStatus}, nil
}
