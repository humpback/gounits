package convert

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

func ConvertToString(v interface{}) (string, error) {

	p := reflect.ValueOf(v)
	switch p.Interface().(type) {
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(p.Int(), 10), nil
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(p.Uint(), 10), nil
	case float32, float64:
		return fmt.Sprintf("%v", p), nil
	case []byte:
		return string(p.Bytes()), nil
	case string:
		return url.QueryEscape(p.String()), nil
	}
	return "", errors.New("can't convert, unsupport data type.")
}
