package orderapp

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"example_shop/common/db"
	model2 "example_shop/internal/model"
	"example_shop/internal/ticket_service/app/shared"
	kitexuser "example_shop/kitex_gen/user"
	"fmt"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CancelOrder 取消订单用例：释放 LOCKED 的区间占用并将订单置为 CANCELLED。
//
// 取消与退款的区别：
// - 取消：发生在“未出票”阶段（PENDING_PAY / PAYING），核心动作是释放锁座
// - 退票：发生在“已出票”阶段（ISSUED），核心动作是释放已售占用 + 计算退款金额
//
// 状态机约束：
// - 仅允许 PENDING_PAY / PAYING -> CANCELLED
// - 使用事务 + 行锁，避免并发下出现“同时支付确认/同时取消”的竞态
func (s *Service) CancelOrder(ctx context.Context, req *kitexuser.CancelOrderReq) (*kitexuser.BaseResp, error) {
	if strings.TrimSpace(req.UserId) == "" || strings.TrimSpace(req.OrderId) == "" {
		return baseResp(400, "user_id/order_id不能为空"), nil
	}

	err := db.MysqlDB.Transaction(func(tx *gorm.DB) error {
		var o model2.OrderInfo
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("order_id = ? AND user_id = ?", req.OrderId, req.UserId).First(&o).Error; err != nil {
			return err
		}
		if o.OrderStatus != "PENDING_PAY" && o.OrderStatus != "PAYING" {
			return fmt.Errorf("当前订单状态不允许取消")
		}
		if err := tx.Model(&model2.SeatSegmentOccupancy{}).
			Where("order_id = ? AND status = ?", o.ID, "LOCKED").
			Updates(map[string]interface{}{"status": "CANCELLED", "lock_expire_time": sql.NullTime{Valid: false}}).Error; err != nil {
			return err
		}
		if err := tx.Model(&model2.OrderInfo{}).Where("order_id = ?", o.ID).Updates(map[string]interface{}{"order_status": "CANCELLED"}).Error; err != nil {
			return err
		}
		writeOrderAuditLog(tx, o.ID, "CANCEL_ORDER", req.UserId, o.OrderStatus, "CANCELLED", nil)
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return baseResp(404, "订单不存在"), nil
		}
		return baseResp(400, err.Error()), nil
	}
	return baseResp(200, "success"), nil
}

