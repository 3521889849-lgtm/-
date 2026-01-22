// Package handler 实现客服系统RPC服务的处理逻辑
// 包含客服管理、排班管理、会话管理、快捷回复等核心业务功能
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"example_shop/service/customer/crypto"             // 消息加密模块
	"example_shop/service/customer/dal"                // 数据访问层
	"example_shop/service/customer/kitex_gen/customer" // Kitex生成的客服服务接口定义
	"example_shop/service/customer/model"              // 数据模型定义
	"example_shop/service/customer/nlp"                // NLP分类模块

	"gorm.io/gorm" // ORM框架
)

// CustomerServiceHandler 客服服务处理器
// 实现了customer.CustomerService接口的所有方法
type CustomerServiceHandler struct{}

// 缓存键定义
const (
	cacheKeyShiftConfigAll  = "customer:shift_config:all:v1"  // 班次配置缓存键
	cacheKeyConvCategoryAll = "customer:conv_category:all:v1" // 会话分类缓存键
	cacheKeyConvTagAll      = "customer:conv_tag:all:v1"      // 会话标签缓存键

	// 缓存过期时间配置
	cacheTTLShiftConfig       = 30 * time.Minute // 班次配置缓存30分钟
	cacheTTLConvCategory      = 60 * time.Minute // 会话分类缓存60分钟
	cacheTTLConvTag           = 60 * time.Minute // 会话标签缓存60分钟
	cacheTTLQuickReplyList    = 5 * time.Minute  // 快捷回复列表缓存5分钟
	cacheTTLConversationLists = 10 * time.Second // 会话列表缓存10秒
)

// NewCustomerServiceHandler 创建客服服务处理器实例
func NewCustomerServiceHandler() *CustomerServiceHandler {
	return &CustomerServiceHandler{}
}

// getShiftConfigAllCached 获取所有班次配置（带缓存）
// 优先从缓存读取，缓存未命中时从数据库查询并写入缓存
// 返回值: 班次配置列表, 错误信息
func getShiftConfigAllCached(ctx context.Context) ([]*customer.ShiftConfig, error) {
	// 尝试从缓存获取
	var cached []*customer.ShiftConfig
	if ok, _ := dal.CacheGetJSON(ctx, cacheKeyShiftConfigAll, &cached); ok && cached != nil {
		return cached, nil
	}

	// 缓存未命中，从数据库查询
	var shifts []model.ShiftConfig
	if err := dal.DB.WithContext(ctx).Model(&model.ShiftConfig{}).Order("shift_id asc").Find(&shifts).Error; err != nil {
		return nil, err
	}

	// 转换为API响应格式
	out := make([]*customer.ShiftConfig, 0, len(shifts))
	for _, s := range shifts {
		out = append(out, &customer.ShiftConfig{
			ShiftId:   s.ShiftID,
			ShiftName: s.ShiftName,
			StartTime: normalizeTimeForAPI(s.StartTime),
			EndTime:   normalizeTimeForAPI(s.EndTime),
			MinStaff:  int32(s.MinStaff),
			IsHoliday: s.IsHoliday,
			CreateBy:  s.CreateBy,
		})
	}

	// 写入缓存
	_ = dal.CacheSetJSON(ctx, cacheKeyShiftConfigAll, out, cacheTTLShiftConfig)
	return out, nil
}

// getConvCategoryAllCached 获取所有会话分类（带缓存）
// 优先从缓存读取，缓存未命中时从数据库查询并写入缓存
// 返回值: 会话分类列表, 错误信息
func getConvCategoryAllCached(ctx context.Context) ([]*customer.ConvCategory, error) {
	// 尝试从缓存获取
	var cached []*customer.ConvCategory
	if ok, _ := dal.CacheGetJSON(ctx, cacheKeyConvCategoryAll, &cached); ok && cached != nil {
		return cached, nil
	}

	// 缓存未命中，从数据库查询（按排序号和分类ID升序）
	var cats []model.ConvCategory
	if err := dal.DB.WithContext(ctx).Model(&model.ConvCategory{}).Order("sort_no asc").Order("category_id asc").Find(&cats).Error; err != nil {
		return nil, err
	}

	// 转换为API响应格式
	out := make([]*customer.ConvCategory, 0, len(cats))
	for _, c := range cats {
		out = append(out, &customer.ConvCategory{
			CategoryId:   c.CategoryID,
			CategoryName: c.CategoryName,
			SortNo:       int32(c.SortNo),
			CreateBy:     c.CreateBy,
		})
	}

	// 写入缓存
	_ = dal.CacheSetJSON(ctx, cacheKeyConvCategoryAll, out, cacheTTLConvCategory)
	return out, nil
}

// GetCustomerService 获取客服信息
func (h *CustomerServiceHandler) GetCustomerService(ctx context.Context, req *customer.GetCustomerServiceReq) (*customer.GetCustomerServiceResp, error) {
	var cs model.CustomerService
	err := dal.DB.Where("cs_id = ?", req.CsId).First(&cs).Error
	if err != nil {
		return &customer.GetCustomerServiceResp{
			BaseResp: &customer.BaseResp{
				Code: 500,
				Msg:  "查询失败: " + err.Error(),
			},
		}, nil
	}

	return &customer.GetCustomerServiceResp{
		BaseResp: &customer.BaseResp{
			Code: 0,
			Msg:  "success",
		},
		CustomerService: &customer.CustomerAgent{
			CsId:          cs.CsID,
			CsName:        cs.CsName,
			DeptId:        cs.DeptID,
			TeamId:        cs.TeamID,
			SkillTags:     cs.SkillTags,
			Status:        cs.Status,
			CurrentStatus: cs.CurrentStatus,
		},
	}, nil
}

// ListCustomerService 查询客服列表
func (h *CustomerServiceHandler) ListCustomerService(ctx context.Context, req *customer.ListCustomerServiceReq) (*customer.ListCustomerServiceResp, error) {
	var csList []model.CustomerService
	var total int64

	query := dal.DB.Model(&model.CustomerService{})
	if req.DeptId != "" {
		query = query.Where("dept_id = ?", req.DeptId)
	}

	err := query.Count(&total).Error
	if err != nil {
		return &customer.ListCustomerServiceResp{
			BaseResp: &customer.BaseResp{
				Code: 500,
				Msg:  "查询失败: " + err.Error(),
			},
		}, nil
	}

	offset := (req.Page - 1) * req.PageSize
	err = query.Offset(int(offset)).Limit(int(req.PageSize)).Find(&csList).Error
	if err != nil {
		return &customer.ListCustomerServiceResp{
			BaseResp: &customer.BaseResp{
				Code: 500,
				Msg:  "查询失败: " + err.Error(),
			},
		}, nil
	}

	var result []*customer.CustomerAgent
	for _, cs := range csList {
		result = append(result, &customer.CustomerAgent{
			CsId:          cs.CsID,
			CsName:        cs.CsName,
			DeptId:        cs.DeptID,
			TeamId:        cs.TeamID,
			SkillTags:     cs.SkillTags,
			Status:        cs.Status,
			CurrentStatus: cs.CurrentStatus,
		})
	}

	return &customer.ListCustomerServiceResp{
		BaseResp: &customer.BaseResp{
			Code: 0,
			Msg:  "success",
		},
		CustomerServices: result,
		Total:            total,
	}, nil
}

// CreateShiftConfig 创建班次配置
// 用于新增客服排班的班次模板，如早班、晚班、节假日班等
// 参数: 班次名称、开始时间、结束时间、最小人数、是否节假日
// 返回: 新创建的班次ID
func (h *CustomerServiceHandler) CreateShiftConfig(ctx context.Context, req *customer.CreateShiftConfigReq) (*customer.CreateShiftConfigResp, error) {
	// 参数校验
	if req == nil || req.Shift == nil {
		return &customer.CreateShiftConfigResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "shift is required"},
		}, nil
	}
	if strings.TrimSpace(req.Shift.ShiftName) == "" {
		return &customer.CreateShiftConfigResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "shift_name is required"},
		}, nil
	}
	if strings.TrimSpace(req.Shift.StartTime) == "" || strings.TrimSpace(req.Shift.EndTime) == "" {
		return &customer.CreateShiftConfigResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "start_time and end_time are required"},
		}, nil
	}
	// 校验时间范围有效性
	if !isValidShiftTimeRange(req.Shift.StartTime, req.Shift.EndTime) {
		return &customer.CreateShiftConfigResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "invalid time range"},
		}, nil
	}

	// 构建班次配置模型
	now := time.Now()
	shift := model.ShiftConfig{
		ShiftName:  strings.TrimSpace(req.Shift.ShiftName),
		StartTime:  normalizeTimeForDB(req.Shift.StartTime),
		EndTime:    normalizeTimeForDB(req.Shift.EndTime),
		MinStaff:   int(req.Shift.MinStaff),
		IsHoliday:  req.Shift.IsHoliday,
		CreateTime: now,
		UpdateTime: now,
		CreateBy:   strings.TrimSpace(req.Shift.CreateBy),
	}
	if shift.CreateBy == "" {
		shift.CreateBy = "ADMIN" // 默认创建者为ADMIN
	}

	// 写入数据库
	if err := dal.DB.Create(&shift).Error; err != nil {
		return &customer.CreateShiftConfigResp{
			BaseResp: &customer.BaseResp{Code: 500, Msg: "create failed: " + err.Error()},
		}, nil
	}

	// 清除班次配置缓存
	_ = dal.CacheDel(ctx, cacheKeyShiftConfigAll)

	return &customer.CreateShiftConfigResp{
		BaseResp: &customer.BaseResp{Code: 0, Msg: "success"},
		ShiftId:  shift.ShiftID,
	}, nil
}

// ListShiftConfig 查询班次配置列表
// 支持按是否节假日和班次名称过滤
// 从缓存获取全量数据后在内存中筛选
func (h *CustomerServiceHandler) ListShiftConfig(ctx context.Context, req *customer.ListShiftConfigReq) (*customer.ListShiftConfigResp, error) {
	// 解析筛选条件
	isHoliday := int8(-1) // -1表示不筛选
	shiftName := ""
	if req != nil {
		if req.IsHoliday == 0 || req.IsHoliday == 1 {
			isHoliday = req.IsHoliday
		}
		shiftName = strings.TrimSpace(req.ShiftName)
	}

	// 获取全量班次配置（带缓存）
	all, err := getShiftConfigAllCached(ctx)
	if err != nil {
		return &customer.ListShiftConfigResp{
			BaseResp: &customer.BaseResp{Code: 500, Msg: "query failed: " + err.Error()},
		}, nil
	}

	// 内存过滤
	respShifts := make([]*customer.ShiftConfig, 0, len(all))
	for _, s := range all {
		// 按节假日标记筛选
		if isHoliday == 0 || isHoliday == 1 {
			if s.IsHoliday != isHoliday {
				continue
			}
		}
		// 按班次名称模糊匹配
		if shiftName != "" && !strings.Contains(s.ShiftName, shiftName) {
			continue
		}
		respShifts = append(respShifts, s)
	}

	return &customer.ListShiftConfigResp{
		BaseResp: &customer.BaseResp{Code: 0, Msg: "success"},
		Shifts:   respShifts,
		Total:    int64(len(respShifts)),
	}, nil
}

// UpdateShiftConfig 更新班次配置
// 根据班次ID更新班次的各项属性
// 更新成功后自动清除班次配置缓存
func (h *CustomerServiceHandler) UpdateShiftConfig(ctx context.Context, req *customer.UpdateShiftConfigReq) (*customer.UpdateShiftConfigResp, error) {
	// 参数校验
	if req == nil || req.Shift == nil {
		return &customer.UpdateShiftConfigResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "shift is required"}}, nil
	}
	if req.Shift.ShiftId <= 0 {
		return &customer.UpdateShiftConfigResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "shift_id is required"}}, nil
	}
	if strings.TrimSpace(req.Shift.ShiftName) == "" {
		return &customer.UpdateShiftConfigResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "shift_name is required"}}, nil
	}
	if strings.TrimSpace(req.Shift.StartTime) == "" || strings.TrimSpace(req.Shift.EndTime) == "" {
		return &customer.UpdateShiftConfigResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "start_time and end_time are required"}}, nil
	}
	if req.Shift.MinStaff < 0 {
		return &customer.UpdateShiftConfigResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "min_staff must be >= 0"}}, nil
	}
	if req.Shift.IsHoliday != 0 && req.Shift.IsHoliday != 1 {
		return &customer.UpdateShiftConfigResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "is_holiday must be 0 or 1"}}, nil
	}
	if !isValidShiftTimeRange(req.Shift.StartTime, req.Shift.EndTime) {
		return &customer.UpdateShiftConfigResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "invalid time range"}}, nil
	}

	// 查询现有班次配置
	var cur model.ShiftConfig
	if err := dal.DB.WithContext(ctx).Where("shift_id = ?", req.Shift.ShiftId).First(&cur).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &customer.UpdateShiftConfigResp{BaseResp: &customer.BaseResp{Code: 404, Msg: "shift not found"}}, nil
		}
		return &customer.UpdateShiftConfigResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "query failed: " + err.Error()}}, nil
	}

	// 更新字段
	cur.ShiftName = strings.TrimSpace(req.Shift.ShiftName)
	cur.StartTime = normalizeTimeForDB(req.Shift.StartTime)
	cur.EndTime = normalizeTimeForDB(req.Shift.EndTime)
	cur.MinStaff = int(req.Shift.MinStaff)
	cur.IsHoliday = req.Shift.IsHoliday
	cur.UpdateTime = time.Now()

	// 保存到数据库
	if err := dal.DB.WithContext(ctx).Save(&cur).Error; err != nil {
		return &customer.UpdateShiftConfigResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "update failed: " + err.Error()}}, nil
	}

	// 清除班次配置缓存
	_ = dal.CacheDel(ctx, cacheKeyShiftConfigAll)

	return &customer.UpdateShiftConfigResp{BaseResp: &customer.BaseResp{Code: 0, Msg: "success"}}, nil
}

// DeleteShiftConfig 删除班次配置
// 删除前会检查该班次是否已被排班使用，如已使用则禁止删除（返回409冲突）
// 删除成功后自动清除班次配置缓存
func (h *CustomerServiceHandler) DeleteShiftConfig(ctx context.Context, req *customer.DeleteShiftConfigReq) (*customer.DeleteShiftConfigResp, error) {
	// 参数校验
	if req == nil || req.ShiftId <= 0 {
		return &customer.DeleteShiftConfigResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "shift_id is required"}}, nil
	}

	// 检查班次是否已被排班使用
	var used int64
	if err := dal.DB.WithContext(ctx).Model(&model.Schedule{}).Where("shift_id = ?", req.ShiftId).Count(&used).Error; err != nil {
		return &customer.DeleteShiftConfigResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "check failed: " + err.Error()}}, nil
	}
	if used > 0 {
		return &customer.DeleteShiftConfigResp{BaseResp: &customer.BaseResp{Code: 409, Msg: "shift is used by schedule"}}, nil
	}

	// 执行删除
	res := dal.DB.WithContext(ctx).Where("shift_id = ?", req.ShiftId).Delete(&model.ShiftConfig{})
	if res.Error != nil {
		return &customer.DeleteShiftConfigResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "delete failed: " + res.Error.Error()}}, nil
	}
	if res.RowsAffected == 0 {
		return &customer.DeleteShiftConfigResp{BaseResp: &customer.BaseResp{Code: 404, Msg: "shift not found"}}, nil
	}

	// 清除班次配置缓存
	_ = dal.CacheDel(ctx, cacheKeyShiftConfigAll)

	return &customer.DeleteShiftConfigResp{BaseResp: &customer.BaseResp{Code: 0, Msg: "success"}}, nil
}

// AssignSchedule 批量分配排班
// 将指定班次分配给多个客服人员的某一天
// 如果客服在该日期已有排班，会返回冲突的客服ID列表
// 分配成功后会更新客服的当前状态为工作中
func (h *CustomerServiceHandler) AssignSchedule(ctx context.Context, req *customer.AssignScheduleReq) (*customer.AssignScheduleResp, error) {
	// 参数校验
	if req == nil {
		return &customer.AssignScheduleResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "req is required"}}, nil
	}
	if strings.TrimSpace(req.ScheduleDate) == "" || req.ShiftId <= 0 || len(req.CsIds) == 0 {
		return &customer.AssignScheduleResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "schedule_date, shift_id and cs_ids are required"}}, nil
	}

	var conflicts []string // 冲突的客服ID列表
	err := dal.DB.Transaction(func(tx *gorm.DB) error {
		// 验证班次是否存在
		var shift model.ShiftConfig
		if err := tx.Where("shift_id = ?", req.ShiftId).First(&shift).Error; err != nil {
			return fmt.Errorf("shift not found: %w", err)
		}

		// 检查是否有排班冲突（同一天同一客服已有正常排班）
		if err := tx.Model(&model.Schedule{}).
			Where("schedule_date = ? AND status = 0 AND cs_id IN ?", req.ScheduleDate, req.CsIds).
			Distinct("cs_id").
			Pluck("cs_id", &conflicts).Error; err != nil {
			return err
		}
		if len(conflicts) > 0 {
			return nil // 有冲突，不继续执行
		}

		// 批量创建排班记录
		now := time.Now()
		schedules := make([]model.Schedule, 0, len(req.CsIds))
		for _, csID := range req.CsIds {
			csID = strings.TrimSpace(csID)
			if csID == "" {
				continue
			}
			schedules = append(schedules, model.Schedule{
				CsID:         csID,
				ShiftID:      req.ShiftId,
				ScheduleDate: req.ScheduleDate,
				Status:       0, // 0=正常排班
				CreateTime:   now,
				UpdateTime:   now,
			})
		}
		if len(schedules) == 0 {
			return nil
		}
		if err := tx.Create(&schedules).Error; err != nil {
			return err
		}
		// 更新客服当前状态为工作中(1)
		if err := tx.Model(&model.CustomerService{}).Where("cs_id IN ?", req.CsIds).Update("current_status", 1).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return &customer.AssignScheduleResp{BaseResp: &customer.BaseResp{Code: 500, Msg: err.Error()}}, nil
	}
	// 返回冲突信息
	if len(conflicts) > 0 {
		return &customer.AssignScheduleResp{
			BaseResp:      &customer.BaseResp{Code: 409, Msg: "schedule conflict"},
			ConflictCsIds: conflicts,
		}, nil
	}
	return &customer.AssignScheduleResp{BaseResp: &customer.BaseResp{Code: 0, Msg: "success"}}, nil
}

