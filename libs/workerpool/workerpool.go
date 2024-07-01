package workerpool

import (
		"log"
		"sync"
)

// WorkerPool defines the contract for a worker pool implementation.
type WorkerPool interface {
    Run()
    AddTask(task func())
    Stop()
}

// workerPool implements the WorkerPool interface.
type workerPool struct {
    maxWorker   int
    queuedTaskC chan func()
    wg          sync.WaitGroup
}

// NewWorkerPool creates an instance of WorkerPool.
func NewWorkerPool(maxWorker int, taskQueueSize int) WorkerPool {
    return &workerPool{
        maxWorker:   maxWorker,
        queuedTaskC: make(chan func(), taskQueueSize),
    }
}

// Run starts the worker pool and spawns workers.
func (wp *workerPool) Run() {
    for i := 0; i < wp.maxWorker; i++ {
        wID := i + 1
        log.Printf("[WorkerPool] Worker %d has been spawned", wID)

        go func(workerID int) {
            for task := range wp.queuedTaskC {
                log.Printf("[WorkerPool] Worker %d start processing task", workerID)
                task()
                log.Printf("[WorkerPool] Worker %d finish processing task", workerID)
                wp.wg.Done()
            }
        }(wID)
    }
}

// AddTask adds a task to the worker pool.
func (wp *workerPool) AddTask(task func()) {
    wp.wg.Add(1)
    wp.queuedTaskC <- task
}

// Stop gracefully shuts down the worker pool, waiting for all tasks to complete.
func (wp *workerPool) Stop() {
    wp.wg.Wait()
    close(wp.queuedTaskC)
}

// GetTotalQueuedTask returns the number of queued tasks (if you really need this method).
func (wp *workerPool) GetTotalQueuedTask() int {
    return len(wp.queuedTaskC)
}