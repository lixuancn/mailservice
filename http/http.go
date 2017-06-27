package http

import (
	"net/http"
	"net/smtp"
	"mailservice/config"
	"strings"
	"github.com/go-errors/errors"
	"log"
	"github.com/scorredoira/email"
	"io/ioutil"
)

func init(){
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(config.VERSION))
	})
	http.HandleFunc("/sender/mail", func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(32 << 20)
		if r.MultipartForm == nil {
			http.Error(w, "获取POST数据失败", http.StatusBadRequest)
			return
		}
		token := GetPostValue(w, r, "token")
		globalConfig := config.GetConfig()
		if globalConfig.Http.Token != token{
			http.Error(w, "Token验证失败", http.StatusForbidden)
			return
		}
		tos := GetPostValue(w, r, "tos")
		if len(tos) <= 0{
			http.Error(w, "tos是必填参数", http.StatusBadRequest)
			return
		}
		tos = strings.Replace(tos, ",", ";", -1)
		toList := strings.Split(tos, ";")
		if len(toList) <= 0{
			http.Error(w, "无效的收件人", http.StatusBadRequest)
			return
		}
		ccs := GetPostValue(w, r, "ccs")
		ccs = strings.Replace(ccs, ",", ";", -1)
		var ccList []string
		if len(ccs) > 0{
			ccList = strings.Split(ccs, ";")
			if len(ccList) <= 0{
				http.Error(w, "无效的抄送人", http.StatusBadRequest)
				return
			}
		}
		subject := GetPostValue(w, r, "subject")
		if len(subject) <= 0{
			http.Error(w, "subject是必填参数", http.StatusBadRequest)
			return
		}
		content := GetPostValue(w, r, "content")
		if len(content) <= 0{
			http.Error(w, "content是必填参数", http.StatusBadRequest)
			return
		}
		//获取附件
		fileList := r.MultipartForm.File["attachment"]
		//准备发送
		m := email.NewHTMLMessage(subject, content)
		m.From = globalConfig.Smtp.From
		m.To = toList
		m.Cc = ccList
		if len(fileList) > 0{
			for _, file := range fileList{
				if len(file.Filename) <= 0{
					http.Error(w, "附件"+file.Filename+"没有名字", http.StatusBadRequest)
					return
				}
				f, err := file.Open()
				if err != nil{
					http.Error(w, "打开附件"+file.Filename+"失败, err: " + err.Error(), http.StatusBadRequest)
					return
				}
				fileContent, err := ioutil.ReadAll(f)
				if err != nil{
					http.Error(w, "读取附件"+file.Filename+"失败, err: " + err.Error(), http.StatusBadRequest)
					return
				}
				err = m.AttachBuffer(file.Filename, fileContent, false)
				if err != nil{
					http.Error(w, "发送附件"+file.Filename+"失败, err: " + err.Error(), http.StatusBadRequest)
					return
				}
			}
		}
		err := email.Send(globalConfig.Smtp.Addr, smtp.PlainAuth("", globalConfig.Smtp.From.Address, globalConfig.Smtp.Username, globalConfig.Smtp.Password), m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}else{
			http.Error(w, "success", http.StatusOK)
			return
		}
	})
}

func Start(){
	port := config.GetConfig().Http.Listen
	if port == ""{
		errors.New("配置文件未设置HTTP监听端口")
		return
	}
	s := &http.Server{
		Addr: port,
	}
	log.Println("HTTP服务已启动,监听" + port + "端口")
	log.Fatalln(s.ListenAndServe())
}

func GetPostValue(w http.ResponseWriter, r *http.Request, name string)string{
	valueList := r.MultipartForm.Value[name]
	value := ""
	if len(valueList) > 0{
		value = valueList[0]
	}
	return value
}