// ApplyLeaveTransfer 申请请假或换班
// apply_type: 0=请假, 1=换班
// 请假时只需指定日期和班次，换班还需指定目标客服ID
// 申请提交后需等待审批
func (h *CustomerServiceHandler) ApplyLeaveTransfer(ctx context.Context, req *customer.ApplyLeaveTransferReq) (*customer.ApplyLeaveTransferResp, error) {
	// 参数校验
	if req == nil {
		return &customer.ApplyLeaveTransferResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "req is required"}}, nil
	}
	if strings.TrimSpace(req.CsId) == "" || strings.TrimSpace(req.TargetDate) == "" || req.ShiftId <= 0 {
		return &customer.ApplyLeaveTransferResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "cs_id, target_date and shift_id are required"}}, nil
	}
	if req.ApplyType != 0 && req.ApplyType != 1 {
		return &customer.ApplyLeaveTransferResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "invalid apply_type"}}, nil
	}
	// 换班时必须指定目标客服
	if req.ApplyType == 1 && strings.TrimSpace(req.TargetCsId) == "" {
		return &customer.ApplyLeaveTransferResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "target_cs_id is required for transfer"}}, nil
	}

	// 创建申请记录
	now := time.Now()
	apply := model.LeaveTransfer{
		CsID:           strings.TrimSpace(req.CsId),
		ApplyType:      req.ApplyType,
		TargetDate:     req.TargetDate,
		ShiftID:        req.ShiftId,
		TargetCsID:     strings.TrimSpace(req.TargetCsId),
		ApprovalStatus: 0, // 0=待审批
		Reason:         strings.TrimSpace(req.Reason),
		CreateTime:     now,
		UpdateTime:     now,
	}

	if err := dal.DB.Create(&apply).Error; err != nil {
		return &customer.ApplyLeaveTransferResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "create failed: " + err.Error()}}, nil
	}
	return &customer.ApplyLeaveTransferResp{BaseResp: &customer.BaseResp{Code: 0, Msg: "success"}, ApplyId: apply.ApplyID}, nil
}

// ApproveLeaveTransfer 审批请假/换班申请
// approval_status: 1=通过, 2=拒绝
// 通过后会自动更新排班表和客服状态
// 请假通过: 标记客服为请假状态
// 换班通过: 创建目标客服的排班记录，更新原客服状态
func (h *CustomerServiceHandler) ApproveLeaveTransfer(ctx context.Context, req *customer.ApproveLeaveTransferReq) (*customer.ApproveLeaveTransferResp, error) {
	// 参数校验
	if req == nil || req.ApplyId <= 0 {
		return &customer.ApproveLeaveTransferResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "apply_id is required"}}, nil
	}
	if req.ApprovalStatus != 1 && req.ApprovalStatus != 2 {
		return &customer.ApproveLeaveTransferResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "invalid approval_status"}}, nil
	}

	err := dal.DB.Transaction(func(tx *gorm.DB) error {
		// 查询申请记录
		var apply model.LeaveTransfer
		if err := tx.Where("apply_id = ?", req.ApplyId).First(&apply).Error; err != nil {
			return err
		}

		// 更新审批状态
		now := time.Now()
		update := map[string]interface{}{
			"approval_status": req.ApprovalStatus,
			"approver_id":     strings.TrimSpace(req.ApproverId),
			"approval_time":   now,
			"update_time":     now,
		}
		if err := tx.Model(&model.LeaveTransfer{}).Where("apply_id = ?", req.ApplyId).Updates(update).Error; err != nil {
			return err
		}

		// 拒绝则不需要后续处理
		if req.ApprovalStatus != 1 {
			return nil
		}

		// 审批通过 - 处理请假(apply_type=0)
		if apply.ApplyType == 0 {
			var s model.Schedule
			findErr := tx.Where("cs_id = ? AND schedule_date = ? AND shift_id = ?", apply.CsID, apply.TargetDate, apply.ShiftID).First(&s).Error
			if findErr != nil {
				if !errorsIsRecordNotFound(findErr) {
					return findErr
				}
				// 排班记录不存在，创建新的请假记录
				s = model.Schedule{
					CsID:         apply.CsID,
					ShiftID:      apply.ShiftID,
					ScheduleDate: apply.TargetDate,
					Status:       1, // 1=请假
					CreateTime:   now,
					UpdateTime:   now,
				}
				if err := tx.Create(&s).Error; err != nil {
					return err
				}
			} else {
				// 更新现有排班记录为请假状态
				if err := tx.Model(&model.Schedule{}).Where("schedule_id = ?", s.ScheduleID).Updates(map[string]interface{}{
					"status":        1,
					"replace_cs_id": "",
					"update_time":   now,
				}).Error; err != nil {
					return err
				}
			}
			// 更新客服状态为请假(2)
			if err := tx.Model(&model.CustomerService{}).Where("cs_id = ?", apply.CsID).Update("current_status", 2).Error; err != nil {
				return err
			}
			return nil
		}

		// 审批通过 - 处理换班(apply_type=1)
		target := strings.TrimSpace(apply.TargetCsID)
		if target == "" {
			return fmt.Errorf("target_cs_id is required")
		}

		// 检查目标客服是否有冲突排班
		var conflictCount int64
		if err := tx.Model(&model.Schedule{}).Where("cs_id = ? AND schedule_date = ? AND status = 0", target, apply.TargetDate).Count(&conflictCount).Error; err != nil {
			return err
		}
		if conflictCount > 0 {
			return fmt.Errorf("target customer has conflict schedule")
		}

		// 更新原客服的排班记录为换班状态
		var s model.Schedule
		findErr := tx.Where("cs_id = ? AND schedule_date = ? AND shift_id = ?", apply.CsID, apply.TargetDate, apply.ShiftID).First(&s).Error
		if findErr != nil {
			if !errorsIsRecordNotFound(findErr) {
				return findErr
			}
			// 创建新的换班记录
			s = model.Schedule{
				CsID:         apply.CsID,
				ShiftID:      apply.ShiftID,
				ScheduleDate: apply.TargetDate,
				Status:       2, // 2=换班
				ReplaceCsID:  target,
				CreateTime:   now,
				UpdateTime:   now,
			}
			if err := tx.Create(&s).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&model.Schedule{}).Where("schedule_id = ?", s.ScheduleID).Updates(map[string]interface{}{
				"status":        2,
				"replace_cs_id": target,
				"update_time":   now,
			}).Error; err != nil {
				return err
			}
		}

		// 为目标客服创建排班记录
		other := model.Schedule{
			CsID:         target,
			ShiftID:      apply.ShiftID,
			ScheduleDate: apply.TargetDate,
			Status:       0, // 0=正常排班
			ReplaceCsID:  apply.CsID,
			CreateTime:   now,
			UpdateTime:   now,
		}
		if err := tx.Create(&other).Error; err != nil {
			return err
		}
		// 更新两个客服的状态为工作中(1)
		if err := tx.Model(&model.CustomerService{}).Where("cs_id IN ?", []string{apply.CsID, target}).Update("current_status", 1).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		code := int32(500)
		msg := err.Error()
		if strings.Contains(msg, "conflict") {
			code = 409
		}
		return &customer.ApproveLeaveTransferResp{BaseResp: &customer.BaseResp{Code: code, Msg: msg}}, nil
	}
	return &customer.ApproveLeaveTransferResp{BaseResp: &customer.BaseResp{Code: 0, Msg: "success"}}, nil
}

// GetLeaveTransfer 获取请假/换班申请详情
// 返回申请的完整信息，包括申请人、目标日期、班次、审批状态等
func (h *CustomerServiceHandler) GetLeaveTransfer(ctx context.Context, req *customer.GetLeaveTransferReq) (*customer.GetLeaveTransferResp, error) {
	// 参数校验
	if req == nil || req.ApplyId <= 0 {
		return &customer.GetLeaveTransferResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "apply_id is required"}}, nil
	}

	// 定义查询结果结构
	type row struct {
		ApplyID        int64
		CsID           string
		CsName         string
		DeptID         string
		TeamID         string
		ApplyType      int8
		TargetDate     string
		ShiftID        int64
		ShiftName      string
		TargetCsID     string
		TargetCsName   string
		ApprovalStatus int8
		ApproverID     string
		ApprovalTime   *time.Time
		Reason         string
		CreateTime     time.Time
	}

	var out row
	// 联表查询获取完整信息
	q := dal.DB.WithContext(ctx).
		Table("t_leave_transfer as lt").
		Select(`
			lt.apply_id as apply_id,
			lt.cs_id as cs_id,
			cs.cs_name as cs_name,
			cs.dept_id as dept_id,
			cs.team_id as team_id,
			lt.apply_type as apply_type,
			lt.target_date as target_date,
			lt.shift_id as shift_id,
			sc.shift_name as shift_name,
			lt.target_cs_id as target_cs_id,
			tcs.cs_name as target_cs_name,
			lt.approval_status as approval_status,
			lt.approver_id as approver_id,
			lt.approval_time as approval_time,
			lt.reason as reason,
			lt.create_time as create_time
		`).
		Joins("LEFT JOIN t_customer_service cs ON cs.cs_id = lt.cs_id").
		Joins("LEFT JOIN t_customer_service tcs ON tcs.cs_id = lt.target_cs_id").
		Joins("LEFT JOIN t_shift_config sc ON sc.shift_id = lt.shift_id").
		Where("lt.apply_id = ?", req.ApplyId).
		Limit(1)

	if err := q.Scan(&out).Error; err != nil {
		return &customer.GetLeaveTransferResp{BaseResp: &customer.BaseResp{Code: 500, Msg: err.Error()}}, nil
	}
	if out.ApplyID == 0 {
		return &customer.GetLeaveTransferResp{BaseResp: &customer.BaseResp{Code: 404, Msg: "not found"}}, nil
	}

	// 格式化审批时间
	approvalTime := ""
	if out.ApprovalTime != nil {
		approvalTime = out.ApprovalTime.Format("2006-01-02 15:04:05")
	}
	return &customer.GetLeaveTransferResp{
		BaseResp: &customer.BaseResp{Code: 0, Msg: "success"},
		Item: &customer.LeaveTransferItem{
			ApplyId:        out.ApplyID,
			CsId:           out.CsID,
			CsName:         out.CsName,
			DeptId:         out.DeptID,
			TeamId:         out.TeamID,
			ApplyType:      out.ApplyType,
			TargetDate:     out.TargetDate,
			ShiftId:        out.ShiftID,
			ShiftName:      out.ShiftName,
			TargetCsId:     out.TargetCsID,
			TargetCsName:   out.TargetCsName,
			ApprovalStatus: out.ApprovalStatus,
			ApproverId:     out.ApproverID,
			ApprovalTime:   approvalTime,
			Reason:         out.Reason,
			CreateTime:     out.CreateTime.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

// ListLeaveTransfer 查询请假/换班申请列表
// 支持按审批状态筛选，支持关键词模糊搜索
// 返回分页的申请列表
func (h *CustomerServiceHandler) ListLeaveTransfer(ctx context.Context, req *customer.ListLeaveTransferReq) (*customer.ListLeaveTransferResp, error) {
	// 解析查询参数
	approvalStatus := int8(-1) // -1表示不筛选
	keyword := ""
	page := int32(1)
	pageSize := int32(20)
	if req != nil {
		approvalStatus = req.ApprovalStatus
		keyword = strings.TrimSpace(req.Keyword)
		if req.Page > 0 {
			page = req.Page
		}
		if req.PageSize > 0 {
			pageSize = req.PageSize
		}
	}
	// 限制每页最大100条
	if pageSize > 100 {
		pageSize = 100
	}
	offset := int((page - 1) * pageSize)

	// 定义查询结果结构
	type row struct {
		ApplyID        int64
		CsID           string
		CsName         string
		DeptID         string
		TeamID         string
		ApplyType      int8
		TargetDate     string
		ShiftID        int64
		ShiftName      string
		TargetCsID     string
		TargetCsName   string
		ApprovalStatus int8
		ApproverID     string
		ApprovalTime   *time.Time
		Reason         string
		CreateTime     time.Time
	}

	// 构建基础查询
	base := func(db *gorm.DB) *gorm.DB {
		q := db.WithContext(ctx).
			Table("t_leave_transfer as lt").
			Joins("LEFT JOIN t_customer_service cs ON cs.cs_id = lt.cs_id").
			Joins("LEFT JOIN t_customer_service tcs ON tcs.cs_id = lt.target_cs_id").
			Joins("LEFT JOIN t_shift_config sc ON sc.shift_id = lt.shift_id")

		// 按审批状态筛选
		if approvalStatus == 0 || approvalStatus == 1 || approvalStatus == 2 {
			q = q.Where("lt.approval_status = ?", approvalStatus)
		}
		// 关键词模糊搜索
		if keyword != "" {
			like := "%" + keyword + "%"
			q = q.Where("(lt.cs_id LIKE ? OR cs.cs_name LIKE ? OR lt.target_cs_id LIKE ? OR tcs.cs_name LIKE ? OR lt.reason LIKE ?)", like, like, like, like, like)
		}
		return q
	}

	// 查询总数
	var total int64
	if err := base(dal.DB).Count(&total).Error; err != nil {
		return &customer.ListLeaveTransferResp{BaseResp: &customer.BaseResp{Code: 500, Msg: err.Error()}}, nil
	}

	// 查询列表数据
	var rows []row
	if err := base(dal.DB).
		Select(`
			lt.apply_id as apply_id,
			lt.cs_id as cs_id,
			cs.cs_name as cs_name,
			cs.dept_id as dept_id,
			cs.team_id as team_id,
			lt.apply_type as apply_type,
			lt.target_date as target_date,
			lt.shift_id as shift_id,
			sc.shift_name as shift_name,
			lt.target_cs_id as target_cs_id,
			tcs.cs_name as target_cs_name,
			lt.approval_status as approval_status,
			lt.approver_id as approver_id,
			lt.approval_time as approval_time,
			lt.reason as reason,
			lt.create_time as create_time
		`).
		Order("lt.apply_id desc").
		Offset(offset).
		Limit(int(pageSize)).
		Scan(&rows).Error; err != nil {
		return &customer.ListLeaveTransferResp{BaseResp: &customer.BaseResp{Code: 500, Msg: err.Error()}}, nil
	}

	// 转换为API响应格式
	items := make([]*customer.LeaveTransferItem, 0, len(rows))
	for _, r := range rows {
		approvalTime := ""
		if r.ApprovalTime != nil {
			approvalTime = r.ApprovalTime.Format("2006-01-02 15:04:05")
		}
		items = append(items, &customer.LeaveTransferItem{
			ApplyId:        r.ApplyID,
			CsId:           r.CsID,
			CsName:         r.CsName,
			DeptId:         r.DeptID,
			TeamId:         r.TeamID,
			ApplyType:      r.ApplyType,
			TargetDate:     r.TargetDate,
			ShiftId:        r.ShiftID,
			ShiftName:      r.ShiftName,
			TargetCsId:     r.TargetCsID,
			TargetCsName:   r.TargetCsName,
			ApprovalStatus: r.ApprovalStatus,
			ApproverId:     r.ApproverID,
			ApprovalTime:   approvalTime,
			Reason:         r.Reason,
			CreateTime:     r.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	return &customer.ListLeaveTransferResp{
		BaseResp: &customer.BaseResp{Code: 0, Msg: "success"},
		Items:    items,
		Total:    total,
	}, nil
}

// assignResult 客服分配结果（内部使用）
type assignResult struct {
	ConvID  string // 会话ID
	CsID    string // 客服ID
	CsName  string // 客服名称
	IsNew   bool   // 是否新创建
	ErrCode int32  // 错误码
	ErrMsg  string // 错误信息
}

// assignCustomerInternal 内部客服分配公共逻辑
// 用于 CreateConversation 和 AssignCustomer 的代码复用
// 参数: userID-用户ID, userNickname-用户昵称, source-来源
// 返回: 分配结果
func assignCustomerInternal(userID, userNickname, source string) *assignResult {
	if source == "" {
		source = "Web"
	}

	// 检查用户是否已有进行中的会话
	var existingConv model.Conversation
	if err := dal.DB.Where("user_id = ? AND status = ?", userID, model.ConvStatusOngoing).First(&existingConv).Error; err == nil {
		// 用户已有进行中会话，返回现有会话信息
		var csInfo model.CustomerService
		dal.DB.Where("cs_id = ?", existingConv.CsID).First(&csInfo)
		return &assignResult{
			ConvID: existingConv.ConvID,
			CsID:   existingConv.CsID,
			CsName: csInfo.CsName,
			IsNew:  false,
		}
	}

	// 获取当前时间
	now := time.Now()
	today := now.Format("2006-01-02")
	currentTime := now.Format("15:04:05")

	// 查询当前时段在岗的客服
	type onDutyCS struct {
		CsID   string
		CsName string
	}
	var onDutyList []onDutyCS

	sql := `
		SELECT DISTINCT cs.cs_id, cs.cs_name
		FROM t_customer_service cs
		INNER JOIN t_schedule s ON cs.cs_id = s.cs_id
		INNER JOIN t_shift_config sf ON s.shift_id = sf.shift_id
		WHERE s.schedule_date = ?
		  AND s.status = 0
		  AND cs.status = 1
		  AND (
		      (sf.start_time <= sf.end_time AND ? >= sf.start_time AND ? <= sf.end_time)
		      OR
		      (sf.start_time > sf.end_time AND (? >= sf.start_time OR ? <= sf.end_time))
		  )
	`
	if err := dal.DB.Raw(sql, today, currentTime, currentTime, currentTime, currentTime).Scan(&onDutyList).Error; err != nil {
		return &assignResult{ErrCode: 500, ErrMsg: "failed to query on-duty staff: " + err.Error()}
	}

	if len(onDutyList) == 0 {
		return &assignResult{ErrCode: 404, ErrMsg: "no customer service available at this time"}
	}

	// 统计各客服当前进行中的会话数
	csIDs := make([]string, len(onDutyList))
	for i, cs := range onDutyList {
		csIDs[i] = cs.CsID
	}

	type convCount struct {
		CsID  string
		Count int64
	}
	var counts []convCount
	dal.DB.Model(&model.Conversation{}).
		Select("cs_id, COUNT(*) as count").
		Where("cs_id IN ? AND status = ?", csIDs, model.ConvStatusOngoing).
		Group("cs_id").
		Scan(&counts)

	countMap := make(map[string]int64)
	for _, c := range counts {
		countMap[c.CsID] = c.Count
	}

	// 选择会话数最少的客服（负载均衡）
	minCount := int64(999999)
	for _, cs := range onDutyList {
		cnt := countMap[cs.CsID]
		if cnt < minCount {
			minCount = cnt
		}
	}

	// 收集所有会话数等于最小值的客服
	var candidates []onDutyCS
	for _, cs := range onDutyList {
		if countMap[cs.CsID] == minCount {
			candidates = append(candidates, cs)
		}
	}

	// 从候选客服中随机选择一个
	var selectedCS onDutyCS
	if len(candidates) == 1 {
		selectedCS = candidates[0]
	} else {
		idx := int(now.UnixNano() % int64(len(candidates)))
		selectedCS = candidates[idx]
	}

	// 创建新会话
	convID := fmt.Sprintf("CONV-%d-%s", now.UnixNano(), userID[:min(8, len(userID))])
	conv := model.Conversation{
		ConvID:         convID,
		UserID:         userID,
		UserNickname:   userNickname,
		CsID:           selectedCS.CsID,
		Source:         source,
		StartTime:      now,
		LastMsgTime:    now,
		Status:         model.ConvStatusOngoing,
		MsgType:        0,
		IsManualAdjust: 0,
		CategoryID:     0,
		Tags:           "",
		IsCore:         0,
		Version:        0,
		CreateTime:     now,
		UpdateTime:     now,
	}

	if err := dal.DB.Create(&conv).Error; err != nil {
		return &assignResult{ErrCode: 500, ErrMsg: "failed to create conversation: " + err.Error()}
	}

	// 插入系统欢迎消息
	welcomeMsg := model.ConvMessage{
		ConvID:       convID,
		SenderType:   2,
		SenderID:     "SYSTEM",
		MsgContent:   fmt.Sprintf("您好，客服 %s 为您服务。请描述您的问题。", selectedCS.CsName),
		IsQuickReply: 0,
		SendTime:     now,
	}
	dal.DB.Create(&welcomeMsg)

	return &assignResult{
		ConvID: convID,
		CsID:   selectedCS.CsID,
		CsName: selectedCS.CsName,
		IsNew:  true,
	}
}

// CreateConversation 创建会话
// 为用户创建新会话或返回已存在的进行中会话
func (h *CustomerServiceHandler) CreateConversation(ctx context.Context, req *customer.CreateConversationReq) (*customer.CreateConversationResp, error) {
	if req == nil || strings.TrimSpace(req.UserId) == "" {
		return &customer.CreateConversationResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "user_id is required"},
		}, nil
	}

	result := assignCustomerInternal(
		strings.TrimSpace(req.UserId),
		strings.TrimSpace(req.UserNickname),
		strings.TrimSpace(req.Source),
	)

	if result.ErrCode != 0 {
		return &customer.CreateConversationResp{
			BaseResp: &customer.BaseResp{Code: result.ErrCode, Msg: result.ErrMsg},
		}, nil
	}

	msg := "success"
	if !result.IsNew {
		msg = "existing conversation found"
	}

	return &customer.CreateConversationResp{
		BaseResp: &customer.BaseResp{Code: 0, Msg: msg},
		ConvId:   result.ConvID,
		CsId:     result.CsID,
		IsNew:    result.IsNew,
	}, nil
}

// TransferConversation 转接会话
// 将会话从当前客服转接给目标客服
// 实现逻辑：
// 1. 校验会话状态（状态机验证）
// 2. 使用乐观锁防止并发冲突
// 3. 记录转接历史和上下文
// 4. 更新会话的客服ID和状态
func (h *CustomerServiceHandler) TransferConversation(ctx context.Context, req *customer.TransferConversationReq) (*customer.TransferConversationResp, error) {
	// 参数校验
	if req == nil {
		return &customer.TransferConversationResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "request is required"},
		}, nil
	}

	convID := strings.TrimSpace(req.ConvId)
	targetCsID := strings.TrimSpace(req.ToCsId)

	if convID == "" || targetCsID == "" {
		return &customer.TransferConversationResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "conv_id and to_cs_id are required"},
		}, nil
	}

	now := time.Now()

	// 使用事务确保数据一致性
	err := dal.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 查询并锁定会话记录
		var conv model.Conversation
		if err := tx.Where("conv_id = ?", convID).First(&conv).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("conversation not found")
			}
			return err
		}

		// 2. 状态机验证：检查是否可以转接
		if !model.CanTransitionTo(conv.Status, model.ConvStatusTransferred) {
			return fmt.Errorf("当前会话状态(%s)不允许转接", model.ConvStatusName(conv.Status))
		}

		// 3. 检查目标客服是否存在且在线
		var targetCS model.CustomerService
		if err := tx.Where("cs_id = ? AND status = 1", targetCsID).First(&targetCS).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("目标客服不存在或不在线")
			}
			return err
		}

		// 不能转接给自己
		if targetCsID == conv.CsID {
			return fmt.Errorf("不能转接给当前客服")
		}

		// 4. 获取原客服信息
		var fromCS model.CustomerService
		tx.Where("cs_id = ?", conv.CsID).First(&fromCS)

		// 5. 使用乐观锁更新会话状态
		oldVersion := conv.Version
		result := tx.Model(&model.Conversation{}).
			Where("conv_id = ? AND version = ?", convID, oldVersion).
			Updates(map[string]interface{}{
				"cs_id":       targetCsID,
				"status":      model.ConvStatusTransferred,
				"version":     oldVersion + 1,
				"update_time": now,
			})

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("会话状态已变更，请重试")
		}

		// 6. 创建转接记录
		transfer := model.ConvTransfer{
			ConvID:         convID,
			FromCsID:       conv.CsID,
			FromCsName:     fromCS.CsName,
			ToCsID:         targetCsID,
			ToCsName:       targetCS.CsName,
			TransferReason: "", // 当前请求中没有这个字段，后续可扩展
			ContextRemark:  "", // 当前请求中没有这个字段，后续可扩展
			TransferTime:   now,
			Status:         1, // 已接受
		}
		if err := tx.Create(&transfer).Error; err != nil {
			return err
		}

		// 7. 插入系统通知消息
		sysMsg := model.ConvMessage{
			ConvID:       convID,
			SenderType:   2, // 系统消息
			SenderID:     "SYSTEM",
			MsgContent:   fmt.Sprintf("会话已转接给客服 %s，请稍等...", targetCS.CsName),
			IsQuickReply: 0,
			SendTime:     now,
		}
		return tx.Create(&sysMsg).Error
	})

	if err != nil {
		return &customer.TransferConversationResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: err.Error()},
		}, nil
	}

	return &customer.TransferConversationResp{
		BaseResp: &customer.BaseResp{Code: 0, Msg: "success"},
	}, nil
}

