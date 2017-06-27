package config

import (
	"github.com/go-errors/errors"
	"io/ioutil"
	"strings"
	"sync"
	"encoding/json"
	"os"
	"net/mail"
)

const VERSION = "1.0.0"

type GlobalConfig struct {
	Debug bool `json:"debug"`
	Http *HttpConfig `json:"http"`
	Smtp *SmtpConfig `json:"smtp"`
}

type HttpConfig struct {
	Listen string `json:"listen"`
	Token string `json:"token"`
}

type SmtpConfig struct {
	Addr string `json:"addr"`
	Username string `json:"username"`
	Password string `json:"password"`
	From mail.Address `json:"from"`
}

var globalConfig *GlobalConfig
var configLock = new(sync.RWMutex)

func Parse(configFile string)error{
	println("aaa")
	configLock.Lock()
	defer configLock.Unlock()
	if configFile == "" {
		return errors.New("请传入配置文件参数-c")
	}
	_, err := os.Stat(configFile)
	if err == nil && os.IsExist(err){
		return errors.New("配置文件'" + configFile + "'不存在")
	}
	configContentByte, err := ioutil.ReadFile(configFile)
	if err != nil{
		return errors.New("读取配置文件'" + configFile + "'错误: " + err.Error())
	}
	configContent := strings.TrimSpace(string(configContentByte))

	err = json.Unmarshal([]byte(configContent), &globalConfig)
	if err != nil{
		return errors.New("配置文件'" + configFile + "'逆向Json时错误: " + err.Error())
	}
	return nil
}

func GetConfig()*GlobalConfig{
	configLock.RLock()
	defer configLock.RUnlock()
	return globalConfig
}