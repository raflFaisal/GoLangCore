package main

import (
    "log"
    "time"

    "github.com/raflFaisal/GoLangCore/libs/workerpool"
)

func main() {
    // 3 workers and 10 taskQueueSize
    wp := workerpool.NewWorkerPool(3, 20)
    wp.Run()

    wp.AddTask(func() {
        log.Println("Task 1 is running")
        time.Sleep(2 * time.Second)
        log.Println("Task 1 is done")
    })

    wp.AddTask(func() {
        log.Println("Task 2 is running")
        time.Sleep(1 * time.Second)
        log.Println("Task 2 is done")
    })

    // Add more tasks as needed
    wp.Stop() // Gracefully shutdown the worker pool
}