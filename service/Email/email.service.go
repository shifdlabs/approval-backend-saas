package emailService

type EmailService interface {
	// Notify approver that a document needs their approval
	SendApprovalRequest(toEmail, toName, fromName, documentSubject, documentURL string) error
	// Notify author that all approvers have approved
	SendDocumentApproved(toEmail, toName, documentSubject, documentURL string) error
	// Notify author that a step rejected the document
	SendDocumentRejected(toEmail, toName, documentSubject, rejectedBy, reason, documentURL string) error
}
