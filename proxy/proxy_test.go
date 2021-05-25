package proxy

import (
	"testing"
)

func TestService_Proxy(t *testing.T) {

	// https://api.igoogle.ink/app/v1/ping
	// test path  /app/v1/ping

	// 解开下面注释测试
	//c := &Config{
	//	ProxySchema: SchemaHTTPS,
	//	ProxyHost:   "api.igoogle.ink",
	//	ProxyPort:   "",
	//	ServerPort:  ":2233",
	//	Key:         "5urivxGzAqOzdJotjbK7AOmayYYnyHlP",
	//}
	//
	//handler := NewHandler(c)
	//
	//if err := handler.ListenAndServe(); err != nil {
	//	log.Fatal("Proxy Start Err：", err)
	//}
}
