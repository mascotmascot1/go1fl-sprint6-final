package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/handlers"
	"github.com/go-chi/chi/v5"
)

// Server wraps the configured HTTP server and its associated logger.
type Server struct {
	HTTP   *http.Server
	Logger *log.Logger
}

// NewServer creates and returns a new Server instance.
//
// It sets up the HTTP router, registers request handlers,
// and configures server settings such as address, timeouts, and error logging.
func NewServer(logger *log.Logger) *Server {
	mux := chi.NewRouter()
	mux.Get("/", handlers.RootHandle)
	mux.Post("/upload", handlers.UploadConvertHandle(logger))

	srv := &http.Server{
		Addr:         ":8080",
		ErrorLog:     logger,
		Handler:      mux,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 15,
	}

	return &Server{
		HTTP:   srv,
		Logger: logger,
	}
}
