package json

import (
	"encoding/json"
)

func Serialize(data interface{}) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func Deserialize(data string, target interface{}) error {
	b := []byte(data)
	return json.Unmarshal(b, target)
}
