package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func ViperInit() error {
	var err error
	viper.SetConfigName("config")
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&Cfg)
	if err != nil {
		return err
	}
	fmt.Println("配置", Cfg)
	return nil
}
