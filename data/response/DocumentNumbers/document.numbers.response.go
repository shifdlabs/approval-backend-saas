package documentnumbers

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type DocumentNumbersResponse struct {
	Id                  *uuid.UUID `json:"id"`
	DocumentNumber      string     `json:"documentNumber"`
	NumberingFormatName string     `json:"numberingFormatName"`
	NumberingGroupName  string     `json:"numberingGroupName"`
	NumberingFormatId   *uuid.UUID `json:"numberingFormatId"`
	CreatedAt           time.Time  `json:"createdAt"`
}
