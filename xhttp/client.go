package xhttp

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type Client struct {
	HttpClient *http.Client
	bodySize   int // body size limit(MB), default is 10MB
	err        error
}

// NewClient , default tls.Config{InsecureSkipVerify: true}
func NewClient() (client *Client) {
	client = &Client{
		HttpClient: &http.Client{
			Timeout: 60 * time.Second,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: defaultTransportDialContext(&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}),
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				DisableKeepAlives:     true,
				ForceAttemptHTTP2:     true,
			},
		},
		bodySize: 10, // default is 10MB
	}
	return client
}

func (c *Client) SetTransport(transport *http.Transport) (client *Client) {
	c.HttpClient.Transport = transport
	return c
}

func (c *Client) SetTLSConfig(tlsCfg *tls.Config) (client *Client) {
	c.HttpClient.Transport.(*http.Transport).TLSClientConfig = tlsCfg
	return c
}

func (c *Client) SetTimeout(timeout time.Duration) (client *Client) {
	c.HttpClient.Timeout = timeout
	return c
}

// set body size (MB), default is 10MB
func (c *Client) SetBodySize(sizeMB int) (client *Client) {
	c.bodySize = sizeMB
	return c
}

func (c *Client) Req(typeStr ...RequestType) *Request {
	tp := TypeJSON
	if len(typeStr) == 1 {
		tpp := typeStr[0]
		if _, ok := types[tpp]; ok {
			tp = tpp
		}
	}
	r := &Request{
		client:        c,
		Header:        make(http.Header),
		requestType:   tp,
		unmarshalType: string(tp),
	}
	r.Header.Set("Content-Type", types[tp])
	return r
}

//func (c *Client) EndStruct(ctx context.Context, r *Request, v any) (res *http.Response, err error) {
//	res, bs, err := c.EndBytes(ctx, r)
//	if err != nil {
//		return nil, err
//	}
//	if res.StatusCode != http.StatusOK {
//		return res, fmt.Errorf("StatusCode(%d) != 200", res.StatusCode)
//	}
//
//	switch r.unmarshalType {
//	case string(TypeJSON):
//		err = sonic.Unmarshal(bs, &v)
//		if err != nil {
//			return nil, fmt.Errorf("json.Unmarshal(%s, %+v)：%w", string(bs), v, err)
//		}
//		return res, nil
//	case string(TypeXML):
//		err = xml.Unmarshal(bs, &v)
//		if err != nil {
//			return nil, fmt.Errorf("xml.Unmarshal(%s, %+v)：%w", string(bs), v, err)
//		}
//		return res, nil
//	default:
//		return nil, errors.New("unmarshalType Type Wrong")
//	}
//}
//
//func (c *Client) EndBytes(ctx context.Context, r *Request) (res *http.Response, bs []byte, err error) {
//	if c.err != nil {
//		return nil, nil, c.err
//	}
//	var (
//		body io.Reader
//		bw   *multipart.Writer
//	)
//	// multipart-form-data
//	if r.requestType == TypeMultipartFormData {
//		body = &bytes.Buffer{}
//		bw = multipart.NewWriter(body.(io.Writer))
//	}
//
//	reqFunc := func() (err error) {
//		switch r.method {
//		case GET:
//			// do nothing
//		case POST, PUT, DELETE, PATCH:
//			switch r.requestType {
//			case TypeJSON:
//				if r.jsonByte != nil {
//					body = strings.NewReader(string(r.jsonByte))
//				}
//			case TypeForm, TypeFormData, TypeUrlencoded:
//				body = strings.NewReader(r.formString)
//			case TypeMultipartFormData:
//				for k, v := range r.multipartBodyMap {
//					// file 参数
//					if file, ok := v.(*bm.File); ok {
//						fw, e := bw.CreateFormFile(k, file.Name)
//						if e != nil {
//							return e
//						}
//						_, _ = fw.Write(file.Content)
//						continue
//					}
//					// text 参数
//					vs, ok2 := v.(string)
//					if ok2 {
//						_ = bw.WriteField(k, vs)
//					} else if ss := util.ConvertToString(v); ss != "" {
//						_ = bw.WriteField(k, ss)
//					}
//				}
//				_ = bw.Close()
//				r.Header.Set("Content-Type", bw.FormDataContentType())
//			case TypeXML:
//				body = strings.NewReader(r.formString)
//			default:
//				return errors.New("Request type Error ")
//			}
//		default:
//			return errors.New("Only support GET and POST and PUT and DELETE ")
//		}
//
//		// request
//		req, err := http.NewRequestWithContext(ctx, r.method, r.url, body)
//		if err != nil {
//			return err
//		}
//		req.Header = r.Header
//		res, err = c.HttpClient.Do(req)
//		if err != nil {
//			return err
//		}
//		defer res.Body.Close()
//		bs, err = io.ReadAll(io.LimitReader(res.Body, int64(c.bodySize<<20))) // default 10MB change the size you want
//		if err != nil {
//			return err
//		}
//		return nil
//	}
//
//	if err = reqFunc(); err != nil {
//		return nil, nil, err
//	}
//	return res, bs, nil
//}
