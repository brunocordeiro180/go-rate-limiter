package ratelimiter

import (
	"context"
	"errors"
	"testing"

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

func TestCheck_IPLimit(t *testing.T) {
	ctx := context.TODO()
	repo := new(MockRepository)
	limiter := NewRateLimiter(repo, 5, 10, 60)

	repo.On("Increment", ctx, "192.168.0.1", 1).Return(3, nil)

	ok := limiter.Check(ctx, "192.168.0.1", false)

	assert.True(t, ok)
	repo.AssertExpectations(t)
}

func TestCheck_TokenLimitExceeded(t *testing.T) {
	ctx := context.TODO()
	repo := new(MockRepository)
	limiter := NewRateLimiter(repo, 5, 2, 60)

	repo.On("Increment", ctx, "auth-token-abc", 1).Return(3, nil)
	repo.On("Increment", ctx, "auth-token-abc:blocked", 60).Return(0, nil)

	ok := limiter.Check(ctx, "auth-token-abc", true)

	assert.False(t, ok)
	repo.AssertExpectations(t)
}

func TestCheck_RepositoryError(t *testing.T) {
	ctx := context.TODO()
	repo := new(MockRepository)
	limiter := NewRateLimiter(repo, 5, 5, 60)

	repo.On("Increment", ctx, "some-key", 1).Return(0, errors.New("db error"))

	ok := limiter.Check(ctx, "some-key", false)

	assert.False(t, ok)
	repo.AssertExpectations(t)
}

func TestIsBlocked_True(t *testing.T) {
	ctx := context.TODO()
	repo := new(MockRepository)
	limiter := NewRateLimiter(repo, 5, 5, 60)

	repo.On("Get", ctx, "some-key:blocked").Return(1, nil)

	blocked := limiter.IsBlocked(ctx, "some-key")

	assert.True(t, blocked)
	repo.AssertExpectations(t)
}

func TestIsBlocked_False(t *testing.T) {
	ctx := context.TODO()
	repo := new(MockRepository)
	limiter := NewRateLimiter(repo, 5, 5, 60)

	repo.On("Get", ctx, "some-key:blocked").Return(0, nil)

	blocked := limiter.IsBlocked(ctx, "some-key")

	assert.False(t, blocked)
	repo.AssertExpectations(t)
}

func TestIsBlocked_Error(t *testing.T) {
	ctx := context.TODO()
	repo := new(MockRepository)
	limiter := NewRateLimiter(repo, 5, 5, 60)

	repo.On("Get", ctx, "some-key:blocked").Return(0, errors.New("db error"))

	blocked := limiter.IsBlocked(ctx, "some-key")

	assert.False(t, blocked)
	repo.AssertExpectations(t)
}
