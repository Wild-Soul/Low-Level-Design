## Objective

Implement a task scheduler that distributes tasks to workers using Golang interfaces, goroutines, and channels. Each task has a unique ID and a duration. Workers can execute tasks concurrently.

#### Task Scheduler with Golang Features :

##### Requirements
- Task Interface:
    - Define an interface Task with a method Execute() that returns the task ID, duration, and any execution error.


- Scheduler Interface:
    - Define an interface Scheduler with the following methods:
    - AddTask(task Task) error: Adds a task to the scheduler. Handles edge cases like duplicate task IDs or invalid task definitions.
    - StartExecution(ctx context.Context) error: Starts the execution of tasks using worker goroutines. Uses a channel for task distribution and ensures proper context handling for task cancellations.
    - GetResults() ([]Result, error): Retrieves the results of task executions. Handles errors gracefully if no results are available or the execution was interrupted.


- Concurrency Management:
    - Use goroutines for worker threads and channels to communicate tasks and results between the scheduler and workers. Ensure tasks are processed concurrently without data races or memory leaks.

- Memory & Context Handling:
    - Ensure efficient memory allocation and cleanup after task execution.
    - Use context.Context to handle cancellations or timeouts during task execution.

- Error Propagation:
    - Handle errors in task execution gracefully and propagate them to the caller.


##### Testing:
- Write comprehensive test cases covering:
    - Positive Cases: Valid tasks executed successfully.
    - Negative Cases: Invalid task inputs or task execution errors.
    - Edge Cases: Handling of duplicate task IDs, empty task queues, and context timeouts.
    - Memory Safety: Ensure no goroutines are leaked and channels are properly closed.

