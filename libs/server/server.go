package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
    "log"
)

var (
	DefaultServer Server
)

// func init() {
// 	DefaultServer = NewServer(env.Port)
// }

type Server struct {
	Port   int
	server *http.Server
}

type HealthcheckResponse struct {
	Status string `json:"status"`
}

func NewServer(p int) Server {
	log.Printf("Setting-up HTTP server")
	return Server{
		Port:   p,
	}
}

func healthcheckHandler(w http.ResponseWriter, req *http.Request) {
	response := HealthcheckResponse{Status: "OK"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) RegisterHandlers() {
	http.HandleFunc("/health", healthcheckHandler)
}

func (s *Server) Start() {
	addr := fmt.Sprintf(":%d", s.Port)
	s.server = &http.Server{
		Addr:    addr,
		Handler: http.DefaultServeMux,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Could not listen on address", err.Error())
		}
	}()

	log.Printf("Http server started")
}

func (s *Server) Stop() {
	// wait for any requests beign served
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown failed", err.Error())
	}
}
