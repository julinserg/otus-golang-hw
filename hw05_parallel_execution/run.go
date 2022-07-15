package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	countTaskRun := 0
	countTaskWhithError := 0
	numWorker := n
	wg := &sync.WaitGroup{}

	if len(tasks) < n {
		numWorker = len(tasks)
	}
	chanNotifyDone := make(chan error, numWorker)

	worker := func(i int, notify chan<- error) {
		defer wg.Done()
		err := tasks[i]()
		notify <- err
	}
	wg.Add(numWorker)
	for i := 0; i < numWorker; i++ {
		go worker(i, chanNotifyDone)
		countTaskRun++
	}

	isErrorExist := false
	for err := range chanNotifyDone {
		if err != nil && m > 0 {
			countTaskWhithError++
			if countTaskWhithError == m {
				isErrorExist = true
				break
			}
		}

		if countTaskRun < len(tasks) {
			wg.Add(1)
			go worker(countTaskRun, chanNotifyDone)
			countTaskRun++
		} else {
			break
		}
	}
	wg.Wait()
	close(chanNotifyDone)
	if isErrorExist {
		return ErrErrorsLimitExceeded
	}
	return nil
}
