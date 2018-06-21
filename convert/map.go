package convert

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"
)

var (
	ERR_CONVERT_INVALID = errors.New("value convert invalid.")
)

func ConvertKVStringSliceToMap(values []string) map[string]string {

	ret := make(map[string]string, len(values))
	for _, value := range values {
		kv := strings.SplitN(value, "=", 2)
		if len(kv) == 1 {
			ret[kv[0]] = ""
		} else {
			ret[kv[0]] = kv[1]
		}
	}
	return ret
}

func ConvertMapToKVStringSlice(values map[string]string) []string {

	ret := make([]string, len(values))
	i := 0
	for key, value := range values {
		ret[i] = key + "=" + value
		i++
	}
	return ret
}

func ConvertMapValueToObjectSlice(m map[string]interface{}, key string) ([]interface{}, error) {

	if value, ret := m[key]; ret {
		if reflect.TypeOf(value).Kind() == reflect.Slice {
			return value.([]interface{}), nil
		}
	}
	return nil, ERR_CONVERT_INVALID
}

func ConvertMapValueToString(m map[string]interface{}, key string) (string, error) {

	if object, ret := m[key]; ret {
		if value, ok := object.(string); ok {
			return value, nil
		}
	}
	return "", ERR_CONVERT_INVALID
}

func ConvertMapValueToStringSlice(m map[string]interface{}, key string) ([]string, error) {

	if object, ret := m[key]; ret {
		if value, ok := object.([]string); ok {
			return value, nil
		}
	}
	return nil, ERR_CONVERT_INVALID
}

func ConvertMapValueToInt(m map[string]interface{}, key string) (int, error) {

	if object, ret := m[key]; ret {
		if value, ok := object.(int); ok {
			return value, nil
		}
	}
	return 0, ERR_CONVERT_INVALID
}

func ConvertMapValueToIntSlice(m map[string]interface{}, key string) ([]int, error) {

	if object, ret := m[key]; ret {
		if value, ok := object.([]int); ok {
			return value, nil
		}
	}
	return nil, ERR_CONVERT_INVALID
}

func ConvertMapValueToInt64(m map[string]interface{}, key string) (int64, error) {

	if object, ret := m[key]; ret {
		if value, ok := object.(int64); ok {
			return value, nil
		}
	}
	return 0, ERR_CONVERT_INVALID
}

func ConvertMapValueToInt64Slice(m map[string]interface{}, key string) ([]int64, error) {

	if object, ret := m[key]; ret {
		if value, ok := object.([]int64); ok {
			return value, nil
		}
	}
	return nil, ERR_CONVERT_INVALID
}

func ConvertMapValueToFloat64(m map[string]interface{}, key string) (float64, error) {

	if object, ret := m[key]; ret {
		if value, ok := object.(float64); ok {
			return value, nil
		}
	}
	return 0.000000, ERR_CONVERT_INVALID
}

func ConvertMapValueToFloat64Slice(m map[string]interface{}, key string) ([]float64, error) {

	if object, ret := m[key]; ret {
		if value, ok := object.([]float64); ok {
			return value, nil
		}
	}
	return nil, ERR_CONVERT_INVALID
}

func ConvertMapValueToBool(m map[string]interface{}, key string) (bool, error) {

	if object, ret := m[key]; ret {
		if value, ok := object.(bool); ok {
			return value, nil
		}
	}
	return false, ERR_CONVERT_INVALID
}

func ConvertMapValueToTime(m map[string]interface{}, key string) (time.Time, error) {

	if object, ret := m[key]; ret {
		if value, ok := object.(time.Time); ok {
			return value, nil
		}
	}
	return time.Time{}, ERR_CONVERT_INVALID
}

func ConvertMapValueToJSONObject(m map[string]interface{}, key string, v interface{}) error {

	if object, ret := m[key]; ret {
		if value, ok := object.(string); ok {
			return json.Unmarshal([]byte(value), v)
		}
	}
	return ERR_CONVERT_INVALID
}

func InterfaceToMap(obj interface{}) map[string]interface{} {

	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}
