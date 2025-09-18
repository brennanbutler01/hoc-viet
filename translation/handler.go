package translation

import (
	"context"
	"net/http"

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
	Word string `path:"word" maxLength:"30" example:"world" doc:"Word to translate"`
}) (*Output, error) {
	result, err := h.service.Translate(input.Word)
	if err != nil {
		return nil, huma.Error500InternalServerError("Translation failed", err)
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
