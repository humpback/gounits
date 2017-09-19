package httpx

import "golang.org/x/net/proxy"

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	DefaultClient = NewClient()

	DefaultTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
		DisableKeepAlives:     false,
		MaxIdleConns:          http.DefaultTransport.(*http.Transport).MaxIdleConns,
		MaxIdleConnsPerHost:   http.DefaultMaxIdleConnsPerHost,
		IdleConnTimeout:       120 * time.Second,
		TLSHandshakeTimeout:   http.DefaultTransport.(*http.Transport).TLSHandshakeTimeout,
		ExpectContinueTimeout: http.DefaultTransport.(*http.Transport).ExpectContinueTimeout,
	}

	DefaultPool = &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 200<<10))
		},
	}
)

type basicAuth struct {
	UserName string
	Password string
}

type httpReader struct {
	io.Reader
	Headers map[string][]string
}

type HttpClient struct {
	c       *http.Client
	pool    *sync.Pool
	auth    basicAuth
	cookies []*http.Cookie
	headers map[string]string
}

func NewClient() *HttpClient {

	client := &http.Client{
		Transport: http.DefaultTransport,
	}
	return NewWithClient(client)
}

func NewWithClient(client *http.Client) *HttpClient {

	if client == nil {
		client = http.DefaultClient
	}

	return &HttpClient{
		c:       client,
		pool:    nil,
		auth:    basicAuth{},
		cookies: make([]*http.Cookie, 0),
		headers: make(map[string]string),
	}
}

func (client *HttpClient) Close() {

	if client.c != nil {
		if transport, ok := client.c.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}
	}
}

func (client *HttpClient) UsePool(pool *sync.Pool) *HttpClient {

	if pool == nil {
		pool = DefaultPool
	}
	client.pool = pool
	return client
}

func (client *HttpClient) GetTransport() *http.Transport {

	if client.c != nil {
		if transport, ok := client.c.Transport.(*http.Transport); ok {
			return transport
		}
	}
	return http.DefaultTransport.(*http.Transport)
}

func (client *HttpClient) SetTransport(transport *http.Transport) *HttpClient {

	if client.c != nil {
		if transport == nil {
			transport = DefaultTransport
		}
		client.c.Transport = transport
	}
	return client
}

func (client *HttpClient) SetBasicAuth(username string, password string) *HttpClient {

	client.auth = basicAuth{
		UserName: username,
		Password: password,
	}
	return client
}

func (client *HttpClient) SetHeader(key string, value string) *HttpClient {

	client.headers[key] = value
	return client
}

func (client *HttpClient) SetHeaders(headers map[string]string) *HttpClient {

	for key, value := range headers {
		client.headers[key] = value
	}
	return client
}

func (client *HttpClient) SetCookie(cookie *http.Cookie) *HttpClient {

	client.cookies = append(client.cookies, cookie)
	return client
}

func (client *HttpClient) SetCookies(cookies []*http.Cookie) *HttpClient {

	client.cookies = append(client.cookies, cookies...)
	return client
}

func (client *HttpClient) SetProxy(proxy *url.URL) *HttpClient {

	client.GetTransport().Proxy = http.ProxyURL(proxy)
	return client
}

func (client *HttpClient) SetSocks5(network string, addr string, auth *proxy.Auth, forward proxy.Dialer) *HttpClient {

	dialer, _ := proxy.SOCKS5(network, addr, auth, forward)
	client.GetTransport().Dial = dialer.Dial
	return client
}

func (client *HttpClient) SetTLSClientConfig(tlsConfig *tls.Config) *HttpClient {

	client.GetTransport().TLSClientConfig = tlsConfig
	return client
}

func Head(path string, query url.Values, headers map[string][]string) (*Response, error) {

	return DefaultClient.Head(path, query, headers)
}

func Options(path string, query url.Values, headers map[string][]string) (*Response, error) {

	return DefaultClient.Options(path, query, headers)
}

