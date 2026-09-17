package vocabulary

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Share one repository instance per file within this process.
// File replacement is atomic; multi-process writers require a database or file locking.
type Repository struct {
	filename string
	mutex    sync.Mutex
}

func NewRepository(filename string) *Repository { return &Repository{filename: filename} }

func (r *Repository) GetAll() ([]Word, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.readWords()
}

func (r *Repository) readWords() ([]Word, error) {
	words := []Word{}
	data, err := os.ReadFile(r.filename)
	if os.IsNotExist(err) {
		return words, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read vocabulary file: %w", err)
	}
	if err := json.Unmarshal(data, &words); err != nil {
		return nil, fmt.Errorf("parse vocabulary file: %w", err)
	}
	if words == nil {
		words = []Word{}
	}
	return words, nil
}

func (r *Repository) Save(word Word) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	words, err := r.readWords()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(append(words, word), "", "  ")
	if err != nil {
		return fmt.Errorf("encode vocabulary: %w", err)
	}
	// Stage beside the destination so rename cannot cross filesystem boundaries.
	file, err := os.CreateTemp(filepath.Dir(r.filename), ".vocabulary-*.json")
	if err != nil {
		return fmt.Errorf("create staged vocabulary: %w", err)
	}
	defer os.Remove(file.Name())
	if _, err = file.Write(data); err != nil {
		file.Close()
		return fmt.Errorf("write staged vocabulary: %w", err)
	}
	if err = file.Sync(); err != nil {
		file.Close()
		return fmt.Errorf("sync staged vocabulary: %w", err)
	}
	if err = file.Close(); err != nil {
		return fmt.Errorf("close staged vocabulary: %w", err)
	}
	if err = os.Rename(file.Name(), r.filename); err != nil {
		return fmt.Errorf("replace vocabulary file: %w", err)
	}
	return nil
}
