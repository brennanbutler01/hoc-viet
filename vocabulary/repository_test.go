package vocabulary

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestRepositoryEmptyAndPersistentVocabulary(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "words.json")
	repo := NewRepository(filename)
	words, err := repo.GetAll()
	if err != nil || words == nil || len(words) != 0 {
		t.Fatalf("expected empty array: %v %v", words, err)
	}
	if err := repo.Save(Word{Word: "hello", Translation: "xin chào"}); err != nil {
		t.Fatal(err)
	}
	words, err = NewRepository(filename).GetAll()
	if err != nil || len(words) != 1 || words[0].Translation != "xin chào" {
		t.Fatalf("vocabulary did not persist: %v %v", words, err)
	}
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("unexpected file permissions: %o", info.Mode().Perm())
	}
}

func TestRepositoryConcurrentSavesRetainEveryWord(t *testing.T) {
	directory := t.TempDir()
	repo := NewRepository(filepath.Join(directory, "words.json"))
	var pending sync.WaitGroup
	for i := 0; i < 40; i++ {
		pending.Add(1)
		go func(index int) {
			defer pending.Done()
			if err := repo.Save(Word{Word: fmt.Sprintf("word-%d", index), Translation: "từ"}); err != nil {
				t.Error(err)
			}
		}(i)
	}
	pending.Wait()
	words, err := repo.GetAll()
	if err != nil || len(words) != 40 {
		t.Fatalf("lost writes: count=%d error=%v", len(words), err)
	}
	files, err := os.ReadDir(directory)
	if err != nil || len(files) != 1 {
		t.Fatalf("staged files were not cleaned up: %v %v", files, err)
	}
}

func TestRepositoryPreservesMalformedExistingFile(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "words.json")
	original := []byte("damaged existing vocabulary")
	if err := os.WriteFile(filename, original, 0600); err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(filename)
	if err := repo.Save(Word{Word: "hello", Translation: "xin chào"}); err == nil {
		t.Fatal("should refuse to overwrite invalid data")
	}
	data, err := os.ReadFile(filename)
	if err != nil || string(data) != string(original) {
		t.Fatal("existing data changed after failed save")
	}
}
