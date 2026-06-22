package helper

import (
	"encoding/json"

	"gorm.io/datatypes"
)

func ToJSON(v interface{}) datatypes.JSON {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return datatypes.JSON(jsonData)
}
