package proxy

import (
	"io"
	"net/http"
)

var srv *Service

type ProxyHandler struct{}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get(HEADER_CONTENT_TYPE)
	if ct != "" {
		w.Header().Set(HEADER_CONTENT_TYPE, ct)
	}
	defer r.Body.Close()
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}
	srv.Proxy(r.Context(), w, r)
}
