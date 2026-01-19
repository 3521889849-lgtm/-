package service

import (
	"errors"
	"example_shop/common/db"
	"example_shop/common/encrypt"
	"example_shop/common/model/hotel_admin"
	"time"
)

// GuestService 客人信息管理服务
// 负责处理酒店客人信息的创建、更新、查询等核心业务逻辑，
// 包括实名制验证、敏感信息加密存储、在住客人查询、敏感信息脱敏显示等。
type GuestService struct{}

// ListInHouseGuestsReq 在住客人列表查询请求
type ListInHouseGuestsReq struct {
	Page     int     `json:"page"`                // 页码
	PageSize int     `json:"page_size"`           // 每页数量
	BranchID *uint64 `json:"branch_id,omitempty"` // 分店ID，可选
	Province *string `json:"province,omitempty"`  // 省份，可选
	City     *string `json:"city,omitempty"`      // 城市，可选
	District *string `json:"district,omitempty"`  // 区县，可选
	Name     *string `json:"name,omitempty"`      // 姓名，可选
	Phone    *string `json:"phone,omitempty"`     // 手机号，可选
	IDNumber *string `json:"id_number,omitempty"` // 身份证号，可选
	RoomNo   *string `json:"room_no,omitempty"`   // 房间号，可选
}

// InHouseGuestInfo 在住客人信息
type InHouseGuestInfo struct {
	ID           uint64    `json:"id"`
	GuestID      uint64    `json:"guest_id"`
	Name         string    `json:"name"`
	IDType       string    `json:"id_type"`
	IDNumber     string    `json:"id_number"`
	Phone        string    `json:"phone"`
	Province     *string   `json:"province,omitempty"`
	Address      *string   `json:"address,omitempty"`
	Ethnicity    *string   `json:"ethnicity,omitempty"`
	CheckInTime  time.Time `json:"check_in_time"`
	CheckOutTime time.Time `json:"check_out_time"`
	// 订单信息
	OrderID     uint64 `json:"order_id"`
	OrderNo     string `json:"order_no"`
	GuestSource string `json:"guest_source"`
	// 房间信息
	RoomID       uint64 `json:"room_id"`
	RoomNo       string `json:"room_no"`
	RoomTypeID   uint64 `json:"room_type_id"`
	RoomTypeName string `json:"room_type_name"`
	// 财务信息
	OrderAmount       float64 `json:"order_amount"`
	DepositReceived   float64 `json:"deposit_received"`
	OutstandingAmount float64 `json:"outstanding_amount"`
}

