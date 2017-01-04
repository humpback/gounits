package container

import (
	"sync"
)

type ArrayStack struct {
	mutex *sync.RWMutex
	items []interface{}
}

func NewArrayStack(maxcount int) *ArrayStack {

	if maxcount <= 0 {
		maxcount = DEFAULT_STACK_MAXCOUNT
	}

	return &ArrayStack{
		mutex: new(sync.RWMutex),
		items: make([]interface{}, 0, maxcount),
	}
}

func (stack *ArrayStack) Push(value interface{}) error {

	stack.mutex.Lock()
	defer stack.mutex.Unlock()
	if len(stack.items) < cap(stack.items) {
		stack.items = append(stack.items, value)
		return nil
	}
	return ERR_STACK_FULL
}

func (stack *ArrayStack) Pop() (interface{}, error) {

	stack.mutex.Lock()
	defer stack.mutex.Unlock()
	size := len(stack.items)
	if size == 0 {
		return nil, ERR_STACK_EMPTY
	}
	pop := stack.items[size-1]
	stack.items = stack.items[:size-1]
	return pop, nil
}

func (stack *ArrayStack) Top() (interface{}, error) {

	stack.mutex.RLock()
	defer stack.mutex.RUnlock()
	size := len(stack.items)
	if size == 0 {
		return nil, ERR_STACK_EMPTY
	}
	return stack.items[size-1], nil
}

func (stack *ArrayStack) Clear() {

	stack.mutex.Lock()
	defer stack.mutex.Unlock()
	stack.items = stack.items[0:0]
}

func (stack *ArrayStack) Len() int {

	stack.mutex.RLock()
	defer stack.mutex.RUnlock()
	return len(stack.items)
}

func (stack *ArrayStack) Empty() bool {

	stack.mutex.RLock()
	defer stack.mutex.RUnlock()
	return len(stack.items) == 0
}
