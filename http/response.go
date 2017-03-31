package http

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

const ResponseBodyAllSize int64 = 0

type Response struct {
	rawurl     string
	body       io.ReadCloser
	header     http.Header
	statuscode int
}

func (resp *Response) Bytes() ([]byte, error) {

	return ioutil.ReadAll(resp.body)
}

func (resp *Response) String() string {

	buf, err := resp.Bytes()
	if err != nil {
		return ""
	}
	return string(buf)
}

func (resp *Response) JSON(object interface{}) error {

	buf, err := resp.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, object)
}

func (resp *Response) JSONMapper(data interface{}) error {

	dec := json.NewDecoder(resp.body)
	for {
		if err := dec.Decode(data); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}

func (resp *Response) RawURL() string {

	return resp.rawurl
}

func (resp *Response) Header(key string) string {

	return resp.header.Get(key)
}

func (resp *Response) Headers() http.Header {

	return resp.header
}

func (resp *Response) StatusCode() int {

	return resp.statuscode
}

func (resp *Response) Close() error {

	io.Copy(ioutil.Discard, resp.body)
	return resp.body.Close()
}
