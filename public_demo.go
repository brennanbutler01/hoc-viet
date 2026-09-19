package main

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const publicTranslationRequestsPerMinute = 60

type publicDemoLimiter struct {
	mutex       sync.Mutex
	windowStart time.Time
	requests    int
}

func (l *publicDemoLimiter) allow(now time.Time) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.windowStart.IsZero() || now.Sub(l.windowStart) >= time.Minute {
		l.windowStart = now
		l.requests = 0
	}
	if l.requests >= publicTranslationRequestsPerMinute {
		return false
	}
	l.requests++
	return true
}

func newPublicDemoMiddleware(next http.Handler) http.Handler {
	limiter := &publicDemoLimiter{}
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Content-Type-Options", "nosniff")
		writer.Header().Set("Referrer-Policy", "no-referrer")
		if request.Method != http.MethodGet && request.Method != http.MethodHead {
			writer.Header().Set("Allow", "GET, HEAD")
			http.Error(writer, "The hosted demo is read-only.", http.StatusMethodNotAllowed)
			return
		}
		if strings.HasPrefix(request.URL.Path, "/translation/") && !limiter.allow(time.Now()) {
			writer.Header().Set("Retry-After", strconv.Itoa(60))
			http.Error(writer, "The hosted demo translation limit has been reached. Try again later.", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(writer, request)
	})
}
