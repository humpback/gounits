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
	"container/list"
	"reflect"
	"sync"
)

var QueryFunc func(value interface{}) bool

type ListQueue struct {
	mutex    *sync.RWMutex
	maxcount int
	list     *list.List
}

func NewListQueue(maxcount int) *ListQueue {

	if maxcount <= 0 {
		maxcount = DEFAULT_QUEUE_MAXCOUNT
	}

	return &ListQueue{
		mutex:    new(sync.RWMutex),
		list:     list.New(),
		maxcount: maxcount,
	}
}

func (queue *ListQueue) Push(value interface{}) error {

	queue.mutex.Lock()
	defer queue.mutex.Unlock()
	if queue.list.Len() < queue.maxcount {
		if e := queue.list.PushFront(value); e != nil {
			return nil
		}
	}
	return ERR_QUEUE_FULL
}

func (queue *ListQueue) Pop() (interface{}, error) {

	queue.mutex.Lock()
	defer queue.mutex.Unlock()
	if e := queue.list.Back(); e != nil {
		queue.list.Remove(e)
		return e.Value, nil
	}
	return nil, ERR_QUEUE_EMPTY
}

func (queue *ListQueue) Contains(value interface{}) bool {

	queue.mutex.RLock()
	defer queue.mutex.RUnlock()
	e := queue.list.Front()
	for e != nil {
		if e.Value == value {
			return true
		} else {
			e = e.Next()
		}
	}
	return false
}

func (queue *ListQueue) Query(queryFunc interface{}) interface{} {

	queue.mutex.RLock()
	defer queue.mutex.RUnlock()
	e := queue.list.Front()
	for e != nil {
		if reflect.TypeOf(queryFunc) == reflect.TypeOf(QueryFunc) {
			if queryFunc.(func(value interface{}) bool)(e.Value) {
				return e.Value
			}
		} else {
			return nil
		}
		e = e.Next()
	}
	return nil
}

func (queue *ListQueue) Clear() {

	queue.mutex.Lock()
	defer queue.mutex.Unlock()
	for queue.list.Len() > 0 {
		queue.list.Remove(queue.list.Back())
	}
}

func (queue *ListQueue) Len() int {

	queue.mutex.RLock()
	defer queue.mutex.RUnlock()
	return queue.list.Len()
}

func (queue *ListQueue) Empty() bool {

	queue.mutex.RLock()
	defer queue.mutex.RUnlock()
	return queue.list.Len() == 0
}
