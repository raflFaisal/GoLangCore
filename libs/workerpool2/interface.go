// interface.go 
package workerpool2

import "context"

// WorkerLauncher is an interface for launching workers.
type WorkerLauncher interface {
 LaunchWorker(in chan Request, stopCh chan struct{})
}

// Dispatcher is an interface for managing the worker pool.
type Dispatcher interface {
 AddWorker(w WorkerLauncher)
 RemoveWorker(minWorkers int)
 LaunchWorker(id int, w WorkerLauncher)
 ScaleWorkers(minWorkers, maxWorkers, loadThreshold int)
 MakeRequest(Request)
 Stop(ctx context.Context)
}

// Dispatcher
// The dispatcher is responsible for managing the workers and distributing the incoming requests among them. It can dynamically add or remove workers based on the current load and ensures a graceful shutdown of all workers.

// AddWorker: Adds a new worker to the pool and increments the worker count. The worker is launched to start processing requests.
// RemoveWorker: Removes a worker from the pool if there are more than the minimum required workers. The worker is signaled to stop via the stopCh channel.
// ScaleWorkers: Dynamically adjusts the number of workers based on the load. If the load exceeds a threshold and there are fewer than the maximum allowed workers, a new worker is added. If the load is below the threshold and there are more than the minimum required workers, a worker is removed.
// LaunchWorker: Launches a worker and increments the worker count. This is typically used for the initial set of workers.
// MakeRequest: Adds a request to the input channel. If the channel is full, the request is dropped, and a message is logged.
// Stop: Gracefully stops all workers. It waits for all workers to finish processing their current requests. If the timeout is reached, it forcefully stops all workers.