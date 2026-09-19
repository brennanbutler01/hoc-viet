package translation

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	service *Service
}

func NewHandler() *Handler {
	return &Handler{
		service: NewService(),
	}
}

func (h *Handler) GetTranslation(ctx context.Context, input *struct {
	Word string `path:"word" minLength:"1" maxLength:"30" example:"world" doc:"Word to translate"`
}) (*Output, error) {
	if strings.TrimSpace(input.Word) == "" {
		return nil, huma.Error400BadRequest("A word is required.")
	}
	result, err := h.service.Translate(ctx, input.Word)
	if err != nil {
		// Omit request URLs so submitted text is not written to runtime logs.
		var requestError *url.Error
		if errors.As(err, &requestError) {
			log.Printf("translation provider transport: %v", requestError.Err)
		} else {
			log.Printf("translation provider: %v", err)
		}
		return nil, huma.Error502BadGateway("Translation provider is unavailable. Please try again later.")
	}

	return &Output{Body: result}, nil
}

func RegisterRoutes(api huma.API) {
	handler := NewHandler()

	huma.Register(api, huma.Operation{
		OperationID: "get-translation",
		Method:      http.MethodGet,
		Path:        "/translation/{word}",
		Summary:     "Get a translation",
		Description: "Get a translation for a word from English to Vietnamese.",
		Tags:        []string{"Translations"},
	}, handler.GetTranslation)
}
