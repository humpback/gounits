package container

import (
	"runtime"
)

type TaskHandleFunc func(context interface{})

type Task struct {
	Context    interface{}
	HandleFunc TaskHandleFunc
}

func NewTask(context interface{}, handleFunc TaskHandleFunc) *Task {

	return &Task{
		Context:    context,
		HandleFunc: handleFunc,
	}
}

type LocalQueue struct {
	taskChan chan *Task
	stopCh   <-chan struct{}
}

func NewLocalQueue(goNum int, queueSize int, stopCh <-chan struct{}) *LocalQueue {

	localQueue := &LocalQueue{
		taskChan: make(chan *Task, queueSize),
		stopCh:   stopCh,
	}

	if goNum == 0 {
		goNum = 1
	}

	for i := 0; i < goNum; i++ {
		go localQueue.consumeTask()
	}
	return localQueue
}

func (localQueue *LocalQueue) Add(task *Task) {

	if task != nil {
		go func() {
			localQueue.taskChan <- task
		}()
	}
}

func (localQueue *LocalQueue) consumeTask() {

	for {
		select {
		case task := <-localQueue.taskChan:
			{
				task.HandleFunc(task.Context)
				runtime.Gosched()
			}
		case <-localQueue.stopCh:
			{
				return
			}
		}
	}
}
