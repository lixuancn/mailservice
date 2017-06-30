package http

import (
	"strings"
	"io/ioutil"
	"net/http"
	"mailservice/email"
	"mailservice/config"
	"net/mail"
)

func HttpMail(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	if r.MultipartForm == nil {
		http.Error(w, "获取POST数据失败", http.StatusBadRequest)
		return
	}
	token, err := GetPostValue(w, r, "token", false)
	if err != nil{
		return
	}
	globalConfig := config.Get()
	if globalConfig.Http.Token != token{
		http.Error(w, "Token验证失败", http.StatusForbidden)
		return
	}
	fromAddress, err := GetPostValue(w, r, "fromaddress", true)
	if err != nil{
		return
	}
	fromName, err := GetPostValue(w, r, "fromname", true)
	if err != nil{
		return
	}
	from := mail.Address{Name:fromName, Address:fromAddress}
	tos, err := GetPostValue(w, r, "tos", true)
	if err != nil{
		return
	}
	tos = strings.Replace(tos, ",", ";", -1)
	toList := strings.Split(tos, ";")
	if len(toList) <= 0{
		http.Error(w, "无效的收件人", http.StatusBadRequest)
		return
	}
	ccs, _ := GetPostValue(w, r, "ccs", false)
	ccs = strings.Replace(ccs, ",", ";", -1)
	var ccList []string
	if len(ccs) > 0{
		ccList = strings.Split(ccs, ";")
		if len(ccList) <= 0{
			http.Error(w, "无效的抄送人", http.StatusBadRequest)
			return
		}
	}
	subject, err := GetPostValue(w, r, "subject", true)
	if err != nil{
		return
	}
	content, err := GetPostValue(w, r, "content", true)
	if err != nil{
		return
	}
	//获取附件
	fileList := r.MultipartForm.File["attachment"]
	attachmentList := make(map[string]*email.Attachment)
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
			attachment := email.Attachment{
				Filename: file.Filename,
				Data: fileContent,
				Inline: false,
			}
			attachmentList[file.Filename] = &attachment
			if err != nil{
				http.Error(w, "发送附件"+file.Filename+"失败, err: " + err.Error(), http.StatusBadRequest)
				return
			}
		}
	}
	//准备发送
	m := email.New()
	m.From = from
	m.To = toList
	m.Cc = ccList
	m.Subject = subject
	m.Body = content
	m.Attachments = attachmentList
	err = m.Send()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}else{
		http.Error(w, "success", http.StatusOK)
		return
	}
}