package service

import (
	"errors"
	"example_shop/common/db"
	"example_shop/common/model/hotel_admin"
)

// RoomBindingService 房源关联绑定管理服务
// 负责处理房源与房源之间的关联绑定关系的创建、删除等核心业务逻辑，
// 包括主房源与关联房源的绑定校验、绑定关系去重检查等。
type RoomBindingService struct{}

// CreateRoomBinding 创建关联房绑定
func (s *RoomBindingService) CreateRoomBinding(mainRoomID uint64, relatedRoomID uint64, bindingDesc *string) error {
	// 检查主房源是否存在
	var mainRoom hotel_admin.RoomInfo                                     // 声明房源实体变量，用于存储查询到的主房源信息
	if err := db.MysqlDB.First(&mainRoom, mainRoomID).Error; err != nil { // 通过主房源ID查询房源信息，如果查询失败则说明主房源不存在
		return errors.New("主房源不存在") // 返回错误信息，表示主房源不存在
	}

	// 检查关联房源是否存在
	var relatedRoom hotel_admin.RoomInfo                                        // 声明房源实体变量，用于存储查询到的关联房源信息
	if err := db.MysqlDB.First(&relatedRoom, relatedRoomID).Error; err != nil { // 通过关联房源ID查询房源信息，如果查询失败则说明关联房源不存在
		return errors.New("关联房源不存在") // 返回错误信息，表示关联房源不存在
	}

	// 检查是否是自己关联自己
	if mainRoomID == relatedRoomID { // 如果主房源ID等于关联房源ID（自己关联自己）
		return errors.New("不能将房源关联到自己") // 返回错误信息，表示不能将房源关联到自己
	}

	// 检查是否已经存在绑定关系
	var existBinding hotel_admin.RelatedRoomBinding                                                                                               // 声明绑定实体变量，用于存储查询到的已存在绑定信息
	result := db.MysqlDB.Where("main_room_id = ? AND related_room_id = ? AND deleted_at IS NULL", mainRoomID, relatedRoomID).First(&existBinding) // 查询是否已存在该主房源和关联房源的绑定关系（排除已删除绑定）
	if result.Error == nil {                                                                                                                      // 如果查询成功（未报错），说明该绑定关系已存在
		return errors.New("该关联关系已存在") // 返回错误信息，表示该关联关系已存在
	}

	// 检查反向绑定是否已存在
	result = db.MysqlDB.Where("main_room_id = ? AND related_room_id = ? AND deleted_at IS NULL", relatedRoomID, mainRoomID).First(&existBinding) // 查询是否已存在反向绑定关系（关联房源作为主房源，主房源作为关联房源），如果查询失败则说明反向绑定不存在
	if result.Error == nil {                                                                                                                     // 如果查询成功（未报错），说明反向绑定关系已存在
		return errors.New("该关联关系已存在（反向）") // 返回错误信息，表示该关联关系已存在（反向）
	}

	binding := &hotel_admin.RelatedRoomBinding{ // 创建房源绑定实体对象指针
		MainRoomID:    mainRoomID,    // 设置主房源ID
		RelatedRoomID: relatedRoomID, // 设置关联房源ID
		BindingDesc:   bindingDesc,   // 设置绑定描述（从参数中获取，可为空）
	}

	return db.MysqlDB.Create(binding).Error // 将绑定关系保存到数据库，返回保存操作的结果（成功为nil，失败为error）
}

// DeleteRoomBinding 删除关联房绑定
func (s *RoomBindingService) DeleteRoomBinding(bindingID uint64) error {
	var binding hotel_admin.RelatedRoomBinding                          // 声明绑定实体变量，用于存储查询到的绑定信息
	if err := db.MysqlDB.First(&binding, bindingID).Error; err != nil { // 通过绑定ID查询绑定信息，如果查询失败则说明绑定关系不存在
		return errors.New("绑定关系不存在") // 返回错误信息，表示绑定关系不存在
	}

	return db.MysqlDB.Delete(&binding).Error // 执行软删除操作（设置deleted_at字段），根据绑定ID删除绑定记录，返回删除操作的结果（成功为nil，失败为error）
}

