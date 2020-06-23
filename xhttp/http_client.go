package gopay

import (
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/iGoogle-ink/goutil"
)

const (
	POST           = "POST"
	GET            = "GET"
	TypeJSON       = "json"
	TypeXML        = "xml"
	TypeUrlencoded = "urlencoded"
	TypeForm       = "form"
	TypeFormData   = "form-data"
)

var Types = map[string]string{
	TypeJSON:       "application/json",
	TypeXML:        "application/xml",
	TypeForm:       "application/x-www-form-urlencoded",
	TypeFormData:   "application/x-www-form-urlencoded",
	TypeUrlencoded: "application/x-www-form-urlencoded",
}

type Client struct {
	HttpClient    *http.Client
	Transport     *http.Transport
	Header        http.Header
	Timeout       time.Duration
	Url           string
	Method        string
	RequestType   string
	FormString    string
	ContentType   string
	UnmarshalType string
	Types         map[string]string
	JsonByte      []byte
	Errors        []error
	mu            sync.RWMutex
}

// NewHttpClient , default tls.Config{InsecureSkipVerify: true}
func NewHttpClient() (client *Client) {
	client = &Client{
		HttpClient: &http.Client{
			Timeout: 60 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
				DisableKeepAlives: true,
			},
		},
		Transport:     &http.Transport{},
		Header:        make(http.Header),
		RequestType:   TypeUrlencoded,
		UnmarshalType: TypeJSON,
		Errors:        make([]error, 0),
	}
	return client
}

func (c *Client) SetTLSConfig(tlsCfg *tls.Config) (client *Client) {
	c.mu.Lock()
	c.Transport.TLSClientConfig = tlsCfg
	c.mu.Unlock()
	return c
}

func (c *Client) SetTimeout(timeout time.Duration) (client *Client) {
	c.mu.Lock()
	c.Timeout = timeout
	c.mu.Unlock()
	return c
}

func (c *Client) Post(url string) (client *Client) {
	c.mu.Lock()
	c.Method = POST
	c.Url = url
	c.mu.Unlock()
	return c
}

func (c *Client) Type(typeStr string) (client *Client) {
	if _, ok := Types[typeStr]; ok {
		c.mu.Lock()
		c.RequestType = typeStr
		c.mu.Unlock()
	} else {
		c.Errors = append(c.Errors, errors.New("Type func: incorrect type \""+typeStr+"\""))
	}
	return c
}

func (c *Client) Get(url string) (client *Client) {
	c.mu.Lock()
	c.Method = GET
	c.Url = url
	c.mu.Unlock()
	return c
}

func (c *Client) SendStruct(v interface{}) (client *Client) {
	bs, err := json.Marshal(v)
	if err != nil {
		c.Errors = append(c.Errors, err)
		return c
	}
	c.mu.Lock()
	c.JsonByte = bs
	c.mu.Unlock()
	return c
}

func (c *Client) SendString(v string) (client *Client) {
	c.mu.Lock()
	c.FormString = v
	c.mu.Unlock()
	return c
}

func (c *Client) EndStruct(v interface{}) (res *http.Response, errs []error) {
	res, bs, errs := c.EndBytes()
	if errs != nil && len(errs) > 0 {
		c.Errors = append(c.Errors, errs...)
		return nil, c.Errors
	}
	c.mu.RLock()
	defer c.mu.RUnlock()

	switch c.UnmarshalType {
	case TypeJSON:
		err := json.Unmarshal(bs, &v)
		if err != nil {
			c.Errors = append(c.Errors, fmt.Errorf("json.Unmarshal(%s)：%w", string(bs), err))
			return nil, c.Errors
		}
		return res, nil
	case TypeXML:
		err := xml.Unmarshal(bs, &v)
		if err != nil {
			c.Errors = append(c.Errors, fmt.Errorf("xml.Unmarshal(%s)：%w", string(bs), err))
			return nil, c.Errors
		}
		return res, nil
	default:
		c.Errors = append(c.Errors, errors.New("UnmarshalType Type Wrong"))
		return nil, c.Errors
	}
}

func (c *Client) EndBytes() (res *http.Response, bs []byte, errs []error) {
	if len(c.Errors) > 0 {
		return nil, nil, c.Errors
	}
	var reader = strings.NewReader(goutil.NULL)

	req, err := func() (*http.Request, error) {
		c.mu.RLock()
		defer c.mu.RUnlock()

		switch c.Method {
		case GET:
			//todo: nothing
		case POST:
			switch c.RequestType {
			case TypeJSON:
				if c.JsonByte != nil {
					reader = strings.NewReader(string(c.JsonByte))
				}
				c.ContentType = Types[TypeJSON]
			case TypeForm, TypeFormData, TypeUrlencoded:
				reader = strings.NewReader(c.FormString)
				c.ContentType = Types[TypeForm]
			case TypeXML:
				reader = strings.NewReader(c.FormString)
				c.ContentType = Types[TypeXML]
				c.UnmarshalType = TypeXML
			default:
				return nil, errors.New("Request type Error ")
			}
		default:
			return nil, errors.New("Only support Get and Post ")
		}

		req, err := http.NewRequest(c.Method, c.Url, reader)
		if err != nil {
			return nil, err
		}
		req.Header = c.Header
		req.Header.Set("Content-Type", c.ContentType)
		c.HttpClient.Transport = c.Transport
		return req, nil
	}()
	if err != nil {
		c.Errors = append(c.Errors, err)
		return nil, nil, c.Errors
	}
	if c.Timeout != time.Duration(0) {
		c.HttpClient.Timeout = c.Timeout
	}
	res, err = c.HttpClient.Do(req)
	if err != nil {
		c.Errors = append(c.Errors, err)
		return nil, nil, c.Errors
	}
	defer res.Body.Close()
	bs, err = ioutil.ReadAll(io.LimitReader(res.Body, int64(3<<20))) // default 3MB change the size you want
	if err != nil {
		c.Errors = append(c.Errors, err)
		return nil, nil, c.Errors
	}
	return res, bs, nil
}
