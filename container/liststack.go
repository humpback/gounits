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
	"sync"
)

type ListStack struct {
	mutex    *sync.RWMutex
	maxcount int
	list     *list.List
}

func NewListStack(maxcount int) *ListStack {

	if maxcount <= 0 {
		maxcount = DEFAULT_STACK_MAXCOUNT
	}

	return &ListStack{
		mutex:    new(sync.RWMutex),
		list:     list.New(),
		maxcount: maxcount,
	}
}

func (stack *ListStack) Push(value interface{}) error {

	stack.mutex.Lock()
	defer stack.mutex.Unlock()
	if stack.list.Len() < stack.maxcount {
		if e := stack.list.PushBack(value); e != nil {
			return nil
		}
	}
	return ERR_STACK_FULL
}

func (stack *ListStack) Pop() (interface{}, error) {

	stack.mutex.Lock()
	defer stack.mutex.Unlock()
	if e := stack.list.Back(); e != nil {
		stack.list.Remove(e)
		return e.Value, nil
	}
	return nil, ERR_STACK_EMPTY
}

func (stack *ListStack) Top() (interface{}, error) {

	stack.mutex.RLock()
	defer stack.mutex.RUnlock()
	if e := stack.list.Back(); e != nil {
		return e.Value, nil
	}
	return nil, ERR_STACK_EMPTY
}

func (stack *ListStack) Clear() {

	stack.mutex.Lock()
	defer stack.mutex.Unlock()
	for stack.list.Len() > 0 {
		stack.list.Remove(stack.list.Back())
	}
}

func (stack *ListStack) Len() int {

	stack.mutex.RLock()
	defer stack.mutex.RUnlock()
	return stack.list.Len()
}

func (stack *ListStack) Empty() bool {

	stack.mutex.RLock()
	defer stack.mutex.RUnlock()
	return stack.list.Len() == 0
}
