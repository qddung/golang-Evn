package user_service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/homework/lab/internal/models/dto/api/user"
	"github.com/homework/lab/internal/models/entity"
	"github.com/homework/lab/internal/repository/user/mocks"
	helper_mocks "github.com/homework/lab/pkg/helpers/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testErr = errors.New("test error")

func TestService_Register(t *testing.T) {
	testCases := []struct {
		name         string
		setupRepo    func(ctx context.Context, info *user.UserRegister) (*mocks.UserRepository, *helper_mocks.HashHelper)
		expectedFunc func(t *testing.T, info *user.UserInfo, registerInput *user.UserRegister, err error)
	}{
		{
			// Duplicate email
			name: "Duplicate email",
			setupRepo: func(ctx context.Context, info *user.UserRegister) (*mocks.UserRepository, *helper_mocks.HashHelper) {
				repo := mocks.NewUserRepository(t)
				hasher := helper_mocks.NewHashHelper(t)
				repo.On("GetUserByUserName", ctx, info.UserName).Return(nil, nil)
				repo.On("GetUserByEmail", ctx, info.Email).Return(&entity.User{}, nil)
				return repo, hasher
			},
			expectedFunc: func(t *testing.T, info *user.UserInfo, registerInput *user.UserRegister, err error) {
				assert.Equal(t, err, EmailExistError)
			},
		},
		{
			// Duplicate username
			name: "Duplicate username",
			setupRepo: func(ctx context.Context, info *user.UserRegister) (*mocks.UserRepository, *helper_mocks.HashHelper) {
				repo := mocks.NewUserRepository(t)
				hasher := helper_mocks.NewHashHelper(t)
				repo.On("GetUserByUserName", ctx, info.UserName).Return(&entity.User{}, nil)
				return repo, hasher
			},
			expectedFunc: func(t *testing.T, info *user.UserInfo, registerInput *user.UserRegister, err error) {
				assert.Equal(t, err, UserNameExistError)
			},
		},
		{
			name: "Create user error",
			setupRepo: func(ctx context.Context, info *user.UserRegister) (*mocks.UserRepository, *helper_mocks.HashHelper) {
				repo := mocks.NewUserRepository(t)
				hasher := helper_mocks.NewHashHelper(t)
				expectedHashPass := "hash"
				hasher.On("HashPassword", info.Password).Return(expectedHashPass, nil)
				repo.On("GetUserByUserName", ctx, info.UserName).Return(nil, nil)
				repo.On("GetUserByEmail", ctx, info.Email).Return(nil, nil)
				repo.On("CreateUser", ctx, &entity.User{
					DisplayName: info.DisplayName,
					Email:       info.Email,
					Password:    expectedHashPass,
					UserName:    info.UserName,
				}).Return(testErr)
				return repo, hasher
			},
			expectedFunc: func(t *testing.T, info *user.UserInfo, registerInput *user.UserRegister, err error) {
				assert.Equal(t, err, testErr)
			},
		},
		{
			name: "Create user successfully",
			setupRepo: func(ctx context.Context, info *user.UserRegister) (*mocks.UserRepository, *helper_mocks.HashHelper) {
				repo := mocks.NewUserRepository(t)
				hasher := helper_mocks.NewHashHelper(t)
				expectedHashPass := "hash"
				hasher.On("HashPassword", info.Password).Return(expectedHashPass, nil)
				repo.On("GetUserByUserName", ctx, info.UserName).Return(nil, nil)
				repo.On("GetUserByEmail", ctx, info.Email).Return(nil, nil)
				user := &entity.User{
					DisplayName: info.DisplayName,
					Email:       info.Email,
					Password:    expectedHashPass,
					UserName:    info.UserName,
				}
				repo.On("CreateUser", ctx, user).Run(func(args mock.Arguments) {
					now := time.Now()
					user.Id = uuid.NewString()
					user.CreatedAt = now
					user.UpdatedAt = now
				}).Return(nil)
				return repo, hasher
			},
			expectedFunc: func(t *testing.T, info *user.UserInfo, registerInput *user.UserRegister, err error) {
				assert.NoError(t, err)
				assert.Equal(t, info.DisplayName, registerInput.DisplayName)
				assert.Equal(t, info.Email, registerInput.Email)
				assert.Equal(t, info.UserName, registerInput.UserName)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(testItem *testing.T) {
			testItem.Parallel()
			ctx := context.Background()
			u := user.UserRegister{
				DisplayName: "test",
				Email:       "test@example.com",
				Password:    "123131242",
				UserName:    "testuser",
			}
			repo, hashMock := tc.setupRepo(ctx, &u)
			service := NewUserService(repo, hashMock)
			info, err := service.Register(ctx, u)

			tc.expectedFunc(testItem, info, &u, err)
		})
	}
}
