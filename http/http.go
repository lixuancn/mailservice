package http

import (
	"net/http"
	"mailservice/config"
	"github.com/go-errors/errors"
	"log"
)

func init(){
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	http.HandleFunc("/sender/mail", HttpMail)
}

func Start(){
	listen := config.Get().Http.Listen
	if listen == ""{
		errors.New("配置文件未设置HTTP监听端口")
		return
	}
	s := &http.Server{
		Addr: listen,
	}
	log.Println("HTTP服务已启动,监听" + listen + "端口")
	log.Fatalln(s.ListenAndServe())
}

func GetPostValue(w http.ResponseWriter, r *http.Request, name string, isMust bool)(string, error){
	valueList := r.MultipartForm.Value[name]
	value := ""
	if len(valueList) > 0{
		value = valueList[0]
	}
	if isMust && len(value) <= 0{
		msg := name + "是必填参数"
		http.Error(w, msg, http.StatusBadRequest)
		return value, errors.New(msg)
	}
	return value, nil
}