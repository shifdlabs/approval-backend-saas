package emailService

import (
	"fmt"

	"github.com/resend/resend-go/v2"
)

// ── Real implementation using Resend ──────────────────────────────────────────

type EmailServiceImpl struct {
	client      *resend.Client
	fromEmail   string
	frontendURL string
}

func NewEmailService(apiKey, fromEmail, frontendURL string) EmailService {
	if apiKey == "" {
		return &NoopEmailService{}
	}
	return &EmailServiceImpl{
		client:      resend.NewClient(apiKey),
		fromEmail:   fromEmail,
		frontendURL: frontendURL,
	}
}

func (e *EmailServiceImpl) SendApprovalRequest(toEmail, toName, fromName, documentSubject, documentURL string) error {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="id">
<head><meta charset="UTF-8"/></head>
<body style="margin:0;padding:0;background:#f4f6fb;font-family:sans-serif;">
  <div style="max-width:600px;margin:32px auto;background:#fff;border-radius:14px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,.08);">
    <div style="background:#1e3a5f;padding:28px 32px;">
      <p style="margin:0;color:#fff;font-size:20px;font-weight:700;">Shifd Approval</p>
    </div>
    <div style="padding:32px;">
      <h2 style="margin:0 0 8px;color:#1e3a5f;font-size:18px;">Permintaan Persetujuan</h2>
      <p style="color:#555;margin:0 0 20px;">Halo <strong>%s</strong>,</p>
      <p style="color:#555;margin:0 0 20px;"><strong>%s</strong> meminta persetujuan Anda untuk dokumen berikut:</p>
      <div style="background:#f0f4ff;border-left:4px solid #1e3a5f;border-radius:0 8px 8px 0;padding:14px 18px;margin-bottom:28px;">
        <p style="margin:0;color:#1e3a5f;font-weight:600;font-size:15px;">%s</p>
      </div>
      <a href="%s" style="display:inline-block;background:#1e3a5f;color:#fff;padding:13px 28px;border-radius:9px;text-decoration:none;font-weight:700;font-size:14px;">
        Lihat &amp; Setujui Dokumen →
      </a>
    </div>
    <div style="padding:20px 32px;border-top:1px solid #eee;">
      <p style="margin:0;color:#aaa;font-size:12px;">Email ini dikirim otomatis oleh Shifd Approval. Jangan balas email ini.</p>
    </div>
  </div>
</body>
</html>`, toName, fromName, documentSubject, documentURL)

	_, err := e.client.Emails.Send(&resend.SendEmailRequest{
		From:    e.fromEmail,
		To:      []string{toEmail},
		Subject: fmt.Sprintf("Permintaan Persetujuan: %s", documentSubject),
		Html:    html,
	})
	return err
}

func (e *EmailServiceImpl) SendDocumentApproved(toEmail, toName, documentSubject, documentURL string) error {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="id">
<head><meta charset="UTF-8"/></head>
<body style="margin:0;padding:0;background:#f4f6fb;font-family:sans-serif;">
  <div style="max-width:600px;margin:32px auto;background:#fff;border-radius:14px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,.08);">
    <div style="background:#1e7a4a;padding:28px 32px;">
      <p style="margin:0;color:#fff;font-size:20px;font-weight:700;">Shifd Approval</p>
    </div>
    <div style="padding:32px;">
      <h2 style="margin:0 0 8px;color:#1e7a4a;font-size:18px;">✓ Dokumen Disetujui</h2>
      <p style="color:#555;margin:0 0 20px;">Halo <strong>%s</strong>,</p>
      <p style="color:#555;margin:0 0 20px;">Dokumen Anda telah disetujui oleh semua pihak:</p>
      <div style="background:#f0fff6;border-left:4px solid #1e7a4a;border-radius:0 8px 8px 0;padding:14px 18px;margin-bottom:28px;">
        <p style="margin:0;color:#1e7a4a;font-weight:600;font-size:15px;">%s</p>
      </div>
      <a href="%s" style="display:inline-block;background:#1e7a4a;color:#fff;padding:13px 28px;border-radius:9px;text-decoration:none;font-weight:700;font-size:14px;">
        Lihat Dokumen →
      </a>
    </div>
    <div style="padding:20px 32px;border-top:1px solid #eee;">
      <p style="margin:0;color:#aaa;font-size:12px;">Email ini dikirim otomatis oleh Shifd Approval. Jangan balas email ini.</p>
    </div>
  </div>
</body>
</html>`, toName, documentSubject, documentURL)

	_, err := e.client.Emails.Send(&resend.SendEmailRequest{
		From:    e.fromEmail,
		To:      []string{toEmail},
		Subject: fmt.Sprintf("Dokumen Disetujui: %s", documentSubject),
		Html:    html,
	})
	return err
}

