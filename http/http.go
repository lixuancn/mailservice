package http

import (
	"net/http"
	"mailservice/config"
	"github.com/go-errors/errors"
	"log"
	"encoding/json"
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
		Output(w, nil, http.StatusBadRequest, msg)
		return value, errors.New(msg)
	}
	return value, nil
}

type Result struct{
	Data interface{}
	ErrorCode int
	ErrorMsg string
}

func ResultJson(data interface{}, errCode int, errMsg string)([]byte, error){
	result := Result{Data:data, ErrorCode:errCode, ErrorMsg:errMsg}
	return json.Marshal(result)
}

func Output(w http.ResponseWriter, data interface{}, errCode int, errMsg string)error{
	j, err := ResultJson(data, errCode, errMsg)
	if err != nil{
		return err
	}
	w.Write(j)
	return nil
}