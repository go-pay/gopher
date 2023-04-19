package web

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopher/ecode"
	"github.com/go-pay/gopher/proxy"
)

var (
	httpCli = new(http.Client)
	pLog    = log.New(os.Stdout, "[PROXY] ", log.Ldate|log.Ltime)
)

// GinProxy gin request proxy and get rsp
func GinProxy[Rsp any](c *gin.Context, method, host, uri string) (rspParam Rsp, err error) {
	var (
		req     *http.Request
		reader  *strings.Reader
		rMethod = c.Request.Method
		rHeader = c.Request.Header
		rUri    = c.Request.RequestURI
		pa      = c.Request.Form.Encode()
		rBody   = c.Request.Body
	)
	vo := reflect.ValueOf(rspParam)
	if vo.Kind() != reflect.Ptr {
		err = ecode.NewV2(500, "", "rspParam must be point kind")
		return
	}
	if uri != "" {
		rUri = uri
	}
	if method != "" {
		rMethod = strings.ToUpper(method)
	}
	uri = host + rUri
	// Request
	ct := rHeader.Get(proxy.HEADER_CONTENT_TYPE)
	switch rMethod {
	case proxy.HTTP_METHOD_POST:
		switch ct {
		case proxy.CONTENT_TYPE_JSON:
			jsbs, e := io.ReadAll(rBody)
			if e != nil {
				err = e
				return
			}
			reader = strings.NewReader(string(jsbs))
		case proxy.CONTENT_TYPE_FORM:
			reader = strings.NewReader(pa)
		}
		req, err = http.NewRequestWithContext(c, rMethod, uri, reader)
		if err != nil {
			return
		}
	case proxy.HTTP_METHOD_GET:
		req, err = http.NewRequestWithContext(c, rMethod, uri, nil)
		if err != nil {
			return
		}
	default:
		err = ecode.NewV2(500, "", "only support GET and POST")
		return
	}

	// Request Content
	req.Header = rHeader
	req.Header.Del("Accept-Encoding")
	//xlog.Warnf("reqH: %+v", req.Header)
	httpCli.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, DisableKeepAlives: true}

	resp, e := httpCli.Do(req)
	if e != nil {
		err = e
		return
	}
	defer resp.Body.Close()
	pLog.Printf("| %d | %s      %s\n", resp.StatusCode, rMethod, uri)
	rspBytes, e := io.ReadAll(resp.Body)
	if e != nil {
		err = e
		return
	}
	if resp.StatusCode != 200 {
		err = ecode.NewV2(resp.StatusCode, "", string(rspBytes))
		return
	}
	//xlog.Infof("rspBytes:%v", string(rspBytes))
	res := &HttpRsp[Rsp]{}
	if err = json.Unmarshal(rspBytes, res); err != nil {
		return
	}
	rspParam = res.Data
	//xlog.Infof("rspParam: %+v", rspParam)
	return rspParam, nil
}

// GinPureProxy gin request proxy
func GinPureProxy(c *gin.Context, method, host, uri string) {
	var (
		req     *http.Request
		reader  *strings.Reader
		rMethod = c.Request.Method
		rHeader = c.Request.Header
		rUri    = c.Request.RequestURI
		pa      = c.Request.Form.Encode()
		rBody   = c.Request.Body
		err     error
	)
	if uri != "" {
		rUri = uri
	}
	if method != "" {
		rMethod = strings.ToUpper(method)
	}
	uri = host + rUri
	// Request
	ct := rHeader.Get(proxy.HEADER_CONTENT_TYPE)
	switch rMethod {
	case proxy.HTTP_METHOD_POST:
		switch ct {
		case proxy.CONTENT_TYPE_JSON:
			jsbs, e := io.ReadAll(rBody)
			if e != nil {
				err = e
				JSON(c, "", err)
				return
			}
			reader = strings.NewReader(string(jsbs))
		case proxy.CONTENT_TYPE_FORM:
			reader = strings.NewReader(pa)
		}
		req, err = http.NewRequestWithContext(c, rMethod, uri, reader)
		if err != nil {
			JSON(c, "", err)
			return
		}
	case proxy.HTTP_METHOD_GET:
		req, err = http.NewRequestWithContext(c, rMethod, uri, nil)
		if err != nil {
			JSON(c, "", err)
			return
		}
	default:
		err = ecode.NewV2(500, "", "only support GET and POST")
		JSON(c, "", err)
		return
	}

	// Request Content
	req.Header = rHeader
	req.Header.Del("Accept-Encoding")
	//xlog.Warnf("reqH: %+v", req.Header)
	httpCli.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, DisableKeepAlives: true}

	resp, e := httpCli.Do(req)
	if e != nil {
		err = e
		JSON(c, "", err)
		return
	}
	defer resp.Body.Close()
	pLog.Printf("| %d | %s      %s\n", resp.StatusCode, rMethod, uri)
	rspBytes, e := io.ReadAll(resp.Body)
	if e != nil {
		err = e
		JSON(c, "", err)
		return
	}
	if resp.StatusCode != 200 {
		err = ecode.NewV2(resp.StatusCode, "", string(rspBytes))
		JSON(c, "", err)
		return
	}
	JSON(c, rspBytes, nil)
}

func ClientIP(r *http.Request, rHeader http.Header) string {
	cIp := rHeader.Get("X-Forwarded-For")
	cIp = strings.TrimSpace(strings.Split(cIp, ",")[0])
	if cIp == "" {
		cIp = strings.TrimSpace(rHeader.Get("X-Real-Ip"))
	}
	if cIp != "" {
		return cIp
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}
