package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to string, subject string, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")

	from := smtpUser
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	err := smtp.SendMail(addr, auth, from, []string{to}, message)
	if err != nil {
		return err
	}
	return nil
}

func SendEmailWithAttachment(to, subject, body string, imageData []byte) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")

	from := smtpUser
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	boundary := "boundary123"

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = fmt.Sprintf("multipart/mixed; boundary=%s", boundary)

	var bodyBuffer bytes.Buffer
	for k, v := range headers {
		bodyBuffer.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	bodyBuffer.WriteString("\r\n--" + boundary + "\r\n")
	bodyBuffer.WriteString("Content-Type: text/plain; charset=utf-8\r\n\r\n")
	bodyBuffer.WriteString(body + "\r\n")
	bodyBuffer.WriteString("\r\n--" + boundary + "\r\n")

	bodyBuffer.WriteString("Content-Type: image/png\r\n")
	bodyBuffer.WriteString("Content-Disposition: attachment; filename=\"image.png\"\r\n")
	bodyBuffer.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
	bodyBuffer.WriteString(base64.StdEncoding.EncodeToString(imageData))
	bodyBuffer.WriteString("\r\n--" + boundary + "--\r\n")

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	err := smtp.SendMail(addr, auth, from, []string{to}, bodyBuffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func SendEmailWithPdfAttachment(to, subject, body string, pdfData []byte) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")

	from := smtpUser
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	boundary := "boundary123"

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = fmt.Sprintf("multipart/mixed; boundary=%s", boundary)

	var bodyBuffer bytes.Buffer
	for k, v := range headers {
		bodyBuffer.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	bodyBuffer.WriteString("\r\n--" + boundary + "\r\n")
	bodyBuffer.WriteString("Content-Type: text/plain; charset=utf-8\r\n\r\n")
	bodyBuffer.WriteString(body + "\r\n")
	bodyBuffer.WriteString("\r\n--" + boundary + "\r\n")

	bodyBuffer.WriteString("Content-Type: application/pdf\r\n")
	bodyBuffer.WriteString("Content-Disposition: attachment; filename=\"document.pdf\"\r\n")
	bodyBuffer.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
	bodyBuffer.WriteString(base64.StdEncoding.EncodeToString(pdfData))
	bodyBuffer.WriteString("\r\n--" + boundary + "--\r\n")

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	err := smtp.SendMail(addr, auth, from, []string{to}, bodyBuffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}