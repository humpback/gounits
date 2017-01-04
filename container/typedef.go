/*
* (C) 2001-2017 humpback Inc.
*
* gounits source code
* version: 1.0.0
* author: bobliu0909@gmail.com
* datetime: 2015-10-14
*
 */

package container

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"
)

var DEFAULT_STACK_MAXCOUNT = 1024
var DEFAULT_QUEUE_MAXCOUNT = 1024

var (
	ERR_STACK_FULL  = errors.New("stack is full.")
	ERR_STACK_EMPTY = errors.New("stack is empty.")
	ERR_QUEUE_FULL  = errors.New("queue is full.")
	ERR_QUEUE_EMPTY = errors.New("queue is empty.")
)

func Contains(obj interface{}, target interface{}) bool {

	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}

func ConvertJSONMapperArray(m map[string]interface{}, key string) []interface{} {

	if _, ret := m[key]; ret {
		return m[key].([]interface{})
	}
	return nil
}

func ReadJSONMapperString(m map[string]interface{}, key string) string {

	if _, ret := m[key]; ret {
		return m[key].(string)
	}
	return ""
}

func ReadJSONMapperInt(m map[string]interface{}, key string) int {

	if _, ret := m[key]; ret {
		return (int)(m[key].(float64))
	}
	return 0
}

func ReadJSONMapperFloat64(m map[string]interface{}, key string) float64 {

	if _, ret := m[key]; ret {
		return m[key].(float64)
	}
	return 0.0000000
}

func ReadJSONMapperInt64(m map[string]interface{}, key string) int64 {

	if _, ret := m[key]; ret {
		return (int64)(m[key].(float64))
	}
	return 0
}

func ReadJSONMapperBool(m map[string]interface{}, key string) bool {

	if _, ret := m[key]; ret {
		return m[key].(bool)
	}
	return false
}

func ReadJSONMapperDateTime(m map[string]interface{}, key string) time.Time {

	v := time.Time{}
	if _, ret := m[key]; ret {
		t, err := time.ParseInLocation("2006-01-02T15:04:05.000Z", m[key].(string), time.Local)
		if err != nil {
			return time.Time{}
		}
		v = t
	}
	return v
}

func ReadJSONMapperObject(m map[string]interface{}, key string, i interface{}) error {

	if _, ret := m[key]; ret {
		return json.Unmarshal([]byte(m[key].(string)), i)
	}
	return fmt.Errorf("read key not found.")
}
