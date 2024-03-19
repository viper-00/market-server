package mail

import (
	"testing"
)

func TestMail(t *testing.T) {
	msg := []byte(
		"From: " + "no-reply@predictmarket.xyz" + "\r\n" +
			"To: " + "zhongmingyang2000@gmail.com" + "\r\n" +
			"Subject: Log in to Predictmarket\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"utf-8\"\r\n" +
			"\r\n" +
			"<html><body>" +
			"<h1>Predictmarket</h1>" +
			"<p>Click the link below to log in to Predictmarket.</p>" +
			"<p>This link will expire in 20 minutes.</p>" +
			"<p><a href=\"" + "sdfsdf" + "\">Log in to Predictmarket</a></p>" +
			"<p>Confirming this request will securely log you in using " + "zhongmingyang2000@gmail.com" + ".</p>" +
			"<h2>- Predictmarket Team</h2>" +
			"</body></html>\r\n")

	err := sendCoreMail("smtpout.secureserver.net", 587, "no-reply@predictmarket.xyz", "Viper123&", "zhongmingyang2000@gmail.com", msg)
	if err != nil {
		t.Log(err.Error())
	}
}
