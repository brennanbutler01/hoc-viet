package vocabulary

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Repository struct {
	filename string
	mutex    *sync.Mutex
}

func NewRepository(filename string) *Repository {
	return &Repository{
		filename: filename,
		mutex:    &sync.Mutex{},
	}
}

func (r *Repository) GetAll() ([]Word, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var words []Word
	data, err := os.ReadFile(r.filename)

	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read words file: %w", err)
	}

	if err == nil {
		if err := json.Unmarshal(data, &words); err != nil {
			return nil, fmt.Errorf("failed to parse words file: %w", err)
		}
	}

	return words, nil
}

func (r *Repository) Save(word Word) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	words, err := r.getAllUnsafe()
	if err != nil {
		return err
	}

	words = append(words, word)

	data, err := json.MarshalIndent(words, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal words: %w", err)
	}

	if err := os.WriteFile(r.filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write words file: %w", err)
	}

	return nil
}

// getAllUnsafe is a helper method that doesn't lock (assumes caller has lock)
func (r *Repository) getAllUnsafe() ([]Word, error) {
	var words []Word
	data, err := os.ReadFile(r.filename)

	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read words file: %w", err)
	}

	if err == nil {
		if err := json.Unmarshal(data, &words); err != nil {
			return nil, fmt.Errorf("failed to parse words file: %w", err)
		}
	}

	return words, nil
}
