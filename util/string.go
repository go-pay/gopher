package util

import (
	"net/url"
	"strconv"
	"strings"
	"sync"
)

const (
	NULL = ""
)

var (
	bfPool = sync.Pool{
		New: func() any {
			return &strings.Builder{}
		},
	}
)

// JoinInts format int64 slice like:n1,n2,n3.
func JoinInts(is []int64) string {
	if len(is) == 0 {
		return NULL
	}
	if len(is) == 1 {
		return strconv.FormatInt(is[0], 10)
	}
	buf, ok := bfPool.Get().(*strings.Builder)
	if ok && buf != nil {
		for _, i := range is {
			buf.WriteString(strconv.FormatInt(i, 10))
			buf.WriteByte(',')
		}
		s := buf.String()
		if len(s) > 0 {
			s = s[:len(s)-1]
		}
		buf.Reset()
		bfPool.Put(buf)
		return s
	}
	return NULL
}

// SplitInts split string into int64 slice.
func SplitInts(s string) ([]int64, error) {
	if s == NULL {
		return nil, nil
	}
	sArr := strings.Split(s, ",")
	res := make([]int64, 0, len(sArr))
	for _, sc := range sArr {
		i, err := strconv.ParseInt(sc, 10, 64)
		if err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}

func FormatURLParam(body map[string]any) (urlParam string) {
	v := url.Values{}
	for key, value := range body {
		v.Add(key, value.(string))
	}
	return v.Encode()
}
