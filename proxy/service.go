package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/go-pay/xlog"
)

type Service struct {
	httpCli *http.Client
	Schema  SchemaType // SchemaHTTP or SchemaHTTPS
	r       *http.Request
	Host    string
	Port    string
	Key     string
	log     *log.Logger
}

func NewHandler(c *Config) (handler *Handler) {
	srv = &Service{
		httpCli: new(http.Client),
		Schema:  c.ProxySchema,
		Host:    c.ProxyHost,
		Port:    c.ProxyPort,
		Key:     c.Key,
		log:     log.New(os.Stdout, string(xlog.Magenta)+" [PROXY] "+string(xlog.Reset), log.Ldate|log.Lmicroseconds),
	}
	return &Handler{c: c}
}

// Proxy
func (s *Service) Proxy(c context.Context, w http.ResponseWriter, r *http.Request) {
	var (
		req     *http.Request
		reader  *strings.Reader
		err     error
		rMethod = r.Method
		rHeader = r.Header
		rUri    = r.RequestURI
		pa      = r.Form.Encode()
		rBody   = r.Body
	)
	s.r = r
	// 验证 Key
	key := rHeader.Get(HEADER_CONTENT_KEY)
	if s.Key != key {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("[%s] invalid key", key))
		return
	}
	uri := string(s.Schema) + s.Host + s.Port + rUri
	// Request
	m := strings.ToUpper(r.Method)
	ct := rHeader.Get(HEADER_CONTENT_TYPE)
	switch m {
	case HTTP_METHOD_POST:
		switch ct {
		case CONTENT_TYPE_JSON:
			jsbs, e := io.ReadAll(io.LimitReader(rBody, int64(4<<20))) // default 4MB, change the size you want;
			if e != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = io.WriteString(w, e.Error())
				return
			}
			reader = strings.NewReader(string(jsbs))
		case CONTENT_TYPE_FORM:
			reader = strings.NewReader(pa)
		case CONTENT_TYPE_XML:
			xmlbs, e := io.ReadAll(io.LimitReader(rBody, int64(4<<20)))
			if e != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = io.WriteString(w, e.Error())
				return
			}
			reader = strings.NewReader(string(xmlbs))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = io.WriteString(w, "request type error")
			return
		}
		req, err = http.NewRequest(HTTP_METHOD_POST, uri, reader)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = io.WriteString(w, err.Error())
			return
		}
	case HTTP_METHOD_GET:
		req, err = http.NewRequest(HTTP_METHOD_GET, uri, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = io.WriteString(w, err.Error())
			return
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, "only support GET and POST")
		return
	}

	// Request Content
	req.Header = rHeader
	s.httpCli.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, DisableKeepAlives: true}

	resp, err := s.httpCli.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, err.Error())
		return
	}
	defer resp.Body.Close()
	s.log.Printf("| %d | %s | %s      %s\n", resp.StatusCode, s.clientIP(), rMethod, r.RequestURI)
	rsp, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, err.Error())
		return
	}
	for k := range resp.Header {
		w.Header().Set(k, resp.Header.Get(k))
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(rsp)
}

func (s *Service) clientIP() string {
	rHeader := s.r.Header
	clientIP := rHeader.Get("X-Forwarded-For")
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(rHeader.Get("X-Real-Ip"))
	}
	if clientIP != "" {
		return clientIP
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(s.r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}
