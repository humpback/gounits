package http

import "golang.org/x/net/proxy"

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
)

type HttpClient struct {
	c                   *http.Client
	tlsconfig           *tls.Config
	tlshandshaketimeout time.Duration
	proxy               *url.URL
	auth                basicAuth
	cookies             []*http.Cookie
	headers             map[string]string
}

type basicAuth struct {
	username string
	password string
}

var DefaultClient = NewClient()

func NewClient() *HttpClient {

	client := &http.Client{
		Transport: &http.Transport{},
	}
	return NewWithClient(client)
}

func NewWithClient(client *http.Client) *HttpClient {

	return &HttpClient{
		c:                   client,
		tlshandshaketimeout: time.Second * 35,
		auth:                basicAuth{},
		cookies:             make([]*http.Cookie, 0),
		headers:             make(map[string]string),
	}
}

func NewWithTimeout(timeout time.Duration) *HttpClient {

	client := NewClient()
	transport := client.GetTransport()
	transport.MaxIdleConnsPerHost = 5
	transport.Dial = func(network, addr string) (net.Conn, error) {
		dial := net.Dialer{
			Timeout:   timeout,
			KeepAlive: 45 * time.Second,
		}
		conn, err := dial.Dial(network, addr)
		if err != nil {
			return conn, err
		}
		err = conn.SetDeadline(time.Now().Add(timeout))
		return conn, err
	}
	return client.SetTransport(transport)
}

func (client *HttpClient) Close() {

	if client.c != nil {
		if transport, ok := client.c.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}
	}
}

func (client *HttpClient) GetTransport() *http.Transport {

	if client.c != nil {
		if transport, ok := client.c.Transport.(*http.Transport); ok {
			return transport
		}
	}
	return &http.Transport{}
}

func (client *HttpClient) SetTransport(transport *http.Transport) *HttpClient {

	if client.c != nil {
		client.c.Transport = transport
	}
	return client
}

func (client *HttpClient) GetProxy() *url.URL {

	return client.proxy
}

func (client *HttpClient) SetProxy(proxy *url.URL) *HttpClient {

	transport := client.GetTransport()
	transport.Proxy = http.ProxyURL(proxy)
	client.SetTransport(transport)
	client.proxy = proxy
	return client
}

func (client *HttpClient) SOCKS5(network string, addr string, auth *proxy.Auth, forward proxy.Dialer) (*HttpClient, error) {

	dialer, err := proxy.SOCKS5(network, addr, auth, forward)
	if err != nil {
		return nil, err
	}
	transport := client.GetTransport()
	transport.Dial = dialer.Dial
	client.SetTransport(transport)
	return client, nil
}

func (client *HttpClient) SetTLSClientConfig(tlsConfig *tls.Config) *HttpClient {

	transport := client.GetTransport()
	transport.TLSClientConfig = tlsConfig
	if transport.TLSClientConfig != nil {
		transport.TLSHandshakeTimeout = client.tlshandshaketimeout
	}
	client.SetTransport(transport)
	return client
}

func (client *HttpClient) SetTLSHandshakeTimeout(timeout time.Duration) *HttpClient {

	transport := client.GetTransport()
	if transport.TLSClientConfig != nil {
		transport.TLSHandshakeTimeout = client.tlshandshaketimeout
		client.SetTransport(transport)
	}
	client.tlshandshaketimeout = timeout
	return client
}

func (client *HttpClient) SetBasicAuth(username string, password string) *HttpClient {

	client.auth.username = username
	client.auth.password = password
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
