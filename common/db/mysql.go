// Package db 提供数据库连接管理
//
// 本文件实现MySQL数据库的初始化、连接池配置和表结构迁移
package db

import (
	"example_shop/common/config"
	"example_shop/common/model/hotel_admin"
	"fmt"
	"net/url"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MysqlDB 全局MySQL数据库连接实例
//
// 程序启动时初始化，全局共享使用
// 基于GORM框架，提供ORM功能
var MysqlDB *gorm.DB

// MysqlInit 初始化MySQL数据库连接
//
// 功能：
//   - 读取配置创建数据库连接
//   - 配置连接池参数
//   - 自动迁移数据表结构
//   - 设置外键约束
//
// 返回：
//   - error: 初始化失败时返回错误信息
func MysqlInit() error {
	// 读取MySQL配置
	mysqlCfg := config.Cfg.MysqlInit

	// 配置时区：使用Asia/Shanghai（北京时间）
	// URL编码是必需的，防止特殊字符引起的DSN解析错误
	loc := url.QueryEscape("Asia/Shanghai")

	// 构建DSN（Data Source Name）连接字符串
	// 参数说明：
	//   - charset=utf8mb4: 使用utf8mb4字符集（支持emoji等4字节字符）
	//   - parseTime=True: 自动解析DATE/DATETIME为time.Time
	//   - loc=...: 设置时区
	//   - timeout=5s: 连接超时5秒
	//   - readTimeout=10s: 读取超时10秒
	//   - writeTimeout=10s: 写入超时10秒
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=%s&timeout=5s&readTimeout=10s&writeTimeout=10s",
		mysqlCfg.User,
		mysqlCfg.Password,
		mysqlCfg.Host,
		mysqlCfg.Port,
		mysqlCfg.Database,
		loc,
	)

	// 配置GORM选项
	gormConfig := &gorm.Config{
		// 禁用自动创建外键约束
		// 原因：外键约束在某些情况下会影响性能，且增加维护复杂度
		// 我们通过应用层逻辑保证数据一致性
		DisableForeignKeyConstraintWhenMigrating: true,

		// 是否跳过默认事务
		// false: 每次Create/Update/Delete都自动开启事务（安全但稍慢）
		// true: 不自动开启事务（需要手动管理，性能更好）
		SkipDefaultTransaction: false,
	}

	// 打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("MySQL连接失败: %w", err)
	}

	// 获取底层的 *sql.DB 实例，用于配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取SQL.DB实例失败: %w", err)
	}

	// ========== 配置连接池参数 ==========

	// SetMaxIdleConns 设置最大空闲连接数
	// 空闲连接：已建立但暂时未使用的连接
	// 作用：保持一定数量的空闲连接，提高响应速度
	// 推荐值：10-20（根据并发量调整）
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置最大打开连接数
	// 打开连接：包括正在使用和空闲的所有连接
	// 作用：限制同时连接数，防止数据库过载
	// 推荐值：50-200（根据数据库配置和并发量调整）
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置连接最大生存时间
	// 作用：定期关闭旧连接，避免连接泄漏和数据库端超时
	// 推荐值：1-4小时
	sqlDB.SetConnMaxLifetime(time.Hour * 30)

	// SetConnMaxIdleTime 设置连接最大空闲时间
	// 作用：关闭长时间未使用的空闲连接，释放资源
	// 推荐值：5-15分钟
	sqlDB.SetConnMaxIdleTime(time.Minute * 10)

	// ========== 数据表自动迁移 ==========

	// 临时禁用外键检查，避免迁移时的外键约束问题
	// 迁移完成后会重新启用
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// AutoMigrate 自动迁移表结构
	// 功能：
	//   - 检查表是否存在，不存在则创建
	//   - 检查字段是否存在，不存在则添加
	//   - 不会删除已存在的表或字段
	//   - 不会修改已存在字段的类型（需要手动修改）
	//
	// 迁移顺序说明：
	//   按照表的依赖关系排序，先创建父表，再创建子表
	//   避免外键约束导致的创建失败
	err = db.AutoMigrate(
		// ========== 第一层：基础表（无外键依赖） ==========
		&hotel_admin.HotelBranch{},        // 酒店分店表
		&hotel_admin.RoomTypeDict{},       // 房型字典表
		&hotel_admin.FacilityDict{},       // 设施字典表
		&hotel_admin.CancellationPolicy{}, // 退订政策表
		&hotel_admin.Role{},               // 角色表
		&hotel_admin.Permission{},         // 权限表
		&hotel_admin.SystemConfig{},       // 系统配置表
		&hotel_admin.ChannelConfig{},      // 渠道配置表
		&hotel_admin.MemberRights{},       // 会员权益表

		// ========== 第二层：依赖基础表的表 ==========
		&hotel_admin.RoomInfo{},               // 房源信息表
		&hotel_admin.UserAccount{},            // 用户账号表
		&hotel_admin.RolePermissionRelation{}, // 角色权限关联表

		// ========== 第三层：依赖第二层的表 ==========
		&hotel_admin.RoomFacilityRelation{}, // 房间设施关联表
		&hotel_admin.RoomImage{},            // 房间图片表
		&hotel_admin.RelatedRoomBinding{},   // 关联房间绑定表
		&hotel_admin.RoomStatusDetail{},     // 房态明细表
		&hotel_admin.GuestInfo{},            // 客人信息表
		&hotel_admin.Member{},               // 会员表
		&hotel_admin.OrderMain{},            // 订单主表
		&hotel_admin.OrderExtension{},       // 订单扩展表
		&hotel_admin.Blacklist{},            // 黑名单表
		&hotel_admin.MemberPointsRecord{},   // 积分记录表
		&hotel_admin.FinancialFlow{},        // 财务流水表
		&hotel_admin.SalesStatistics{},      // 销售统计表
		&hotel_admin.ShiftChangeRecord{},    // 交接班记录表
		&hotel_admin.OperationLog{},         // 操作日志表
		&hotel_admin.ChannelSyncLog{},       // 渠道同步日志表
	)

	if err != nil {
		// 迁移失败，恢复外键检查并返回错误
		db.Exec("SET FOREIGN_KEY_CHECKS = 1")
		return fmt.Errorf("数据表迁移失败: %w", err)
	}

	// 恢复外键检查
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	// 保存到全局变量
	MysqlDB = db

	return nil
}
