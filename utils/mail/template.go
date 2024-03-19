package mail

import "market/global"

func UserLoginTemplate(to string, loginUrl string) []byte {
	msg := []byte(
		"From: " + global.MARKET_CONFIG.Smtp.From + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: Log in to Predictmarket\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"utf-8\"\r\n" +
			"\r\n" +
			"<html><body>" +
			"<h1>Predictmarket</h1>" +
			"<p>Click the button below to log in to Predictmarket.</p>" +
			"<p>This button will expire in 20 minutes.</p>" +
			"<p><a href=\"" + loginUrl + "\">Log in to Predictmarket</a></p>" +
			"<p>Confirming this request will securely log you in using " + to + ".</p>" +
			"<h2>- Predictmarket Team</h2>" +
			"</body></html>\r\n")
	return msg
}
