package main

import (
	"net/http"

	"hoc-viet/translation"
	"hoc-viet/vocabulary"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

func main() {
	router := chi.NewMux()
	api := humachi.New(router, huma.DefaultConfig("Translation API", "1.0.0"))

	translation.RegisterRoutes(api)
	vocabulary.RegisterRoutes(api)

	http.ListenAndServe("localhost:8888", router)
}
