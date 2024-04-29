package validator

import (
	"encoding/json"
	"strings"
)

type Validator interface {
	Validate(data string) bool
}

type JSONValidator struct{}

func NewJSONValidator() *JSONValidator {
	return &JSONValidator{}
}

func (v *JSONValidator) Validate(data string) bool {
	data = strings.TrimSpace(data)

	if !strings.HasPrefix(data, "{") || !strings.HasSuffix(data, "}") {
		return false
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(data), &obj); err != nil {
		return false
	}

	_, messageExists := obj["message"]
	_, timestampExists := obj["timestamp"]
	return messageExists && timestampExists
}
