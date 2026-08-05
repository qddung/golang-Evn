package user_service

import (
	"context"
	"errors"
	"testing"

	"github.com/homework/lab/internal/models/dto/api/user"
	"github.com/homework/lab/internal/models/entity"
	"github.com/homework/lab/internal/repository/user/mocks"
	"github.com/stretchr/testify/assert"
)

var testErr = errors.New("test error")

func TestService_GetLinkFromCode(t *testing.T) {
	testCases := []struct {
		name string

		setupRepo    func(ctx context.Context) *mocks.UserRepository
		expectedFunc func(t *testing.T, info *user.UserInfo, err error)
	}{
		{
			// Duplicate email
			name: "Duplicate email",
			setupRepo: func(ctx context.Context) *mocks.UserRepository {
				repo := &mocks.UserRepository{}
				repo.On("GetUserByUserName", ctx, "testuser").Return(nil, nil)
				repo.On("GetUserByEmail", ctx, "testemail").Return(entity.User{}, nil)
				return repo
			},
			expectedFunc: func(t *testing.T, info *user.UserInfo, err error) {
				assert.Equal(t, err, EmailExistError)
			},
			// Duplicate username
		},
		// {
		// 	name: "Duplicate username",
		// 	setupRepo: func(ctx context.Context) *mocks.UserRepository {
		// 		repo := &mocks.UserRepository{}
		// 		repo.On("GetUserByUserName", ctx, "testuser").Return(entity.User{}, UserNameExistError)
		// 		repo.On("GetUserByEmail", ctx, "testemail").Return(nil, nil)
		// 		return repo
		// 	},
		// 	expectedFunc: func(t *testing.T, info *user.UserInfo, err error) {
		// 		assert.Equal(t, err, UserNameExistError)
		// 	},
		// },
		// {
		// 	name: "Create user error",
		// 	setupRepo: func(ctx context.Context) *mocks.UserRepository {
		// 		repo := &mocks.UserRepository{}
		// 		repo.On("GetUserByUserName", ctx, "testuser").Return(nil, nil)
		// 		repo.On("GetUserByEmail", ctx, "testemail").Return(nil, nil)
		// 		repo.On("CreateUser", ctx, &entity.User{}).Return(testErr)
		// 		return repo
		// 	},
		// 	expectedFunc: func(t *testing.T, info *user.UserInfo, err error) {
		// 		assert.Equal(t, err, testErr)
		// 	},
		// },
		// {
		// 	name: "Create user successfully",
		// 	setupRepo: func(ctx context.Context) *mocks.UserRepository {
		// 		repo := &mocks.UserRepository{}
		// 		repo.On("GetUserByUserName", ctx, "testuser").Return(nil, nil)
		// 		repo.On("GetUserByEmail", ctx, "testemail").Return(nil, nil)
		// 		repo.On("CreateUser", ctx, &entity.User{}).Return(nil)
		// 		return repo
		// 	},
		// 	expectedFunc: func(t *testing.T, info *user.UserInfo, err error) {
		// 		assert.NoError(t, err)
		// 	},
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo := tc.setupRepo(ctx)
			service := NewUserService(repo)
			info, err := service.Register(ctx, user.UserRegister{})

			tc.expectedFunc(t, info, err)
		})
	}

}
