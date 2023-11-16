package misc

import "encoding/json"

func Pick(data any, fields ...string) map[string]any {
	bytes, _ := json.Marshal(&data)
	fullMap := make(map[string]any)
	json.Unmarshal(bytes, &fullMap)
	result := make(map[string]any, len(fields))
	for _, f := range fields {
		if v, ok := fullMap[f]; ok {
			result[f] = v
		}
	}
	return result
}
