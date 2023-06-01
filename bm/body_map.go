package bm

import (
	"encoding/xml"
	"errors"
	"io"
	"net/url"
	"sort"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/go-pay/gopher/util"
)

type BodyMap map[string]any

type xmlMapMarshal struct {
	XMLName xml.Name
	Value   any `xml:",cdata"`
}

type xmlMapUnmarshal struct {
	XMLName xml.Name
	Value   string `xml:",cdata"`
}

type File struct {
	Name    string `json:"name"`
	Content []byte `json:"content"`
}

// 设置参数
func (bm BodyMap) Set(key string, value any) BodyMap {
	bm[key] = value
	return bm
}

func (bm BodyMap) SetBodyMap(key string, value func(bm BodyMap)) BodyMap {
	_bm := make(BodyMap)
	value(_bm)
	bm[key] = _bm
	return bm
}

// 设置 FormFile
func (bm BodyMap) SetFormFile(key string, file *File) BodyMap {
	bm[key] = file
	return bm
}

// 获取参数转换string
func (bm BodyMap) GetString(key string) string {
	if bm == nil {
		return util.NULL
	}
	value, ok := bm[key]
	if !ok {
		return util.NULL
	}
	v, ok := value.(string)
	if !ok {
		return convertToString(value)
	}
	return v
}

// 获取原始参数
func (bm BodyMap) GetInterface(key string) any {
	if bm == nil {
		return nil
	}
	return bm[key]
}

// 删除参数
func (bm BodyMap) Remove(key string) {
	delete(bm, key)
}

// 置空BodyMap
func (bm BodyMap) Reset() {
	for k := range bm {
		delete(bm, k)
	}
}

func (bm BodyMap) JsonBody() (jb string) {
	bs, err := sonic.Marshal(bm)
	if err != nil {
		return ""
	}
	jb = string(bs)
	return jb
}

func (bm BodyMap) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	if len(bm) == 0 {
		return nil
	}
	start.Name = xml.Name{Space: util.NULL, Local: "xml"}
	if err = e.EncodeToken(start); err != nil {
		return
	}
	for k := range bm {
		if v := bm.GetString(k); v != util.NULL {
			_ = e.Encode(xmlMapMarshal{XMLName: xml.Name{Local: k}, Value: v})
		}
	}
	return e.EncodeToken(start.End())
}

func (bm *BodyMap) UnmarshalXML(d *xml.Decoder, _ xml.StartElement) (err error) {
	for {
		var e xmlMapUnmarshal
		err = d.Decode(&e)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		bm.Set(e.XMLName.Local, e.Value)
	}
}

// ("bar=baz&foo=quux") sorted by key.
func (bm BodyMap) EncodeURLParams() string {
	if bm == nil {
		return util.NULL
	}
	var (
		buf  strings.Builder
		keys []string
	)
	for k := range bm {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if v := bm.GetString(k); v != util.NULL {
			buf.WriteString(url.QueryEscape(k))
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
			buf.WriteByte('&')
		}
	}
	if buf.Len() <= 0 {
		return util.NULL
	}
	return buf.String()[:buf.Len()-1]
}

func (bm BodyMap) CheckEmptyError(keys ...string) error {
	var emptyKeys []string
	for _, k := range keys {
		if v := bm.GetString(k); v == util.NULL {
			emptyKeys = append(emptyKeys, k)
		}
	}
	if len(emptyKeys) > 0 {
		return errors.New(strings.Join(emptyKeys, ", ") + " : cannot be empty")
	}
	return nil
}

func convertToString(v any) (str string) {
	if v == nil {
		return util.NULL
	}
	var (
		bs  []byte
		err error
	)
	if bs, err = sonic.Marshal(v); err != nil {
		return util.NULL
	}
	str = string(bs)
	return
}
