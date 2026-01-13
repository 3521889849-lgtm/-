package model

import "time"

// SysOperationLog 系统操作日志表（合规审计）
type SysOperationLog struct {
	ID               int64     `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	UserID           int       `gorm:"not null;index:idx_user_id;comment:关联用户ID（sys_user.id）" json:"user_id"`
	OperationModule  string    `gorm:"type:varchar(50);not null;index:idx_operation_module;comment:操作模块（如：房源管理、订单处理）" json:"operation_module"`
	OperationType    string    `gorm:"type:varchar(20);not null;comment:操作类型（如：添加、修改、删除、查询）" json:"operation_type"`
	OperationContent string    `gorm:"type:text;comment:操作内容（如：修改房源8101状态为维修房）" json:"operation_content"`
	OperationIP      string    `gorm:"type:varchar(50);comment:操作IP" json:"operation_ip"`
	OperationTime    time.Time `gorm:"not null;index:idx_operation_time;comment:操作时间" json:"operation_time"`
	CreatedAt        time.Time `gorm:"not null;comment:创建时间" json:"created_at"`

	// 关联关系
	User *SysUser `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
}

// TableName 指定表名
func (SysOperationLog) TableName() string {
	return "sys_operation_log"
}
