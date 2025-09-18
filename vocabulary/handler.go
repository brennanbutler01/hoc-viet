package vocabulary

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	service *Service
}

func NewHandler() *Handler {
	repo := NewRepository("words.json")
	service := NewService(repo)
	
	return &Handler{
		service: service,
	}
}

// AddWord handles POST /words
func (h *Handler) AddWord(ctx context.Context, input *struct {
	Body AddWordRequest
}) (*AddWordResponse, error) {
	word, err := h.service.AddWord(input.Body)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid word data", err)
	}

	return &AddWordResponse{Body: *word}, nil
}

// GetWords handles GET /words
func (h *Handler) GetWords(ctx context.Context, input *struct{}) (*struct {
	Body []Word
}, error) {
	words, err := h.service.GetAllWords()
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to retrieve words", err)
	}

	return &struct {
		Body []Word
	}{Body: words}, nil
}

// RegisterRoutes registers all vocabulary routes
func RegisterRoutes(api huma.API) {
	handler := NewHandler()

	huma.Register(api, huma.Operation{
		OperationID: "add-word",
		Method:      http.MethodPost,
		Path:        "/words",
		Summary:     "Add a new studied word",
		Description: "Add a new word to your vocabulary collection.",
		Tags:        []string{"Vocabulary"},
	}, handler.AddWord)

	huma.Register(api, huma.Operation{
		OperationID: "get-words",
		Method:      http.MethodGet,
		Path:        "/words",
		Summary:     "Get all studied words",
		Description: "Retrieve all words from your vocabulary collection.",
		Tags:        []string{"Vocabulary"},
	}, handler.GetWords)
}