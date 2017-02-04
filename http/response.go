package http

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type Response struct {
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

func (resp *Response) JSONMapper() (map[string]interface{}, error) {

	var mapper map[string]interface{}
	dec := json.NewDecoder(resp.body)
	for {
		if err := dec.Decode(mapper); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return mapper, nil
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

	io.CopyN(ioutil.Discard, resp.body, 512)
	return resp.body.Close()
}
