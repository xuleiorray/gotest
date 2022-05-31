package utils

import "encoding/json"

func ToJSON(any interface{}) string {
	data, err := json.Marshal(any)
	if err != nil {
		log.Errorf("json marshal happens error, error msg: %s", err.Error())
		return ""
	}
	return string(data)
}
