package config

var Cfg = new(Config)

type Config struct {
	Mysql
	Redis
	Coupon
}

type Mysql struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type Redis struct {
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
