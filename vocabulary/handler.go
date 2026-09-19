package vocabulary

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(repo *Repository) *Handler {
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
		if errors.Is(err, ErrInvalidWord) {
			return nil, huma.Error400BadRequest(err.Error())
		}
		log.Printf("save vocabulary: %v", err)
		return nil, huma.Error500InternalServerError("Could not save the word.")
	}

	return &AddWordResponse{Body: *word}, nil
}

// GetWords handles GET /words
func (h *Handler) GetWords(ctx context.Context, input *struct{}) (*struct {
	Body []Word
}, error) {
	words, err := h.service.GetAllWords()
	if err != nil {
		log.Printf("read vocabulary: %v", err)
		return nil, huma.Error500InternalServerError("Could not retrieve vocabulary.")
	}

	return &struct {
		Body []Word
	}{Body: words}, nil
}

// RegisterRoutes registers all vocabulary routes
func RegisterRoutes(api huma.API, repository *Repository) {
	handler := NewHandler(repository)
	registerGetWords(api, handler)

	huma.Register(api, huma.Operation{
		OperationID: "add-word",
		Method:      http.MethodPost,
		Path:        "/words",
		Summary:     "Add a new studied word",
		Description: "Add a new word to your vocabulary collection.",
		Tags:        []string{"Vocabulary"},
	}, handler.AddWord)
}

func RegisterReadOnlyRoutes(api huma.API, repository *Repository) {
	registerGetWords(api, NewHandler(repository))
}

func registerGetWords(api huma.API, handler *Handler) {
	huma.Register(api, huma.Operation{
		OperationID: "get-words",
		Method:      http.MethodGet,
		Path:        "/words",
		Summary:     "Get all studied words",
		Description: "Retrieve all words from your vocabulary collection.",
		Tags:        []string{"Vocabulary"},
	}, handler.GetWords)
}
