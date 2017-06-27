package main

import (
	"runtime"
	"log"
	"flag"
	"fmt"
	"os"
	"mailservice/config"
	"mailservice/http"
)

func init(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	//接收命令行参数
	flagConfigFile := flag.String("c", "config.json", "configuration file")
	flagIsVersion := flag.Bool("v", false, "show version")
	flagIsHelp := flag.Bool("h", false, "help")
	flag.Parse()
	flagFuncVersion(*flagIsVersion)
	flagFuncHelp(*flagIsHelp)
	flagFuncConfigFile(*flagConfigFile)

}

func flagFuncVersion(isVersion bool){
	if isVersion {
		fmt.Println("邮件SMTP服务, 当前版本: " + config.VERSION)
		fmt.Println("作者: lixuan-it@360.cn")
		os.Exit(0)
	}
}

func flagFuncHelp(isHelp bool){
	if isHelp{
		flag.Usage()
		os.Exit(0)
	}
}

func flagFuncConfigFile(configFile string){
	err := config.Parse(configFile)
	if err != nil{
		log.Fatalln(err)
	}
}

func main(){
	http.Start()
}