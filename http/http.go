package http

import (
	"net/http"
	"mailservice/config"
	"github.com/go-errors/errors"
	"log"
	"encoding/json"
	"mime/multipart"
	"strconv"
)

func init(){
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	http.HandleFunc("/mail", HttpMail)
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

/**
 * 获取文件列表
 */
func GetFileList(r *http.Request, name string)[]*multipart.FileHeader{
	//如果该名字下面是个数组的话, 比如form表单提交的
	fileList, ok := r.MultipartForm.File[name];
	if ok{
		return fileList
	}
	//如果该名字找不到, 数组是name[0], name[1]之类的
	fileList = make([]*multipart.FileHeader, 0)
	count := 0
	for k, v := range r.MultipartForm.File{
		key := name + "[" + strconv.Itoa(count) + "]"
		if k != key{
			continue;
		}
		fileList = append(fileList, v[0])
		count++
	}
	return fileList
}