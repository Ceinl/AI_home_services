package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Cein13/AbeTTS/server/api"
)

type Server struct{
	Port string
	Mux *http.ServeMux
	Api *api.Api 
}

func NewServer(API_KEY string, port string) *Server {
	return &Server{
		Port: port,
		Mux: http.NewServeMux(),
		Api: api.NewApi(API_KEY),
	}
}

func (s *Server) StartWithGracefulShutdown() {
	srv := &http.Server{
		Addr:    ":" + s.Port,
		Handler: s.Mux,
	}

	// Сигнали завершення
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Println("Server started on port :" + s.Port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
	log.Println("Server shut down gracefully")
}


func (s *Server) RegisterApi() {
//	s.Mux.HandleFunc("/api/summarize", s.Api.POSTSummarize)
//	s.Mux.HandleFunc("/api/systemPrompt", s.Api.POSTSystemPrompt)
}
