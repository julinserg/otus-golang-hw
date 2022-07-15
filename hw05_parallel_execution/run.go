package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {

	wg := &sync.WaitGroup{}
	wgMonitor1 := &sync.WaitGroup{}
	wgMonitor2 := &sync.WaitGroup{}
	muTask := &sync.Mutex{}
	chanErrorNotify := make(chan error, n)

	worker := func() {
		defer wg.Done()
		for {
			muTask.Lock()
			if len(tasks) > 0 {
				taskForRun := tasks[0]
				tasks = tasks[1:]
				//fmt.Println("tasks size", len(tasks))
				muTask.Unlock()

				err := taskForRun()
				chanErrorNotify <- err
			} else {
				muTask.Unlock()
				return
			}
		}
	}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go worker()
	}
	wgMonitor1.Add(1)
	go func() {
		defer wgMonitor1.Done()
		wg.Wait()
		close(chanErrorNotify)
	}()

	isStopSendError := false
	wgMonitor2.Add(1)
	go func() {
		defer wgMonitor2.Done()
		countTaskWhithError := 0
		for err := range chanErrorNotify {
			if err != nil && m > 0 && !isStopSendError {
				countTaskWhithError++
				if countTaskWhithError == m {
					isStopSendError = true
					muTask.Lock()
					tasks = nil
					muTask.Unlock()
				}
			}
		}
	}()

	wgMonitor1.Wait()
	wgMonitor2.Wait()
	if isStopSendError {
		return ErrErrorsLimitExceeded
	}
	return nil
}
