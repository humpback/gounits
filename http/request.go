package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

func Head(path string, query url.Values, headers map[string][]string) (*Response, error) {

	return DefaultClient.Head(path, query, headers)
}

func Options(path string, query url.Values, headers map[string][]string) (*Response, error) {

	return DefaultClient.Options(path, query, headers)
}

func Get(path string, query url.Values, headers map[string][]string) (*Response, error) {

	return DefaultClient.Get(path, query, headers)
}

func Post(path string, query url.Values, body io.Reader, headers map[string][]string) (*Response, error) {

	return DefaultClient.Post(path, query, body, headers)
}

func Put(path string, query url.Values, body io.Reader, headers map[string][]string) (*Response, error) {

	return DefaultClient.Put(path, query, body, headers)
}

func Patch(path string, query url.Values, body io.Reader, headers map[string][]string) (*Response, error) {

	return DefaultClient.Patch(path, query, body, headers)
}

func Delete(path string, query url.Values, headers map[string][]string) (*Response, error) {

	return DefaultClient.Delete(path, query, headers)
}

func (client *HttpClient) Head(path string, query url.Values, headers map[string][]string) (*Response, error) {

	return client.doSendRequest("HEAD", path, query, nil, headers)
}

func (client *HttpClient) Options(path string, query url.Values, headers map[string][]string) (*Response, error) {

	return client.doSendRequest("OPTIONS", path, query, nil, headers)
}

func (client *HttpClient) Get(path string, query url.Values, headers map[string][]string) (*Response, error) {

	return client.doSendRequest("GET", path, query, nil, headers)
}

func (client *HttpClient) Post(path string, query url.Values, body io.Reader, headers map[string][]string) (*Response, error) {

	return client.doSendRequest("POST", path, query, body, headers)
}

func (client *HttpClient) PostJSON(path string, query url.Values, object interface{}, headers map[string][]string) (*Response, error) {

	body, err := makeJsonBody(object, headers)
	if err != nil {
		return nil, err
	}
	return client.doSendRequest("POST", path, query, body, headers)
}

func (client *HttpClient) Put(path string, query url.Values, body io.Reader, headers map[string][]string) (*Response, error) {

	return client.doSendRequest("PUT", path, query, body, headers)
}

func (client *HttpClient) PutJSON(path string, query url.Values, object interface{}, headers map[string][]string) (*Response, error) {

	body, err := makeJsonBody(object, headers)
	if err != nil {
		return nil, err
	}
	return client.doSendRequest("PUT", path, query, body, headers)
}

func (client *HttpClient) Patch(path string, query url.Values, body io.Reader, headers map[string][]string) (*Response, error) {

	return client.doSendRequest("PATCH", path, query, body, headers)
}

func (client *HttpClient) PatchJSON(path string, query url.Values, object interface{}, headers map[string][]string) (*Response, error) {

	body, err := makeJsonBody(object, headers)
	if err != nil {
		return nil, err
	}
	return client.doSendRequest("PATCH", path, query, body, headers)
}

func (client *HttpClient) Delete(path string, query url.Values, headers map[string][]string) (*Response, error) {

	return client.doSendRequest("DELETE", path, query, nil, headers)
}

func (client *HttpClient) doSendRequest(method string, path string, query url.Values, body io.Reader, headers map[string][]string) (*Response, error) {

	resp := &Response{
		body:       nil,
		statuscode: -1,
	}

	ispayload := (method == "POST" || method == "PUT" || method == "PATCH")
	if ispayload && body == nil {
		body = bytes.NewReader([]byte{})
	}

	request, err := client.newRequest(method, path, query, body, headers)
	if err != nil {
		return nil, err
	}

	if ispayload && request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "text/plain")
	}

	response, err := client.c.Do(request)
	if err != nil {
		return nil, err
	}

	resp.body = response.Body
	resp.header = response.Header
	resp.statuscode = response.StatusCode
	return resp, nil
}

func (client *HttpClient) newRequest(method, path string, query url.Values, body io.Reader, headers map[string][]string) (*http.Request, error) {

	rawurl := &url.URL{
		Path: path,
	}
	if len(query) > 0 {
		rawurl.RawQuery = query.Encode()
	}

	request, err := http.NewRequest(method, rawurl.String(), body)
	if err != nil {
		return nil, err
	}

	for key, value := range client.headers {
		request.Header.Set(key, value)
	}

	if headers != nil {
		for key, value := range headers {
			request.Header[key] = value
		}
	}

	if len(client.cookies) != 0 {
		for _, cookie := range client.cookies {
			request.AddCookie(cookie)
		}
	}

	if client.auth.username != "" && client.auth.password != "" {
		request.SetBasicAuth(client.auth.username, client.auth.password)
	}
	return request, nil
}

func makeJsonBody(object interface{}, headers map[string][]string) (io.Reader, error) {

	var body io.Reader
	if object != nil {
		var err error
		body, err = encodeJsonData(object)
		if err != nil {
			return nil, err
		}
		if headers == nil {
			headers = make(map[string][]string)
		}
		headers["Content-Type"] = []string{"application/json;charset=utf-8"}
	}
	return body, nil
}

func encodeJsonData(object interface{}) (*bytes.Buffer, error) {

	buffer := bytes.NewBuffer(nil)
	if object != nil {
		if err := json.NewEncoder(buffer).Encode(object); err != nil {
			return nil, err
		}
	}
	return buffer, nil
}
