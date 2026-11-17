package translate

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/imroc/req/v3"
)

func TestRetryOneshotRecoversAfterRateLimit(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requests.Add(1) == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		fmt.Fprint(w, `{"translations":[{"text":"안녕하세요"}]}`)
	}))
	defer server.Close()
	client := req.C()
	response, err := retryOneshot(context.Background(), func(ctx context.Context) (oneshotResponse, error) {
		request := client.R().SetContext(ctx).DisableAutoReadResponse().SetBodyBytes([]byte(`{"text":["こんにちは"]}`))
		return postOneshot(request, server.URL)
	})
	if err != nil || response.status != http.StatusOK || response.body.Get("translations.0.text").String() != "안녕하세요" {
		t.Fatalf("got response %+v and error %v, want a successful Korean translation", response, err)
	}
	if got := requests.Load(); got != 2 {
		t.Fatalf("got %d upstream requests, want 2", got)
	}
}

func TestRetryOneshotStopsAtLimit(t *testing.T) {
	for _, tt := range []struct {
		name       string
		retryAfter string
		attempts   int
		cancel     bool
	}{
		{name: "missing Retry-After", attempts: 1},
		{name: "malformed Retry-After", retryAfter: "soon", attempts: 1},
		{name: "persistent 429", retryAfter: "0", attempts: 3},
		{name: "wait longer than budget", retryAfter: "60", attempts: 1},
		{name: "cancelled", retryAfter: "0", attempts: 1, cancel: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			attempts := 0
			response, err := retryOneshot(ctx, func(context.Context) (oneshotResponse, error) {
				attempts++
				if tt.cancel {
					cancel()
				}
				return oneshotResponse{status: http.StatusTooManyRequests, retryAfter: tt.retryAfter}, nil
			})
			if err != nil || response.status != http.StatusTooManyRequests || response.retryAfter != tt.retryAfter || attempts != tt.attempts {
				t.Fatalf("got response %+v, error %v and %d attempts, want 429 with Retry-After %q and %d attempts",
					response, err, attempts, tt.retryAfter, tt.attempts)
			}
		})
	}
}

func TestOneshotRetryDelay(t *testing.T) {
	now := time.Date(2026, 11, 17, 0, 0, 0, 0, time.UTC)
	for _, tt := range []struct {
		header string
		want   time.Duration
		valid  bool
	}{
		{"7", 7 * time.Second, true},
		{"0", 0, true},
		{now.Add(5 * time.Second).Format(http.TimeFormat), 5 * time.Second, true},
		{now.Add(-5 * time.Second).Format(http.TimeFormat), 0, true},
		{"184467440737095516160", oneshotTimeout, true},
		{"", 0, false},
		{"soon", 0, false},
		{"-1", 0, false},
	} {
		t.Run(tt.header, func(t *testing.T) {
			got, valid := oneshotRetryDelay(tt.header, now)
			if got != tt.want || valid != tt.valid {
				t.Fatalf("got (%s, %v), want (%s, %v)", got, valid, tt.want, tt.valid)
			}
		})
	}
}

func TestRetryOneshotPreservesTransportError(t *testing.T) {
	wantErr := errors.New("connection reset")
	attempts := 0
	response, err := retryOneshot(context.Background(), func(context.Context) (oneshotResponse, error) {
		attempts++
		if attempts == 1 {
			return oneshotResponse{status: http.StatusTooManyRequests, retryAfter: "0"}, nil
		}
		return oneshotResponse{}, wantErr
	})
	if !errors.Is(err, wantErr) || response.status != 0 || attempts != 2 {
		t.Fatalf("got %+v, %v after %d attempts; want original error after 2 attempts", response, err, attempts)
	}
}

func TestRetryOneshotCancelledBeforeRequest(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := retryOneshot(ctx, func(context.Context) (oneshotResponse, error) {
		t.Fatal("cancelled context must not send an upstream request")
		return oneshotResponse{}, nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v, want context.Canceled", err)
	}
}
