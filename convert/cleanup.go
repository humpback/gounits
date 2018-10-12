package convert

import "fmt"

func CleanupInterfaceArray(in []interface{}) []interface{} {

	res := make([]interface{}, len(in))
	for i, v := range in {
		res[i] = CleanupMapValue(v)
	}
	return res
}

func CleanupInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {

	res := make(map[string]interface{})
	for k, v := range in {
		res[fmt.Sprintf("%v", k)] = CleanupMapValue(v)
	}
	return res
}

func CleanupMapValue(v interface{}) interface{} {

	switch v := v.(type) {
	case []interface{}:
		return CleanupInterfaceArray(v)
	case map[interface{}]interface{}:
		return CleanupInterfaceMap(v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}
