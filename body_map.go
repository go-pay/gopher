package goutil

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"sync"
)

type BodyMap struct {
	sm sync.Map
}

// 设置参数
func (bm *BodyMap) Set(key string, value interface{}) {
	bm.sm.Store(key, value)
}

// 获取参数
func (bm *BodyMap) Get(key string) string {
	if bm == nil {
		return null
	}
	var (
		value interface{}
		ok    bool
		v     string
	)
	if value, ok = bm.sm.Load(key); !ok {
		return null
	}
	if v, ok = value.(string); ok {
		return v
	}
	return convertToString(value)
}

func convertToString(v interface{}) (str string) {
	if v == nil {
		return null
	}
	var (
		bs  []byte
		err error
	)
	if bs, err = json.Marshal(v); err != nil {
		return null
	}
	str = string(bs)
	return
}

// 删除参数
func (bm *BodyMap) Remove(key string) {
	bm.sm.Delete(key)
}

// Map长度
func (bm *BodyMap) Len() (len int) {
	bm.sm.Range(func(key, value interface{}) bool {
		len++
		return true
	})
	return
}

// 遍历Map
func (bm *BodyMap) Range(f func(key, value interface{}) bool) {
	bm.sm.Range(f)
}

type xmlMapEntry struct {
	XMLName xml.Name
	Value   interface{} `xml:",chardata"`
}

func (bm *BodyMap) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	if bm.Len() == 0 {
		return nil
	}
	if err = e.EncodeToken(start); err != nil {
		return
	}
	bm.sm.Range(func(key, value interface{}) bool {
		e.Encode(xmlMapEntry{XMLName: xml.Name{Local: key.(string)}, Value: value})
		return true
	})
	return e.EncodeToken(start.End())
}

func (bm *BodyMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	for {
		var e xmlMapEntry
		err = d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return
		}
		bm.Set(e.XMLName.Local, e.Value)
	}
	return
}
