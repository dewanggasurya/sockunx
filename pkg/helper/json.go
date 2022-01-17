package helper

import "encoding/json"

// ToJSON data to json
func ToJSON(data interface{}) []byte {
	bytes, _ := json.Marshal(data)
	return bytes
}

// ToJSONString data to json
func ToJSONString(data interface{}) string {
	return string(ToJSON(data))
}

// FromJSON data to json
func FromJSON(data []byte, result interface{}) error {
	e := json.Unmarshal(data, result)
	return e
}

// FromJSONString data to json
func FromJSONString(data string, result interface{}) error {
	return FromJSON([]byte(data), result)
}
