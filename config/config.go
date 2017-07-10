package config

const GLOBAL_VERSION = "1.0.0"
const HTTP_LISTEN = "0.0.0.0:4000"
const HTTP_TOKEN = ""
const SMTP_ADDR = "127.0.0.1:25"
const SMTP_USERNAME = ""
const SMTP_PASSWORD = ""

type GlobalConfig struct {
	Http HttpConfig
	Smtp SmtpConfig
}

type HttpConfig struct {
	Listen string
	Token string
}

type SmtpConfig struct {
	Addr string
	Username string
	Password string
}

var globalConfig GlobalConfig

func Get()GlobalConfig{
	globalConfig.Http = HttpConfig{
		Listen: HTTP_LISTEN,
		Token: HTTP_TOKEN,
	}
	globalConfig.Smtp = SmtpConfig{
		Addr: SMTP_ADDR,
		Username: SMTP_USERNAME,
		Password: SMTP_PASSWORD,
	}
	return globalConfig
}