// GetRoomBindings 获取房源的关联房列表
func (s *RoomBindingService) GetRoomBindings(roomID uint64) ([]hotel_admin.RelatedRoomBinding, error) {
	var bindings []hotel_admin.RelatedRoomBinding                                  // 声明绑定列表变量，用于存储查询到的绑定信息列表
	if err := db.MysqlDB.Where("main_room_id = ? AND deleted_at IS NULL", roomID). // 添加筛选条件：主房源ID匹配且未被删除
											Preload("RelatedRoom").             // 预加载关联房源信息关联数据（JOIN查询关联房源信息）
											Find(&bindings).Error; err != nil { // 执行查询并获取符合条件的绑定列表（包含所有预加载的关联数据），如果查询失败则返回错误
		return nil, err // 返回nil和错误信息
	}
	return bindings, nil // 返回绑定列表和无错误
}

// BatchCreateRoomBindings 批量创建关联房绑定
func (s *RoomBindingService) BatchCreateRoomBindings(mainRoomID uint64, relatedRoomIDs []uint64, bindingDesc *string) error {
	if len(relatedRoomIDs) == 0 { // 如果关联房源ID列表为空（长度为0）
		return errors.New("关联房源ID列表不能为空") // 返回错误信息，表示关联房源ID列表不能为空
	}

	// 检查主房源是否存在
	var mainRoom hotel_admin.RoomInfo                                     // 声明房源实体变量，用于存储查询到的主房源信息
	if err := db.MysqlDB.First(&mainRoom, mainRoomID).Error; err != nil { // 通过主房源ID查询房源信息，如果查询失败则说明主房源不存在
		return errors.New("主房源不存在") // 返回错误信息，表示主房源不存在
	}

	// 检查所有关联房源是否存在
	var count int64                                                                          // 声明计数变量，用于存储存在的关联房源数量
	db.MysqlDB.Model(&hotel_admin.RoomInfo{}).Where("id IN ?", relatedRoomIDs).Count(&count) // 统计关联房源ID列表中存在的房源数量（使用IN子句）
	if int(count) != len(relatedRoomIDs) {                                                   // 如果存在的房源数量不等于请求的房源数量，说明部分关联房源不存在
		return errors.New("部分关联房源不存在") // 返回错误信息，表示部分关联房源不存在
	}

	// 检查是否包含自己
	for _, id := range relatedRoomIDs { // 遍历关联房源ID列表
		if id == mainRoomID { // 如果关联房源ID等于主房源ID（自己关联自己）
			return errors.New("不能将房源关联到自己") // 返回错误信息，表示不能将房源关联到自己
		}
	}

	// 检查是否已存在绑定关系
	var existBindings []hotel_admin.RelatedRoomBinding                                                                // 声明绑定列表变量，用于存储查询到的已存在绑定信息
	db.MysqlDB.Where("main_room_id = ? AND related_room_id IN ? AND deleted_at IS NULL", mainRoomID, relatedRoomIDs). // 查询是否已存在该主房源和关联房源列表中的绑定关系（使用IN子句，排除已删除绑定）
																Find(&existBindings) // 执行查询并获取已存在的绑定列表
	if len(existBindings) > 0 { // 如果已存在的绑定列表不为空（长度大于0），说明部分关联关系已存在
		return errors.New("部分关联关系已存在") // 返回错误信息，表示部分关联关系已存在
	}

	// 批量创建
	bindings := make([]hotel_admin.RelatedRoomBinding, 0, len(relatedRoomIDs)) // 创建绑定记录切片，初始容量为关联房源ID列表长度
	for _, relatedRoomID := range relatedRoomIDs {                             // 遍历关联房源ID列表，为每个关联房源创建绑定记录
		bindings = append(bindings, hotel_admin.RelatedRoomBinding{ // 将新的绑定关系添加到列表中
			MainRoomID:    mainRoomID,    // 设置主房源ID
			RelatedRoomID: relatedRoomID, // 设置关联房源ID
			BindingDesc:   bindingDesc,   // 设置绑定描述（从参数中获取，可为空）
		})
	}

	return db.MysqlDB.Create(&bindings).Error // 批量创建绑定关系，返回创建操作的结果（成功为nil，失败为error）
}
