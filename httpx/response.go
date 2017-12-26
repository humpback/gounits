package httpx

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
)

const ResponseBodyAllSize int64 = 0

type HttpResponse struct {
	rawurl     string
	body       io.ReadCloser
	header     http.Header
	status     string
	statuscode int
}

func (resp *HttpResponse) Body() io.ReadCloser {

	return resp.Body()
}

func (resp *HttpResponse) Bytes() ([]byte, error) {

	return ioutil.ReadAll(resp.body)
}

func (resp *HttpResponse) String() string {

	buf, err := resp.Bytes()
	if err != nil {
		return ""
	}
	return string(buf)
}

func (resp *HttpResponse) JSON(object interface{}) error {

	buf, err := resp.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, object)
}

func (resp *HttpResponse) XML(object interface{}) error {

	buf, err := resp.Bytes()
	if err != nil {
		return err
	}
	return xml.Unmarshal(buf, object)
}

func (resp *HttpResponse) JSONMapper(data interface{}) error {

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

func (resp *HttpResponse) XMLMapper(data interface{}) error {

	dec := xml.NewDecoder(resp.body)
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

func (resp *HttpResponse) RawURL() string {

	return resp.rawurl
}

func (resp *HttpResponse) Header(key string) string {

	return resp.header.Get(key)
}

func (resp *HttpResponse) Headers() http.Header {

	return resp.header
}

func (resp *HttpResponse) Status() string {

	return resp.status
}

func (resp *HttpResponse) StatusCode() int {

	return resp.statuscode
}

func (resp *HttpResponse) Close() error {

	io.Copy(ioutil.Discard, resp.body)
	return resp.body.Close()
}
