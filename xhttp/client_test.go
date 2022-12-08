package xhttp

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/go-pay/gopher/bm"
	"github.com/go-pay/gopher/xlog"
)

type HttpGet struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

var ctx = context.Background()

func TestHttpGet(t *testing.T) {
	xlog.Level = xlog.DebugLevel
	client := NewClient()
	client.Timeout = 10 * time.Second
	// test
	_, bs, err := client.Get("http://www.baidu.com").EndBytes(ctx)
	if err != nil {
		xlog.Error(err)
		return
	}
	xlog.Debug(string(bs))

	//rsp := new(HttpGet)
	//_, err = client.Type(TypeJSON).Get("http://api.igoogle.ink/app/v1/ping").EndStruct(ctx, rsp)
	//if err != nil {
	//	xlog.Error(err)
	//	return
	//}
	//xlog.Debug(rsp)
}

func TestHttpUploadFile(t *testing.T) {
	xlog.Level = xlog.DebugLevel
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
	_, err = client.Type(TypeMultipartFormData).
		Post("http://localhost:2233/admin/v1/oss/uploadImage").
		SendMultipartBodyMap(bmm).
		EndStruct(ctx, rsp)
	if err != nil {
		xlog.Error(err)
		return
	}
	xlog.Debugf("%+v", rsp)
}
