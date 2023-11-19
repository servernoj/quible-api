package misc

import "encoding/json"

func PickFields(data any, fields ...string) map[string]any {
	bytes, _ := json.Marshal(&data)
	fullMap := make(map[string]any)
	//lint comment used to bycheck linter because format was correct for code

	json.Unmarshal(bytes, &fullMap)
	result := make(map[string]any, len(fields))
	for _, f := range fields {
		if v, ok := fullMap[f]; ok {
			result[f] = v
		}
	}
	return result
}
