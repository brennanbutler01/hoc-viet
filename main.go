package main

import (
	"context"
	"errors"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"github.com/go-chi/chi/v5"
	"hoc-viet/translation"
	"hoc-viet/vocabulary"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func newRouter(filename string) http.Handler {
	router := chi.NewMux()
	api := humachi.New(router, huma.DefaultConfig("Translation API", "1.0.0"))
	translation.RegisterRoutes(api)
	vocabulary.RegisterRoutes(api, vocabulary.NewRepository(filename))
	return router
}

func run() error {
	address := os.Getenv("LISTEN_ADDR")
	if address == "" {
		address = "127.0.0.1:8888"
	}
	filename := os.Getenv("VOCABULARY_FILE")
	if filename == "" {
		filename = "words.json"
	}
	server := &http.Server{Addr: address, Handler: newRouter(filename), ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout: 15 * time.Second, WriteTimeout: 20 * time.Second, IdleTimeout: 60 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	failures := make(chan error, 1)
	go func() { failures <- server.ListenAndServe() }()
	log.Printf("Serving local vocabulary API at %s", address)
	select {
	case err := <-failures:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		return server.Shutdown(shutdown)
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