// EndConversation 结束会话
// 将会话状态置为已结束，计算会话时长
// 实现逻辑：
// 1. 校验会话状态（状态机验证）
// 2. 使用乐观锁防止并发冲突
// 3. 记录结束时间和会话时长
// 4. 自动分类：根据会话消息内容进行NLP智能分类
// 5. 插入系统消息通知
func (h *CustomerServiceHandler) EndConversation(ctx context.Context, req *customer.EndConversationReq) (*customer.EndConversationResp, error) {
	// 参数校验
	if req == nil || strings.TrimSpace(req.ConvId) == "" {
		return &customer.EndConversationResp{
			BaseResp:        &customer.BaseResp{Code: 400, Msg: "conv_id is required"},
			DurationSeconds: 0,
		}, nil
	}

	convID := strings.TrimSpace(req.ConvId)
	endReason := strings.TrimSpace(req.EndReason)
	now := time.Now()
	var durationSeconds int32

	// 使用事务确保数据一致性
	err := dal.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 查询会话记录
		var conv model.Conversation
		if err := tx.Where("conv_id = ?", convID).First(&conv).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("会话不存在")
			}
			return err
		}

		// 2. 状态机验证：检查是否可以结束
		if !model.CanTransitionTo(conv.Status, model.ConvStatusEnded) {
			return fmt.Errorf("当前会话状态(%s)不允许结束", model.ConvStatusName(conv.Status))
		}

		// 3. 计算会话时长
		durationSeconds = int32(now.Sub(conv.StartTime).Seconds())

		// 4. 自动分类：只有未分类(category_id=0)且非手动调整的会话才执行自动分类
		autoClassifiedCategoryID := conv.CategoryID
		if conv.CategoryID == 0 && conv.IsManualAdjust == 0 {
			// 4.1 获取会话消息（用户和客服的消息，最多50条）
			var messages []model.ConvMessage
			tx.Where("conv_id = ? AND sender_type IN (0, 1)", convID).
				Order("send_time asc").
				Limit(50).
				Find(&messages)

			if len(messages) > 0 {
				// 4.2 获取所有消息分类维度
				var categories []model.MsgCategory
				tx.Model(&model.MsgCategory{}).Order("sort_no asc").Find(&categories)

				if len(categories) > 0 {
					// 4.3 创建NLP分类器并添加分类规则
					classifier := nlp.NewClassifier()
					for _, cat := range categories {
						keywords := nlp.ParseKeywordsJSON(cat.Keywords)
						classifier.AddCategory(cat.CategoryID, cat.CategoryName, keywords)
					}

					// 4.4 提取消息内容并执行分类
					var msgContents []string
					for _, m := range messages {
						msgContents = append(msgContents, m.MsgContent)
					}
					allContent := strings.Join(msgContents, " ")
					classifyResult := classifier.Classify(allContent)

					// 4.5 如果置信度足够高（不需要人工确认），则记录分类结果
					if classifyResult.CategoryID > 0 && !classifyResult.NeedManual {
						autoClassifiedCategoryID = classifyResult.CategoryID
					}
				}
			}
		}

		// 5. 使用乐观锁更新会话状态（包含自动分类结果）
		oldVersion := conv.Version
		updateFields := map[string]interface{}{
			"status":      model.ConvStatusEnded,
			"end_time":    now,
			"version":     oldVersion + 1,
			"update_time": now,
		}
		// 如果自动分类成功，更新分类ID
		if autoClassifiedCategoryID > 0 && conv.CategoryID == 0 {
			updateFields["category_id"] = autoClassifiedCategoryID
		}

		result := tx.Model(&model.Conversation{}).
			Where("conv_id = ? AND version = ?", convID, oldVersion).
			Updates(updateFields)

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("会话状态已变更，请重试")
		}

		// 6. 插入系统结束消息
		endMsgContent := "会话已结束，感谢您的咨询。"
		if endReason != "" {
			endMsgContent = fmt.Sprintf("会话已结束（%s），感谢您的咨询。", endReason)
		}
		sysMsg := model.ConvMessage{
			ConvID:       convID,
			SenderType:   2, // 系统消息
			SenderID:     "SYSTEM",
			MsgContent:   endMsgContent,
			IsQuickReply: 0,
			SendTime:     now,
		}
		return tx.Create(&sysMsg).Error
	})

	if err != nil {
		return &customer.EndConversationResp{
			BaseResp:        &customer.BaseResp{Code: 400, Msg: err.Error()},
			DurationSeconds: 0,
		}, nil
	}

	return &customer.EndConversationResp{
		BaseResp:        &customer.BaseResp{Code: 0, Msg: "success"},
		DurationSeconds: durationSeconds,
	}, nil
}

// AbandonConversationResponse 放弃会话响应（内部用）
type AbandonConversationResponse struct {
	Code            int32  // 响应码
	Msg             string // 响应消息
	DurationSeconds int32  // 会话时长(秒)
}

// AbandonConversation 放弃会话（用户中途退出）
// 将会话状态置为已放弃
func (h *CustomerServiceHandler) AbandonConversation(convID string) *AbandonConversationResponse {
	convID = strings.TrimSpace(convID)
	if convID == "" {
		return &AbandonConversationResponse{Code: 400, Msg: "conv_id is required"}
	}

	now := time.Now()
	var durationSeconds int32

	err := dal.DB.Transaction(func(tx *gorm.DB) error {
		var conv model.Conversation
		if err := tx.Where("conv_id = ?", convID).First(&conv).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("会话不存在")
			}
			return err
		}

		// 状态机验证
		if !model.CanTransitionTo(conv.Status, model.ConvStatusAbandoned) {
			return fmt.Errorf("当前会话状态(%s)不允许放弃", model.ConvStatusName(conv.Status))
		}

		durationSeconds = int32(now.Sub(conv.StartTime).Seconds())

		// 使用乐观锁更新
		oldVersion := conv.Version
		result := tx.Model(&model.Conversation{}).
			Where("conv_id = ? AND version = ?", convID, oldVersion).
			Updates(map[string]interface{}{
				"status":      model.ConvStatusAbandoned,
				"end_time":    now,
				"version":     oldVersion + 1,
				"update_time": now,
			})

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("会话状态已变更，请重试")
		}

		return nil
	})

	if err != nil {
		return &AbandonConversationResponse{Code: 400, Msg: err.Error()}
	}

	return &AbandonConversationResponse{
		Code:            0,
		Msg:             "success",
		DurationSeconds: durationSeconds,
	}
}

// ListScheduleGrid 查询排班表格数据（用于类似Excel的排班视图）
func (h *CustomerServiceHandler) ListScheduleGrid(ctx context.Context, req *customer.ListScheduleGridReq) (*customer.ListScheduleGridResp, error) {
	if req == nil {
		return &customer.ListScheduleGridResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "req is required"}}, nil
	}

	startDate := strings.TrimSpace(req.StartDate)
	endDate := strings.TrimSpace(req.EndDate)
	if startDate == "" || endDate == "" {
		return &customer.ListScheduleGridResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "start_date and end_date are required"}}, nil
	}

	st, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return &customer.ListScheduleGridResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "start_date must be YYYY-MM-DD"}}, nil
	}
	et, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return &customer.ListScheduleGridResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "end_date must be YYYY-MM-DD"}}, nil
	}
	if et.Before(st) {
		st, et = et, st
		startDate, endDate = st.Format("2006-01-02"), et.Format("2006-01-02")
	}
	if et.Sub(st) > 62*24*time.Hour {
		et = st.AddDate(0, 0, 62)
		endDate = et.Format("2006-01-02")
	}

	dates := make([]string, 0, 64)
	for d := st; !d.After(et); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format("2006-01-02"))
	}

	var csList []model.CustomerService
	query := dal.DB.WithContext(ctx).Model(&model.CustomerService{})
	if strings.TrimSpace(req.DeptId) != "" {
		query = query.Where("dept_id = ?", strings.TrimSpace(req.DeptId))
	}
	if strings.TrimSpace(req.TeamId) != "" {
		query = query.Where("team_id = ?", strings.TrimSpace(req.TeamId))
	}
	if err := query.Order("cs_id asc").Find(&csList).Error; err != nil {
		return &customer.ListScheduleGridResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "查询客服失败: " + err.Error()}}, nil
	}

	shifts, err := getShiftConfigAllCached(ctx)
	if err != nil {
		return &customer.ListScheduleGridResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "查询班次失败: " + err.Error()}}, nil
	}

	csIDs := make([]string, 0, len(csList))
	for _, c := range csList {
		csIDs = append(csIDs, c.CsID)
	}

	var schedules []model.Schedule
	if len(csIDs) > 0 {
		if err := dal.DB.WithContext(ctx).Model(&model.Schedule{}).
			Where("schedule_date BETWEEN ? AND ? AND cs_id IN ?", startDate, endDate, csIDs).
			Find(&schedules).Error; err != nil {
			return &customer.ListScheduleGridResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "查询排班失败: " + err.Error()}}, nil
		}
	}

	customers := make([]*customer.CustomerAgent, 0, len(csList))
	for _, cs := range csList {
		customers = append(customers, &customer.CustomerAgent{
			CsId:          cs.CsID,
			CsName:        cs.CsName,
			DeptId:        cs.DeptID,
			TeamId:        cs.TeamID,
			SkillTags:     cs.SkillTags,
			Status:        cs.Status,
			CurrentStatus: cs.CurrentStatus,
		})
	}

	cells := make([]*customer.ScheduleCell, 0, len(schedules))
	for _, s := range schedules {
		cells = append(cells, &customer.ScheduleCell{
			CsId:         s.CsID,
			ScheduleDate: normalizeDateForAPI(s.ScheduleDate),
			ShiftId:      s.ShiftID,
			Status:       s.Status,
		})
	}

	return &customer.ListScheduleGridResp{
		BaseResp:  &customer.BaseResp{Code: 0, Msg: "success"},
		Dates:     dates,
		Customers: customers,
		Shifts:    shifts,
		Cells:     cells,
	}, nil
}

