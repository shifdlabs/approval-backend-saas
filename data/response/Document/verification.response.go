package document

type VerificationResponse struct {
	Subject            string `json:"subject"`
	Body               string `json:"body"`
	DocumentNumber     string `json:"documentNumber"`
	OrganizationName   string `json:"organizationName"`
	ApprovalDate       string `json:"approvalDate"`
	LastApproverName   string `json:"lastApproverName"`
	Type               int    `json:"type"`
	CompanyLogoUrl     string `json:"companyLogoUrl"`
	CompanyAddress     string `json:"companyAddress"`
	CompanyCity        string `json:"companyCity"`
	CompanyPhone       string `json:"companyPhone"`
	CompanyEmail       string `json:"companyEmail"`
	CompanyDescription string `json:"companyDescription"`
}