// RefundOrder 退票用例：将 SOLD 的区间占用取消，并根据规则计算退款金额。
//
// 退款规则来源：
// - 优先从 TicketRuleConfig（rule_type=REFUND_FEE）读取可配置的阶梯手续费
// - 找不到配置则使用内置默认规则（见 loadRefundFeeRule）
//
// 关键约束：
// - 仅允许 ISSUED -> REFUNDED
// - 已发车不允许退票（如果能查到对应站点的 departure_time）
// - 在同一事务内：
//   1) SOLD 占用 -> CANCELLED（释放库存）
//   2) 更新订单退款字段与状态
func (s *Service) RefundOrder(ctx context.Context, req *kitexuser.RefundOrderReq) (*kitexuser.RefundOrderResp, error) {
	if strings.TrimSpace(req.UserId) == "" || strings.TrimSpace(req.OrderId) == "" {
		return &kitexuser.RefundOrderResp{BaseResp: baseResp(400, "user_id/order_id不能为空")}, nil
	}

	var refundAmount float64
	var outStatus string
	err := db.MysqlDB.Transaction(func(tx *gorm.DB) error {
		var o model2.OrderInfo
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("order_id = ? AND user_id = ?", req.OrderId, req.UserId).First(&o).Error; err != nil {
			return err
		}
		if o.OrderStatus != "ISSUED" {
			return fmt.Errorf("订单状态不允许退票")
		}
		var depStop model2.TrainStationPass
		if err := tx.Where("train_id = ? AND station_name = ?", o.TrainID, o.DepartureStation).First(&depStop).Error; err == nil {
			if depStop.DepartureTime.Valid {
				if time.Now().After(depStop.DepartureTime.Time) {
					return fmt.Errorf("已发车，无法退票")
				}
				hoursBefore := depStop.DepartureTime.Time.Sub(time.Now()).Hours()
				rule := loadRefundFeeRule("")
				rate, ok := refundFeeRate(rule, hoursBefore)
				if !ok {
					return fmt.Errorf("距离发车过近，不允许退票")
				}
				refundAmount = shared.Money2(o.TotalAmount * (1 - rate))
			} else {
				rule := loadRefundFeeRule("")
				rate, _ := refundFeeRate(rule, 24)
				refundAmount = shared.Money2(o.TotalAmount * (1 - rate))
			}
		} else {
			refundAmount = shared.Money2(o.TotalAmount * 0.90)
		}
		if refundAmount < 0 {
			refundAmount = 0
		}

		if err := tx.Model(&model2.SeatSegmentOccupancy{}).
			Where("order_id = ? AND status = ?", o.ID, "SOLD").
			Updates(map[string]any{"status": "CANCELLED"}).Error; err != nil {
			return err
		}
		if err := tx.Model(&model2.OrderInfo{}).Where("order_id = ?", o.ID).Updates(map[string]any{
			"order_status":  "REFUNDED",
			"refund_amount": refundAmount,
			"refund_time":   nullTime(time.Now()),
			"refund_status": "REFUNDED",
		}).Error; err != nil {
			return err
		}
		writeOrderAuditLog(tx, o.ID, "REFUND_ORDER", req.UserId, "ISSUED", "REFUNDED", map[string]any{"reason": req.GetReason(), "refund_amount": refundAmount})
		outStatus = "REFUNDED"
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &kitexuser.RefundOrderResp{BaseResp: baseResp(404, "订单不存在")}, nil
		}
		return &kitexuser.RefundOrderResp{BaseResp: baseResp(400, err.Error())}, nil
	}
	return &kitexuser.RefundOrderResp{BaseResp: baseResp(200, "success"), RefundAmount: shared.Money2(refundAmount), RefundStatus: "REFUNDED", OrderStatus: outStatus}, nil
}

type refundTier struct {
	HoursBefore int     `json:"hours_before"`
	FeeRate     float64 `json:"fee_rate"`
}

type refundFeeRule struct {
	Tiers      []refundTier `json:"tiers"`
	DefaultFee float64      `json:"default_fee"`
	NoRefundLT float64      `json:"no_refund_lt_hours"`
}

func loadRefundFeeRule(trainType string) refundFeeRule {
	var cfg model2.TicketRuleConfig
	err := db.ReadDB().Where("rule_type = ? AND status = ?", "REFUND_FEE", "ENABLED").Order("version DESC").First(&cfg).Error
	if err == nil && cfg.RuleParam != nil {
		var rule refundFeeRule
		if jsonErr := json.Unmarshal([]byte(*cfg.RuleParam), &rule); jsonErr == nil && len(rule.Tiers) > 0 {
			sort.Slice(rule.Tiers, func(i, j int) bool { return rule.Tiers[i].HoursBefore > rule.Tiers[j].HoursBefore })
			return rule
		}
	}
	return refundFeeRule{
		Tiers: []refundTier{
			{HoursBefore: 48, FeeRate: 0.05},
			{HoursBefore: 24, FeeRate: 0.10},
			{HoursBefore: 2, FeeRate: 0.20},
			{HoursBefore: 0, FeeRate: 0.50},
		},
		DefaultFee: 0.20,
		NoRefundLT: 0,
	}
}

func refundFeeRate(rule refundFeeRule, hoursBefore float64) (float64, bool) {
	if rule.NoRefundLT > 0 && hoursBefore < rule.NoRefundLT {
		return 0, false
	}
	for _, t := range rule.Tiers {
		if hoursBefore >= float64(t.HoursBefore) {
			return t.FeeRate, true
		}
	}
	return rule.DefaultFee, true
}
