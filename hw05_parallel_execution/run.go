package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	chTasks := getChanTasks(tasks)
	chErrors := make(chan struct{}, len(tasks))
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	chStop := make(chan struct{}, n)

	totalErrors := 0

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

					taskHandler(task, &mu, &totalErrors, m, chStop)
				case <-chStop:
					return
				}
			}
		}()
	}

	wg.Wait()
	close(chErrors)
	close(chStop)

	if totalErrors > 0 || m == 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func getChanTasks(tasks []Task) <-chan Task {
	ch := make(chan Task, len(tasks))

	for _, task := range tasks {
		ch <- task
	}

	close(ch)

	return ch
}

func taskHandler(task Task, mu *sync.Mutex, totalErrors *int, m int, chStop chan struct{}) {
	mu.Lock()
	if *totalErrors >= m {
		mu.Unlock()
		return
	}
	mu.Unlock()

	if task() != nil {
		mu.Lock()
		*totalErrors++
		if *totalErrors >= m {
			mu.Unlock()
			chStop <- struct{}{}
			return
		}
		mu.Unlock()
	}
}
