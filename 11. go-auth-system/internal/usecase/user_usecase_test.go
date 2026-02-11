package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-auth-system/internal/domain"
	"go-auth-system/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := usecase.NewUserUsecase(mockRepo, time.Second*2, "secret")
	// Expect GetByEmail to be called and return error (user not found)
	// -> Which is good for registration
	mockRepo.On("GetByEmail", mock.Anything, "test@gmail.com").Return(nil, errors.New("not found"))
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	err := uc.Register(context.Background(), "test@gmail.com", "password123", "Test User")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
