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
	"strings"
	"syscall"
	"time"
)

func newRouter(filename string) http.Handler {
	return newConfiguredRouter(filename, false)
}

func newPublicDemoRouter(filename string) http.Handler {
	return newConfiguredRouter(filename, true)
}

func newConfiguredRouter(filename string, publicDemo bool) http.Handler {
	router := chi.NewMux()
	title := "Học Việt API"
	if publicDemo {
		title = "Học Việt Public Demo API"
	}
	api := humachi.New(router, huma.DefaultConfig(title, "1.0.0"))
	translation.RegisterRoutes(api)
	repository := vocabulary.NewRepository(filename)
	if publicDemo {
		vocabulary.RegisterReadOnlyRoutes(api, repository)
	} else {
		vocabulary.RegisterRoutes(api, repository)
	}
	router.Get("/health", func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"status":"ok"}`))
	})
	if publicDemo {
		return newPublicDemoMiddleware(router)
	}
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
	publicDemo := strings.EqualFold(os.Getenv("PUBLIC_DEMO"), "true")
	handler := newRouter(filename)
	if publicDemo {
		handler = newPublicDemoRouter(filename)
	}
	server := &http.Server{Addr: address, Handler: handler, ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout: 15 * time.Second, WriteTimeout: 20 * time.Second, IdleTimeout: 60 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	failures := make(chan error, 1)
	go func() { failures <- server.ListenAndServe() }()
	if publicDemo {
		log.Printf("Serving Học Việt public demo API at %s", address)
	} else {
		log.Printf("Serving Học Việt local API at %s", address)
	}
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
