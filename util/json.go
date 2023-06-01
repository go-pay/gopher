package util

import (
	"github.com/bytedance/sonic"
)

func MarshalString(v any) string {
	bs, err := sonic.Marshal(v)
	if err != nil {
		return ""
	}
	return string(bs)
}

func MarshalBytes(v any) []byte {
	bs, err := sonic.Marshal(v)
	if err != nil {
		return nil
	}
	return bs
}

func UnmarshalString(jsonStr string, v any) error {
	return sonic.Unmarshal([]byte(jsonStr), v)
}

func UnmarshalBytes(bs []byte, v any) error {
	return sonic.Unmarshal(bs, v)
}
