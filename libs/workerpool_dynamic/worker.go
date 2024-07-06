// worker.go

package workerpool

import (
 "context"
 "fmt"
 "sync"
 "time"
)

// Reference: https://medium.com/@souravchoudhary0306/mastering-concurrent-processing-a-step-by-step-guide-to-building-a-scalable-worker-pool-in-go-54093074c612


// Worker
// The Worker struct represents a worker that processes requests. Each worker runs in its own goroutine and listens for incoming requests on a channel.

// LaunchWorker: Launches the worker in a separate goroutine. The worker processes incoming requests until the input channel is closed or it receives a stop signal.
// processRequest: Processes a single request. It retries the request up to the specified maximum retries if an error occurs or if the request times out.

// Worker represents a worker that processes requests.
type Worker struct {
 Id         int
 Wg         *sync.WaitGroup
 ReqHandler map[int]RequestHandler
}

// LaunchWorker launches the worker to process incoming requests.
// It runs in a separate goroutine, continuously listening for incoming requests on the input channel.
// The worker gracefully stops when either the input channel is closed or it receives a stop signal.
func (w *Worker) LaunchWorker(in chan Request, stopCh chan struct{}) {
 go func() {
  defer w.Wg.Done()
  for {
   select {
   case msg, open := <-in:
    if !open {
     // If the channel is closed, stop processing and return
     // if we skip close channel check then after closing channel, 
     // worker keep reading empty values from closed channel.
     fmt.Println("Stopping worker:", w.Id)
     return
    }
    w.processRequest(msg)
    time.Sleep(1 * time.Microsecond) // Small delay to prevent tight loop
   case <-stopCh:
    fmt.Println("Stopping worker:", w.Id)
    return
   }
  }
 }()
}

// processRequest processes a single request.
func (w *Worker) processRequest(msg Request) {
//  fmt.Printf("Worker %d processing request: %v\n", w.Id, msg)
 var handler RequestHandler
 var ok bool
 if handler, ok = w.ReqHandler[msg.Type]; !ok {
  fmt.Println("Handler not implemented: workerID:", w.Id)
 } else {
  if msg.Timeout == 0 {
   msg.Timeout = time.Duration(10 * time.Millisecond) // Default timeout
  }
  for attempt := 0; attempt <= msg.MaxRetries; attempt++ {
   var err error
   done := make(chan struct{})
   ctx, cancel := context.WithTimeout(context.Background(), msg.Timeout)
   defer cancel()

   go func() {
    err = handler(msg.Data)
    close(done)
   }()

   select {
   case <-done:
    if err == nil {
     return // Successfully processed
    }
    fmt.Printf("Worker %d: Error processing request: %v\n", w.Id, err)
   case <-ctx.Done():
    fmt.Printf("Worker %d: Timeout processing request: %v\n", w.Id, msg.Data)
   }
   fmt.Printf("Worker %d: Retry %d for request %v\n", w.Id, attempt, msg.Data)
  }
  fmt.Printf("Worker %d: Failed to process request %v after %d retries\n", w.Id, msg.Data, msg.MaxRetries)
 }
}
