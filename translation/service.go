package translation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Service struct {
	apiURL string
}

func NewService() *Service {
	return &Service{
		apiURL: "https://api.mymemory.translated.net/get",
	}
}

func (s *Service) Translate(word string) (*MyMemoryResponse, error) {
	encodedWord := url.QueryEscape(word)
	fullURL := fmt.Sprintf("%s?q=%s&langpair=en|vi", s.apiURL, encodedWord)

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to contact translation service: %w", err)
	}
	defer resp.Body.Close()

	var result MyMemoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
