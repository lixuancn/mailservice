package email

import (
	"net/mail"
	"mailservice/config"
	"net/smtp"
	"bytes"
	"strings"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"mime"
	"time"
)

type Message struct{
	From mail.Address
	To []string
	Cc []string
	Subject string
	Body string
	BodyContentType string
	Attachments map[string]*Attachment

	Bcc             []string
	ReplyTo         string
}

type Attachment struct {
	Filename string
	Data     []byte
	Inline   bool
}

func New()*Message{
	message := &Message{}
	//"text/html"
	message.BodyContentType = "text/plain"
	//message.Attachments = make(map[string]*Attachment)
	return message
}

func NewAttachmentList()map[string]*Attachment{
	return make(map[string]*Attachment)
}

func NewAttachment(filename string, fileContent []byte, inline bool)*Attachment{
	return &Attachment{
		Filename: filename,
		Data: fileContent,
		Inline: inline,
	}
}

func (this *Message)Send()error{
	globalConfig := config.Get()
	auth := smtp.PlainAuth("", this.From.Address, globalConfig.Smtp.Username, globalConfig.Smtp.Password)
	return smtp.SendMail(globalConfig.Smtp.Addr, auth, this.From.Address, this.GetToList(), this.GetMsg())
}

func (this *Message) GetToList()[]string {
	toList := this.To
	for _, cc := range this.Cc {
		toList = append(toList, cc)
	}
	for _, bcc := range this.Bcc {
		toList = append(toList, bcc)
	}
	return toList
}

func (this *Message)GetMsg()[]byte{
	buf := bytes.NewBuffer(nil)
	division := "\r\n"
	buf.WriteString("From: " + this.From.String() + division)
	buf.WriteString("Date: " + time.Now().Format(time.RFC822) + division)
	buf.WriteString("To: " + strings.Join(this.To, ",") + division)
	if len(this.Cc) > 0{
		buf.WriteString("Cc: " + strings.Join(this.Cc, ",") + division)
	}

	var coder = base64.StdEncoding
	var subject = "=?UTF-8?B?" + coder.EncodeToString([]byte(this.Subject)) + "?="
	buf.WriteString("Subject: " + subject + division)

	if len(this.ReplyTo) > 0 {
		buf.WriteString("Reply-To: " + this.ReplyTo + division)
	}
	buf.WriteString("MIME-Version: 1.0" + division)
	boundary := "f46d043c813270fc6b04c2d223da"
	if len(this.Attachments) > 0 {
		buf.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + division)
		buf.WriteString(division + "--" + boundary + division)
	}
	buf.WriteString(fmt.Sprintf("Content-Type: %s; charset=utf-8" + division + division, this.BodyContentType))
	buf.WriteString(this.Body)
	buf.WriteString(division)

	if len(this.Attachments) > 0 {
		for _, attachment := range this.Attachments {
			buf.WriteString(division + division + "--" + boundary + division)
			if attachment.Inline {
				buf.WriteString("Content-Type: message/rfc822" + division)
				buf.WriteString(fmt.Sprintf("Content-Disposition: inline; filename=\"%s\"" + division + division, attachment.Filename))
				buf.Write(attachment.Data)
			}else {
				ext := filepath.Ext(attachment.Filename)
				mimeType := mime.TypeByExtension(ext)
				if mimeType != "" {
					buf.WriteString(fmt.Sprintf("Content-Type: %s" + division, mimeType))
				}else{
					buf.WriteString("Content-Type: application/octet-stream" + division)
				}
				buf.WriteString("Content-Transfer-Encoding: base64" + division)
				buf.WriteString("Content-Disposition: attachment; filename=\"=?UTF-8?B?")
				buf.WriteString(coder.EncodeToString([]byte(attachment.Filename)))
				buf.WriteString("?=\"" + division + division)
				b := make([]byte, base64.StdEncoding.EncodedLen(len(attachment.Data)))
				base64.StdEncoding.Encode(b, attachment.Data)
				for i, l := 0, len(b); i < l; i++ {
					buf.WriteByte(b[i])
					if (i+1)%76 == 0 {
						buf.WriteString(division)
					}
				}
			}
			buf.WriteString(division + "--" + boundary)
		}
		buf.WriteString("--")
	}
	return buf.Bytes()
}