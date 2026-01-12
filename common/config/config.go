package config

var Cfg = new(Config)

type Config struct {
	MysqlInit
	RedisInit
	Coupon
}

type MysqlInit struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type RedisInit struct {
	Host     string
	Port     int
	Password string
	Database int
}

type Coupon struct {
	AntiBrushLimit  int
	AntiBrushExpire int
	AesKey          string
}