func Get(path string, query url.Values, headers map[string][]string) (*Response, error) {

	return DefaultClient.Get(path, query, headers)
}

func Put(path string, query url.Values, data io.Reader, headers map[string][]string) (*Response, error) {

	return DefaultClient.Put(path, query, data, headers)
}

func PutJSON(path string, query url.Values, object interface{}, headers map[string][]string) (*Response, error) {

	return DefaultClient.PutJSON(path, query, object, headers)
}

func PutXML(path string, query url.Values, object interface{}, headers map[string][]string) (*Response, error) {

	return DefaultClient.PutXML(path, query, object, headers)
}

func Post(path string, query url.Values, data io.Reader, headers map[string][]string) (*Response, error) {

	return DefaultClient.Post(path, query, data, headers)
}

func PostJSON(path string, query url.Values, object interface{}, headers map[string][]string) (*Response, error) {

	return DefaultClient.PostJSON(path, query, object, headers)
}

func PostXML(path string, query url.Values, object interface{}, headers map[string][]string) (*Response, error) {

	return DefaultClient.PostXML(path, query, object, headers)
}

func Patch(path string, query url.Values, data io.Reader, headers map[string][]string) (*Response, error) {

	return DefaultClient.Patch(path, query, data, headers)
}

func PatchJSON(path string, query url.Values, object interface{}, headers map[string][]string) (*Response, error) {

	return DefaultClient.PatchJSON(path, query, object, headers)
}

func PatchXML(path string, query url.Values, object interface{}, headers map[string][]string) (*Response, error) {

	return DefaultClient.PatchXML(path, query, object, headers)
}

func Delete(path string, query url.Values, headers map[string][]string) (*Response, error) {

	return DefaultClient.Delete(path, query, headers)
}

func (client *HttpClient) Head(path string, query url.Values, headers map[string][]string) (*Response, error) {

	return client.sendRequest(&Request{
		Method:  http.MethodHead,
		RawURL:  path,
		Query:   query,
		Data:    nil,
		Headers: headers,
	})
}

func (client *HttpClient) Options(path string, query url.Values, headers map[string][]string) (*Response, error) {

	return client.sendRequest(&Request{
		Method:  http.MethodOptions,
		RawURL:  path,
		Query:   query,
		Data:    nil,
		Headers: headers,
	})
}

func (client *HttpClient) Get(path string, query url.Values, headers map[string][]string) (*Response, error) {

	return client.sendRequest(&Request{
		Method:  http.MethodGet,
		RawURL:  path,
		Query:   query,
		Data:    nil,
		Headers: headers,
	})
}

func (client *HttpClient) Put(path string, query url.Values, data io.Reader, headers map[string][]string) (*Response, error) {

	return client.sendRequest(&Request{
		Method:  http.MethodPut,
		RawURL:  path,
		Query:   query,
		Data:    data,
		Headers: headers,
	})
}

func (client *HttpClient) PutJSON(path string, query url.Values, object interface{}, headers map[string][]string) (*Response, error) {

	data, err := client.encodeJson(object, headers)
	defer client.putBuffer(data.Reader.(*bytes.Buffer))
	if err != nil {
		return nil, err
	}

	return client.sendRequest(&Request{
		Method:  http.MethodPut,
		RawURL:  path,
		Query:   query,
		Data:    data,
		Headers: data.Headers,
	})
}

func (client *HttpClient) PutXML(path string, query url.Values, object interface{}, headers map[string][]string) (*Response, error) {

	data, err := client.encodeXml(object, headers)
	defer client.putBuffer(data.Reader.(*bytes.Buffer))
	if err != nil {
		return nil, err
	}

	return client.sendRequest(&Request{
		Method:  http.MethodPut,
		RawURL:  path,
		Query:   query,
		Data:    data,
		Headers: data.Headers,
	})
}

