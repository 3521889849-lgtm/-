package service

import (
	"errors"
	"example_shop/common/db"
	"example_shop/common/model/hotel_admin"

	"gorm.io/gorm"
)

// RoomFacilityService 房源设施关联管理服务
// 负责处理房源与设施的关联关系的设置、删除等核心业务逻辑，
// 包括房源设施的全量更新（先删除旧关联，再创建新关联）、设施有效性检查等。
type RoomFacilityService struct{}

// SetRoomFacilities 设置房源的设施（先删除旧的，再创建新的）
func (s *RoomFacilityService) SetRoomFacilities(roomID uint64, facilityIDs []uint64) error {
	// 检查房源是否存在
	var roomInfo hotel_admin.RoomInfo                                 // 声明房源实体变量，用于存储查询到的房源信息
	if err := db.MysqlDB.First(&roomInfo, roomID).Error; err != nil { // 通过房源ID查询房源信息，如果查询失败则说明房源不存在
		return errors.New("房源不存在") // 返回错误信息，表示房源不存在
	}

	// 检查所有设施是否存在
	if len(facilityIDs) > 0 { // 如果设施ID列表不为空（长度大于0），则需要验证所有设施是否存在
		var count int64                                                                                                  // 声明计数变量，用于存储存在的设施数量
		db.MysqlDB.Model(&hotel_admin.FacilityDict{}).Where("id IN ? AND deleted_at IS NULL", facilityIDs).Count(&count) // 统计设施ID列表中存在的设施数量（使用IN子句，排除已删除设施）
		if int(count) != len(facilityIDs) {                                                                              // 如果存在的设施数量不等于请求的设施数量，说明部分设施不存在
			return errors.New("部分设施不存在") // 返回错误信息，表示部分设施不存在
		}
	}

	// 使用事务
	return db.MysqlDB.Transaction(func(tx *gorm.DB) error { // 开启数据库事务，传入事务处理函数，返回事务执行结果
		// 软删除旧的关联关系
		if err := tx.Model(&hotel_admin.RoomFacilityRelation{}). // 创建房源设施关联模型的查询构建器（使用事务连接）
										Where("room_id = ? AND deleted_at IS NULL", roomID).            // 添加筛选条件：房源ID匹配且未被删除
										Delete(&hotel_admin.RoomFacilityRelation{}).Error; err != nil { // 执行软删除操作（设置deleted_at字段），删除该房源的所有旧关联关系，如果删除失败则返回错误
			return err // 返回数据库操作错误
		}

		// 创建新的关联关系
		if len(facilityIDs) > 0 { // 如果设施ID列表不为空（长度大于0），则需要创建新的关联关系
			relations := make([]hotel_admin.RoomFacilityRelation, 0, len(facilityIDs)) // 创建房源设施关联记录切片，初始容量为设施ID列表长度
			for _, facilityID := range facilityIDs {                                   // 遍历设施ID列表，为每个设施创建关联记录
				// 检查是否已存在（避免重复）
				var exist hotel_admin.RoomFacilityRelation                                                                     // 声明关联实体变量，用于存储查询到的已存在关联信息
				result := tx.Where("room_id = ? AND facility_id = ? AND deleted_at IS NULL", roomID, facilityID).First(&exist) // 查询是否已存在该房源和设施的关联关系（使用事务连接，排除已删除关联）
				if result.Error != nil {                                                                                       // 如果查询失败（未找到已存在的关联），说明该关联关系不存在，可以创建
					relations = append(relations, hotel_admin.RoomFacilityRelation{ // 将新的关联关系添加到列表中
						RoomID:     roomID,     // 设置房源ID
						FacilityID: facilityID, // 设置设施ID
					})
				}
			}
			if len(relations) > 0 { // 如果关联关系列表不为空（有需要创建的关联关系）
				if err := tx.Create(&relations).Error; err != nil { // 批量创建关联关系（使用事务连接），如果创建失败则返回错误
					return err // 返回数据库操作错误
				}
			}
		}

		return nil // 返回nil表示事务执行成功（旧关联删除和新关联创建都成功）
	})
}

