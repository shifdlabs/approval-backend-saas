package response

import (
	user "Microservice/data/request/User"
)

type PreviewImportResponse struct {
	Users    []user.ImportedUserData `json:"users"`
	Errors   []ImportError           `json:"errors"`
	Warnings []ImportWarning         `json:"warnings"`
}

type ImportError struct {
	Row     int    `json:"row"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ImportWarning struct {
	Row     int    `json:"row"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

type BulkImportResponse struct {
	SuccessCount int           `json:"successCount"`
	FailedCount  int           `json:"failedCount"`
	Errors       []ImportError `json:"errors"`
}
