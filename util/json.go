package util

import "encoding/json"

func MarshalString(v any) string {
	bs, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(bs)
}

func MarshalBytes(v any) []byte {
	bs, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return bs
}

func UnmarshalString(jsonStr string, v any) error {
	return json.Unmarshal([]byte(jsonStr), v)
}

func UnmarshalBytes(bs []byte, v any) error {
	return json.Unmarshal(bs, v)
}
