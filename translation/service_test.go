package translation

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestTranslateEncodesQueryAndReadsVietnamese(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") != "bread & butter" || r.URL.Query().Get("langpair") != "en|vi" {
			t.Error("translation query was not preserved")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"responseData":{"translatedText":"bánh mì và bơ"},"responseStatus":200}`))
	}))
	defer server.Close()
	service := &Service{apiURL: server.URL, client: server.Client()}
	result, err := service.Translate(context.Background(), " bread & butter ")
	if err != nil {
		t.Fatal(err)
	}
	if result.ResponseData.TranslatedText != "bánh mì và bơ" {
		t.Fatalf("unexpected translation: %q", result.ResponseData.TranslatedText)
	}
}

func TestTranslateRejectsUnusableResponses(t *testing.T) {
	cases := []struct {
		name   string
		status int
		body   string
	}{
		{"HTTP failure", 503, `{"responseData":{"translatedText":"wrong"},"responseStatus":200}`},
		{"invalid JSON", 200, `<html>bad gateway</html>`},
		{"provider failure", 200, `{"responseData":{"translatedText":"quota reached"},"responseStatus":403}`},
		{"quota finished", 200, `{"responseData":{"translatedText":"quota reached"},"responseStatus":200,"quotaFinished":true}`},
		{"empty translation", 200, `{"responseData":{"translatedText":" "},"responseStatus":200}`},
		{"oversized response", 200, strings.Repeat("x", maximumResponseBytes+1)},
	}
	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(item.status); w.Write([]byte(item.body)) }))
			defer server.Close()
			service := &Service{apiURL: server.URL, client: server.Client()}
			if result, err := service.Translate(context.Background(), "hello"); err == nil || result != nil {
				t.Fatal("unusable response was accepted")
			}
		})
	}
}

func TestTranslateRespectsCancellationAndTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { <-r.Context().Done() }))
	defer server.Close()
	service := &Service{apiURL: server.URL, client: &http.Client{Timeout: 30 * time.Millisecond}}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := service.Translate(ctx, "hello"); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
	started := time.Now()
	if _, err := service.Translate(context.Background(), "hello"); err == nil {
		t.Fatal("expected timeout")
	}
	if time.Since(started) > time.Second {
		t.Fatal("upstream timeout did not bound the request")
	}
}

func TestTranslateRejectsInvalidInputBeforeNetwork(t *testing.T) {
	service := NewService()
	for _, word := range []string{"", "   ", strings.Repeat("a", 31)} {
		if _, err := service.Translate(context.Background(), word); err == nil {
			t.Fatal("invalid word accepted")
		}
	}
}
