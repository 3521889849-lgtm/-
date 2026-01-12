package db

import (
	"database/sql"
	"example_shop/common/config"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MysqlDB *sql.DB

func MysqlInit() error {
	MysqlS := config.Cfg.MysqlInit
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local)",
		MysqlS.User,
		MysqlS.Password,
		MysqlS.Host,
		MysqlS.Port,
		MysqlS.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return err
	}
	sqlDB, _ := db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了可以重新使用连接的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour * 30)
	return nil
}
