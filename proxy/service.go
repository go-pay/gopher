package proxy

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Service struct {
	httpCli *http.Client
	Schema  SchemaType // SchemaHTTP or SchemaHTTPS
	Host    string
	Port    string
	Key     string
}

func New(c *Config) {
	srv = &Service{
		httpCli: new(http.Client),
		Schema:  c.ProxySchema,
		Host:    c.ProxyHost,
		Port:    c.ProxyPort,
		Key:     c.Key,
	}
}

func (s *Service) Proxy(c context.Context, method string, header http.Header, path string, params url.Values, body io.ReadCloser) (res []byte, err error) {
	var (
		req    *http.Request
		reader *strings.Reader
	)
	// 验证 Key
	key := header.Get(HEADER_CONTENT_KEY)
	if s.Key != key {
		return nil, fmt.Errorf("[%s] invalid key", key)
	}
	uri := fmt.Sprintf("%s", s.Schema) + s.Host + s.Port + path
	// Request
	m := strings.ToUpper(method)
	ct := header.Get(HEADER_CONTENT_TYPE)
	switch m {
	case HTTP_METHOD_POST:
		switch ct {
		case CONTENT_TYPE_JSON:
			jsbs, err := ioutil.ReadAll(io.LimitReader(body, int64(4<<20))) // default 4MB, change the size you want;
			if err != nil {
				return nil, err
			}
			reader = strings.NewReader(string(jsbs))
		case CONTENT_TYPE_FORM:
			reader = strings.NewReader(params.Encode())
		case CONTENT_TYPE_XML:
			xmlbs, err := ioutil.ReadAll(io.LimitReader(body, int64(4<<20)))
			if err != nil {
				return nil, err
			}
			reader = strings.NewReader(string(xmlbs))
		default:
			return nil, errors.New("request type error")
		}
		req, err = http.NewRequest(HTTP_METHOD_POST, uri, reader)
		if err != nil {
			return nil, err
		}
	case HTTP_METHOD_GET:
		pa := params.Encode()
		uri = uri + "?" + pa
		req, err = http.NewRequest(HTTP_METHOD_GET, uri, nil)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("only support GET and POST")
	}

	// Request Content
	req.Header = header
	s.httpCli.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, DisableKeepAlives: true}

	resp, err := s.httpCli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return res, nil
}