// UpsertScheduleCell 设置某个客服某天的班次（用于表格单元格编辑）
func (h *CustomerServiceHandler) UpsertScheduleCell(ctx context.Context, req *customer.UpsertScheduleCellReq) (*customer.UpsertScheduleCellResp, error) {
	if req == nil {
		return &customer.UpsertScheduleCellResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "req is required"}}, nil
	}
	csID := strings.TrimSpace(req.CsId)
	scheduleDate := strings.TrimSpace(req.ScheduleDate)
	if csID == "" || scheduleDate == "" {
		return &customer.UpsertScheduleCellResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "cs_id and schedule_date are required"}}, nil
	}
	if _, err := time.Parse("2006-01-02", scheduleDate); err != nil {
		return &customer.UpsertScheduleCellResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "schedule_date must be YYYY-MM-DD"}}, nil
	}

	shiftID := req.ShiftId
	if shiftID < 0 {
		return &customer.UpsertScheduleCellResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "shift_id must be >= 0"}}, nil
	}

	err := dal.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if shiftID > 0 {
			var shift model.ShiftConfig
			if err := tx.Where("shift_id = ?", shiftID).First(&shift).Error; err != nil {
				return fmt.Errorf("shift not found: %w", err)
			}
		}

		var existing model.Schedule
		findErr := tx.Where("cs_id = ? AND schedule_date = ? AND status = 0", csID, scheduleDate).First(&existing).Error
		if findErr != nil && !errorsIsRecordNotFound(findErr) {
			return findErr
		}

		now := time.Now()
		if shiftID == 0 {
			if errorsIsRecordNotFound(findErr) {
				return nil
			}
			return tx.Where("schedule_id = ?", existing.ScheduleID).Delete(&model.Schedule{}).Error
		}

		if errorsIsRecordNotFound(findErr) {
			s := model.Schedule{
				CsID:         csID,
				ShiftID:      shiftID,
				ScheduleDate: scheduleDate,
				Status:       0,
				CreateTime:   now,
				UpdateTime:   now,
			}
			if err := tx.Create(&s).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&model.Schedule{}).Where("schedule_id = ?", existing.ScheduleID).Updates(map[string]interface{}{
				"shift_id":      shiftID,
				"update_time":   now,
				"replace_cs_id": nil,
				"status":        0,
			}).Error; err != nil {
				return err
			}
		}
		return tx.Model(&model.CustomerService{}).Where("cs_id = ?", csID).Update("current_status", 1).Error
	})
	if err != nil {
		code := int32(500)
		msg := err.Error()
		if strings.Contains(msg, "Duplicate") || strings.Contains(msg, "duplicate") {
			code = 409
		}
		return &customer.UpsertScheduleCellResp{BaseResp: &customer.BaseResp{Code: code, Msg: msg}}, nil
	}
	return &customer.UpsertScheduleCellResp{BaseResp: &customer.BaseResp{Code: 0, Msg: "success"}}, nil
}

// AutoSchedule 自动排班
// 根据日期范围和班次配置，自动为客服人员分配排班
// 规则：
// 1. 获取所有有效客服和班次
// 2. 遍历日期范围，每天轮询分配班次
// 3. 尽量满足每个班次的最少人数要求
// 4. 跳过已有排班的记录
func (h *CustomerServiceHandler) AutoSchedule(ctx context.Context, req *customer.AutoScheduleReq) (*customer.AutoScheduleResp, error) {
	if req == nil {
		return &customer.AutoScheduleResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "req is required"}}, nil
	}

	startDate := strings.TrimSpace(req.StartDate)
	endDate := strings.TrimSpace(req.EndDate)
	if startDate == "" || endDate == "" {
		return &customer.AutoScheduleResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "start_date and end_date are required"}}, nil
	}

	st, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return &customer.AutoScheduleResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "start_date format error"}}, nil
	}
	et, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return &customer.AutoScheduleResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "end_date format error"}}, nil
	}
	if et.Before(st) {
		return &customer.AutoScheduleResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "end_date must be after start_date"}}, nil
	}
	// 限制自动排班最大跨度为31天，防止误操作
	if et.Sub(st) > 31*24*time.Hour {
		return &customer.AutoScheduleResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "max range is 31 days"}}, nil
	}

	// 1. 获取班次配置
	shifts, err := getShiftConfigAllCached(ctx)
	if err != nil || len(shifts) == 0 {
		return &customer.AutoScheduleResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "no shift config found"}}, nil
	}

	// 2. 获取客服列表
	var csList []model.CustomerService
	query := dal.DB.WithContext(ctx).Model(&model.CustomerService{}).Where("status = 1")
	if strings.TrimSpace(req.DeptId) != "" {
		query = query.Where("dept_id = ?", req.DeptId)
	}
	if strings.TrimSpace(req.TeamId) != "" {
		query = query.Where("team_id = ?", req.TeamId)
	}
	if err := query.Find(&csList).Error; err != nil {
		return &customer.AutoScheduleResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "query cs failed"}}, nil
	}
	if len(csList) == 0 {
		return &customer.AutoScheduleResp{BaseResp: &customer.BaseResp{Code: 404, Msg: "no customer service found"}}, nil
	}

	totalScheduled := int64(0)
	now := time.Now()

	// 3. 开始排班
	// 使用简单的轮询算法：每天将客服轮流分配到不同班次
	// 优化点：可以增加"上一个班次"的记录，避免 晚班->早班 的情况
	csIndex := 0

	err = dal.DB.Transaction(func(tx *gorm.DB) error {
		for d := st; !d.After(et); d = d.AddDate(0, 0, 1) {
			dateStr := d.Format("2006-01-02")

			// 获取当天已有排班，避免重复
			var existCsIds []string
			tx.Model(&model.Schedule{}).Where("schedule_date = ?", dateStr).Pluck("cs_id", &existCsIds)
			existMap := make(map[string]bool)
			for _, id := range existCsIds {
				existMap[id] = true
			}

			// 遍历所有班次，尝试填满最小人数
			// 如果还有剩余客服，继续轮询分配

			// 简化版逻辑：
			// 每天，将所有未排班的客服，按顺序分配给各个班次
			for _, cs := range csList {
				if existMap[cs.CsID] {
					continue
				}

				// 简单的轮班规则：(csIndex + dayIndex) % shiftCount
				shift := shifts[csIndex%len(shifts)]
				csIndex++

				s := model.Schedule{
					CsID:         cs.CsID,
					ShiftID:      shift.ShiftId,
					ScheduleDate: dateStr,
					Status:       0,
					CreateTime:   now,
					UpdateTime:   now,
				}
				if err := tx.Create(&s).Error; err != nil {
					return err
				}
				totalScheduled++
			}
		}

		// 更新参与排班的客服状态为工作中
		if totalScheduled > 0 {
			var ids []string
			for _, cs := range csList {
				ids = append(ids, cs.CsID)
			}
			tx.Model(&model.CustomerService{}).Where("cs_id IN ?", ids).Update("current_status", 1)
		}

		return nil
	})

	if err != nil {
		return &customer.AutoScheduleResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "auto schedule failed: " + err.Error()}}, nil
	}

	return &customer.AutoScheduleResp{
		BaseResp:      &customer.BaseResp{Code: 0, Msg: "success"},
		ScheduleCount: totalScheduled,
	}, nil
}

// AssignCustomer 自动分配客服
// 根据当前排班情况和负载均衡策略，为用户分配一个在岗客服
// 分配策略：
// 1. 查询当前时间段排班的客服（基于 t_schedule 和 t_shift_config）
// 2. 排除请假/调班的客服（status != 0）
// 3. 统计各客服进行中会话数，选择负载最低的客服
// 4. 创建新会话并返回
// AssignCustomer 自动分配客服
// 为用户分配客服并创建会话，复用 assignCustomerInternal 公共逻辑
func (h *CustomerServiceHandler) AssignCustomer(ctx context.Context, req *customer.AssignCustomerReq) (*customer.AssignCustomerResp, error) {
	if req == nil || strings.TrimSpace(req.UserId) == "" {
		return &customer.AssignCustomerResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "user_id is required"},
		}, nil
	}

	result := assignCustomerInternal(
		strings.TrimSpace(req.UserId),
		strings.TrimSpace(req.UserNickname),
		strings.TrimSpace(req.Source),
	)

	if result.ErrCode != 0 {
		return &customer.AssignCustomerResp{
			BaseResp: &customer.BaseResp{Code: result.ErrCode, Msg: result.ErrMsg},
		}, nil
	}

	msg := "success"
	if !result.IsNew {
		msg = "existing conversation found"
	}

	return &customer.AssignCustomerResp{
		BaseResp: &customer.BaseResp{Code: 0, Msg: msg},
		CsId:     result.CsID,
		CsName:   result.CsName,
		ConvId:   result.ConvID,
	}, nil
}

// ListConversation 查询会话列表（用于客服会话管理与记录查询）
func (h *CustomerServiceHandler) ListConversation(ctx context.Context, req *customer.ListConversationReq) (*customer.ListConversationResp, error) {
	if req == nil {
		return &customer.ListConversationResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "req is required"}}, nil
	}
	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 200 {
		pageSize = 20
	}

	csID := strings.TrimSpace(req.CsId)
	kw := strings.TrimSpace(req.Keyword)
	cacheKey := fmt.Sprintf(
		"customer:conversation:list:v1:cs=%s:kw=%s:st=%d:p=%d:ps=%d",
		url.QueryEscape(csID),
		url.QueryEscape(kw),
		req.Status,
		page,
		pageSize,
	)
	{
		var cached customer.ListConversationResp
		if ok, _ := dal.CacheGetJSON(ctx, cacheKey, &cached); ok && cached.BaseResp != nil {
			return &cached, nil
		}
	}

	var total int64
	q := dal.DB.WithContext(ctx).Model(&model.Conversation{})
	if csID != "" {
		q = q.Where("cs_id = ?", csID)
	}
	if kw != "" {
		like := "%" + kw + "%"
		q = q.Where("conv_id LIKE ? OR user_id LIKE ? OR user_nickname LIKE ?", like, like, like)
	}
	if req.Status >= 0 {
		q = q.Where("status = ?", req.Status)
	}
	if err := q.Count(&total).Error; err != nil {
		return &customer.ListConversationResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "查询会话失败: " + err.Error()}}, nil
	}

	var convs []model.Conversation
	offset := (page - 1) * pageSize
	if err := q.Order("update_time desc").Offset(int(offset)).Limit(int(pageSize)).Find(&convs).Error; err != nil {
		return &customer.ListConversationResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "查询会话失败: " + err.Error()}}, nil
	}

	convIDs := make([]string, 0, len(convs))
	for _, c := range convs {
		convIDs = append(convIDs, c.ConvID)
	}

	type lastMsgRow struct {
		ConvID     string    `gorm:"column:conv_id"`
		MsgContent string    `gorm:"column:msg_content"`
		SendTime   time.Time `gorm:"column:send_time"`
	}
	lastMsgMap := map[string]lastMsgRow{}
	if len(convIDs) > 0 {
		var rows []lastMsgRow
		raw := `
SELECT m.conv_id, m.msg_content, m.send_time
FROM t_conv_message m
INNER JOIN (
  SELECT conv_id, MAX(send_time) AS max_time
  FROM t_conv_message
  WHERE conv_id IN ?
  GROUP BY conv_id
) t ON m.conv_id = t.conv_id AND m.send_time = t.max_time
`
		if err := dal.DB.WithContext(ctx).Raw(raw, convIDs).Scan(&rows).Error; err == nil {
			for _, r := range rows {
				lastMsgMap[r.ConvID] = r
			}
		}
	}

	catNameMap := map[int64]string{}
	if cats, err := getConvCategoryAllCached(ctx); err == nil {
		for _, c := range cats {
			catNameMap[c.CategoryId] = c.CategoryName
		}
	}

	items := make([]*customer.ConversationItem, 0, len(convs))
	for _, c := range convs {
		lm, ok := lastMsgMap[c.ConvID]
		lastMsg := ""
		lastTime := c.UpdateTime.Format("2006-01-02 15:04:05")
		if ok {
			lastMsg = lm.MsgContent
			lastTime = lm.SendTime.Format("2006-01-02 15:04:05")
		}
		items = append(items, &customer.ConversationItem{
			ConvId:       c.ConvID,
			UserId:       c.UserID,
			UserNickname: c.UserNickname,
			CsId:         c.CsID,
			Source:       c.Source,
			Status:       c.Status,
			LastMsg:      lastMsg,
			LastTime:     lastTime,
			CategoryId:   c.CategoryID,
			CategoryName: catNameMap[c.CategoryID],
			Tags:         c.Tags,
			IsCore:       c.IsCore,
		})
	}

	resp := &customer.ListConversationResp{
		BaseResp:      &customer.BaseResp{Code: 0, Msg: "success"},
		Conversations: items,
		Total:         total,
	}
	_ = dal.CacheSetJSON(ctx, cacheKey, resp, cacheTTLConversationLists)
	return resp, nil
}

// ListConversationHistory 查询历史会话记录
// 查询已关闭或转接的会话，且会话中有用户发送的消息
// 支持按客服ID、关键词和状态筛选
func (h *CustomerServiceHandler) ListConversationHistory(ctx context.Context, req *customer.ListConversationHistoryReq) (*customer.ListConversationResp, error) {
	// 参数校验
	if req == nil {
		return &customer.ListConversationResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "req is required"}}, nil
	}
	// 分页参数处理
	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 200 {
		pageSize = 20
	}

	csID := strings.TrimSpace(req.CsId)
	kw := strings.TrimSpace(req.Keyword)
	// 生成缓存键
	cacheKey := fmt.Sprintf(
		"customer:conversation:history:list:v1:cs=%s:kw=%s:st=%d:p=%d:ps=%d",
		url.QueryEscape(csID),
		url.QueryEscape(kw),
		req.Status,
		page,
		pageSize,
	)
	// 尝试从缓存获取
	{
		var cached customer.ListConversationResp
		if ok, _ := dal.CacheGetJSON(ctx, cacheKey, &cached); ok && cached.BaseResp != nil {
			return &cached, nil
		}
	}

	// 构建查询: 已关闭(0)或转接(2)的会话，且有用户消息
	var total int64
	q := dal.DB.WithContext(ctx).
		Model(&model.Conversation{}).
		Where("status IN ?", []int8{0, 2}).
		Where("EXISTS (SELECT 1 FROM t_conv_message m WHERE m.conv_id = t_conversation.conv_id AND m.sender_type = ?)", int8(0))
	if csID != "" {
		q = q.Where("cs_id = ?", csID)
	}
	if kw != "" {
		like := "%" + kw + "%"
		q = q.Where("conv_id LIKE ? OR user_id LIKE ? OR user_nickname LIKE ?", like, like, like)
	}
	if req.Status == 0 || req.Status == 2 {
		q = q.Where("status = ?", req.Status)
	}
	if err := q.Count(&total).Error; err != nil {
		return &customer.ListConversationResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "查询会话失败: " + err.Error()}}, nil
	}

	// 查询会话列表
	var convs []model.Conversation
	offset := (page - 1) * pageSize
	if err := q.Order("update_time desc").Offset(int(offset)).Limit(int(pageSize)).Find(&convs).Error; err != nil {
		return &customer.ListConversationResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "查询会话失败: " + err.Error()}}, nil
	}

	// 获取会话ID列表
	convIDs := make([]string, 0, len(convs))
	for _, c := range convs {
		convIDs = append(convIDs, c.ConvID)
	}

	// 查询每个会话的最后一条消息
	type lastMsgRow struct {
		ConvID     string    `gorm:"column:conv_id"`
		MsgContent string    `gorm:"column:msg_content"`
		SendTime   time.Time `gorm:"column:send_time"`
	}
	lastMsgMap := map[string]lastMsgRow{}
	if len(convIDs) > 0 {
		var rows []lastMsgRow
		raw := `
SELECT m.conv_id, m.msg_content, m.send_time
FROM t_conv_message m
INNER JOIN (
  SELECT conv_id, MAX(send_time) AS max_time
  FROM t_conv_message
  WHERE conv_id IN ? AND sender_type IN (0, 1)
  GROUP BY conv_id
) t ON m.conv_id = t.conv_id AND m.send_time = t.max_time
`
		if err := dal.DB.WithContext(ctx).Raw(raw, convIDs).Scan(&rows).Error; err == nil {
			for _, r := range rows {
				lastMsgMap[r.ConvID] = r
			}
		}
	}

	// 获取分类名称映射
	catNameMap := map[int64]string{}
	// 这里假设 getConvCategoryAllCached 已经在其他地方定义，或者直接查询数据库
	var cats []model.ConvCategory
	_ = dal.DB.WithContext(ctx).Model(&model.ConvCategory{}).Find(&cats).Error
	for _, c := range cats {
		catNameMap[c.CategoryID] = c.CategoryName
	}

	// 转换为API响应格式
	items := make([]*customer.ConversationItem, 0, len(convs))
	for _, c := range convs {
		lm, ok := lastMsgMap[c.ConvID]
		lastMsg := ""
		lastTime := c.UpdateTime.Format("2006-01-02 15:04:05")
		if ok {
			lastMsg = lm.MsgContent
			lastTime = lm.SendTime.Format("2006-01-02 15:04:05")
		}
		items = append(items, &customer.ConversationItem{
			ConvId:       c.ConvID,
			UserId:       c.UserID,
			UserNickname: c.UserNickname,
			CsId:         c.CsID,
			Source:       c.Source,
			Status:       c.Status,
			LastMsg:      lastMsg,
			LastTime:     lastTime,
			CategoryId:   c.CategoryID,
			CategoryName: catNameMap[c.CategoryID],
			Tags:         c.Tags,
			IsCore:       c.IsCore,
		})
	}

	resp := &customer.ListConversationResp{
		BaseResp:      &customer.BaseResp{Code: 0, Msg: "success"},
		Conversations: items,
		Total:         total,
	}
	// 写入缓存
	_ = dal.CacheSetJSON(ctx, cacheKey, resp, cacheTTLConversationLists)
	return resp, nil
}

