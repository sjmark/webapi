package utils

import "encoding/json"

func GetJsonStruct(object interface{}) string {

	jsonobject, err := json.Marshal(object)

	if err != nil {
		return ""
	}

	return string(jsonobject)
}