func (e *EmailServiceImpl) SendDocumentRejected(toEmail, toName, documentSubject, rejectedBy, reason, documentURL string) error {
	reasonSection := ""
	if reason != "" {
		reasonSection = fmt.Sprintf(`
      <p style="color:#555;margin:0 0 8px;">Alasan penolakan:</p>
      <div style="background:#fff5f5;border-left:4px solid #c0392b;border-radius:0 8px 8px 0;padding:12px 16px;margin-bottom:24px;">
        <p style="margin:0;color:#555;font-size:14px;">%s</p>
      </div>`, reason)
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="id">
<head><meta charset="UTF-8"/></head>
<body style="margin:0;padding:0;background:#f4f6fb;font-family:sans-serif;">
  <div style="max-width:600px;margin:32px auto;background:#fff;border-radius:14px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,.08);">
    <div style="background:#c0392b;padding:28px 32px;">
      <p style="margin:0;color:#fff;font-size:20px;font-weight:700;">Shifd Approval</p>
    </div>
    <div style="padding:32px;">
      <h2 style="margin:0 0 8px;color:#c0392b;font-size:18px;">✗ Dokumen Ditolak</h2>
      <p style="color:#555;margin:0 0 20px;">Halo <strong>%s</strong>,</p>
      <p style="color:#555;margin:0 0 20px;">Dokumen Anda ditolak oleh <strong>%s</strong>:</p>
      <div style="background:#fff5f5;border-left:4px solid #c0392b;border-radius:0 8px 8px 0;padding:14px 18px;margin-bottom:20px;">
        <p style="margin:0;color:#c0392b;font-weight:600;font-size:15px;">%s</p>
      </div>
      %s
      <a href="%s" style="display:inline-block;background:#c0392b;color:#fff;padding:13px 28px;border-radius:9px;text-decoration:none;font-weight:700;font-size:14px;">
        Lihat Dokumen →
      </a>
    </div>
    <div style="padding:20px 32px;border-top:1px solid #eee;">
      <p style="margin:0;color:#aaa;font-size:12px;">Email ini dikirim otomatis oleh Shifd Approval. Jangan balas email ini.</p>
    </div>
  </div>
</body>
</html>`, toName, rejectedBy, documentSubject, reasonSection, documentURL)

	_, err := e.client.Emails.Send(&resend.SendEmailRequest{
		From:    e.fromEmail,
		To:      []string{toEmail},
		Subject: fmt.Sprintf("Dokumen Ditolak: %s", documentSubject),
		Html:    html,
	})
	return err
}

// ── No-op fallback (used when RESEND_API_KEY is not set) ─────────────────────

type NoopEmailService struct{}

func (n *NoopEmailService) SendApprovalRequest(_, _, _, _, _ string) error  { return nil }
func (n *NoopEmailService) SendDocumentApproved(_, _, _, _ string) error    { return nil }
func (n *NoopEmailService) SendDocumentRejected(_, _, _, _, _, _ string) error { return nil }
