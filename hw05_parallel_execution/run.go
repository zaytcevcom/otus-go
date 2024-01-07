package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n < 0 {
		return errors.New("invalid count workers")
	}

	chTasks := make(chan Task)
	chStop := make(chan struct{}, n)
	wg := sync.WaitGroup{}

	var totalErrors int32

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for {
				select {
				case task, ok := <-chTasks:
					if !ok {
						return
					}

					taskHandler(task, &totalErrors, m, chStop)
				case <-chStop:
					return
				}
			}
		}()
	}

	for _, task := range tasks {
		select {
		case chTasks <- task:
		case <-chStop:
			break
		}
	}

	close(chTasks)
	wg.Wait()
	close(chStop)

	if totalErrors > 0 || m == 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func taskHandler(task Task, totalErrors *int32, m int, chStop chan struct{}) {
	if atomic.LoadInt32(totalErrors) >= int32(m) {
		return
	}

	if task() != nil {
		atomic.AddInt32(totalErrors, 1)

		if atomic.LoadInt32(totalErrors) >= int32(m) {
			chStop <- struct{}{}
		}
	}
}
