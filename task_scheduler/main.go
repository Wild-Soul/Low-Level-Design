package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/heymarket/core/scheduler"
)

type SleepTask struct {
	Id int
}

func (st *SleepTask) Execute(ctx context.Context) (int, time.Duration, error) {
	timeNow := time.Now()
	fmt.Printf("[Task %d]: started execution\n", st.Id)

	select {
	case <-ctx.Done():
		fmt.Printf("[Task %d]: cancelled execution\n", st.Id)
		return 0, time.Since(timeNow), fmt.Errorf("[Task %d]: Cancelled execution", st.Id)
	default:
		fmt.Printf("[Task %d]: sleeping for %v second\n", st.Id, (1 * time.Second).Seconds())
		time.Sleep(1 * time.Second)
	}

	return st.Id, time.Since(timeNow), nil
}

func (st *SleepTask) GetId() int {
	return st.Id
}

func main() {
	fmt.Println("Starting main")
	scheduler := scheduler.NewScheduler()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for i := range 100 {
			err := scheduler.AddTask(&SleepTask{
				Id: i,
			})
			if err != nil {
				fmt.Println("Failed to add task:", err)
			}
		}

		// Signal that no more tasks will be added.
		scheduler.CloseTasks()
	}()

	err := scheduler.StartExecution(context.Background())
	if err != nil {
		log.Fatal("Failed to start scheduler:", err)
	}

	go func() {
		defer wg.Done()
		results, err := scheduler.GetResults()
		if err != nil {
			log.Fatal("Failed to start scheduler:", err)
		}

		for _, result := range results {
			fmt.Printf("[Task %d]: took: %v\n", result.Id, result.Duration.Seconds())
		}
	}()

	wg.Wait()
	fmt.Println("Execution Done!!")
}