func (client *HttpClient) Post(path string, query url.Values, data io.Reader, headers map[string][]string) (*Response, error) {

	return client.sendRequest(&Request{
		Method:  http.MethodPost,
		RawURL:  path,
		Query:   query,
		Data:    data,
		Headers: headers,
	})
}

func (client *HttpClient) PostJSON(path string, query url.Values, object interface{}, headers map[string][]string) (*Response, error) {

	data, err := client.encodeJson(object, headers)
	defer client.putBuffer(data.Reader.(*bytes.Buffer))
	if err != nil {
		return nil, err
	}

	return client.sendRequest(&Request{
		Method:  http.MethodPost,
		RawURL:  path,
		Query:   query,
		Data:    data,
		Headers: data.Headers,
	})
}

func (client *HttpClient) PostXML(path string, query url.Values, object interface{}, headers map[string][]string) (*Response, error) {

	data, err := client.encodeXml(object, headers)
	defer client.putBuffer(data.Reader.(*bytes.Buffer))
	if err != nil {
		return nil, err
	}

	return client.sendRequest(&Request{
		Method:  http.MethodPost,
		RawURL:  path,
		Query:   query,
		Data:    data,
		Headers: data.Headers,
	})
}

func (client *HttpClient) Patch(path string, query url.Values, data io.Reader, headers map[string][]string) (*Response, error) {

	return client.sendRequest(&Request{
		Method:  http.MethodPatch,
		RawURL:  path,
		Query:   query,
		Data:    data,
		Headers: headers,
	})
}

func (client *HttpClient) PatchJSON(path string, query url.Values, object interface{}, headers map[string][]string) (*Response, error) {

	data, err := client.encodeJson(object, headers)
	defer client.putBuffer(data.Reader.(*bytes.Buffer))
	if err != nil {
		return nil, err
	}

	return client.sendRequest(&Request{
		Method:  http.MethodPatch,
		RawURL:  path,
		Query:   query,
		Data:    data,
		Headers: data.Headers,
	})
}

func (client *HttpClient) PatchXML(path string, query url.Values, object interface{}, headers map[string][]string) (*Response, error) {

	data, err := client.encodeXml(object, headers)
	defer client.putBuffer(data.Reader.(*bytes.Buffer))
	if err != nil {
		return nil, err
	}

	return client.sendRequest(&Request{
		Method:  http.MethodPatch,
		RawURL:  path,
		Query:   query,
		Data:    data,
		Headers: data.Headers,
	})
}

func (client *HttpClient) Delete(path string, query url.Values, headers map[string][]string) (*Response, error) {

	return client.sendRequest(&Request{
		Method:  http.MethodDelete,
		RawURL:  path,
		Query:   query,
		Data:    nil,
		Headers: headers,
	})
}

func (client *HttpClient) encodeJson(object interface{}, headers map[string][]string) (*httpReader, error) {

	buf := client.getBuffer()
	if err := json.NewEncoder(buf).Encode(object); err != nil {
		return nil, err
	}

	if headers == nil {
		headers = make(map[string][]string)
	}

	headers["Content-Type"] = []string{"application/json;charset=utf-8"}
	return &httpReader{
		Reader:  buf,
		Headers: headers,
	}, nil
}

func (client *HttpClient) encodeXml(object interface{}, headers map[string][]string) (*httpReader, error) {

	buf := client.getBuffer()
	if err := xml.NewEncoder(buf).Encode(object); err != nil {
		return nil, err
	}

	if headers == nil {
		headers = make(map[string][]string)
	}

	headers["Content-Type"] = []string{"application/xml;charset=utf-8"}
	return &httpReader{
		Reader:  buf,
		Headers: headers,
	}, nil
}

func (client *HttpClient) getBuffer() *bytes.Buffer {

	if client.pool == nil {
		return bytes.NewBuffer(nil)
	}

	buf := client.pool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func (client *HttpClient) putBuffer(buf *bytes.Buffer) {

	if client.pool != nil {
		client.pool.Put(buf)
	}
}
