package main

import (
	"bytes"
	"encoding/json"
	"hoc-viet/vocabulary"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVocabularyHTTPWorkflow(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "words.json")
	handler := newRouter(filename)
	request := func(method, path, body string) *httptest.ResponseRecorder {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		handler.ServeHTTP(recorder, req)
		return recorder
	}
	empty := request("GET", "/words", "")
	if empty.Code != http.StatusOK || strings.TrimSpace(empty.Body.String()) != "[]" {
		t.Fatalf("unexpected empty response: %d %s", empty.Code, empty.Body.String())
	}
	saved := request("POST", "/words", `{"word":" hello ","translation":" xin chào "}`)
	if saved.Code != http.StatusOK {
		t.Fatalf("save failed: %d %s", saved.Code, saved.Body.String())
	}
	var word vocabulary.Word
	if err := json.Unmarshal(saved.Body.Bytes(), &word); err != nil {
		t.Fatal(err)
	}
	if word.Word != "hello" || word.Translation != "xin chào" || word.DateCreated.IsZero() {
		t.Fatalf("unexpected saved word: %+v", word)
	}
	listed := request("GET", "/words", "")
	var words []vocabulary.Word
	if err := json.Unmarshal(listed.Body.Bytes(), &words); err != nil || len(words) != 1 {
		t.Fatalf("list failed: %v %v", words, err)
	}
	for _, body := range []string{`{"word":" ","translation":"xin chào"}`, `{"word":"hello","translation":" "}`} {
		if response := request("POST", "/words", body); response.Code != 400 {
			t.Fatalf("expected invalid input rejection: %d", response.Code)
		}
	}
	if err := os.WriteFile(filename, []byte("invalid JSON"), 0600); err != nil {
		t.Fatal(err)
	}
	failed := request("POST", "/words", `{"word":"hello","translation":"xin chào"}`)
	if failed.Code != 500 {
		t.Fatalf("storage error should be 500, got %d", failed.Code)
	}
	if strings.Contains(failed.Body.String(), filename) {
		t.Fatal("response leaked internal storage path")
	}
	if response := request("GET", "/docs", ""); response.Code != 200 {
		t.Fatalf("docs unavailable: %d", response.Code)
	}
}
