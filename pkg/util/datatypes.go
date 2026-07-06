package util

import (
	"encoding/json"

	"gorm.io/datatypes"
)

func ToJSON(v any) datatypes.JSON {
	b, _ := json.Marshal(v)
	return datatypes.JSON(b)
}
