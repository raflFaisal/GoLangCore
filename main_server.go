// main.go
package main

import (
 "fmt"
 "time"
 "github.com/raflFaisal/GoLangCore/libs/server"
)

func main() {
 s := server.NewServer(9093)
 s.RegisterHandlers()
 s.Start()
 time.Sleep(10 * time.Second)
 // curl -X GET localhost:9093/health

 s.Stop()
 fmt.Println("Exiting main!")
}