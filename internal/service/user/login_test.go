package user_service

import (
	"context"
	"errors"
	"testing"

	userModel "github.com/homework/lab/internal/models/dto/api/user"
	"github.com/homework/lab/internal/models/entity"
	"github.com/homework/lab/internal/repository/user/mocks"
	helper_mocks "github.com/homework/lab/pkg/helpers/mocks"
	jwt_pkg "github.com/homework/lab/pkg/jwt"
	"github.com/stretchr/testify/assert"
)

func SetupRepoForLogin(t *testing.T) (*mocks.UserRepository, jwt_pkg.JwtGenerator, *helper_mocks.HashHelper) {
	jwtMock := jwt_pkg.NewMockJwt()
	hasher := helper_mocks.NewHashHelper(t)
	repo := mocks.NewUserRepository(t)
	return repo, jwtMock.JwtGenarate, hasher
}

var internalErr = errors.New("internal error")

func TestService_Login(t *testing.T) {
	loginInput := &userModel.UserLogin{UserName: "testuser", Password: "123131242"}
	testCases := []struct {
		name         string
		setupRepo    func(ctx context.Context, info *userModel.UserLogin) (*mocks.UserRepository, jwt_pkg.JwtGenerator, *helper_mocks.HashHelper)
		expectedFunc func(t *testing.T, token string, err error)
	}{
		{
			name: "Not found username",
			setupRepo: func(ctx context.Context, info *userModel.UserLogin) (*mocks.UserRepository, jwt_pkg.JwtGenerator, *helper_mocks.HashHelper) {
				repo, jwtGen, hasher := SetupRepoForLogin(t)
				repo.On("GetUserByUserName", ctx, info.UserName).Return(nil, nil)
				return repo, jwtGen, hasher
			},
			expectedFunc: func(t *testing.T, token string, err error) {
				assert.Equal(t, err, ServiceErr.UserNameNotExistError)
				assert.Equal(t, token, string(""))
			},
		},
		{
			name: "Password not match",
			setupRepo: func(ctx context.Context, info *userModel.UserLogin) (*mocks.UserRepository, jwt_pkg.JwtGenerator, *helper_mocks.HashHelper) {
				repo, jwtGen, hasher := SetupRepoForLogin(t)
				repo.On("GetUserByUserName", ctx, info.UserName).Return(&entity.User{}, nil)
				hasher.On("CheckPasswordHash", loginInput.Password, "").Return(false)
				return repo, jwtGen, hasher
			},
			expectedFunc: func(t *testing.T, token string, err error) {
				assert.Equal(t, err, ServiceErr.PasswordError)
				assert.Equal(t, token, string(""))
			},
		},

		{
			name: "Get user error",
			setupRepo: func(ctx context.Context, info *userModel.UserLogin) (*mocks.UserRepository, jwt_pkg.JwtGenerator, *helper_mocks.HashHelper) {
				repo, jwtGen, hasher := SetupRepoForLogin(t)
				repo.On("GetUserByUserName", ctx, info.UserName).Return(nil, internalErr)
				return repo, jwtGen, hasher
			},
			expectedFunc: func(t *testing.T, token string, err error) {
				assert.Equal(t, err, internalErr)
				assert.Equal(t, token, string(""))
			},
		},
		{
			name: "Success",
			setupRepo: func(ctx context.Context, info *userModel.UserLogin) (*mocks.UserRepository, jwt_pkg.JwtGenerator, *helper_mocks.HashHelper) {
				repo, jwtGen, hasher := SetupRepoForLogin(t)
				repo.On("GetUserByUserName", ctx, info.UserName).Return(&entity.User{}, nil)
				hasher.On("CheckPasswordHash", loginInput.Password, "").Return(true)
				return repo, jwtGen, hasher
			},
			expectedFunc: func(t *testing.T, token string, err error) {
				assert.Equal(t, err, nil)
				assert.NotEmpty(t, token)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(testItem *testing.T) {
			testItem.Parallel()
			ctx := context.Background()

			repo, jwtGen, hasher := tc.setupRepo(ctx, loginInput)
			service := NewUserService(repo, hasher, jwtGen)
			token, err := service.Login(ctx, *loginInput)
			tc.expectedFunc(testItem, token, err)
		})
	}
}
