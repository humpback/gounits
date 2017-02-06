package container

import (
	"errors"
	"reflect"
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