// ListConversationMessage 查询会话消息列表
// 支持分页查询，可按发送时间正序或倒序排列
func (h *CustomerServiceHandler) ListConversationMessage(ctx context.Context, req *customer.ListConversationMessageReq) (*customer.ListConversationMessageResp, error) {
	// 参数校验
	if req == nil || strings.TrimSpace(req.ConvId) == "" {
		return &customer.ListConversationMessageResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "conv_id is required"}}, nil
	}
	// 分页参数处理
	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 500 {
		pageSize = 50
	}

	convID := strings.TrimSpace(req.ConvId)
	var total int64
	q := dal.DB.WithContext(ctx).Model(&model.ConvMessage{}).Where("conv_id = ?", convID)
	// 查询总数
	if err := q.Count(&total).Error; err != nil {
		return &customer.ListConversationMessageResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "查询消息失败: " + err.Error()}}, nil
	}

	// 计算偏移量
	offset := (page - 1) * pageSize
	order := "send_time desc"
	if req.OrderAsc == 1 {
		order = "send_time asc"
	}
	// 分页查询消息
	var msgs []model.ConvMessage
	if err := q.Order(order).Offset(int(offset)).Limit(int(pageSize)).Find(&msgs).Error; err != nil {
		return &customer.ListConversationMessageResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "查询消息失败: " + err.Error()}}, nil
	}

	// 转换为响应格式
	out := make([]*customer.ConvMessageItem, 0, len(msgs))
	for _, m := range msgs {
		out = append(out, &customer.ConvMessageItem{
			MsgId:        m.MsgID,
			ConvId:       m.ConvID,
			SenderType:   m.SenderType,
			SenderId:     m.SenderID,
			MsgContent:   m.MsgContent,
			IsQuickReply: m.IsQuickReply,
			QuickReplyId: m.QuickReplyID,
			SendTime:     m.SendTime.Format("2006-01-02 15:04:05"),
		})
	}

	return &customer.ListConversationMessageResp{
		BaseResp: &customer.BaseResp{Code: 0, Msg: "success"},
		Messages: out,
		Total:    total,
	}, nil
}

// SendConversationMessage 发送会话消息（写入消息表并刷新会话更新时间）
// 权限控制：客服只能回复自己负责的会话，除非已进行转接
func (h *CustomerServiceHandler) SendConversationMessage(ctx context.Context, req *customer.SendConversationMessageReq) (*customer.SendConversationMessageResp, error) {
	if req == nil || strings.TrimSpace(req.ConvId) == "" {
		return &customer.SendConversationMessageResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "conv_id is required"}}, nil
	}
	content := strings.TrimSpace(req.MsgContent)
	if content == "" {
		return &customer.SendConversationMessageResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "msg_content is required"}}, nil
	}

	convID := strings.TrimSpace(req.ConvId)
	senderID := strings.TrimSpace(req.SenderId)
	if senderID == "" {
		senderID = "UNKNOWN"
	}

	// 查询会话信息
	var conv model.Conversation
	if err := dal.DB.WithContext(ctx).Where("conv_id = ?", convID).First(&conv).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &customer.SendConversationMessageResp{BaseResp: &customer.BaseResp{Code: 404, Msg: "会话不存在"}}, nil
		}
		return &customer.SendConversationMessageResp{BaseResp: &customer.BaseResp{Code: 500, Msg: err.Error()}}, nil
	}

	// 权限验证：客服只能回复自己负责的会话
	// sender_type: 0-用户, 1-客服, 2-系统
	if req.SenderType == 1 {
		// 客服发送消息，验证是否为该会话的负责客服
		if conv.CsID != senderID {
			return &customer.SendConversationMessageResp{
				BaseResp: &customer.BaseResp{
					Code: 403,
					Msg:  "权限不足：您不是该会话的负责客服，无法发送消息",
				},
			}, nil
		}
	} else if req.SenderType == 0 {
		// 用户发送消息，验证是否为该会话的用户
		if conv.UserID != senderID {
			return &customer.SendConversationMessageResp{
				BaseResp: &customer.BaseResp{
					Code: 403,
					Msg:  "权限不足：您不是该会话的用户，无法发送消息",
				},
			}, nil
		}
	}
	// sender_type == 2 系统消息不限制

	// 检查会话状态，只有进行中或已转接的会话才能发送消息
	// 转接后，会话状态为Transferred，但cs_id已更新为新客服，新客服应能继续服务
	if conv.Status != model.ConvStatusOngoing && conv.Status != model.ConvStatusTransferred {
		return &customer.SendConversationMessageResp{
			BaseResp: &customer.BaseResp{
				Code: 400,
				Msg:  "会话已结束，无法发送消息",
			},
		}, nil
	}

	var msgID int64
	err := dal.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		m := model.ConvMessage{
			ConvID:       convID,
			SenderType:   req.SenderType,
			SenderID:     senderID,
			MsgContent:   content,
			IsQuickReply: req.IsQuickReply,
			QuickReplyID: req.QuickReplyId,
			SendTime:     now,
		}
		if err := tx.Create(&m).Error; err != nil {
			return err
		}
		msgID = m.MsgID
		// 同时更新会话的update_time和last_msg_time，用于会话超时检测
		return tx.Model(&model.Conversation{}).Where("conv_id = ?", convID).Updates(map[string]interface{}{
			"update_time":   now,
			"last_msg_time": now,
		}).Error
	})
	if err != nil {
		return &customer.SendConversationMessageResp{BaseResp: &customer.BaseResp{Code: 500, Msg: err.Error()}}, nil
	}
	return &customer.SendConversationMessageResp{BaseResp: &customer.BaseResp{Code: 0, Msg: "success"}, MsgId: msgID}, nil
}

// ListQuickReply 查询快捷回复（用于客服会话管理右侧面板）
func (h *CustomerServiceHandler) ListQuickReply(ctx context.Context, req *customer.ListQuickReplyReq) (*customer.ListQuickReplyResp, error) {
	if req == nil {
		return &customer.ListQuickReplyResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "req is required"}}, nil
	}
	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 200 {
		pageSize = 50
	}

	kw := strings.TrimSpace(req.Keyword)
	cacheKey := fmt.Sprintf(
		"customer:quick_reply:list:v1:kw=%s:rt=%d:pub=%d:p=%d:ps=%d",
		url.QueryEscape(kw),
		req.ReplyType,
		req.IsPublic,
		page,
		pageSize,
	)
	{
		var cached customer.ListQuickReplyResp
		if ok, _ := dal.CacheGetJSON(ctx, cacheKey, &cached); ok && cached.BaseResp != nil {
			return &cached, nil
		}
	}

	var total int64
	q := dal.DB.WithContext(ctx).Model(&model.QuickReply{})
	if kw != "" {
		q = q.Where("reply_content LIKE ?", "%"+kw+"%")
	}
	if req.ReplyType >= 0 {
		q = q.Where("reply_type = ?", req.ReplyType)
	}
	if req.IsPublic >= 0 {
		q = q.Where("is_public = ?", req.IsPublic)
	}
	if err := q.Count(&total).Error; err != nil {
		return &customer.ListQuickReplyResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "查询快捷回复失败: " + err.Error()}}, nil
	}

	var replies []model.QuickReply
	offset := (page - 1) * pageSize
	if err := q.Order("update_time desc").Offset(int(offset)).Limit(int(pageSize)).Find(&replies).Error; err != nil {
		return &customer.ListQuickReplyResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "查询快捷回复失败: " + err.Error()}}, nil
	}

	out := make([]*customer.QuickReplyItem, 0, len(replies))
	for _, r := range replies {
		out = append(out, &customer.QuickReplyItem{
			ReplyId:      r.ReplyID,
			ReplyType:    r.ReplyType,
			ReplyContent: r.ReplyContent,
			CreateBy:     r.CreateBy,
			IsPublic:     r.IsPublic,
			UpdateTime:   r.UpdateTime.Format("2006-01-02 15:04:05"),
		})
	}

	resp := &customer.ListQuickReplyResp{BaseResp: &customer.BaseResp{Code: 0, Msg: "success"}, Replies: out, Total: total}
	_ = dal.CacheSetJSON(ctx, cacheKey, resp, cacheTTLQuickReplyList)
	return resp, nil
}

// CreateConvCategory 新增会话分类
func (h *CustomerServiceHandler) CreateConvCategory(ctx context.Context, req *customer.CreateConvCategoryReq) (*customer.CreateConvCategoryResp, error) {
	if req == nil || strings.TrimSpace(req.CategoryName) == "" {
		return &customer.CreateConvCategoryResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "category_name is required"}}, nil
	}
	now := time.Now()
	c := model.ConvCategory{
		CategoryName: strings.TrimSpace(req.CategoryName),
		SortNo:       int(req.SortNo),
		CreateBy:     strings.TrimSpace(req.CreateBy),
		CreateTime:   now,
		UpdateTime:   now,
	}
	if c.CreateBy == "" {
		c.CreateBy = "ADMIN"
	}
	if err := dal.DB.WithContext(ctx).Create(&c).Error; err != nil {
		code := int32(500)
		msg := err.Error()
		if strings.Contains(msg, "Duplicate") || strings.Contains(msg, "duplicate") {
			code = 409
		}
		return &customer.CreateConvCategoryResp{BaseResp: &customer.BaseResp{Code: code, Msg: msg}}, nil
	}
	_ = dal.CacheDel(ctx, cacheKeyConvCategoryAll)
	return &customer.CreateConvCategoryResp{BaseResp: &customer.BaseResp{Code: 0, Msg: "success"}, CategoryId: c.CategoryID}, nil
}

// ListConvCategory 查询会话分类列表
func (h *CustomerServiceHandler) ListConvCategory(ctx context.Context, req *customer.ListConvCategoryReq) (*customer.ListConvCategoryResp, error) {
	out, err := getConvCategoryAllCached(ctx)
	if err != nil {
		return &customer.ListConvCategoryResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "查询分类失败: " + err.Error()}}, nil
	}
	return &customer.ListConvCategoryResp{BaseResp: &customer.BaseResp{Code: 0, Msg: "success"}, Categories: out, Total: int64(len(out))}, nil
}

// UpdateConversationClassify 更新会话分类/标签/核心标记
func (h *CustomerServiceHandler) UpdateConversationClassify(ctx context.Context, req *customer.UpdateConversationClassifyReq) (*customer.UpdateConversationClassifyResp, error) {
	if req == nil || strings.TrimSpace(req.ConvId) == "" {
		return &customer.UpdateConversationClassifyResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "conv_id is required"}}, nil
	}
	if req.IsCore != 0 && req.IsCore != 1 {
		return &customer.UpdateConversationClassifyResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "is_core must be 0 or 1"}}, nil
	}
	update := map[string]interface{}{
		"category_id": req.CategoryId,
		"tags":        strings.TrimSpace(req.Tags),
		"is_core":     req.IsCore,
		"update_time": time.Now(),
	}
	if err := dal.DB.WithContext(ctx).Model(&model.Conversation{}).Where("conv_id = ?", strings.TrimSpace(req.ConvId)).Updates(update).Error; err != nil {
		return &customer.UpdateConversationClassifyResp{BaseResp: &customer.BaseResp{Code: 500, Msg: err.Error()}}, nil
	}
	return &customer.UpdateConversationClassifyResp{BaseResp: &customer.BaseResp{Code: 0, Msg: "success"}}, nil
}

// isValidShiftTimeRange 校验班次时间范围是否有效
// 支持跨夜班次（如 22:00:00 到 06:00:00）
// 返回true表示时间范围有效
func isValidShiftTimeRange(start, end string) bool {
	start = strings.TrimSpace(start)
	end = strings.TrimSpace(end)
	// 解析时间格式 HH:mm:ss
	st, err1 := time.Parse("15:04:05", trimDatePrefix(start))
	et, err2 := time.Parse("15:04:05", trimDatePrefix(end))
	if err1 != nil || err2 != nil {
		return false
	}
	// 正常情况: 结束时间大于开始时间
	if et.After(st) {
		return true
	}
	// 跨夜情况: 结束时间加24小时后应大于开始时间
	et = et.Add(24 * time.Hour)
	return et.After(st)
}

// normalizeTimeForDB 将时间字符串标准化为数据库存储格式
// 如果输入只有时间部分（如 "09:00:00"），会添加默认日期前缀
// 用于班次配置的开始/结束时间存储
func normalizeTimeForDB(t string) string {
	t = strings.TrimSpace(t)
	if t == "" {
		return t
	}
	// 如果已包含日期部分，直接返回
	if strings.Contains(t, "-") {
		return t
	}
	// 添加默认日期前缀
	return "2000-01-01 " + t
}

// normalizeTimeForAPI 将数据库时间转换为API响应格式
// 提取时间部分（HH:mm:ss），去除日期和时区信息
// 用于班次配置的开始/结束时间展示
func normalizeTimeForAPI(t string) string {
	t = strings.TrimSpace(t)
	if t == "" {
		return t
	}
	// 处理ISO 8601格式 (2000-01-01T09:00:00Z)
	if idx := strings.LastIndex(t, "T"); idx >= 0 && idx+1 < len(t) {
		t = t[idx+1:]
		if j := strings.IndexAny(t, "Z+-"); j > 0 {
			t = t[:j]
		}
		if len(t) >= 8 {
			return t[:8]
		}
		return t
	}
	// 处理普通格式 (2000-01-01 09:00:00)
	t = trimDatePrefix(t)
	if j := strings.IndexAny(t, "Z+-"); j > 0 {
		t = t[:j]
	}
	if len(t) >= 8 {
		return t[:8]
	}
	return t
}

// normalizeDateForAPI 将日期字符串标准化为API响应格式
// 提取日期部分（YYYY-MM-DD），去除时间部分
// 用于排班日期的展示
func normalizeDateForAPI(d string) string {
	d = strings.TrimSpace(d)
	if d == "" {
		return d
	}
	// 截取前10个字符（YYYY-MM-DD）
	if len(d) >= 10 {
		return d[:10]
	}
	return d
}

// trimDatePrefix 去除时间字符串中的日期前缀
// 将 "2000-01-01 09:00:00" 转换为 "09:00:00"
func trimDatePrefix(t string) string {
	if idx := strings.LastIndex(t, " "); idx >= 0 {
		return t[idx+1:]
	}
	return t
}

// errorsIsRecordNotFound 判断错误是否为GORM的记录未找到错误
// 用于区分"记录不存在"和其他数据库错误
func errorsIsRecordNotFound(err error) bool {
	return err != nil && errors.Is(err, gorm.ErrRecordNotFound)
}

// ============ 会话标签管理 ============

// getConvTagAllCached 获取所有会话标签（带缓存）
func getConvTagAllCached(ctx context.Context) ([]*customer.ConvTag, error) {
	var cached []*customer.ConvTag
	if ok, _ := dal.CacheGetJSON(ctx, cacheKeyConvTagAll, &cached); ok && cached != nil {
		return cached, nil
	}

	var tags []model.ConvTag
	if err := dal.DB.WithContext(ctx).Model(&model.ConvTag{}).Order("sort_no asc").Order("tag_id asc").Find(&tags).Error; err != nil {
		return nil, err
	}

	out := make([]*customer.ConvTag, 0, len(tags))
	for _, t := range tags {
		out = append(out, &customer.ConvTag{
			TagId:    t.TagID,
			TagName:  t.TagName,
			TagColor: t.TagColor,
			SortNo:   int32(t.SortNo),
		})
	}

	_ = dal.CacheSetJSON(ctx, cacheKeyConvTagAll, out, cacheTTLConvTag)
	return out, nil
}

// CreateConvTag 创建会话标签
func (h *CustomerServiceHandler) CreateConvTag(ctx context.Context, req *customer.CreateConvTagReq) (*customer.CreateConvTagResp, error) {
	if req == nil || strings.TrimSpace(req.TagName) == "" {
		return &customer.CreateConvTagResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "tag_name is required"}}, nil
	}

	now := time.Now()
	tag := model.ConvTag{
		TagName:    strings.TrimSpace(req.TagName),
		TagColor:   strings.TrimSpace(req.TagColor),
		SortNo:     int(req.SortNo),
		CreateBy:   strings.TrimSpace(req.CreateBy),
		CreateTime: now,
		UpdateTime: now,
	}
	if tag.TagColor == "" {
		tag.TagColor = "#1890ff" // 默认颜色
	}
	if tag.CreateBy == "" {
		tag.CreateBy = "ADMIN"
	}

	if err := dal.DB.WithContext(ctx).Create(&tag).Error; err != nil {
		code := int32(500)
		msg := err.Error()
		if strings.Contains(msg, "Duplicate") || strings.Contains(msg, "duplicate") {
			code = 409
		}
		return &customer.CreateConvTagResp{BaseResp: &customer.BaseResp{Code: code, Msg: msg}}, nil
	}

	_ = dal.CacheDel(ctx, cacheKeyConvTagAll)
	return &customer.CreateConvTagResp{BaseResp: &customer.BaseResp{Code: 0, Msg: "success"}, TagId: tag.TagID}, nil
}

// ListConvTag 查询会话标签列表
func (h *CustomerServiceHandler) ListConvTag(ctx context.Context, req *customer.ListConvTagReq) (*customer.ListConvTagResp, error) {
	out, err := getConvTagAllCached(ctx)
	if err != nil {
		return &customer.ListConvTagResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "查询标签失败: " + err.Error()}}, nil
	}
	return &customer.ListConvTagResp{BaseResp: &customer.BaseResp{Code: 0, Msg: "success"}, Tags: out, Total: int64(len(out))}, nil
}

// UpdateConvTag 更新会话标签
func (h *CustomerServiceHandler) UpdateConvTag(ctx context.Context, req *customer.UpdateConvTagReq) (*customer.UpdateConvTagResp, error) {
	if req == nil || req.TagId <= 0 {
		return &customer.UpdateConvTagResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "tag_id is required"}}, nil
	}
	if strings.TrimSpace(req.TagName) == "" {
		return &customer.UpdateConvTagResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "tag_name is required"}}, nil
	}

	update := map[string]interface{}{
		"tag_name":    strings.TrimSpace(req.TagName),
		"tag_color":   strings.TrimSpace(req.TagColor),
		"sort_no":     req.SortNo,
		"update_time": time.Now(),
	}
	if update["tag_color"] == "" {
		update["tag_color"] = "#1890ff"
	}

	res := dal.DB.WithContext(ctx).Model(&model.ConvTag{}).Where("tag_id = ?", req.TagId).Updates(update)
	if res.Error != nil {
		return &customer.UpdateConvTagResp{BaseResp: &customer.BaseResp{Code: 500, Msg: res.Error.Error()}}, nil
	}
	if res.RowsAffected == 0 {
		return &customer.UpdateConvTagResp{BaseResp: &customer.BaseResp{Code: 404, Msg: "tag not found"}}, nil
	}

	_ = dal.CacheDel(ctx, cacheKeyConvTagAll)
	return &customer.UpdateConvTagResp{BaseResp: &customer.BaseResp{Code: 0, Msg: "success"}}, nil
}

