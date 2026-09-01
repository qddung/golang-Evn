package user_service

import (
	"context"
	"errors"
	"testing"

	domain_model "github.com/homework/lab/internal/models/domain"
	userModel "github.com/homework/lab/internal/models/dto/api/user"
	"github.com/homework/lab/internal/models/entity"
	repo_mocks "github.com/homework/lab/internal/repository/user/mocks"
	"github.com/homework/lab/internal/test/data/fixture"
	helper_mocks "github.com/homework/lab/pkg/helpers/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestToUpdateUser(t *testing.T) {
	testCases := []struct {
		name        string
		input       *userModel.UpdateUserInput
		setupHasher func() *helper_mocks.HashHelper
		expects     func(t *testing.T, out *domain_model.UpdateUser, err error)
	}{
		{
			name:        "no password",
			input:       &userModel.UpdateUserInput{UserName: "bob", Password: ""},
			setupHasher: func() *helper_mocks.HashHelper { return nil },
			expects: func(t *testing.T, out *domain_model.UpdateUser, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "bob", out.UserName)
				assert.Equal(t, "", out.Password)
			},
		},
		{
			name:  "hash success",
			input: &userModel.UpdateUserInput{UserName: "alice", Password: "secret"},
			setupHasher: func() *helper_mocks.HashHelper {
				h := helper_mocks.NewHashHelper(t)
				h.On("HashPassword", "secret").Return("hashed-secret", nil)
				return h
			},
			expects: func(t *testing.T, out *domain_model.UpdateUser, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "alice", out.UserName)
				assert.Equal(t, "hashed-secret", out.Password)
			},
		},
		{
			name:  "hash fail",
			input: &userModel.UpdateUserInput{UserName: "joe", Password: "bad"},
			setupHasher: func() *helper_mocks.HashHelper {
				h := helper_mocks.NewHashHelper(t)
				h.On("HashPassword", "bad").Return("", errors.New("hash fail"))
				return h
			},
			expects: func(t *testing.T, out *domain_model.UpdateUser, err error) {
				assert.Equal(t, ErrorParse, err)
				assert.Nil(t, out)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			hasher := tc.setupHasher()
			out, err := ToUpdateUser(tc.input, hasher)
			tc.expects(t, out, err)
		})
	}
}

func TestService_UpdateUserInfo(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name      string
		setupRepo func() *repo_mocks.UserRepository
		inputId   string
		input     *userModel.UpdateUserInput
		expectErr error
	}{
		{
			name: "user not found",
			setupRepo: func() *repo_mocks.UserRepository {
				repo := repo_mocks.NewUserRepository(t)
				repo.On("GetUserById", ctx, "u-not").Return((*entity.User)(nil), gorm.ErrRecordNotFound)
				return repo
			},
			inputId:   "u-not",
			input:     &userModel.UpdateUserInput{UserName: "x"},
			expectErr: ServiceErr.UserNameExistError,
		},
		{
			name: "repo get error",
			setupRepo: func() *repo_mocks.UserRepository {
				repo := repo_mocks.NewUserRepository(t)
				repo.On("GetUserById", ctx, "u-err").Return((*entity.User)(nil), errors.New("boom"))
				return repo
			},
			inputId:   "u-err",
			input:     &userModel.UpdateUserInput{UserName: "x"},
			expectErr: ErrorGetUser,
		},
		{
			name: "update fails",
			setupRepo: func() *repo_mocks.UserRepository {
				repo := repo_mocks.NewUserRepository(t)
				repo.On("GetUserById", ctx, "u-1").Return(&entity.User{Base: fixture.GetBaseEntity("u-1")}, nil)
				repo.On("UpdateUser", ctx, mock.MatchedBy(func(x interface{}) bool {
					v, ok := x.(*domain_model.UpdateUser)
					return ok && v.Id == "u-1" && v.UserName == "new" && v.Password == ""
				})).Return(assert.AnError)
				return repo
			},
			inputId:   "u-1",
			input:     &userModel.UpdateUserInput{UserName: "new"},
			expectErr: UpdateUserFailed,
		},
		{
			name: "success",
			setupRepo: func() *repo_mocks.UserRepository {
				repo := repo_mocks.NewUserRepository(t)
				repo.On("GetUserById", ctx, "u-2").Return(&entity.User{Base: fixture.GetBaseEntity("u-2")}, nil)
				repo.On("UpdateUser", ctx, mock.MatchedBy(func(x interface{}) bool {
					v, ok := x.(*domain_model.UpdateUser)
					return ok && v.Id == "u-2" && v.UserName == "ok" && v.Password == ""
				})).Return(nil)
				return repo
			},
			inputId:   "u-2",
			input:     &userModel.UpdateUserInput{UserName: "ok", Password: ""},
			expectErr: nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repo := tc.setupRepo()
			service := NewUserService(repo, nil, nil)
			err := service.UpdateUserInfo(ctx, tc.inputId, tc.input)
			if tc.expectErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tc.expectErr, err)
			}
		})
	}
}
