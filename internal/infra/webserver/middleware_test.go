package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brunocordeiro180/go-rate-limiter/internal/pkg/ratelimiter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Increment(ctx context.Context, key string, value int) (int, error) {
	args := m.Called(ctx, key, value)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) Get(ctx context.Context, key string) (int, error) {
	args := m.Called(ctx, key)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) Reset(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func TestMiddleware_AllowsRequest(t *testing.T) {
	repo := new(MockRepository)
	limiter := ratelimiter.NewRateLimiter(repo, 5, 5, 60)

	repo.On("Get", mock.Anything, "192.0.2.1:blocked").Return(0, nil)
	repo.On("Increment", mock.Anything, "192.0.2.1", 1).Return(1, nil)

	called := false
	handler := RateLimiterMiddleware(limiter, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.0.2.1"
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.True(t, called)
	repo.AssertExpectations(t)
}

func TestMiddleware_BlocksWhenRateLimitExceeded(t *testing.T) {
	repo := new(MockRepository)
	limiter := ratelimiter.NewRateLimiter(repo, 5, 5, 60)

	repo.On("Get", mock.Anything, "token-123:blocked").Return(0, nil)
	repo.On("Increment", mock.Anything, "token-123", 1).Return(6, nil)
	repo.On("Increment", mock.Anything, "token-123:blocked", 60).Return(0, nil)

	handler := RateLimiterMiddleware(limiter, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fail()
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "should-be-overridden"
	req.Header.Set("API_KEY", "token-123")
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusTooManyRequests, resp.Code)
	repo.AssertExpectations(t)
}

func TestMiddleware_BlocksWhenBlocked(t *testing.T) {
	repo := new(MockRepository)
	limiter := ratelimiter.NewRateLimiter(repo, 5, 5, 60)

	repo.On("Get", mock.Anything, "192.0.2.1:blocked").Return(1, nil)

	handler := RateLimiterMiddleware(limiter, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fail()
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.0.2.1"
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusTooManyRequests, resp.Code)
	repo.AssertExpectations(t)
}

func TestMiddleware_HandlesRepositoryError(t *testing.T) {
	repo := new(MockRepository)
	limiter := ratelimiter.NewRateLimiter(repo, 5, 5, 60)

	repo.On("Get", mock.Anything, "192.0.2.1:blocked").Return(0, errors.New("db error"))
	repo.On("Increment", mock.Anything, "192.0.2.1", 1).Return(0, errors.New("db error"))

	handler := RateLimiterMiddleware(limiter, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fail()
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.0.2.1"
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusTooManyRequests, resp.Code)
	repo.AssertExpectations(t)
}
