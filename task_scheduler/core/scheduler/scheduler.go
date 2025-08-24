package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/heymarket/core/task"
)

const (
	NUM_WORKERS = 10
)

type WorkerResult struct {
	Id       int
	Duration time.Duration
}

type WorkerError struct {
	Id    int
	Error error
}

type Scheduler struct {
	// routine safe map[int]struct{}
	lookup *Lookup

	// channel that will contain all pending tasks
	workItems chan task.Task

	// channel for results
	results chan *WorkerResult

	// channel for errors
	exError chan *WorkerError

	// wait group to synchronize all go routines.
	wg *sync.WaitGroup
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		lookup:    NewLookup(),
		workItems: make(chan task.Task, NUM_WORKERS),
		results:   make(chan *WorkerResult, NUM_WORKERS),
		exError:   make(chan *WorkerError, NUM_WORKERS),
		wg:        &sync.WaitGroup{},
	}
}

func (s *Scheduler) AddTask(task task.Task) error {
	if err := s.lookup.Add(task.GetId()); err != nil {
		return err
	}

	s.workItems <- task
	return nil
}

func (s *Scheduler) work(ctx context.Context) {
	defer s.wg.Done()

	for task := range s.workItems {
		timeoutContext, cancel := context.WithTimeout(ctx, 2*time.Second)
		taskID, duration, err := task.Execute(timeoutContext)
		cancel()

		if err != nil {
			s.exError <- &WorkerError{
				Id:    taskID,
				Error: err,
			}
		} else {
			s.results <- &WorkerResult{
				Id:       taskID,
				Duration: duration,
			}
		}
	}
}

func (s *Scheduler) StartExecution(ctx context.Context) error {
	s.wg.Add(NUM_WORKERS)

	for i := 0; i < NUM_WORKERS; i++ {
		fmt.Printf("Spinning up go routine: %d\n", i)
		go s.work(ctx)
	}

	return nil
}

func (s *Scheduler) Wait() {
	s.wg.Wait()
	close(s.results)
	close(s.exError)
}

func (s *Scheduler) CloseTasks() {
	close(s.workItems)
	go s.Wait()
}

func (s *Scheduler) GetResults() ([]*WorkerResult, error) {
	fmt.Println("Getting results")
	results := make([]*WorkerResult, 0)
	for res := range s.results {
		results = append(results, res)
	}
	return results, nil
}
