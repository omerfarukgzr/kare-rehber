package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONB map[string]any

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value any) error {
	if value == nil {
		*j = nil
		return nil
	}
	var src []byte
	switch v := value.(type) {
	case []byte:
		src = v
	case string:
		src = []byte(v)
	default:
		return errors.New("jsonb: unsupported scan type")
	}
	return json.Unmarshal(src, j)
}
