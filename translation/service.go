package translation

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"
)

const maximumResponseBytes = 1024 * 1024

type Service struct {
	apiURL string
	client *http.Client
}

func NewService() *Service {
	return &Service{
		apiURL: "https://api.mymemory.translated.net/get",
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *Service) Translate(ctx context.Context, word string) (*MyMemoryResponse, error) {
	word = strings.TrimSpace(word)
	if word == "" || utf8.RuneCountInString(word) > 30 {
		return nil, fmt.Errorf("word must contain 1 to 30 characters")
	}
	endpoint, err := url.Parse(s.apiURL)
	if err != nil {
		return nil, fmt.Errorf("invalid translation endpoint: %w", err)
	}
	query := endpoint.Query()
	query.Set("q", word)
	query.Set("langpair", "en|vi")
	endpoint.RawQuery = query.Encode()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create translation request: %w", err)
	}
	response, err := s.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("contact translation service: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("translation provider returned HTTP %d", response.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, maximumResponseBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read translation response: %w", err)
	}
	if len(data) > maximumResponseBytes {
		return nil, fmt.Errorf("translation response exceeds size limit")
	}
	var result MyMemoryResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode translation response: %w", err)
	}
	if result.ResponseStatus != http.StatusOK || result.QuotaFinished || strings.TrimSpace(result.ResponseData.TranslatedText) == "" {
		return nil, fmt.Errorf("translation provider did not return a usable translation (status=%d, quotaFinished=%t)", result.ResponseStatus, result.QuotaFinished)
	}
	return &result, nil
}