// DeleteConvTag 删除会话标签
func (h *CustomerServiceHandler) DeleteConvTag(ctx context.Context, req *customer.DeleteConvTagReq) (*customer.DeleteConvTagResp, error) {
	if req == nil || req.TagId <= 0 {
		return &customer.DeleteConvTagResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "tag_id is required"}}, nil
	}

	res := dal.DB.WithContext(ctx).Where("tag_id = ?", req.TagId).Delete(&model.ConvTag{})
	if res.Error != nil {
		return &customer.DeleteConvTagResp{BaseResp: &customer.BaseResp{Code: 500, Msg: res.Error.Error()}}, nil
	}
	if res.RowsAffected == 0 {
		return &customer.DeleteConvTagResp{BaseResp: &customer.BaseResp{Code: 404, Msg: "tag not found"}}, nil
	}

	_ = dal.CacheDel(ctx, cacheKeyConvTagAll)
	return &customer.DeleteConvTagResp{BaseResp: &customer.BaseResp{Code: 0, Msg: "success"}}, nil
}

// ============ 会话统计看板 ============

// GetConversationStats 获取会话统计数据
// 返回：Top标签分布、Top分类分布、处理时长趋势、核心会话占比
func (h *CustomerServiceHandler) GetConversationStats(ctx context.Context, req *customer.GetConversationStatsReq) (*customer.GetConversationStatsResp, error) {
	if req == nil {
		return &customer.GetConversationStatsResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "req is required"}}, nil
	}

	startDate := strings.TrimSpace(req.StartDate)
	endDate := strings.TrimSpace(req.EndDate)
	if startDate == "" {
		startDate = time.Now().AddDate(0, -1, 0).Format("2006-01-02") // 默认最近1个月
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	// 1. 查询时间范围内的会话基础统计
	var totalConversations, coreConversations int64
	baseQuery := dal.DB.WithContext(ctx).Model(&model.Conversation{}).
		Where("start_time >= ? AND start_time <= ?", startDate+" 00:00:00", endDate+" 23:59:59")

	if err := baseQuery.Count(&totalConversations).Error; err != nil {
		return &customer.GetConversationStatsResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "统计总会话失败: " + err.Error()}}, nil
	}

	if err := dal.DB.WithContext(ctx).Model(&model.Conversation{}).
		Where("start_time >= ? AND start_time <= ? AND is_core = 1", startDate+" 00:00:00", endDate+" 23:59:59").
		Count(&coreConversations).Error; err != nil {
		return &customer.GetConversationStatsResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "统计核心会话失败: " + err.Error()}}, nil
	}

	coreRatio := 0.0
	if totalConversations > 0 {
		coreRatio = float64(coreConversations) / float64(totalConversations) * 100
	}

	// 2. 按分类统计
	type catStat struct {
		CategoryID   int64
		CategoryName string
		Count        int64
	}
	var catStats []catStat
	if err := dal.DB.WithContext(ctx).
		Table("t_conversation c").
		Select("c.category_id, COALESCE(cat.category_name, '未分类') as category_name, COUNT(*) as count").
		Joins("LEFT JOIN t_conv_category cat ON c.category_id = cat.category_id").
		Where("c.start_time >= ? AND c.start_time <= ?", startDate+" 00:00:00", endDate+" 23:59:59").
		Group("c.category_id, cat.category_name").
		Order("count DESC").
		Limit(10).
		Scan(&catStats).Error; err != nil {
		return &customer.GetConversationStatsResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "统计分类失败: " + err.Error()}}, nil
	}

	topCategories := make([]*customer.CategoryStat, 0, len(catStats))
	for _, cs := range catStats {
		ratio := 0.0
		if totalConversations > 0 {
			ratio = float64(cs.Count) / float64(totalConversations) * 100
		}
		topCategories = append(topCategories, &customer.CategoryStat{
			CategoryName: cs.CategoryName,
			Count:        cs.Count,
			Ratio:        ratio,
		})
	}

	// 3. 按标签统计（解析tags字段，逗号分隔）
	type convTags struct {
		Tags string
	}
	var allTags []convTags
	if err := dal.DB.WithContext(ctx).
		Table("t_conversation").
		Select("tags").
		Where("start_time >= ? AND start_time <= ? AND tags != '' AND tags IS NOT NULL", startDate+" 00:00:00", endDate+" 23:59:59").
		Scan(&allTags).Error; err != nil {
		return &customer.GetConversationStatsResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "统计标签失败: " + err.Error()}}, nil
	}

	tagCountMap := make(map[string]int64)
	totalTagCount := int64(0)
	for _, ct := range allTags {
		for _, tag := range strings.Split(ct.Tags, ",") {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				tagCountMap[tag]++
				totalTagCount++
			}
		}
	}

	// 转换为切片并排序
	type tagCount struct {
		Name  string
		Count int64
	}
	tagCounts := make([]tagCount, 0, len(tagCountMap))
	for name, count := range tagCountMap {
		tagCounts = append(tagCounts, tagCount{Name: name, Count: count})
	}
	// 按数量降序排序
	for i := 0; i < len(tagCounts)-1; i++ {
		for j := i + 1; j < len(tagCounts); j++ {
			if tagCounts[i].Count < tagCounts[j].Count {
				tagCounts[i], tagCounts[j] = tagCounts[j], tagCounts[i]
			}
		}
	}

	topTags := make([]*customer.TagStat, 0, 10)
	for i, tc := range tagCounts {
		if i >= 10 {
			break
		}
		ratio := 0.0
		if totalConversations > 0 {
			ratio = float64(tc.Count) / float64(totalConversations) * 100
		}
		topTags = append(topTags, &customer.TagStat{
			TagName: tc.Name,
			Count:   tc.Count,
			Ratio:   ratio,
		})
	}

	// 4. 按日期统计平均处理时长
	type durationRow struct {
		Date        string
		AvgDuration float64
		ConvCount   int64
	}
	var durationStats []durationRow
	if err := dal.DB.WithContext(ctx).
		Table("t_conversation").
		Select("DATE(start_time) as date, AVG(TIMESTAMPDIFF(MINUTE, start_time, end_time)) as avg_duration, COUNT(*) as conv_count").
		Where("start_time >= ? AND start_time <= ? AND end_time IS NOT NULL AND end_time > start_time", startDate+" 00:00:00", endDate+" 23:59:59").
		Group("DATE(start_time)").
		Order("date ASC").
		Scan(&durationStats).Error; err != nil {
		return &customer.GetConversationStatsResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "统计时长失败: " + err.Error()}}, nil
	}

	durationTrend := make([]*customer.DurationStat, 0, len(durationStats))
	for _, ds := range durationStats {
		durationTrend = append(durationTrend, &customer.DurationStat{
			Date:               ds.Date,
			AvgDurationMinutes: ds.AvgDuration,
			ConvCount:          ds.ConvCount,
		})
	}

	return &customer.GetConversationStatsResp{
		BaseResp:           &customer.BaseResp{Code: 0, Msg: "success"},
		TopTags:            topTags,
		TopCategories:      topCategories,
		DurationTrend:      durationTrend,
		TotalConversations: totalConversations,
		CoreConversations:  coreConversations,
		CoreRatio:          coreRatio,
	}, nil
}

// ============ 会话记录导出 ============

// ExportConversationItem 导出的会话记录项
type ExportConversationItem struct {
	ConvID       string // 会话ID
	UserID       string // 用户ID
	UserNickname string // 用户昵称
	CsID         string // 客服ID
	CsName       string // 客服名称
	Source       string // 来源
	Status       string // 状态
	StartTime    string // 开始时间
	EndTime      string // 结束时间
	Duration     string // 时长
	Category     string // 分类
	Tags         string // 标签
	IsCore       string // 是否核心
	MsgCount     int64  // 消息数
}

// ExportConversationRequest 导出会话请求
type ExportConversationRequest struct {
	StartDate string // 开始日期
	EndDate   string // 结束日期
	CsID      string // 客服ID（可选）
	Status    int8   // 状态（可选，-1表示全部）
	Keyword   string // 关键词（可选）
}

// GetConversationsForExport 获取用于导出的会话数据
// 返回结构化数据，Gateway层负责生成Excel文件
func (h *CustomerServiceHandler) GetConversationsForExport(ctx context.Context, req *ExportConversationRequest) ([]ExportConversationItem, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	startDate := strings.TrimSpace(req.StartDate)
	endDate := strings.TrimSpace(req.EndDate)
	if startDate == "" {
		startDate = time.Now().AddDate(0, -1, 0).Format("2006-01-02") // 默认最近1个月
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	// 构建查询
	q := dal.DB.WithContext(ctx).Model(&model.Conversation{}).
		Where("start_time >= ? AND start_time <= ?", startDate+" 00:00:00", endDate+" 23:59:59")

	if req.CsID != "" {
		q = q.Where("cs_id = ?", req.CsID)
	}
	if req.Status >= 0 {
		q = q.Where("status = ?", req.Status)
	}
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		q = q.Where("conv_id LIKE ? OR user_id LIKE ? OR user_nickname LIKE ?", like, like, like)
	}

	var convs []model.Conversation
	if err := q.Order("start_time DESC").Limit(5000).Find(&convs).Error; err != nil { // 限制5000条防止内存溢出
		return nil, err
	}

	// 获取客服名称映射
	csNameMap := make(map[string]string)
	var csList []model.CustomerService
	if err := dal.DB.WithContext(ctx).Find(&csList).Error; err == nil {
		for _, cs := range csList {
			csNameMap[cs.CsID] = cs.CsName
		}
	}

	// 获取分类名称映射
	catNameMap := make(map[int64]string)
	if cats, err := getConvCategoryAllCached(ctx); err == nil {
		for _, c := range cats {
			catNameMap[c.CategoryId] = c.CategoryName
		}
	}

	// 统计消息数
	convIDs := make([]string, 0, len(convs))
	for _, c := range convs {
		convIDs = append(convIDs, c.ConvID)
	}

	msgCountMap := make(map[string]int64)
	if len(convIDs) > 0 {
		type msgCount struct {
			ConvID string
			Count  int64
		}
		var counts []msgCount
		dal.DB.WithContext(ctx).Model(&model.ConvMessage{}).
			Select("conv_id, COUNT(*) as count").
			Where("conv_id IN ?", convIDs).
			Group("conv_id").
			Scan(&counts)
		for _, mc := range counts {
			msgCountMap[mc.ConvID] = mc.Count
		}
	}

	// 构建导出数据
	items := make([]ExportConversationItem, 0, len(convs))
	for _, c := range convs {
		endTime := ""
		duration := ""
		if !c.EndTime.IsZero() {
			endTime = c.EndTime.Format("2006-01-02 15:04:05")
			dur := c.EndTime.Sub(c.StartTime)
			duration = fmt.Sprintf("%d分%d秒", int(dur.Minutes()), int(dur.Seconds())%60)
		}

		isCore := "否"
		if c.IsCore == 1 {
			isCore = "是"
		}

		items = append(items, ExportConversationItem{
			ConvID:       c.ConvID,
			UserID:       c.UserID,
			UserNickname: c.UserNickname,
			CsID:         c.CsID,
			CsName:       csNameMap[c.CsID],
			Source:       c.Source,
			Status:       model.ConvStatusName(c.Status),
			StartTime:    c.StartTime.Format("2006-01-02 15:04:05"),
			EndTime:      endTime,
			Duration:     duration,
			Category:     catNameMap[c.CategoryID],
			Tags:         c.Tags,
			IsCore:       isCore,
			MsgCount:     msgCountMap[c.ConvID],
		})
	}

	return items, nil
}

// ============ 用户认证 ============

// Login 用户登录
// 验证用户名和密码，返回用户信息
// Token由Gateway层生成，RPC层只负责验证
func (h *CustomerServiceHandler) Login(ctx context.Context, req *customer.LoginReq) (*customer.LoginResp, error) {
	// 参数校验
	if req == nil {
		return &customer.LoginResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "请求参数不能为空"}}, nil
	}
	userName := strings.TrimSpace(req.UserName)
	password := strings.TrimSpace(req.Password)
	if userName == "" {
		return &customer.LoginResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "用户名不能为空"}}, nil
	}
	if password == "" {
		return &customer.LoginResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "密码不能为空"}}, nil
	}

	// 查询用户
	var user model.SysUser
	if err := dal.DB.WithContext(ctx).Where("user_name = ?", userName).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &customer.LoginResp{BaseResp: &customer.BaseResp{Code: 401, Msg: "用户名或密码错误"}}, nil
		}
		return &customer.LoginResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "查询用户失败: " + err.Error()}}, nil
	}

	// 检查用户状态
	if user.Status != 1 {
		return &customer.LoginResp{BaseResp: &customer.BaseResp{Code: 403, Msg: "账号已被禁用"}}, nil
	}

	// 校验密码
	if !model.CheckPassword(user.Password, password) {
		return &customer.LoginResp{BaseResp: &customer.BaseResp{Code: 401, Msg: "用户名或密码错误"}}, nil
	}

	// 查询角色名称
	var role model.SysRole
	roleName := ""
	if err := dal.DB.WithContext(ctx).Where("role_code = ?", user.RoleCode).First(&role).Error; err == nil {
		roleName = role.RoleName
	}

	// 返回用户信息
	return &customer.LoginResp{
		BaseResp: &customer.BaseResp{Code: 0, Msg: "登录成功"},
		UserInfo: &customer.UserInfo{
			Id:       user.ID,
			UserName: user.UserName,
			RealName: user.RealName,
			Phone:    user.Phone,
			RoleCode: user.RoleCode,
			RoleName: roleName,
			Status:   user.Status,
		},
	}, nil
}

// GetCurrentUser 获取当前用户信息
// 根据用户ID查询用户详细信息
func (h *CustomerServiceHandler) GetCurrentUser(ctx context.Context, req *customer.GetCurrentUserReq) (*customer.GetCurrentUserResp, error) {
	// 参数校验
	if req == nil || req.UserId <= 0 {
		return &customer.GetCurrentUserResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "用户ID不能为空"}}, nil
	}

	// 查询用户
	var user model.SysUser
	if err := dal.DB.WithContext(ctx).Where("id = ?", req.UserId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &customer.GetCurrentUserResp{BaseResp: &customer.BaseResp{Code: 404, Msg: "用户不存在"}}, nil
		}
		return &customer.GetCurrentUserResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "查询用户失败: " + err.Error()}}, nil
	}

	// 检查用户状态
	if user.Status != 1 {
		return &customer.GetCurrentUserResp{BaseResp: &customer.BaseResp{Code: 403, Msg: "账号已被禁用"}}, nil
	}

	// 查询角色名称
	var role model.SysRole
	roleName := ""
	if err := dal.DB.WithContext(ctx).Where("role_code = ?", user.RoleCode).First(&role).Error; err == nil {
		roleName = role.RoleName
	}

	// 返回用户信息
	return &customer.GetCurrentUserResp{
		BaseResp: &customer.BaseResp{Code: 0, Msg: "success"},
		UserInfo: &customer.UserInfo{
			Id:       user.ID,
			UserName: user.UserName,
			RealName: user.RealName,
			Phone:    user.Phone,
			RoleCode: user.RoleCode,
			RoleName: roleName,
			Status:   user.Status,
		},
	}, nil
}

// Register 用户注册
// 仅允许注册客服账号，管理员账号需要通过后台创建
// 注册成功后会同时在t_customer_service表中创建客服信息，以便后台进行排班管理
func (h *CustomerServiceHandler) Register(ctx context.Context, req *customer.RegisterReq) (*customer.RegisterResp, error) {
	// 参数校验
	if req == nil {
		return &customer.RegisterResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "请求参数不能为空"}}, nil
	}
	userName := strings.TrimSpace(req.UserName)
	password := strings.TrimSpace(req.Password)
	realName := strings.TrimSpace(req.RealName)
	phone := strings.TrimSpace(req.Phone)

	if userName == "" {
		return &customer.RegisterResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "用户名不能为空"}}, nil
	}
	if len(userName) < 3 || len(userName) > 32 {
		return &customer.RegisterResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "用户名长度需要3-32个字符"}}, nil
	}
	if password == "" {
		return &customer.RegisterResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "密码不能为空"}}, nil
	}
	if len(password) < 6 {
		return &customer.RegisterResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "密码长度至少6位"}}, nil
	}
	if realName == "" {
		return &customer.RegisterResp{BaseResp: &customer.BaseResp{Code: 400, Msg: "真实姓名不能为空"}}, nil
	}

	// 检查用户名是否已存在
	var existUser model.SysUser
	if err := dal.DB.WithContext(ctx).Where("user_name = ?", userName).First(&existUser).Error; err == nil {
		return &customer.RegisterResp{BaseResp: &customer.BaseResp{Code: 409, Msg: "用户名已存在"}}, nil
	}

	// 加密密码
	hashedPassword, err := model.HashPassword(password)
	if err != nil {
		return &customer.RegisterResp{BaseResp: &customer.BaseResp{Code: 500, Msg: "密码加密失败: " + err.Error()}}, nil
	}

	// 使用事务确保sys_users和t_customer_service两张表数据一致
	var newUserID int64
	now := time.Now()

	err = dal.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 创建客服账号（固定为customer_service角色）
		newUser := model.SysUser{
			UserName: userName,
			Password: hashedPassword,
			RealName: realName,
			Phone:    phone,
			RoleCode: model.RoleCustomerService, // 强制设置为客服角色
			Status:   1,                         // 默认启用
		}
		if err := tx.Create(&newUser).Error; err != nil {
			return fmt.Errorf("创建用户失败: %w", err)
		}
		newUserID = newUser.ID

		// 2. 同步创建客服信息表记录，用于后台排班管理
		// 客服ID格式：CS + 用户ID，如 CS1、CS2 等
		csID := fmt.Sprintf("CS%d", newUser.ID)
		customerService := model.CustomerService{
			CsID:          csID,
			CsName:        realName,  // 使用真实姓名作为客服姓名
			DeptID:        "DEFAULT", // 默认部门，后续可由管理员修改
			TeamID:        "",        // 默认无班组
			SkillTags:     "",        // 默认无技能标签
			Status:        1,         // 1=在职
			CurrentStatus: 0,         // 0=空闲，等待排班
			CreateTime:    now,
			UpdateTime:    now,
		}
		if err := tx.Create(&customerService).Error; err != nil {
			return fmt.Errorf("创建客服信息失败: %w", err)
		}

		return nil
	})

	if err != nil {
		return &customer.RegisterResp{BaseResp: &customer.BaseResp{Code: 500, Msg: err.Error()}}, nil
	}

	return &customer.RegisterResp{
		BaseResp: &customer.BaseResp{Code: 0, Msg: "注册成功"},
		UserId:   newUserID,
	}, nil
}