// ListInHouseGuestsResp 在住客人列表响应
type ListInHouseGuestsResp struct {
	List     []InHouseGuestInfo `json:"list"`
	Total    uint64             `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

// ListInHouseGuests 获取在住客人列表
// 业务功能：查询当前在住的客人信息，支持多条件筛选和分页，用于在住客人管理和统计场景
// 入参说明：
//   - req: 在住客人列表查询请求，支持按分店、省份、城市、区县、姓名、手机号、身份证号、房间号筛选，支持分页
//
// 返回值说明：
//   - *ListInHouseGuestsResp: 符合条件的在住客人列表（包含订单、房间、财务信息，敏感信息已脱敏）及分页信息
//   - error: 查询失败错误
//
// 业务规则：在住客人的判断标准为：订单状态为"CHECKED_IN"（已入住），且离店时间大于等于当前时间
func (s *GuestService) ListInHouseGuests(req ListInHouseGuestsReq) (*ListInHouseGuestsResp, error) {
	// 业务规则：分页参数默认值设置，页码最小为1，每页数量默认200条（在住客人通常数据量大），最大不超过500条
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 200 // 默认每页200条，50页共10000条
	}
	if req.PageSize > 500 {
		req.PageSize = 500
	}

	offset := (req.Page - 1) * req.PageSize

	// 构建查询：通过JOIN订单表、客人表、房源表、房型表关联查询在住客人信息
	// 业务规则：在住客人的判断标准为订单状态为"CHECKED_IN"（已入住），且离店时间大于等于当前时间
	now := time.Now()
	query := db.MysqlDB.Model(&hotel_admin.OrderMain{}).
		Select(`
			guest_info.id as guest_id,
			guest_info.name,
			guest_info.id_type,
			guest_info.id_number,
			guest_info.phone,
			guest_info.province,
			guest_info.address,
			guest_info.ethnicity,
			hotel_order_main.check_in_time,
			hotel_order_main.check_out_time,
			hotel_order_main.id as order_id,
			hotel_order_main.order_no,
			hotel_order_main.guest_source,
			hotel_order_main.room_id,
			room_info.room_no,
			hotel_order_main.room_type_id,
			room_type_dict.room_type_name,
			hotel_order_main.order_amount,
			hotel_order_main.deposit_received,
			hotel_order_main.outstanding_amount
		`).
		Joins("INNER JOIN guest_info ON hotel_order_main.guest_id = guest_info.id").
		Joins("LEFT JOIN room_info ON hotel_order_main.room_id = room_info.id").
		Joins("LEFT JOIN room_type_dict ON hotel_order_main.room_type_id = room_type_dict.id").
		Where("hotel_order_main.deleted_at IS NULL").
		Where("guest_info.deleted_at IS NULL").
		Where("hotel_order_main.order_status = ?", "CHECKED_IN"). // 已入住状态
		Where("hotel_order_main.check_out_time >= ?", now)        // 离店时间大于等于当前时间

	// 分店筛选
	if req.BranchID != nil {
		query = query.Where("hotel_order_main.branch_id = ?", *req.BranchID)
	}

	// 省份筛选
	if req.Province != nil && *req.Province != "" {
		query = query.Where("guest_info.province = ?", *req.Province)
	}

	// 城市筛选（从地址中提取，如果地址格式是"省份城市区县"）
	if req.City != nil && *req.City != "" {
		query = query.Where("guest_info.address LIKE ?", "%"+*req.City+"%")
	}

	// 区县筛选（从地址中提取）
	if req.District != nil && *req.District != "" {
		query = query.Where("guest_info.address LIKE ?", "%"+*req.District+"%")
	}

	// 姓名筛选
	if req.Name != nil && *req.Name != "" {
		query = query.Where("guest_info.name LIKE ?", "%"+*req.Name+"%")
	}

	// 手机号筛选
	if req.Phone != nil && *req.Phone != "" {
		query = query.Where("guest_info.phone LIKE ?", "%"+*req.Phone+"%")
	}

	// 身份证号筛选
	if req.IDNumber != nil && *req.IDNumber != "" {
		query = query.Where("guest_info.id_number LIKE ?", "%"+*req.IDNumber+"%")
	}

	// 房间号筛选
	if req.RoomNo != nil && *req.RoomNo != "" {
		query = query.Where("room_info.room_no LIKE ?", "%"+*req.RoomNo+"%")
	}

	// 获取总数
	var total int64
	countQuery := db.MysqlDB.Model(&hotel_admin.OrderMain{}).
		Joins("INNER JOIN guest_info ON hotel_order_main.guest_id = guest_info.id").
		Joins("LEFT JOIN room_info ON hotel_order_main.room_id = room_info.id").
		Where("hotel_order_main.deleted_at IS NULL").
		Where("guest_info.deleted_at IS NULL").
		Where("hotel_order_main.order_status = ?", "CHECKED_IN").
		Where("hotel_order_main.check_out_time >= ?", now)

	if req.BranchID != nil {
		countQuery = countQuery.Where("hotel_order_main.branch_id = ?", *req.BranchID)
	}
	if req.Province != nil && *req.Province != "" {
		countQuery = countQuery.Where("guest_info.province = ?", *req.Province)
	}
	if req.City != nil && *req.City != "" {
		countQuery = countQuery.Where("guest_info.address LIKE ?", "%"+*req.City+"%")
	}
	if req.District != nil && *req.District != "" {
		countQuery = countQuery.Where("guest_info.address LIKE ?", "%"+*req.District+"%")
	}
	if req.Name != nil && *req.Name != "" {
		countQuery = countQuery.Where("guest_info.name LIKE ?", "%"+*req.Name+"%")
	}
	if req.Phone != nil && *req.Phone != "" {
		countQuery = countQuery.Where("guest_info.phone LIKE ?", "%"+*req.Phone+"%")
	}
	if req.IDNumber != nil && *req.IDNumber != "" {
		countQuery = countQuery.Where("guest_info.id_number LIKE ?", "%"+*req.IDNumber+"%")
	}
	if req.RoomNo != nil && *req.RoomNo != "" {
		countQuery = countQuery.Where("room_info.room_no LIKE ?", "%"+*req.RoomNo+"%")
	}

	if err := countQuery.Distinct("hotel_order_main.id").Count(&total).Error; err != nil {
		return nil, err
	}

	// 排序：按入住时间倒序
	query = query.Order("hotel_order_main.check_in_time DESC") // 添加排序条件，按订单入住时间倒序排列（最新入住的客人排在前面）

	// 分页
	query = query.Offset(offset).Limit(req.PageSize) // 添加分页限制（偏移量、每页数量），用于SQL查询的OFFSET和LIMIT子句

	// 查询结果
	type GuestResult struct { // 定义查询结果结构体，用于存储从数据库查询到的原始数据
		GuestID           uint64    `gorm:"column:guest_id"`           // 客人ID（数据库列名：guest_id）
		Name              string    `gorm:"column:name"`               // 姓名（数据库列名：name）
		IDType            string    `gorm:"column:id_type"`            // 证件类型（数据库列名：id_type）
		IDNumber          string    `gorm:"column:id_number"`          // 证件号（数据库列名：id_number，已解密）
		Phone             string    `gorm:"column:phone"`              // 手机号（数据库列名：phone，已解密）
		Province          *string   `gorm:"column:province"`           // 省份（数据库列名：province，可为空）
		Address           *string   `gorm:"column:address"`            // 地址（数据库列名：address，可为空）
		Ethnicity         *string   `gorm:"column:ethnicity"`          // 民族（数据库列名：ethnicity，可为空）
		CheckInTime       time.Time `gorm:"column:check_in_time"`      // 入住时间（数据库列名：check_in_time）
		CheckOutTime      time.Time `gorm:"column:check_out_time"`     // 离店时间（数据库列名：check_out_time）
		OrderID           uint64    `gorm:"column:order_id"`           // 订单ID（数据库列名：order_id）
		OrderNo           string    `gorm:"column:order_no"`           // 订单号（数据库列名：order_no）
		GuestSource       string    `gorm:"column:guest_source"`       // 客人来源（数据库列名：guest_source）
		RoomID            uint64    `gorm:"column:room_id"`            // 房间ID（数据库列名：room_id）
		RoomNo            *string   `gorm:"column:room_no"`            // 房间号（数据库列名：room_no，可为空）
		RoomTypeID        uint64    `gorm:"column:room_type_id"`       // 房型ID（数据库列名：room_type_id）
		RoomTypeName      *string   `gorm:"column:room_type_name"`     // 房型名称（数据库列名：room_type_name，可为空）
		OrderAmount       float64   `gorm:"column:order_amount"`       // 订单金额（数据库列名：order_amount）
		DepositReceived   float64   `gorm:"column:deposit_received"`   // 已收押金（数据库列名：deposit_received）
		OutstandingAmount float64   `gorm:"column:outstanding_amount"` // 未付金额（数据库列名：outstanding_amount）
	}

	var results []GuestResult                          // 声明查询结果列表变量，用于存储从数据库查询到的原始数据列表
	if err := query.Scan(&results).Error; err != nil { // 执行查询并将结果扫描到结果结构体列表中，如果查询失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	// 转换为返回格式
	guests := make([]InHouseGuestInfo, len(results)) // 创建在住客人信息列表，长度为查询结果数量
	for i, r := range results {                      // 遍历查询结果列表
		guest := InHouseGuestInfo{ // 创建在住客人信息对象
			ID:                r.GuestID,           // 设置ID（使用客人ID）
			GuestID:           r.GuestID,           // 设置客人ID（从查询结果中获取）
			Name:              r.Name,              // 设置姓名（从查询结果中获取）
			IDType:            r.IDType,            // 设置证件类型（从查询结果中获取）
			IDNumber:          r.IDNumber,          // 设置证件号（从查询结果中获取，已解密）
			Phone:             r.Phone,             // 设置手机号（从查询结果中获取，已解密）
			CheckInTime:       r.CheckInTime,       // 设置入住时间（从查询结果中获取）
			CheckOutTime:      r.CheckOutTime,      // 设置离店时间（从查询结果中获取）
			OrderID:           r.OrderID,           // 设置订单ID（从查询结果中获取）
			OrderNo:           r.OrderNo,           // 设置订单号（从查询结果中获取）
			GuestSource:       r.GuestSource,       // 设置客人来源（从查询结果中获取）
			RoomID:            r.RoomID,            // 设置房间ID（从查询结果中获取）
			RoomTypeID:        r.RoomTypeID,        // 设置房型ID（从查询结果中获取）
			OrderAmount:       r.OrderAmount,       // 设置订单金额（从查询结果中获取）
			DepositReceived:   r.DepositReceived,   // 设置已收押金（从查询结果中获取）
			OutstandingAmount: r.OutstandingAmount, // 设置未付金额（从查询结果中获取）
		}

		if r.Province != nil { // 如果查询结果中省份不为空（指针非空），则设置省份
			guest.Province = r.Province // 设置省份（从查询结果中获取）
		}
		if r.Address != nil { // 如果查询结果中地址不为空（指针非空），则设置地址
			guest.Address = r.Address // 设置地址（从查询结果中获取）
		}
		if r.Ethnicity != nil { // 如果查询结果中民族不为空（指针非空），则设置民族
			guest.Ethnicity = r.Ethnicity // 设置民族（从查询结果中获取）
		}
		if r.RoomNo != nil { // 如果查询结果中房间号不为空（指针非空），则设置房间号
			guest.RoomNo = *r.RoomNo // 设置房间号（解引用指针获取值）
		}
		if r.RoomTypeName != nil { // 如果查询结果中房型名称不为空（指针非空），则设置房型名称
			guest.RoomTypeName = *r.RoomTypeName // 设置房型名称（解引用指针获取值）
		}

		// 业务规则：敏感信息脱敏处理，在返回给前端前对身份证号、手机号、姓名进行脱敏（如中间部分用*替代）
		DesensitizeGuestInfo(&guest) // 调用脱敏函数，对客人信息中的敏感字段进行脱敏处理（身份证号、手机号、姓名）

		guests[i] = guest // 将转换后的在住客人信息对象添加到列表中
	}

	return &ListInHouseGuestsResp{ // 返回在住客人列表响应对象
		List:     guests,        // 设置在住客人列表（转换后的在住客人信息列表，敏感信息已脱敏）
		Total:    uint64(total), // 设置总数（转换为uint64类型）
		Page:     req.Page,      // 设置当前页码（从请求中获取）
		PageSize: req.PageSize,  // 设置每页数量（从请求中获取）
	}, nil // 返回响应对象和无错误
}

// CreateGuestInfoReq 创建客人信息请求
type CreateGuestInfoReq struct {
	Name         string     `json:"name"`      // 姓名
	IDType       string     `json:"id_type"`   // 证件类型
	IDNumber     string     `json:"id_number"` // 证件号（明文，需要加密存储）
	Phone        string     `json:"phone"`     // 手机号（明文，需要加密存储）
	Gender       *string    `json:"gender,omitempty"`
	Ethnicity    *string    `json:"ethnicity,omitempty"`
	Province     *string    `json:"province,omitempty"`
	Address      *string    `json:"address,omitempty"`
	CheckInTime  *time.Time `json:"check_in_time,omitempty"`
	CheckOutTime *time.Time `json:"check_out_time,omitempty"`
	RoomID       *uint64    `json:"room_id,omitempty"`
	OrderID      *uint64    `json:"order_id,omitempty"`
	RegisterBy   uint64     `json:"register_by"` // 登记人ID
	IsMember     bool       `json:"is_member"`
	MemberID     *uint64    `json:"member_id,omitempty"`
}

// CreateGuestInfo 创建客人信息（实名制验证和加密存储）
// 业务功能：创建新的客人信息记录，对身份证号和手机号进行格式验证和加密存储，支持会员关联
// 入参说明：
//   - req: 创建客人信息请求，包含姓名、证件类型、证件号、手机号、性别、民族、省份、地址等，可选关联订单、房源、会员
//
// 返回值说明：
//   - *hotel_admin.GuestInfo: 成功创建后的客人完整信息（敏感信息已加密）
//   - error: 身份证号格式不正确、手机号格式不正确、加密失败或数据库操作错误
func (s *GuestService) CreateGuestInfo(req CreateGuestInfoReq) (*hotel_admin.GuestInfo, error) {
	// 业务规则：实名制验证，身份证号必须符合18位身份证格式（如证件类型为身份证）
	if req.IDType == "身份证" || req.IDType == "ID_CARD" { // 如果证件类型为身份证（中文或英文标识）
		if !encrypt.ValidateIDNumber(req.IDNumber) { // 调用身份证号格式验证函数，如果验证失败则返回错误
			return nil, errors.New("身份证号格式不正确") // 返回nil和错误信息，表示身份证号格式不正确
		}
	}
	// 业务规则：手机号必须符合11位手机号格式
	if !encrypt.ValidatePhone(req.Phone) { // 调用手机号格式验证函数，如果验证失败则返回错误
		return nil, errors.New("手机号格式不正确") // 返回nil和错误信息，表示手机号格式不正确
	}

	// 业务规则：敏感信息加密存储，身份证号和手机号必须使用AES-GCM加密算法加密后存储
	encryptedIDNumber, err := encrypt.Encrypt(req.IDNumber) // 调用加密函数，对明文身份证号进行AES-GCM加密
	if err != nil {                                         // 如果加密失败，则返回错误
		return nil, errors.New("加密身份证号失败: " + err.Error()) // 返回nil和错误信息，表示加密身份证号失败
	}
	encryptedPhone, err := encrypt.Encrypt(req.Phone) // 调用加密函数，对明文手机号进行AES-GCM加密
	if err != nil {                                   // 如果加密失败，则返回错误
		return nil, errors.New("加密手机号失败: " + err.Error()) // 返回nil和错误信息，表示加密手机号失败
	}

	// 创建客人信息实体，设置注册时间为当前时间
	guest := hotel_admin.GuestInfo{ // 创建客人实体对象
		Name:         req.Name,          // 设置姓名（从请求中获取）
		IDType:       req.IDType,        // 设置证件类型（从请求中获取）
		IDNumber:     encryptedIDNumber, // 设置证件号（存储加密后的数据，而不是明文）
		Phone:        encryptedPhone,    // 设置手机号（存储加密后的数据，而不是明文）
		Gender:       req.Gender,        // 设置性别（从请求中获取，可为空）
		Ethnicity:    req.Ethnicity,     // 设置民族（从请求中获取，可为空）
		Province:     req.Province,      // 设置省份（从请求中获取，可为空）
		Address:      req.Address,       // 设置地址（从请求中获取，可为空）
		CheckInTime:  req.CheckInTime,   // 设置入住时间（从请求中获取，可为空）
		CheckOutTime: req.CheckOutTime,  // 设置离店时间（从请求中获取，可为空）
		RoomID:       req.RoomID,        // 设置房间ID（从请求中获取，可为空）
		OrderID:      req.OrderID,       // 设置订单ID（从请求中获取，可为空）
		RegisterBy:   req.RegisterBy,    // 设置登记人ID（从请求中获取）
		RegisterTime: time.Now(),        // 设置登记时间为当前时间（自动生成）
		IsMember:     req.IsMember,      // 设置是否为会员（从请求中获取）
		MemberID:     req.MemberID,      // 设置会员ID（从请求中获取，可为空）
	}

	if err := db.MysqlDB.Create(&guest).Error; err != nil { // 将客人信息保存到数据库，如果保存失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	return &guest, nil // 返回客人实体指针和无错误
}

// UpdateGuestInfoReq 更新客人信息请求
type UpdateGuestInfoReq struct {
	ID           uint64     `json:"id"`
	Name         *string    `json:"name,omitempty"`
	IDType       *string    `json:"id_type,omitempty"`
	IDNumber     *string    `json:"id_number,omitempty"` // 如果提供，需要加密
	Phone        *string    `json:"phone,omitempty"`     // 如果提供，需要加密
	Gender       *string    `json:"gender,omitempty"`
	Ethnicity    *string    `json:"ethnicity,omitempty"`
	Province     *string    `json:"province,omitempty"`
	Address      *string    `json:"address,omitempty"`
	CheckInTime  *time.Time `json:"check_in_time,omitempty"`
	CheckOutTime *time.Time `json:"check_out_time,omitempty"`
}

// UpdateGuestInfo 更新客人信息
func (s *GuestService) UpdateGuestInfo(req UpdateGuestInfoReq) error {
	var guest hotel_admin.GuestInfo                                // 声明客人实体变量，用于存储查询到的客人信息
	if err := db.MysqlDB.First(&guest, req.ID).Error; err != nil { // 通过客人ID查询客人信息，如果查询失败则说明客人不存在
		return err // 返回数据库查询错误
	}

	updates := make(map[string]interface{}) // 创建更新字段映射表，用于存储需要更新的字段和值
	if req.Name != nil {                    // 如果请求中提供了姓名（指针非空），则添加到更新映射表
		updates["name"] = *req.Name // 添加姓名到更新映射表（解引用指针获取值）
	}
	if req.IDType != nil { // 如果请求中提供了证件类型（指针非空），则添加到更新映射表
		updates["id_type"] = *req.IDType // 添加证件类型到更新映射表（解引用指针获取值）
	}
	if req.IDNumber != nil { // 如果请求中提供了证件号（指针非空），则需要验证并加密
		// 验证并加密
		if *req.IDType == "身份证" || *req.IDType == "ID_CARD" { // 如果证件类型是身份证（解引用指针获取值）
			if !encrypt.ValidateIDNumber(*req.IDNumber) { // 调用身份证号格式验证函数，如果验证失败则返回错误（解引用指针获取值）
				return errors.New("身份证号格式不正确") // 返回错误信息，表示身份证号格式不正确
			}
		}
		encrypted, err := encrypt.Encrypt(*req.IDNumber) // 调用加密函数，对明文身份证号进行AES-GCM加密（解引用指针获取值）
		if err != nil {                                  // 如果加密失败，则返回错误
			return errors.New("加密身份证号失败: " + err.Error()) // 返回错误信息，表示加密身份证号失败
		}
		updates["id_number"] = encrypted // 添加加密后的证件号到更新映射表
	}
	if req.Phone != nil { // 如果请求中提供了手机号（指针非空），则需要验证并加密
		if !encrypt.ValidatePhone(*req.Phone) { // 调用手机号格式验证函数，如果验证失败则返回错误（解引用指针获取值）
			return errors.New("手机号格式不正确") // 返回错误信息，表示手机号格式不正确
		}
		encrypted, err := encrypt.Encrypt(*req.Phone) // 调用加密函数，对明文手机号进行AES-GCM加密（解引用指针获取值）
		if err != nil {                               // 如果加密失败，则返回错误
			return errors.New("加密手机号失败: " + err.Error()) // 返回错误信息，表示加密手机号失败
		}
		updates["phone"] = encrypted // 添加加密后的手机号到更新映射表
	}
	if req.Gender != nil { // 如果请求中提供了性别（指针非空），则添加到更新映射表
		updates["gender"] = *req.Gender // 添加性别到更新映射表（解引用指针获取值）
	}
	if req.Ethnicity != nil { // 如果请求中提供了民族（指针非空），则添加到更新映射表
		updates["ethnicity"] = *req.Ethnicity // 添加民族到更新映射表（解引用指针获取值）
	}
	if req.Province != nil { // 如果请求中提供了省份（指针非空），则添加到更新映射表
		updates["province"] = *req.Province // 添加省份到更新映射表（解引用指针获取值）
	}
	if req.Address != nil { // 如果请求中提供了地址（指针非空），则添加到更新映射表
		updates["address"] = *req.Address // 添加地址到更新映射表（解引用指针获取值）
	}
	if req.CheckInTime != nil { // 如果请求中提供了入住时间（指针非空），则添加到更新映射表
		updates["check_in_time"] = *req.CheckInTime // 添加入住时间到更新映射表（解引用指针获取值）
	}
	if req.CheckOutTime != nil { // 如果请求中提供了离店时间（指针非空），则添加到更新映射表
		updates["check_out_time"] = *req.CheckOutTime // 添加离店时间到更新映射表（解引用指针获取值）
	}

	return db.MysqlDB.Model(&guest).Updates(updates).Error // 根据客人实体更新客人信息，使用更新映射表中的字段和值，返回更新操作的结果（成功为nil，失败为error）
}

// DesensitizeGuestInfo 对客人信息进行脱敏处理（用于前端显示）
func DesensitizeGuestInfo(guest *InHouseGuestInfo) {
	// 脱敏身份证号
	if guest.IDNumber != "" { // 如果客人信息中身份证号不为空，则进行脱敏处理
		// 尝试解密，如果失败则直接脱敏加密后的字符串
		decrypted, err := encrypt.Decrypt(guest.IDNumber) // 调用解密函数，尝试解密身份证号（可能是加密存储的）
		if err == nil {                                   // 如果解密成功（未报错），说明身份证号是加密存储的
			guest.IDNumber = encrypt.DesensitizeIDNumber(decrypted) // 对解密后的身份证号进行脱敏处理（如中间部分用*替代），并更新客人信息中的身份证号
		} else {
			// 如果解密失败，可能是明文，直接脱敏
			guest.IDNumber = encrypt.DesensitizeIDNumber(guest.IDNumber) // 对明文身份证号进行脱敏处理（如中间部分用*替代），并更新客人信息中的身份证号
		}
	}
	// 脱敏手机号
	if guest.Phone != "" { // 如果客人信息中手机号不为空，则进行脱敏处理
		decrypted, err := encrypt.Decrypt(guest.Phone) // 调用解密函数，尝试解密手机号（可能是加密存储的）
		if err == nil {                                // 如果解密成功（未报错），说明手机号是加密存储的
			guest.Phone = encrypt.DesensitizePhone(decrypted) // 对解密后的手机号进行脱敏处理（如中间部分用*替代），并更新客人信息中的手机号
		} else {
			guest.Phone = encrypt.DesensitizePhone(guest.Phone) // 对明文手机号进行脱敏处理（如中间部分用*替代），并更新客人信息中的手机号
		}
	}
	// 脱敏姓名
	if guest.Name != "" { // 如果客人信息中姓名不为空，则进行脱敏处理
		guest.Name = encrypt.DesensitizeName(guest.Name) // 对姓名进行脱敏处理（如中间部分用*替代），并更新客人信息中的姓名
	}
}
