package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(tasks <-chan Task, results chan<- error, quit chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-quit:
			return
		default:
		}
		select {
		case task, ok := <-tasks:
			if !ok {
				return
			}
			err := task()
			results <- err
		case <-quit:
			return
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	numTask := len(tasks)
	indexTask := 0
	isStopSendError := false
	wg := &sync.WaitGroup{}
	chanTasks := make(chan Task, numTask)
	chanError := make(chan error, numTask)
	chanQuit := make(chan interface{}, n)

	wg.Add(n)
	for i := 0; i < n; i++ {
		go worker(chanTasks, chanError, chanQuit, wg)
	}
	runTask := n
	if numTask < n {
		runTask = numTask
	}
	for i := 0; i < runTask; i++ {
		chanTasks <- tasks[i]
		indexTask++
	}

	countTaskWhithError := 0
	for i := 0; i < numTask; i++ {
		err := <-chanError
		if err != nil && m > 0 {
			countTaskWhithError++
			if countTaskWhithError == m {
				isStopSendError = true
				for j := 0; j < n; j++ {
					chanQuit <- struct{}{}
				}
				close(chanQuit)
				break
			}
		}

		if indexTask < numTask {
			chanTasks <- tasks[indexTask]
			indexTask++
		}
	}
	close(chanTasks)
	wg.Wait()
	if isStopSendError {
		return ErrErrorsLimitExceeded
	}
	return nil
}