// GetRoomFacilities 获取房源的设施列表
func (s *RoomFacilityService) GetRoomFacilities(roomID uint64) ([]hotel_admin.FacilityDict, error) {
	var facilities []hotel_admin.FacilityDict // 声明设施列表变量，用于存储查询到的设施信息列表

	if err := db.MysqlDB.Model(&hotel_admin.FacilityDict{}). // 创建设施模型的查询构建器
									Joins("JOIN room_facility_relation ON facility_dict.id = room_facility_relation.facility_id").                                          // 通过JOIN关联房源设施关联表（内连接，获取关联的设施信息）
									Where("room_facility_relation.room_id = ? AND room_facility_relation.deleted_at IS NULL AND facility_dict.deleted_at IS NULL", roomID). // 添加筛选条件：房源ID匹配、关联关系未被删除、设施未被删除
									Find(&facilities).Error; err != nil {                                                                                                   // 执行查询并获取符合条件的设施列表，如果查询失败则返回错误
		return nil, err // 返回nil和错误信息
	}

	return facilities, nil // 返回设施列表和无错误
}

// AddRoomFacility 为房源添加单个设施
func (s *RoomFacilityService) AddRoomFacility(roomID uint64, facilityID uint64) error {
	// 检查房源是否存在
	var roomInfo hotel_admin.RoomInfo                                 // 声明房源实体变量，用于存储查询到的房源信息
	if err := db.MysqlDB.First(&roomInfo, roomID).Error; err != nil { // 通过房源ID查询房源信息，如果查询失败则说明房源不存在
		return errors.New("房源不存在") // 返回错误信息，表示房源不存在
	}

	// 检查设施是否存在
	var facility hotel_admin.FacilityDict                                 // 声明设施实体变量，用于存储查询到的设施信息
	if err := db.MysqlDB.First(&facility, facilityID).Error; err != nil { // 通过设施ID查询设施信息，如果查询失败则说明设施不存在
		return errors.New("设施不存在") // 返回错误信息，表示设施不存在
	}

	// 检查关联关系是否已存在
	var exist hotel_admin.RoomFacilityRelation                                                                             // 声明关联实体变量，用于存储查询到的已存在关联信息
	result := db.MysqlDB.Where("room_id = ? AND facility_id = ? AND deleted_at IS NULL", roomID, facilityID).First(&exist) // 查询是否已存在该房源和设施的关联关系（排除已删除关联）
	if result.Error == nil {                                                                                               // 如果查询成功（未报错），说明该关联关系已存在
		return errors.New("该设施已关联到此房源") // 返回错误信息，表示该设施已关联到此房源
	}

	relation := &hotel_admin.RoomFacilityRelation{ // 创建房源设施关联实体对象指针
		RoomID:     roomID,     // 设置房源ID
		FacilityID: facilityID, // 设置设施ID
	}

	return db.MysqlDB.Create(relation).Error // 将关联关系保存到数据库，返回保存操作的结果（成功为nil，失败为error）
}

// RemoveRoomFacility 移除房源的单个设施
func (s *RoomFacilityService) RemoveRoomFacility(roomID uint64, facilityID uint64) error {
	var relation hotel_admin.RoomFacilityRelation                                                                                                 // 声明关联实体变量，用于存储查询到的关联信息
	if err := db.MysqlDB.Where("room_id = ? AND facility_id = ? AND deleted_at IS NULL", roomID, facilityID).First(&relation).Error; err != nil { // 查询该房源和设施的关联关系（排除已删除关联），如果查询失败则说明关联关系不存在
		return errors.New("关联关系不存在") // 返回错误信息，表示关联关系不存在
	}

	return db.MysqlDB.Delete(&relation).Error // 执行软删除操作（设置deleted_at字段），根据关联关系删除记录，返回删除操作的结果（成功为nil，失败为error）
}
