package helper

import (
	"fmt"
	"net/smtp"
)

func SendResetPasswordEmail(smtpHost, smtpPort, smtpUser, smtpPassword, smtpFrom, toEmail, resetLink string) error {
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	subject := "[Shifd] Permintaan Reset Kata Sandi"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; color: #333; max-width: 600px; margin: 0 auto; padding: 24px;">
  <h2 style="color: #1976D2;">Reset Kata Sandi</h2>
  <p>Kami menerima permintaan untuk mereset kata sandi akun Anda.</p>
  <p>Klik tombol di bawah untuk mereset kata sandi Anda. Link ini aktif selama <strong>30 menit</strong>.</p>
  <p style="margin: 32px 0;">
    <a href="%s" style="background-color: #1976D2; color: white; padding: 12px 28px; text-decoration: none; border-radius: 6px; font-weight: bold;">
      Reset Kata Sandi
    </a>
  </p>
  <p style="color: #555; font-size: 14px;">Atau salin link berikut ke browser Anda:</p>
  <p style="word-break: break-all; font-size: 13px;"><a href="%s" style="color: #1976D2;">%s</a></p>
  <hr style="margin: 32px 0; border: none; border-top: 1px solid #eee;">
  <p style="color: #999; font-size: 12px;">
    Jika Anda tidak merasa meminta reset kata sandi, abaikan email ini.
    Kata sandi Anda tidak akan berubah.
  </p>
</body>
</html>`, resetLink, resetLink, resetLink)

	msg := []byte(
		"From: " + smtpFrom + "\r\n" +
			"To: " + toEmail + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"\r\n" +
			body,
	)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{toEmail}, msg)
}
