package position

import (
	uuid "github.com/satori/go.uuid"
)

type PositionResponse struct {
	Id   *uuid.UUID `json:"id"`
	Name string     `json:"name"`
}
