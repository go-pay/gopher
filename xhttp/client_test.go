package xhttp

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/iGoogle-ink/gopher/bm"
	"github.com/iGoogle-ink/gopher/xlog"
)

type HttpGet struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func TestHttpGet(t *testing.T) {
	client := NewClient()
	client.Timeout = 10 * time.Second
	// test
	_, bs, errs := client.Get("http://www.baidu.com").EndBytes()
	if len(errs) > 0 {
		xlog.Error(errs[0])
		return
	}
	xlog.Debug(string(bs))

	//rsp := new(HttpGet)
	//_, errs := client.Type(TypeJSON).Get("http://api.igoogle.ink/app/v1/ping").EndStruct(rsp)
	//if len(errs) > 0 {
	//	xlog.Error(errs[0])
	//	return
	//}
	//xlog.Debug(rsp)
}

func TestHttpUploadFile(t *testing.T) {
	fileContent, err := ioutil.ReadFile("logo.png")
	if err != nil {
		xlog.Error(err)
		return
	}
	//xlog.Debug("fileByte：", string(fileContent))

	bmm := make(bm.BodyMap)
	bmm.SetBodyMap("meta", func(bm bm.BodyMap) {
		bm.Set("filename", "123.jpg").
			Set("sha256", "ad4465asd4fgw5q")
	}).SetFormFile("image", &bm.File{Name: "logo.png", Content: fileContent})

	client := NewClient()
	client.Timeout = 10 * time.Second

	rsp := new(HttpGet)
	_, errs := client.Type(TypeMultipartFormData).
		Post("http://localhost:2233/admin/v1/oss/uploadImage").
		SendMultipartBodyMap(bmm).
		EndStruct(rsp)
	if len(errs) > 0 {
		xlog.Error(errs[0])
		return
	}
	xlog.Debugf("%+v", rsp)
}
