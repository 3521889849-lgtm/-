package init

import (
	"example_shop/common/config"
	"example_shop/common/db"
	"log"
)

func init() {

	if err := config.ViperInit(); err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}
	if err := db.MysqlInit(); err != nil {
		log.Fatalf("MySQL连接初始化失败: %v", err)
	}
	if err := db.RedisInit(); err != nil {
		log.Fatalf("Redis连接初始化失败: %v", err)
	}
}
