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
func (this *Message)Send()error{
	globalConfig := config.Get()
	auth := smtp.PlainAuth("", this.From.Address, globalConfig.Smtp.Username, globalConfig.Smtp.Password)
	return smtp.SendMail(globalConfig.Smtp.Addr, auth, this.From.Address, this.ToList(), this.Bytes())
}

func (this *Message) ToList()[]string {
	toList := this.To
	for _, cc := range this.Cc {
		toList = append(toList, cc)
	}
	for _, bcc := range this.Bcc {
		toList = append(toList, bcc)
	}
	return toList
}

// Bytes returns the mail data
func (this *Message) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("From: " + this.From.String() + "\r\n")

	t := time.Now()
	buf.WriteString("Date: " + t.Format(time.RFC822) + "\r\n")

	buf.WriteString("To: " + strings.Join(this.To, ",") + "\r\n")
	if len(this.Cc) > 0 {
		buf.WriteString("Cc: " + strings.Join(this.Cc, ",") + "\r\n")
	}

	//fix  Encode
	var coder = base64.StdEncoding
	var subject = "=?UTF-8?B?" + coder.EncodeToString([]byte(this.Subject)) + "?="
	buf.WriteString("Subject: " + subject + "\r\n")

	if len(this.ReplyTo) > 0 {
		buf.WriteString("Reply-To: " + this.ReplyTo + "\r\n")
	}

	buf.WriteString("MIME-Version: 1.0\r\n")

	boundary := "f46d043c813270fc6b04c2d223da"

	if len(this.Attachments) > 0 {
		buf.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\r\n")
		buf.WriteString("\r\n--" + boundary + "\r\n")
	}

	buf.WriteString(fmt.Sprintf("Content-Type: %s; charset=utf-8\r\n\r\n", this.BodyContentType))
	buf.WriteString(this.Body)
	buf.WriteString("\r\n")

	if len(this.Attachments) > 0 {
		for _, attachment := range this.Attachments {
			buf.WriteString("\r\n\r\n--" + boundary + "\r\n")

			if attachment.Inline {
				buf.WriteString("Content-Type: message/rfc822\r\n")
				buf.WriteString("Content-Disposition: inline; filename=\"" + attachment.Filename + "\"\r\n\r\n")

				buf.Write(attachment.Data)
			} else {
				ext := filepath.Ext(attachment.Filename)
				mimetype := mime.TypeByExtension(ext)
				if mimetype != "" {
					mime := fmt.Sprintf("Content-Type: %s\r\n", mimetype)
					buf.WriteString(mime)
				} else {
					buf.WriteString("Content-Type: application/octet-stream\r\n")
				}
				buf.WriteString("Content-Transfer-Encoding: base64\r\n")

				buf.WriteString("Content-Disposition: attachment; filename=\"=?UTF-8?B?")
				buf.WriteString(coder.EncodeToString([]byte(attachment.Filename)))
				buf.WriteString("?=\"\r\n\r\n")

				b := make([]byte, base64.StdEncoding.EncodedLen(len(attachment.Data)))
				base64.StdEncoding.Encode(b, attachment.Data)

				// write base64 content in lines of up to 76 chars
				for i, l := 0, len(b); i < l; i++ {
					buf.WriteByte(b[i])
					if (i+1)%76 == 0 {
						buf.WriteString("\r\n")
					}
				}
			}

			buf.WriteString("\r\n--" + boundary)
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}