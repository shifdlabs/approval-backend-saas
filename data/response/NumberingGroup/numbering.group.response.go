package numberinggroup

import (
	uuid "github.com/satori/go.uuid"
)

type NumberingGroupResponse struct {
	Id          *uuid.UUID `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	TotalItem   int        `json:"totalItem"`
}