// ============ 快捷回复管理 ============

// CreateQuickReply 创建快捷回复
// 管理员或客服创建预设的快捷回复话术
// 参数:
//   - reply_type: 回复类型（0-通用, 1-售前, 2-售后, 3-投诉等）
//   - reply_content: 回复内容
//   - create_by: 创建人
//   - is_public: 是否公开（0-私有仅创建者可见, 1-公开所有客服可用）
func (h *CustomerServiceHandler) CreateQuickReply(ctx context.Context, req *customer.CreateQuickReplyReq) (*customer.CreateQuickReplyResp, error) {
	// 参数校验
	if req == nil {
		return &customer.CreateQuickReplyResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "请求参数不能为空"},
		}, nil
	}

	replyContent := strings.TrimSpace(req.ReplyContent)
	if replyContent == "" {
		return &customer.CreateQuickReplyResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "回复内容不能为空"},
		}, nil
	}
	if len(replyContent) > 2000 {
		return &customer.CreateQuickReplyResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "回复内容不能超过2000字符"},
		}, nil
	}

	createBy := strings.TrimSpace(req.CreateBy)
	if createBy == "" {
		createBy = "ADMIN"
	}

	// 校验is_public
	isPublic := req.IsPublic
	if isPublic != 0 && isPublic != 1 {
		isPublic = 0 // 默认私有
	}

	// 创建记录
	now := time.Now()
	reply := model.QuickReply{
		ReplyType:    req.ReplyType,
		ReplyContent: replyContent,
		CreateBy:     createBy,
		IsPublic:     isPublic,
		CreateTime:   now,
		UpdateTime:   now,
	}

	if err := dal.DB.WithContext(ctx).Create(&reply).Error; err != nil {
		return &customer.CreateQuickReplyResp{
			BaseResp: &customer.BaseResp{Code: 500, Msg: "创建失败: " + err.Error()},
		}, nil
	}

	return &customer.CreateQuickReplyResp{
		BaseResp: &customer.BaseResp{Code: 0, Msg: "success"},
		ReplyId:  reply.ReplyID,
	}, nil
}

// UpdateQuickReply 更新快捷回复
// 修改已有的快捷回复内容或属性
// 参数:
//   - reply_id: 回复ID（必填）
//   - reply_type: 回复类型
//   - reply_content: 回复内容
//   - is_public: 是否公开
func (h *CustomerServiceHandler) UpdateQuickReply(ctx context.Context, req *customer.UpdateQuickReplyReq) (*customer.UpdateQuickReplyResp, error) {
	// 参数校验
	if req == nil || req.ReplyId <= 0 {
		return &customer.UpdateQuickReplyResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "reply_id不能为空"},
		}, nil
	}

	replyContent := strings.TrimSpace(req.ReplyContent)
	if replyContent == "" {
		return &customer.UpdateQuickReplyResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "回复内容不能为空"},
		}, nil
	}
	if len(replyContent) > 2000 {
		return &customer.UpdateQuickReplyResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "回复内容不能超过2000字符"},
		}, nil
	}

	// 检查记录是否存在
	var existing model.QuickReply
	if err := dal.DB.WithContext(ctx).Where("reply_id = ?", req.ReplyId).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &customer.UpdateQuickReplyResp{
				BaseResp: &customer.BaseResp{Code: 404, Msg: "快捷回复不存在"},
			}, nil
		}
		return &customer.UpdateQuickReplyResp{
			BaseResp: &customer.BaseResp{Code: 500, Msg: "查询失败: " + err.Error()},
		}, nil
	}

	// 校验is_public
	isPublic := req.IsPublic
	if isPublic != 0 && isPublic != 1 {
		isPublic = existing.IsPublic // 保持原值
	}

	// 更新记录
	updates := map[string]interface{}{
		"reply_type":    req.ReplyType,
		"reply_content": replyContent,
		"is_public":     isPublic,
		"update_time":   time.Now(),
	}

	if err := dal.DB.WithContext(ctx).Model(&model.QuickReply{}).Where("reply_id = ?", req.ReplyId).Updates(updates).Error; err != nil {
		return &customer.UpdateQuickReplyResp{
			BaseResp: &customer.BaseResp{Code: 500, Msg: "更新失败: " + err.Error()},
		}, nil
	}

	return &customer.UpdateQuickReplyResp{
		BaseResp: &customer.BaseResp{Code: 0, Msg: "success"},
	}, nil
}

// DeleteQuickReply 删除快捷回复
// 删除指定的快捷回复记录
// 参数:
//   - reply_id: 回复ID（必填）
func (h *CustomerServiceHandler) DeleteQuickReply(ctx context.Context, req *customer.DeleteQuickReplyReq) (*customer.DeleteQuickReplyResp, error) {
	// 参数校验
	if req == nil || req.ReplyId <= 0 {
		return &customer.DeleteQuickReplyResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "reply_id不能为空"},
		}, nil
	}

	// 执行删除
	result := dal.DB.WithContext(ctx).Where("reply_id = ?", req.ReplyId).Delete(&model.QuickReply{})
	if result.Error != nil {
		return &customer.DeleteQuickReplyResp{
			BaseResp: &customer.BaseResp{Code: 500, Msg: "删除失败: " + result.Error.Error()},
		}, nil
	}

	if result.RowsAffected == 0 {
		return &customer.DeleteQuickReplyResp{
			BaseResp: &customer.BaseResp{Code: 404, Msg: "快捷回复不存在"},
		}, nil
	}

	return &customer.DeleteQuickReplyResp{
		BaseResp: &customer.BaseResp{Code: 0, Msg: "success"},
	}, nil
}

// ============ 会话监控 ============

// GetConversationMonitor 获取会话监控数据
// 实时查看所有会话状态、客服在线状态
// 参数:
//   - dept_id: 部门ID筛选（可选）
//   - status_filter: 状态筛选 -1-全部 0-等待 1-进行中
func (h *CustomerServiceHandler) GetConversationMonitor(ctx context.Context, req *customer.GetConversationMonitorReq) (*customer.GetConversationMonitorResp, error) {
	if req == nil {
		req = &customer.GetConversationMonitorReq{StatusFilter: -1}
	}

	// 查询客服状态列表
	var csList []model.CustomerService
	csQuery := dal.DB.WithContext(ctx).Model(&model.CustomerService{}).Where("status = 1") // 在职客服
	if req.DeptId != "" {
		csQuery = csQuery.Where("dept_id = ?", req.DeptId)
	}
	if err := csQuery.Find(&csList).Error; err != nil {
		return &customer.GetConversationMonitorResp{
			BaseResp: &customer.BaseResp{Code: 500, Msg: "查询客服列表失败: " + err.Error()},
		}, nil
	}

	// 统计每个客服的当前会话数和今日处理数
	today := time.Now().Format("2006-01-02")
	csStatusList := make([]*customer.CsStatusInfo, 0, len(csList))
	onlineCsCount := int32(0)

	for _, cs := range csList {
		// 当前进行中会话数
		var currentCount int64
		dal.DB.WithContext(ctx).Model(&model.Conversation{}).
			Where("cs_id = ? AND status IN (?, ?)", cs.CsID, model.ConvStatusWaiting, model.ConvStatusOngoing).
			Count(&currentCount)

		// 今日处理总数
		var todayCount int64
		dal.DB.WithContext(ctx).Model(&model.Conversation{}).
			Where("cs_id = ? AND DATE(start_time) = ?", cs.CsID, today).
			Count(&todayCount)

		// 确定在线状态
		onlineStatus := int8(0) // 默认离线
		if cs.CurrentStatus == 1 {
			if currentCount > 0 {
				onlineStatus = 2 // 忙碌
			} else {
				onlineStatus = 1 // 在线
				onlineCsCount++
			}
		}

		csStatusList = append(csStatusList, &customer.CsStatusInfo{
			CsId:             cs.CsID,
			CsName:           cs.CsName,
			OnlineStatus:     onlineStatus,
			CurrentConvCount: int32(currentCount),
			TodayConvCount:   int32(todayCount),
			LastActiveTime:   cs.UpdateTime.Format("2006-01-02 15:04:05"),
		})
	}

	// 查询会话列表
	var convList []model.Conversation
	convQuery := dal.DB.WithContext(ctx).Model(&model.Conversation{}).
		Where("status IN (?, ?)", model.ConvStatusWaiting, model.ConvStatusOngoing)
	if req.StatusFilter == 0 {
		convQuery = convQuery.Where("status = ?", model.ConvStatusWaiting)
	} else if req.StatusFilter == 1 {
		convQuery = convQuery.Where("status = ?", model.ConvStatusOngoing)
	}
	convQuery.Order("start_time desc").Limit(100).Find(&convList)

	// 统计
	var waitingCount, ongoingCount int64
	dal.DB.WithContext(ctx).Model(&model.Conversation{}).Where("status = ?", model.ConvStatusWaiting).Count(&waitingCount)
	dal.DB.WithContext(ctx).Model(&model.Conversation{}).Where("status = ?", model.ConvStatusOngoing).Count(&ongoingCount)

	// 转换会话列表
	monitorConvList := make([]*customer.MonitorConvItem, 0, len(convList))
	now := time.Now()
	for _, conv := range convList {
		// 获取最后一条消息
		var lastMsg model.ConvMessage
		dal.DB.WithContext(ctx).Model(&model.ConvMessage{}).
			Where("conv_id = ?", conv.ConvID).
			Order("send_time desc").First(&lastMsg)

		// 获取客服名称
		var csName string
		for _, cs := range csList {
			if cs.CsID == conv.CsID {
				csName = cs.CsName
				break
			}
		}

		waitSeconds := int32(0)
		durationSeconds := int32(0)
		if conv.Status == model.ConvStatusWaiting {
			waitSeconds = int32(now.Sub(conv.StartTime).Seconds())
		} else {
			durationSeconds = int32(now.Sub(conv.StartTime).Seconds())
		}

		monitorConvList = append(monitorConvList, &customer.MonitorConvItem{
			ConvId:          conv.ConvID,
			UserId:          conv.UserID,
			UserNickname:    conv.UserNickname,
			CsId:            conv.CsID,
			CsName:          csName,
			Status:          conv.Status,
			WaitSeconds:     waitSeconds,
			DurationSeconds: durationSeconds,
			LastMsg:         lastMsg.MsgContent,
			StartTime:       conv.StartTime.Format("2006-01-02 15:04:05"),
		})
	}

	return &customer.GetConversationMonitorResp{
		BaseResp:      &customer.BaseResp{Code: 0, Msg: "success"},
		CsList:        csStatusList,
		ConvList:      monitorConvList,
		WaitingCount:  int32(waitingCount),
		OngoingCount:  int32(ongoingCount),
		OnlineCsCount: onlineCsCount,
	}, nil
}

// ============ 会话记录导出 ============

// ExportConversations 导出会话记录
// 支持导出为Excel/CSV格式
// 参数:
//   - cs_id: 客服ID筛选（可选）
//   - user_id: 用户ID筛选（可选）
//   - start_date: 开始日期
//   - end_date: 结束日期
//   - status: 状态筛选 -1-全部
//   - keyword: 关键词搜索
//   - export_format: 导出格式 excel/csv
func (h *CustomerServiceHandler) ExportConversations(ctx context.Context, req *customer.ExportConversationsReq) (*customer.ExportConversationsResp, error) {
	if req == nil {
		return &customer.ExportConversationsResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "请求参数不能为空"},
		}, nil
	}

	// 构建查询条件
	query := dal.DB.WithContext(ctx).Model(&model.Conversation{})

	if req.CsId != "" {
		query = query.Where("cs_id = ?", req.CsId)
	}
	if req.UserId != "" {
		query = query.Where("user_id = ?", req.UserId)
	}
	if req.StartDate != "" {
		query = query.Where("DATE(start_time) >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("DATE(start_time) <= ?", req.EndDate)
	}
	if req.Status != -1 {
		query = query.Where("status = ?", req.Status)
	}
	if req.Keyword != "" {
		kw := "%" + req.Keyword + "%"
		query = query.Where("user_nickname LIKE ? OR user_id LIKE ?", kw, kw)
	}

	// 查询数据
	var convList []model.Conversation
	if err := query.Order("start_time desc").Limit(10000).Find(&convList).Error; err != nil {
		return &customer.ExportConversationsResp{
			BaseResp: &customer.BaseResp{Code: 500, Msg: "查询失败: " + err.Error()},
		}, nil
	}

	// 生成CSV数据
	var csvData strings.Builder
	csvData.WriteString("会话ID,用户ID,用户昵称,客服ID,来源渠道,状态,开始时间,结束时间,分类ID,标签,是否核心\n")

	for _, conv := range convList {
		statusName := model.ConvStatusName(conv.Status)
		endTimeStr := ""
		if !conv.EndTime.IsZero() {
			endTimeStr = conv.EndTime.Format("2006-01-02 15:04:05")
		}
		isCore := "否"
		if conv.IsCore == 1 {
			isCore = "是"
		}

		csvData.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%d,%s,%s\n",
			conv.ConvID, conv.UserID, conv.UserNickname, conv.CsID, conv.Source,
			statusName, conv.StartTime.Format("2006-01-02 15:04:05"), endTimeStr,
			conv.CategoryID, conv.Tags, isCore))
	}

	fileName := fmt.Sprintf("会话记录_%s.csv", time.Now().Format("20060102150405"))

	return &customer.ExportConversationsResp{
		BaseResp:   &customer.BaseResp{Code: 0, Msg: "success"},
		FileData:   []byte(csvData.String()),
		FileName:   fileName,
		TotalCount: int64(len(convList)),
	}, nil
}

// ============ 消息分类管理 ============

// MsgAutoClassify 消息自动分类（NLP增强版）
// 采用多策略融合分类：关键词匹配(40%) + TF-IDF相似度(40%) + 语义规则(20%)
// 参数:
//   - conv_id: 会话ID
//   - msg_contents: 消息内容列表
func (h *CustomerServiceHandler) MsgAutoClassify(ctx context.Context, req *customer.MsgAutoClassifyReq) (*customer.MsgAutoClassifyResp, error) {
	if req == nil || req.ConvId == "" || len(req.MsgContents) == 0 {
		return &customer.MsgAutoClassifyResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "参数不完整"},
		}, nil
	}

	// 查询所有分类及关键词
	var categories []model.MsgCategory
	if err := dal.DB.WithContext(ctx).Model(&model.MsgCategory{}).Order("sort_no asc").Find(&categories).Error; err != nil {
		return &customer.MsgAutoClassifyResp{
			BaseResp: &customer.BaseResp{Code: 500, Msg: "查询分类失败: " + err.Error()},
		}, nil
	}

	if len(categories) == 0 {
		return &customer.MsgAutoClassifyResp{
			BaseResp:          &customer.BaseResp{Code: 0, Msg: "success"},
			CategoryId:        0,
			CategoryName:      "未分类",
			Confidence:        0,
			NeedManualConfirm: true,
		}, nil
	}

	// 创建NLP分类器
	classifier := nlp.NewClassifier()
	for _, cat := range categories {
		keywords := nlp.ParseKeywordsJSON(cat.Keywords)
		classifier.AddCategory(cat.CategoryID, cat.CategoryName, keywords)
	}

	// 合并所有消息内容进行分类
	allContent := strings.Join(req.MsgContents, " ")
	result := classifier.Classify(allContent)

	// 更新会话的分类（置信度足够高时自动更新）
	if !result.NeedManual {
		dal.DB.WithContext(ctx).Model(&model.Conversation{}).
			Where("conv_id = ?", req.ConvId).
			Updates(map[string]interface{}{
				"category_id":      result.CategoryID,
				"is_manual_adjust": 0,
				"update_time":      time.Now(),
			})
	}

	return &customer.MsgAutoClassifyResp{
		BaseResp:          &customer.BaseResp{Code: 0, Msg: "success"},
		CategoryId:        result.CategoryID,
		CategoryName:      result.CategoryName,
		Confidence:        result.Confidence,
		NeedManualConfirm: result.NeedManual,
		MatchedKeywords:   result.MatchedKeywords,
	}, nil
}

// AdjustMsgClassify 人工调整分类
// 客服人工修正自动分类结果
// 参数:
//   - conv_id: 会话ID
//   - original_category_id: 原分类ID
//   - new_category_id: 新分类ID
//   - operator_id: 操作人ID
//   - adjust_reason: 调整原因
func (h *CustomerServiceHandler) AdjustMsgClassify(ctx context.Context, req *customer.AdjustMsgClassifyReq) (*customer.AdjustMsgClassifyResp, error) {
	if req == nil || req.ConvId == "" || req.NewCategoryId_ <= 0 {
		return &customer.AdjustMsgClassifyResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "参数不完整"},
		}, nil
	}

	// 验证新分类是否存在
	var newCat model.MsgCategory
	if err := dal.DB.WithContext(ctx).Where("category_id = ?", req.NewCategoryId_).First(&newCat).Error; err != nil {
		return &customer.AdjustMsgClassifyResp{
			BaseResp: &customer.BaseResp{Code: 404, Msg: "目标分类不存在"},
		}, nil
	}

	// 记录调整日志
	adjustLog := model.ClassifyAdjustLog{
		ConvID:             req.ConvId,
		OriginalCategoryID: req.OriginalCategoryId,
		NewCategoryID:      req.NewCategoryId_,
		OperatorID:         req.OperatorId,
		AdjustReason:       req.AdjustReason,
		CreateTime:         time.Now(),
	}
	if err := dal.DB.WithContext(ctx).Create(&adjustLog).Error; err != nil {
		return &customer.AdjustMsgClassifyResp{
			BaseResp: &customer.BaseResp{Code: 500, Msg: "记录调整日志失败: " + err.Error()},
		}, nil
	}

	// 更新会话分类
	updates := map[string]interface{}{
		"category_id":      req.NewCategoryId_,
		"is_manual_adjust": 1,
		"update_time":      time.Now(),
	}
	if err := dal.DB.WithContext(ctx).Model(&model.Conversation{}).Where("conv_id = ?", req.ConvId).Updates(updates).Error; err != nil {
		return &customer.AdjustMsgClassifyResp{
			BaseResp: &customer.BaseResp{Code: 500, Msg: "更新会话分类失败: " + err.Error()},
		}, nil
	}

	return &customer.AdjustMsgClassifyResp{
		BaseResp:    &customer.BaseResp{Code: 0, Msg: "success"},
		AdjustLogId: adjustLog.LogID,
	}, nil
}

