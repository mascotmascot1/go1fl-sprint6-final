package main

import (
	"log"
	"os"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "[SERVER] ", log.LstdFlags)

	s := server.NewServer(logger)

	logger.Printf("Starting server on %s\n", s.HTTP.Addr)

	if err := s.HTTP.ListenAndServe(); err != nil {
		logger.Fatalf("Error starting server: %s", err)
	}
}
