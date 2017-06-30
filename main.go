package main

/**
 * 这个页面基本是github.com/open-falcon/mail-provider的东西没变。
 */

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
	//Go1.5以后这句就可以不要了,默认P数量就是CPU数量了,好像是1.5=。=
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	//接收命令行参数
	flagIsVersion := flag.Bool("v", false, "show version")
	flagIsHelp := flag.Bool("h", false, "help")
	flag.Parse()
	flagFuncVersion(*flagIsVersion)
	flagFuncHelp(*flagIsHelp)

}

func flagFuncVersion(isVersion bool){
	if isVersion {
		fmt.Println("邮件SMTP服务, 当前版本: " + config.GLOBAL_VERSION)
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

func main(){
	http.Start()
}