package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type tcpFunc func(*net.TCPConn, time.Duration) error

type HttpClient struct {
	client *http.Client
}

func NewHTTPClientTimeout(rawurl string, tlsConfig *tls.Config, timeout time.Duration, setUserTimeout tcpFunc) (*HttpClient, *url.URL, error) {

	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, nil, err
	}

	if u.Scheme == "" || u.Scheme == "tcp" {
		if tlsConfig == nil {
			u.Scheme = "http"
		} else {
			u.Scheme = "https"
		}
	}

	client := newHTTPClient(u, tlsConfig, timeout, setUserTimeout)
	return &HttpClient{
		client: client,
	}, u, nil
}

func NewHTTPClient(client *http.Client) *HttpClient {

	if client == nil {
		client = http.DefaultClient
		timeout, _ := time.ParseDuration("30s")
		client.Transport = &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(timeout)
				c, err := net.DialTimeout(netw, addr, timeout)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		}
	}
	return &HttpClient{
		client: client,
	}
}

func newHTTPClient(u *url.URL, tlsConfig *tls.Config, timeout time.Duration, setUserTimeout tcpFunc) *http.Client {

	httpTransport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	switch u.Scheme {
	default: //tcp
		httpTransport.Dial = func(proto, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(proto, addr, timeout)
			if tcpConn, ok := conn.(*net.TCPConn); ok && setUserTimeout != nil {
				setUserTimeout(tcpConn, timeout)
			}
			return conn, err
		}
	case "unix":
		socketPath := u.Path
		unixDial := func(proto, addr string) (net.Conn, error) {
			return net.DialTimeout("unix", socketPath, timeout)
		}
		httpTransport.Dial = unixDial
		u.Scheme = "http"
		u.Host = "unix.sock"
		u.Path = ""
	}
	return &http.Client{Transport: httpTransport}
}

func (c *HttpClient) PostJSON(url string, request interface{}, response interface{}) error {

	if request == nil {
		return ErrHttpRequestInvalid
	}

	buffer := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buffer).Encode(request); err != nil {
		return err
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewReader(buffer.Bytes()))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", resp.Status)
	}

	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return err
		}
	}
	return nil
}

func (c *HttpClient) GetJSON(url string, response interface{}) error {

	if response == nil {
		return ErrHttpResponseInvalid
	}

	resp, err := c.client.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return err
	}
	return nil
}

func (c *HttpClient) PostJSONMapper(url string, request map[string]interface{}, response *map[string]interface{}) error {

	if request == nil {
		return ErrHttpRequestInvalid
	}

	buffer := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buffer).Encode(request); err != nil {
		return err
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewReader(buffer.Bytes()))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", resp.Status)
	}

	if response != nil {
		dec := json.NewDecoder(resp.Body)
		for {
			if err := dec.Decode(response); err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
		}
	}
	return nil
}

func (c *HttpClient) GetJSONMapper(url string, response *map[string]interface{}) error {

	if response == nil {
		return ErrHttpResponseInvalid
	}

	resp, err := c.client.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", resp.Status)
	}

	dec := json.NewDecoder(resp.Body)
	for {
		if err := dec.Decode(response); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}

func (c *HttpClient) PostFile(url string, buf []byte) error {

	//预留方法，后期完善
	return nil
}

func (c *HttpClient) GetFile(fd *os.File, url string, resumable bool) error {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ErrHttpNewRequest
	}

	if resumable { //续传
		st, err := fd.Stat()
		if err != nil {
			return ErrFileStatException
		}
		seek := st.Size()
		fd.Seek(seek, 0)
		req.Header.Add("Range", "bytes="+strconv.FormatInt(seek, 10)+"-")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return ErrHttpRequestFailed
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		if _, err := io.Copy(fd, resp.Body); err != nil {
			return ErrHttpIOCopyFailed
		}
		return nil
	}
	return fmt.Errorf("download file http response error: %s", resp.Status)
}
