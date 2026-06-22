package documentSequence

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/datatypes"
)

type UserLogResponse struct {
	Id       *uuid.UUID     `json:"id"`
	UserID   *uuid.UUID     `json:"userId"`
	UserName string         `json:"userName"`
	Action   string         `json:"action"`
	Module   string         `json:"module"`
	Log      datatypes.JSON `json:"log"`
	LogDate  time.Time      `json:"logDate"`
}
