package json

import (
	"bytes"
	"encoding/json"
)

func EnCodeObjectToBuffer(v interface{}) ([]byte, error) {

	b := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(b).Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func DeCodeBufferToObject(buf []byte, v interface{}) error {

	r := bytes.NewReader(buf)
	return json.NewDecoder(r).Decode(v)
}
