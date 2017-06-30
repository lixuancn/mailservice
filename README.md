# 邮件发送服务

### 功能
能发邮件、抄送、HTML内容、附件等

### 起因
项目需要发送邮件给公司人员

### 参考
- github.com/open-falcon/mail-provider

- github.com/scorredoira/email 这本来是个发邮件的包, 我想自己造个轮子, 就参考了它, 然后抛弃了它。

### 依赖
- Linux下的sendmail

### 原理
使用SMTP服务发送。

### 用法

- Protocol: HTTP

- Method: POST

- 发件箱中文名: 必填, fromname string 

- 发件箱地址: 必填, fromaddress string

- 收件人: 必填, tos string, 多个用逗号分隔

- 抄送人: 选填, ccs string, 多个用逗号分隔

- 主题: 必填, subject string

- 内容: 必填, content string, 可用HTML代码

- 附件: 选填, attachment, 可单文件可数组

- 编码: 选填, body_content_type string, 默认是text/plain, 可选是text/html等

### 实例
##### FORM表单
```html
<form method="post" action="http://****:4000/sender/mail" enctype="multipart/form-data">
    <input type="text" name="fromname" value="**邮件系统"><br>
    <input type="text" name="fromaddress" value="example@example.com"><br>
    <input type="text" name="tos" value="***@qq.com"><br>
    <input type="text" name="ccs" value="***@163.com"><br>
    <input type="text" name="subject" value="搭建邮件服务器-测试附件发送"><br>
    <input type="file" name="attachment"><br>
    <input type="file" name="attachment"><br>
    <textarea name="content">
<p>测试邮件发送服务</p>
<p>附件、抄送等功能</p>
<img src="http://*****.png">
    </textarea><br>
    <input type="submit">
</form>
```