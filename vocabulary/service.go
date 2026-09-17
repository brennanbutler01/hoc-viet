package vocabulary

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

var ErrInvalidWord = errors.New("invalid word data")

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
	request.Word = strings.TrimSpace(request.Word)
	request.Translation = strings.TrimSpace(request.Translation)
	if request.Word == "" || utf8.RuneCountInString(request.Word) > 200 {
		return nil, fmt.Errorf("%w: word must contain 1 to 200 characters", ErrInvalidWord)
	}
	if request.Translation == "" || utf8.RuneCountInString(request.Translation) > 500 {
		return nil, fmt.Errorf("%w: translation must contain 1 to 500 characters", ErrInvalidWord)
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
