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
	hasher_mocks "github.com/homework/lab/pkg/helpers/hasher/mocks"
	jwt_pkg "github.com/homework/lab/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testErr = errors.New("test error")

func SetupRepo(t *testing.T) (*mocks.UserRepository, *hasher_mocks.HashHelper) {
	repo := mocks.NewUserRepository(t)
	hasher := hasher_mocks.NewHashHelper(t)
	return repo, hasher
}

func setupCreateUser(repo *mocks.UserRepository, hasher *hasher_mocks.HashHelper, info *user.UserRegister, ctx context.Context) *entity.User {
	expectedHashPass := "hash"
	hasher.On("HashPassword", info.Password).Return(expectedHashPass, nil)
	repo.On("GetUserByUserName", ctx, info.UserName).Return(nil, nil)
	repo.On("GetUserByEmail", ctx, info.Email).Return(nil, nil)
	return &entity.User{
		DisplayName: info.DisplayName,
		Email:       info.Email,
		Password:    expectedHashPass,
		UserName:    info.UserName,
	}
}

func TestService_Register(t *testing.T) {

	testCases := []struct {
		name        string
		setupRepo   func(ctx context.Context, info *user.UserRegister) (*mocks.UserRepository, *hasher_mocks.HashHelper)
		expectedErr error
	}{
		{
			name: "Failed to GetUserByUserName in userService.Register",
			setupRepo: func(ctx context.Context, info *user.UserRegister) (*mocks.UserRepository, *hasher_mocks.HashHelper) {
				repo, hasher := SetupRepo(t)
				repo.On("GetUserByUserName", ctx, info.UserName).Return(nil, testErr)
				return repo, hasher
			},
			expectedErr: testErr,
		},
		{
			// Duplicate email
			name: "Duplicate email",
			setupRepo: func(ctx context.Context, info *user.UserRegister) (*mocks.UserRepository, *hasher_mocks.HashHelper) {
				repo, hasher := SetupRepo(t)
				repo.On("GetUserByUserName", ctx, info.UserName).Return(nil, nil)
				repo.On("GetUserByEmail", ctx, info.Email).Return(&entity.User{}, nil)
				return repo, hasher
			},

			expectedErr: ServiceErr.EmailExistError,
		},
		{
			// Duplicate username
			name: "Duplicate username",
			setupRepo: func(ctx context.Context, info *user.UserRegister) (*mocks.UserRepository, *hasher_mocks.HashHelper) {
				repo, hasher := SetupRepo(t)
				repo.On("GetUserByUserName", ctx, info.UserName).Return(&entity.User{}, nil)
				return repo, hasher
			},
			expectedErr: ServiceErr.UserNameExistError,
		},
		{
			name: "Create user error",
			setupRepo: func(ctx context.Context, info *user.UserRegister) (*mocks.UserRepository, *hasher_mocks.HashHelper) {
				repo, hasher := SetupRepo(t)
				u := setupCreateUser(repo, hasher, info, ctx)
				repo.On("CreateUser", ctx, u).Return(testErr)
				return repo, hasher
			},
			expectedErr: testErr,
		},
		{
			name: "Create user successfully",
			setupRepo: func(ctx context.Context, info *user.UserRegister) (*mocks.UserRepository, *hasher_mocks.HashHelper) {
				repo, hasher := SetupRepo(t)
				user := setupCreateUser(repo, hasher, info, ctx)
				repo.On("CreateUser", ctx, user).Run(func(args mock.Arguments) {
					now := time.Now()
					user.Id = uuid.NewString()
					user.CreatedAt = now
					user.UpdatedAt = now
				}).Return(nil)
				return repo, hasher
			},
			expectedErr: nil,
		},

		{
			name: "Failed to HashPassword in userService.Register",
			setupRepo: func(ctx context.Context, info *user.UserRegister) (*mocks.UserRepository, *hasher_mocks.HashHelper) {
				repo, hasher := SetupRepo(t)
				// existence checks succeed
				repo.On("GetUserByUserName", ctx, info.UserName).Return(nil, nil)
				repo.On("GetUserByEmail", ctx, info.Email).Return(nil, nil)
				// hasher fails
				hasher.On("HashPassword", info.Password).Return("", testErr)
				return repo, hasher
			},
		},

		{
			name: "Failed to GetUserByEmail in userService.Register",
			setupRepo: func(ctx context.Context, info *user.UserRegister) (*mocks.UserRepository, *hasher_mocks.HashHelper) {
				repo, hasher := SetupRepo(t)
				// username check passes
				repo.On("GetUserByUserName", ctx, info.UserName).Return(nil, nil)
				// GetUserByEmail returns error
				repo.On("GetUserByEmail", ctx, info.Email).Return(nil, testErr)
				return repo, hasher
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(testItem *testing.T) {
			testItem.Parallel()
			ctx := context.Background()
			jwtMock := jwt_pkg.NewMockJwt()
			u := &user.UserRegister{
				DisplayName: "test",
				Email:       "test@example.com",
				Password:    "123131242",
				UserName:    "testuser",
			}
			repo, hashMock := tc.setupRepo(ctx, u)
			service := NewUserService(repo, hashMock, jwtMock.JwtGenarate)
			info, err := service.Register(ctx, u)
			if info != nil {
				assert.Equal(t, info.DisplayName, u.DisplayName)
				assert.Equal(t, info.Email, u.Email)
				assert.Equal(t, info.UserName, u.UserName)
			}
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
			}
		})
	}
}
