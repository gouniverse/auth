package auth

import (
	"bytes"
	"html/template"
	"log"
)

// EmailPasswordChangeTemplate returns the template for the email address verification email
func emailTemplatePasswordChange(name string, url string) string {
	msg := `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html>
<head></head>
<body>
	<p>
		Hello!
	<p>
	<p>
		Someone requested to reset your password. Please click the link bellow to reset it.
	</p>
	<p>
		<a href="{{.URL}}">Change Password</a>
	</p>
	<p>
		If you did not request to change your password no further action is required.
	</p>
	<p>
		Thanks,
		<br />
		The Admin Team
	</p>
	<hr />
	<p>
		If you are having trouble clicking the "Reset Password" link,
		copy and paste the URL below into your web browser:
		{{.URL}}
	</p>
</body>
<html>
`
	data := struct {
		Name string
		URL  string
	}{
		Name: name,
		URL:  url,
	}

	t, err := template.New("template").Parse(msg)
	if err != nil {
		log.Println(err)
		return ""
	}

	var doc bytes.Buffer
	errExecute := t.Execute(&doc, data)

	if errExecute != nil {
		log.Println(errExecute)
		return ""
	}

	s := doc.String()
	return s
}

// emailLoginCodeTemplate returns the template for the lofin code verification email
func emailLoginCodeTemplate(email string, code string) string {
	msg := `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html>
<head></head>
<body>
	<p>
		Hello!
	<p>
	<p>
		Someone requested to login with your emal {{.Email}}. Please use the code below to log in.
	</p>
	<p>
		{{.Code}}
	</p>
	<p>
		If you did not request to login no further action is required.
	</p>
	<p>
		Thanks,
		<br />
		The Admin Team
	</p>
</body>
<html>
`
	data := struct {
		Email string
		Code  string
	}{
		Email: email,
		Code:  code,
	}

	t, err := template.New("template").Parse(msg)
	if err != nil {
		log.Println(err)
		return ""
	}

	var doc bytes.Buffer
	errExecute := t.Execute(&doc, data)

	if errExecute != nil {
		log.Println(errExecute)
		return ""
	}

	s := doc.String()
	return s
}
