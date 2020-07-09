package proxy

import (
	"io"
	"net/http"
)

var srv *Service

type ProxyHandler struct{}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(HEADER_CONTENT_TYPE, CONTENT_TYPE_JSON)
	ct := r.Header.Get(HEADER_CONTENT_TYPE)
	if ct != "" {
		w.Header().Set(HEADER_CONTENT_TYPE, ct)
	}
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}
	res, err := srv.Proxy(r.Context(), r.Method, r.Header, r.URL.Path, r.Form, r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(res))
}
