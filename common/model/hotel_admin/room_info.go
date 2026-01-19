// Package hotel_admin 提供酒店管理系统的数据模型定义
//
// 本文件定义了房源信息的数据模型
//
// 功能说明：
//   - 存储单个房间的基础信息（房号、房型、价格等）
//   - 支持房源的添加、修改、删除功能
//   - 支持房间设施、图片、状态等关联信息管理
//   - 支持退订政策配置
//   - 与订单、客人、财务等模块关联
package hotel_admin

import (
	"time"

	"gorm.io/gorm"
)

// ==================== 房源信息模型 ====================

// RoomInfo 房源信息表
//
// 业务用途：
//   - 管理酒店的所有房间信息
//   - 支持房源的增删改查操作
//   - 作为订单的关联主体（客人预订的是具体房间）
//   - 支持房间状态管理（空闲、已入住、维修等）
//   - 支持动态定价（门市价、日历价）
//
// 设计说明：
//   - 每个房间属于一个分店
//   - 每个房间对应一个房型（房型定义了房间的类别和标准）
//   - 支持软删除，删除后可恢复
//   - 关联多个子表：设施、图片、状态明细等
type RoomInfo struct {
	// ========== 基础字段 ==========
	
	// ID 房源ID，主键，自增
	ID uint64 `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;autoIncrement;comment:房源ID" json:"id"`
	
	// BranchID 所属分店ID，外键关联 hotel_branch 表
	// 用途：实现分店级别的数据隔离
	// 示例：1（锦江之星北京店）、2（锦江之星上海店）
	BranchID uint64 `gorm:"column:branch_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_branch_id;comment:分店ID（外键，关联分店表）" json:"branch_id"`
	
	// RoomTypeID 房型ID，外键关联 room_type_dict 表
	// 用途：定义房间的类型和标准（如：标准间、豪华套房）
	// 同一房型的房间共享相同的设施和服务标准
	RoomTypeID uint64 `gorm:"column:room_type_id;type:BIGINT UNSIGNED;NOT NULL;index:idx_room_type_id;comment:房型ID（外键，关联房型表）" json:"room_type_id"`
	
	// ========== 房间标识 ==========
	
	// RoomNo 房间号，房间的唯一标识
	// 示例："101"、"2A08"、"总统套房A"
	// 用途：前台展示、客人入住、房态图显示
	RoomNo string `gorm:"column:room_no;type:VARCHAR(20);NOT NULL;index:idx_room_no;comment:房间号" json:"room_no"`
	
	// RoomName 房间名称，房间的描述性名称
	// 示例："舒适大床房"、"豪华海景套房"
	// 用途：对外展示、营销推广
	RoomName string `gorm:"column:room_name;type:VARCHAR(100);NOT NULL;comment:房间名称" json:"room_name"`
	
	// ========== 价格信息 ==========
	
	// MarketPrice 门市价，房间的标准价格
	// 单位：元（人民币）
	// 用途：作为定价基准，通常用于线下直接预订
	// 示例：298.00（表示298元/晚）
	MarketPrice float64 `gorm:"column:market_price;type:DECIMAL(10,2);NOT NULL;comment:门市价" json:"market_price"`
	
	// CalendarPrice 日历价，根据日期动态调整的价格
	// 单位：元（人民币）
	// 用途：实现动态定价策略（如：周末价、节假日价、淡旺季价）
	// 示例：258.00（表示今日特价258元/晚）
	CalendarPrice float64 `gorm:"column:calendar_price;type:DECIMAL(10,2);NOT NULL;comment:日历价" json:"calendar_price"`
	
	// ========== 房间规格 ==========
	
	// RoomCount 房间数量，该房源包含的房间数
	// 通常为1，某些套房可能包含多个独立房间
	// 示例：1（单间）、2（套房含两个房间）
	RoomCount uint8 `gorm:"column:room_count;type:TINYINT UNSIGNED;NOT NULL;default:1;comment:房间数量" json:"room_count"`
	
	// Area 房间面积，可选字段
	// 单位：平方米（㎡）
	// 示例：25.50（表示25.5平方米）
	// 用途：房间展示、筛选、排序
	Area *float64 `gorm:"column:area;type:DECIMAL(5,2);comment:面积（平方米）" json:"area,omitempty"`
	
	// BedSpec 床型规格，描述床的类型和数量
	// 示例："1张大床"、"2张单人床"、"1张大床+1张沙发床"
	// 用途：客人选房参考、房态图显示
	BedSpec string `gorm:"column:bed_spec;type:VARCHAR(50);NOT NULL;comment:床型规格" json:"bed_spec"`
	
	// ========== 服务配置 ==========
	
	// HasBreakfast 是否含早餐
	// true: 房费包含早餐
	// false: 不含早餐（可单独购买）
	// 用途：影响房费价格、客人选房决策
	HasBreakfast bool `gorm:"column:has_breakfast;type:BOOLEAN;NOT NULL;default:false;comment:是否含早（0/1）" json:"has_breakfast"`
	
	// HasToiletries 是否提供洗漱用品
	// true: 提供（如：牙刷、牙膏、沐浴液、洗发水）
	// false: 不提供（环保要求或需客人自备）
	// 用途：服务标准说明、客人入住提醒
	HasToiletries bool `gorm:"column:has_toiletries;type:BOOLEAN;NOT NULL;default:false;comment:是否提供洗漱用品（0/1）" json:"has_toiletries"`
	
	// ========== 政策配置 ==========
	
	// CancellationPolicyID 退订政策ID，外键关联 cancellation_policy 表
	// 可选字段，如果为NULL则使用分店默认政策
	// 用途：定义该房间的退订规则（如：提前24小时免费取消）
	CancellationPolicyID *uint64 `gorm:"column:cancellation_policy_id;type:BIGINT UNSIGNED;index:idx_cancellation_policy_id;comment:退订政策ID（外键，关联退订政策表）" json:"cancellation_policy_id,omitempty"`
	
	// ========== 状态控制 ==========
	
	// Status 房间状态，控制房间是否可预订
	// 可选值：
	//   - ACTIVE: 启用，正常可预订
	//   - INACTIVE: 停用，不可预订（如：临时下架）
	//   - MAINTENANCE: 维修中，不可预订（如：设施维修、清洁保养）
	// 用途：房态管理、库存控制
	Status string `gorm:"column:status;type:VARCHAR(20);NOT NULL;default:'ACTIVE';index:idx_status;comment:状态：ACTIVE-启用，INACTIVE-停用，MAINTENANCE-维修" json:"status"`
	
	// ========== 时间戳 ==========
	
	// CreatedAt 创建时间，记录房源信息首次创建的时间
	CreatedAt time.Time `gorm:"column:created_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	
	// UpdatedAt 修改时间，每次修改房源信息时自动更新
	UpdatedAt time.Time `gorm:"column:updated_at;type:DATETIME;NOT NULL;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:修改时间" json:"updated_at"`
	
	// CreatedBy 创建人ID，关联用户账号表
	// 用于审计追溯，记录是谁创建的这个房源
	CreatedBy uint64 `gorm:"column:created_by;type:BIGINT UNSIGNED;NOT NULL;comment:创建人" json:"created_by"`
	
	// DeletedAt 软删除时间，非NULL表示已删除
	// 软删除：数据不会真正删除，只是标记为已删除状态
	// 好处：数据可恢复，历史订单关联完整性保证
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME;index;comment:软删除时间" json:"deleted_at,omitempty"`

	// ========== 关联关系 ==========
	// 以下字段不会存储在数据库中，仅用于GORM的关联查询
	
	// Branch 所属分店信息
	// 多对一关系：多个房间属于同一个分店
	Branch *HotelBranch `gorm:"foreignKey:BranchID;references:ID" json:"branch,omitempty"`
	
	// RoomType 房型信息
	// 多对一关系：多个房间属于同一个房型
	RoomType *RoomTypeDict `gorm:"foreignKey:RoomTypeID;references:ID" json:"room_type,omitempty"`
	
	// CancellationPolicy 退订政策信息
	// 多对一关系：多个房间可以使用同一个退订政策
	CancellationPolicy *CancellationPolicy `gorm:"foreignKey:CancellationPolicyID;references:ID" json:"cancellation_policy,omitempty"`
	
	// RoomFacilityRelations 房间-设施关联关系
	// 一对多关系：一个房间可以有多个设施（如：空调、电视、WiFi）
	RoomFacilityRelations []RoomFacilityRelation `gorm:"foreignKey:RoomID;references:ID" json:"room_facility_relations,omitempty"`
	
	// RoomImages 房间图片
	// 一对多关系：一个房间可以有多张展示图片
	RoomImages []RoomImage `gorm:"foreignKey:RoomID;references:ID" json:"room_images,omitempty"`
	
	// RelatedRoomBindings 关联房间绑定
	// 一对多关系：某些房型可以关联其他房型（如：套房关联主卧和次卧）
	RelatedRoomBindings []RelatedRoomBinding `gorm:"foreignKey:MainRoomID;references:ID" json:"related_room_bindings,omitempty"`
	
	// RoomStatusDetails 房间状态明细
	// 一对多关系：记录房间的历史状态变化（如：每日的入住、退房、清洁状态）
	RoomStatusDetails []RoomStatusDetail `gorm:"foreignKey:RoomID;references:ID" json:"room_status_details,omitempty"`
	
	// Orders 该房间的订单记录
	// 一对多关系：一个房间可以有多个历史订单
	Orders []OrderMain `gorm:"foreignKey:RoomID;references:ID" json:"orders,omitempty"`
	
	// GuestInfos 入住该房间的客人信息
	// 一对多关系：一个房间可以有多个历史入住客人记录
	GuestInfos []GuestInfo `gorm:"foreignKey:RoomID;references:ID" json:"guest_infos,omitempty"`
	
	// FinancialFlows 该房间相关的财务流水
	// 一对多关系：一个房间可以产生多条财务记录（如：房费、加床费、损坏赔偿）
	FinancialFlows []FinancialFlow `gorm:"foreignKey:RoomID;references:ID" json:"financial_flows,omitempty"`
}

// ==================== 表名配置 ====================

// TableName 指定数据库表名
//
// GORM会自动调用此方法获取表名，用于生成SQL语句
//
// 返回：
//   - string: 数据库表名 "room_info"
func (RoomInfo) TableName() string {
	return "room_info"
}
