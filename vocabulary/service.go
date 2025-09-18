package vocabulary

import (
	"fmt"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// AddWord adds a new word to the vocabulary
func (s *Service) AddWord(request AddWordRequest) (*Word, error) {
	if request.Word == "" {
		return nil, fmt.Errorf("word cannot be empty")
	}
	if request.Translation == "" {
		return nil, fmt.Errorf("translation cannot be empty")
	}

	word := Word{
		Word:        request.Word,
		Translation: request.Translation,
		DateCreated: time.Now(),
	}

	if err := s.repo.Save(word); err != nil {
		return nil, fmt.Errorf("failed to save word: %w", err)
	}

	return &word, nil
}

// GetAllWords retrieves all words from vocabulary
func (s *Service) GetAllWords() ([]Word, error) {
	return s.repo.GetAll()
}