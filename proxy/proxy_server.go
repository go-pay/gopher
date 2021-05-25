package proxy

import (
	"io"
	"net/http"
	"sync"
)

var srv *Service

type ProxyHandler struct {
	c        *Config
	handlers map[string]http.Handler
	mu       sync.Mutex
}

func (p *ProxyHandler) Handle(path string, handler http.Handler) {
	p.mu.Lock()
	p.handlers[path] = handler
	p.mu.Unlock()
}

func (p *ProxyHandler) ListenAndServe() error {
	http.Handle("/", p)

	if len(p.handlers) > 0 {
		for path, handler := range p.handlers {
			http.Handle(path, handler)
		}
	}

	// listen and serve
	if err := http.ListenAndServe(p.c.ServerPort, nil); err != nil {
		return err
	}
	return nil
}

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
