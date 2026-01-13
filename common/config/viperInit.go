package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func ViperInit() error {
	var err error
	viper.AddConfigPath("conf")
	viper.AddConfigPath("../conf")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&Cfg)
	if err != nil {
		return err
	}
	fmt.Println("配置动态加载成功", Cfg)
	return nil
}