// GetClassifyStats 获取分类统计数据
// 查询消息分类的统计信息
// 参数:
//   - start_date: 开始日期
//   - end_date: 结束日期
//   - stat_type: 统计类型 day/week/month
func (h *CustomerServiceHandler) GetClassifyStats(ctx context.Context, req *customer.GetClassifyStatsReq) (*customer.GetClassifyStatsResp, error) {
	if req == nil {
		req = &customer.GetClassifyStatsReq{}
	}

	// 默认查询最近30天
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	if req.StartDate != "" {
		if t, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			startDate = t
		}
	}
	if req.EndDate != "" {
		if t, err := time.Parse("2006-01-02", req.EndDate); err == nil {
			endDate = t
		}
	}

	// 查询所有分类
	var categories []model.MsgCategory
	dal.DB.WithContext(ctx).Model(&model.MsgCategory{}).Order("sort_no asc").Find(&categories)

	categoryMap := make(map[int64]string)
	for _, cat := range categories {
		categoryMap[cat.CategoryID] = cat.CategoryName
	}

	// 统计各分类数量
	type catCount struct {
		CategoryID int64 `gorm:"column:category_id"`
		Count      int64 `gorm:"column:count"`
	}
	var catCounts []catCount
	dal.DB.WithContext(ctx).Model(&model.Conversation{}).
		Select("category_id, COUNT(*) as count").
		Where("DATE(start_time) BETWEEN ? AND ?", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Where("category_id > 0").
		Group("category_id").
		Scan(&catCounts)

	// 计算总数和占比
	var totalClassified int64
	for _, cc := range catCounts {
		totalClassified += cc.Count
	}

	categorySummary := make([]*customer.CategoryStat, 0, len(catCounts))
	for _, cc := range catCounts {
		ratio := float64(0)
		if totalClassified > 0 {
			ratio = float64(cc.Count) * 100 / float64(totalClassified)
		}
		categorySummary = append(categorySummary, &customer.CategoryStat{
			CategoryName: categoryMap[cc.CategoryID],
			Count:        cc.Count,
			Ratio:        ratio,
		})
	}

	// 统计人工调整数
	var manualAdjusted int64
	dal.DB.WithContext(ctx).Model(&model.ClassifyAdjustLog{}).
		Where("DATE(create_time) BETWEEN ? AND ?", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Count(&manualAdjusted)

	// 计算自动分类准确率
	autoAccuracy := float64(0)
	if totalClassified > 0 {
		// 准确率 = (总分类数 - 人工调整数) / 总分类数
		autoAccuracy = float64(totalClassified-manualAdjusted) * 100 / float64(totalClassified)
		if autoAccuracy < 0 {
			autoAccuracy = 0
		}
	}

	// 每日统计
	type dailyStat struct {
		Date       string `gorm:"column:date"`
		CategoryID int64  `gorm:"column:category_id"`
		Count      int64  `gorm:"column:count"`
	}
	var dailyStats []dailyStat
	dal.DB.WithContext(ctx).Model(&model.Conversation{}).
		Select("DATE(start_time) as date, category_id, COUNT(*) as count").
		Where("DATE(start_time) BETWEEN ? AND ?", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Where("category_id > 0").
		Group("DATE(start_time), category_id").
		Order("date asc").
		Scan(&dailyStats)

	dailyStatItems := make([]*customer.ClassifyStatItem, 0, len(dailyStats))
	for _, ds := range dailyStats {
		dailyStatItems = append(dailyStatItems, &customer.ClassifyStatItem{
			Date:         ds.Date,
			CategoryId:   ds.CategoryID,
			CategoryName: categoryMap[ds.CategoryID],
			Count:        ds.Count,
		})
	}

	return &customer.GetClassifyStatsResp{
		BaseResp:        &customer.BaseResp{Code: 0, Msg: "success"},
		DailyStats:      dailyStatItems,
		CategorySummary: categorySummary,
		TotalClassified: totalClassified,
		ManualAdjusted:  manualAdjusted,
		AutoAccuracy:    autoAccuracy,
	}, nil
}

// ============ 消息分类维度CRUD ============

// CreateMsgCategory 创建消息分类维度
func (h *CustomerServiceHandler) CreateMsgCategory(ctx context.Context, req *customer.CreateMsgCategoryReq) (*customer.CreateMsgCategoryResp, error) {
	if req == nil || strings.TrimSpace(req.CategoryName) == "" {
		return &customer.CreateMsgCategoryResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "分类名称不能为空"},
		}, nil
	}

	now := time.Now()
	cat := model.MsgCategory{
		CategoryName: strings.TrimSpace(req.CategoryName),
		Keywords:     req.Keywords,
		SortNo:       int(req.SortNo),
		CreateBy:     req.CreateBy,
		CreateTime:   now,
		UpdateTime:   now,
	}

	if err := dal.DB.WithContext(ctx).Create(&cat).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			return &customer.CreateMsgCategoryResp{
				BaseResp: &customer.BaseResp{Code: 400, Msg: "分类名称已存在"},
			}, nil
		}
		return &customer.CreateMsgCategoryResp{
			BaseResp: &customer.BaseResp{Code: 500, Msg: "创建失败: " + err.Error()},
		}, nil
	}

	return &customer.CreateMsgCategoryResp{
		BaseResp:   &customer.BaseResp{Code: 0, Msg: "success"},
		CategoryId: cat.CategoryID,
	}, nil
}

// ListMsgCategory 查询消息分类维度列表
func (h *CustomerServiceHandler) ListMsgCategory(ctx context.Context, req *customer.ListMsgCategoryReq) (*customer.ListMsgCategoryResp, error) {
	var categories []model.MsgCategory
	if err := dal.DB.WithContext(ctx).Model(&model.MsgCategory{}).Order("sort_no asc").Find(&categories).Error; err != nil {
		return &customer.ListMsgCategoryResp{
			BaseResp: &customer.BaseResp{Code: 500, Msg: "查询失败: " + err.Error()},
		}, nil
	}

	result := make([]*customer.MsgCategory, 0, len(categories))
	for _, cat := range categories {
		result = append(result, &customer.MsgCategory{
			CategoryId:   cat.CategoryID,
			CategoryName: cat.CategoryName,
			Keywords:     cat.Keywords,
			SortNo:       int32(cat.SortNo),
		})
	}

	return &customer.ListMsgCategoryResp{
		BaseResp:   &customer.BaseResp{Code: 0, Msg: "success"},
		Categories: result,
		Total:      int64(len(result)),
	}, nil
}

// UpdateMsgCategory 更新消息分类维度
func (h *CustomerServiceHandler) UpdateMsgCategory(ctx context.Context, req *customer.UpdateMsgCategoryReq) (*customer.UpdateMsgCategoryResp, error) {
	if req == nil || req.CategoryId <= 0 {
		return &customer.UpdateMsgCategoryResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "分类ID不能为空"},
		}, nil
	}

	updates := map[string]interface{}{
		"update_time": time.Now(),
	}
	if strings.TrimSpace(req.CategoryName) != "" {
		updates["category_name"] = strings.TrimSpace(req.CategoryName)
	}
	if req.Keywords != "" {
		updates["keywords"] = req.Keywords
	}
	updates["sort_no"] = req.SortNo

	res := dal.DB.WithContext(ctx).Model(&model.MsgCategory{}).Where("category_id = ?", req.CategoryId).Updates(updates)
	if res.Error != nil {
		return &customer.UpdateMsgCategoryResp{
			BaseResp: &customer.BaseResp{Code: 500, Msg: "更新失败: " + res.Error.Error()},
		}, nil
	}
	if res.RowsAffected == 0 {
		return &customer.UpdateMsgCategoryResp{
			BaseResp: &customer.BaseResp{Code: 404, Msg: "分类不存在"},
		}, nil
	}

	return &customer.UpdateMsgCategoryResp{
		BaseResp: &customer.BaseResp{Code: 0, Msg: "success"},
	}, nil
}

// DeleteMsgCategory 删除消息分类维度
func (h *CustomerServiceHandler) DeleteMsgCategory(ctx context.Context, req *customer.DeleteMsgCategoryReq) (*customer.DeleteMsgCategoryResp, error) {
	if req == nil || req.CategoryId <= 0 {
		return &customer.DeleteMsgCategoryResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "分类ID不能为空"},
		}, nil
	}

	res := dal.DB.WithContext(ctx).Where("category_id = ?", req.CategoryId).Delete(&model.MsgCategory{})
	if res.Error != nil {
		return &customer.DeleteMsgCategoryResp{
			BaseResp: &customer.BaseResp{Code: 500, Msg: "删除失败: " + res.Error.Error()},
		}, nil
	}
	if res.RowsAffected == 0 {
		return &customer.DeleteMsgCategoryResp{
			BaseResp: &customer.BaseResp{Code: 404, Msg: "分类不存在"},
		}, nil
	}

	return &customer.DeleteMsgCategoryResp{
		BaseResp: &customer.BaseResp{Code: 0, Msg: "success"},
	}, nil
}

// ============ 消息加密与脱敏 ============

// EncryptMessage 加密消息内容
// 使用AES-256-GCM算法加密敏感消息
func (h *CustomerServiceHandler) EncryptMessage(ctx context.Context, req *customer.EncryptMessageReq) (*customer.EncryptMessageResp, error) {
	if req == nil || req.MsgContent == "" {
		return &customer.EncryptMessageResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "消息内容不能为空"},
		}, nil
	}

	encryptor := crypto.GetEncryptor()
	encrypted, err := encryptor.Encrypt(req.MsgContent)
	if err != nil {
		return &customer.EncryptMessageResp{
			BaseResp: &customer.BaseResp{Code: 500, Msg: "加密失败: " + err.Error()},
		}, nil
	}

	return &customer.EncryptMessageResp{
		BaseResp:         &customer.BaseResp{Code: 0, Msg: "success"},
		EncryptedContent: encrypted,
	}, nil
}

// DecryptMessage 解密消息内容
func (h *CustomerServiceHandler) DecryptMessage(ctx context.Context, req *customer.DecryptMessageReq) (*customer.DecryptMessageResp, error) {
	if req == nil || req.EncryptedContent == "" {
		return &customer.DecryptMessageResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "加密内容不能为空"},
		}, nil
	}

	encryptor := crypto.GetEncryptor()
	decrypted, err := encryptor.Decrypt(req.EncryptedContent)
	if err != nil {
		return &customer.DecryptMessageResp{
			BaseResp: &customer.BaseResp{Code: 500, Msg: "解密失败: " + err.Error()},
		}, nil
	}

	return &customer.DecryptMessageResp{
		BaseResp:   &customer.BaseResp{Code: 0, Msg: "success"},
		MsgContent: decrypted,
	}, nil
}

// DesensitizeMessage 消息脱敏处理
// 对手机号、身份证、银行卡、邮箱等敏感信息进行脱敏
func (h *CustomerServiceHandler) DesensitizeMessage(ctx context.Context, req *customer.DesensitizeMessageReq) (*customer.DesensitizeMessageResp, error) {
	if req == nil || req.MsgContent == "" {
		return &customer.DesensitizeMessageResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "消息内容不能为空"},
		}, nil
	}

	desensitizer := nlp.NewDesensitizer()
	desensitized := desensitizer.Desensitize(req.MsgContent)
	detected := desensitizer.DetectSensitiveInfo(req.MsgContent)

	return &customer.DesensitizeMessageResp{
		BaseResp:            &customer.BaseResp{Code: 0, Msg: "success"},
		DesensitizedContent: desensitized,
		DetectedTypes:       detected,
	}, nil
}

// ============ 数据归档管理 ============

// ArchiveConversations 归档历史会话数据
// 将指定日期之前的会话和消息移动到归档表
func (h *CustomerServiceHandler) ArchiveConversations(ctx context.Context, req *customer.ArchiveConversationsReq) (*customer.ArchiveConversationsResp, error) {
	if req == nil || req.EndDate == "" {
		return &customer.ArchiveConversationsResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "归档截止日期不能为空"},
		}, nil
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return &customer.ArchiveConversationsResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "日期格式错误，应为YYYY-MM-DD"},
		}, nil
	}

	retentionDays := int(req.RetentionDays)
	if retentionDays <= 0 {
		retentionDays = 365 // 默认保留1年
	}

	// 创建归档任务记录
	task := model.ArchiveTask{
		TaskType:   "conv",
		StartDate:  endDate.AddDate(0, -6, 0).Format("2006-01-02"), // 归档6个月前的数据
		EndDate:    req.EndDate,
		Status:     0, // 进行中
		StartTime:  time.Now(),
		OperatorID: req.OperatorId,
	}
	if err := dal.DB.WithContext(ctx).Create(&task).Error; err != nil {
		return &customer.ArchiveConversationsResp{
			BaseResp: &customer.BaseResp{Code: 500, Msg: "创建归档任务失败: " + err.Error()},
		}, nil
	}

	// 异步执行归档（实际生产中应使用消息队列）
	go func() {
		archiveCtx := context.Background()
		var archivedCount int64

		// 查询需要归档的会话
		var convs []model.Conversation
		dal.DB.WithContext(archiveCtx).
			Where("status IN (?, ?) AND end_time < ?", model.ConvStatusEnded, model.ConvStatusAbandoned, endDate).
			Find(&convs)

		for _, conv := range convs {
			// 序列化会话数据
			convData, _ := json.Marshal(conv)

			// 查询消息数量
			var msgCount int64
			dal.DB.WithContext(archiveCtx).Model(&model.ConvMessage{}).
				Where("conv_id = ?", conv.ConvID).Count(&msgCount)

			// 创建归档记录
			archived := model.ArchivedConversation{
				ConvID:        conv.ConvID,
				UserID:        conv.UserID,
				CsID:          conv.CsID,
				ConvData:      string(convData),
				MsgCount:      int(msgCount),
				OriginalDate:  conv.StartTime,
				ArchiveTime:   time.Now(),
				RetentionDays: retentionDays,
			}

			if err := dal.DB.WithContext(archiveCtx).Create(&archived).Error; err == nil {
				archivedCount++
			}
		}

		// 更新任务状态
		dal.DB.WithContext(archiveCtx).Model(&model.ArchiveTask{}).
			Where("task_id = ?", task.TaskID).
			Updates(map[string]interface{}{
				"archived_count": archivedCount,
				"status":         1, // 完成
				"end_time":       time.Now(),
			})
	}()

	return &customer.ArchiveConversationsResp{
		BaseResp:      &customer.BaseResp{Code: 0, Msg: "归档任务已创建"},
		TaskId:        task.TaskID,
		ArchivedCount: 0, // 异步执行，初始为0
	}, nil
}

// GetArchiveTask 查询归档任务状态
func (h *CustomerServiceHandler) GetArchiveTask(ctx context.Context, req *customer.GetArchiveTaskReq) (*customer.GetArchiveTaskResp, error) {
	if req == nil || req.TaskId <= 0 {
		return &customer.GetArchiveTaskResp{
			BaseResp: &customer.BaseResp{Code: 400, Msg: "任务ID不能为空"},
		}, nil
	}

	var task model.ArchiveTask
	if err := dal.DB.WithContext(ctx).Where("task_id = ?", req.TaskId).First(&task).Error; err != nil {
		return &customer.GetArchiveTaskResp{
			BaseResp: &customer.BaseResp{Code: 404, Msg: "任务不存在"},
		}, nil
	}

	return &customer.GetArchiveTaskResp{
		BaseResp:      &customer.BaseResp{Code: 0, Msg: "success"},
		TaskId:        task.TaskID,
		TaskType:      task.TaskType,
		StartDate:     task.StartDate,
		EndDate:       task.EndDate,
		ArchivedCount: task.ArchivedCount,
		DeletedCount:  task.DeletedCount,
		Status:        task.Status,
		ErrorMsg:      task.ErrorMsg,
	}, nil
}

// QueryArchivedConversation 查询已归档会话
func (h *CustomerServiceHandler) QueryArchivedConversation(ctx context.Context, req *customer.QueryArchivedConversationReq) (*customer.QueryArchivedConversationResp, error) {
	if req == nil {
		req = &customer.QueryArchivedConversationReq{Page: 1, PageSize: 20}
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	query := dal.DB.WithContext(ctx).Model(&model.ArchivedConversation{})

	if req.UserId != "" {
		query = query.Where("user_id = ?", req.UserId)
	}
	if req.CsId != "" {
		query = query.Where("cs_id = ?", req.CsId)
	}
	if req.StartDate != "" {
		query = query.Where("original_date >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("original_date <= ?", req.EndDate)
	}

	var total int64
	query.Count(&total)

	var items []model.ArchivedConversation
	offset := (req.Page - 1) * req.PageSize
	query.Offset(int(offset)).Limit(int(req.PageSize)).Order("archive_time desc").Find(&items)

	var result []*customer.ArchivedConvItem
	for _, item := range items {
		result = append(result, &customer.ArchivedConvItem{
			ConvId:       item.ConvID,
			UserId:       item.UserID,
			CsId:         item.CsID,
			MsgCount:     int32(item.MsgCount),
			OriginalDate: item.OriginalDate.Format("2006-01-02"),
			ArchiveTime:  item.ArchiveTime.Format("2006-01-02 15:04:05"),
		})
	}

	return &customer.QueryArchivedConversationResp{
		BaseResp: &customer.BaseResp{Code: 0, Msg: "success"},
		Items:    result,
		Total:    total,
	}, nil
}
