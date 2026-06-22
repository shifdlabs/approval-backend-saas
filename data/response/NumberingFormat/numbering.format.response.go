package numberinggroup

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type NumberingFormatResponse struct {
	Id                 *uuid.UUID `json:"id"`
	Name               string     `json:"name"`
	Format             string     `json:"format"`
	IncrementedByGroup bool       `json:"incrementedByGroup"`
	Separator          string     `json:"separator"`
	CreatedAt          *time.Time `json:"createdAt"`
	Group              string     `json:"group"`
}
