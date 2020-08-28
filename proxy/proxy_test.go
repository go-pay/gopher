package proxy

import (
	"log"
	"net/http"
	"testing"

	"github.com/iGoogle-ink/gotil/xlog"
)

func TestService_Proxy(t *testing.T) {
	// https://api.igoogle.ink/app/v1/ping
	// test path  /app/v1/ping
	c := &Config{
		ProxySchema: SchemaHTTPS,
		ProxyHost:   "api.igoogle.ink",
		ProxyPort:   "",
		ServerPort:  ":2233",
		Key:         "5urivxGzAqOzdJotjbK7AOmayYYnyHlP",
	}

	New(c)

	http.Handle("/", &ProxyHandler{})

	xlog.Info("Proxy Started")
	if err := http.ListenAndServe(c.ServerPort, nil); err != nil {
		log.Fatal("Proxy Start Err：", err)
	}
}
