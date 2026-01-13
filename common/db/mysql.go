package db

import (
	"example_shop/common/config"
	"fmt"
	"net/url"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MysqlDB *gorm.DB

func MysqlInit() error {
	MysqlS := config.Cfg.MysqlInit
	// 在 Windows 上使用 Asia/Shanghai 时区，URL编码后使用
	loc := url.QueryEscape("Asia/Shanghai")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=%s",
		MysqlS.User,
		MysqlS.Password,
		MysqlS.Host,
		MysqlS.Port,
		MysqlS.Database,
		loc,
	)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return err
	}
	fmt.Println("数据库连接成功")
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了可以重新使用连接的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour * 30)
	MysqlDB = db
	return nil
}
