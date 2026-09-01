package user_service

import (
	"context"

	"testing"

	userModel "github.com/homework/lab/internal/models/dto/api/user"
	"github.com/homework/lab/internal/models/entity"
	repo_mocks "github.com/homework/lab/internal/repository/user/mocks"
	"github.com/homework/lab/internal/test/data/fixture"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func SetupRepoForInfo(t *testing.T) *repo_mocks.UserRepository {
	return repo_mocks.NewUserRepository(t)
}

func TestService_GetUserInfo(t *testing.T) {
	testCases := []struct {
		name      string
		setupRepo func(ctx context.Context) *repo_mocks.UserRepository
		expect    func(t *testing.T, info *userModel.UserInfo, err error)
	}{
		{
			name: "not found",
			setupRepo: func(ctx context.Context) *repo_mocks.UserRepository {
				repo := SetupRepoForInfo(t)
				repo.On("GetUserById", ctx, "missing-id").Return(nil, gorm.ErrRecordNotFound)
				return repo
			},
			expect: func(t *testing.T, info *userModel.UserInfo, err error) {
				assert.Nil(t, info)
				assert.Equal(t, ServiceErr.NotFoundUserInfo, err)
			},
		},
		{
			name: "repo error",
			setupRepo: func(ctx context.Context) *repo_mocks.UserRepository {
				repo := SetupRepoForInfo(t)
				repo.On("GetUserById", ctx, "some-id").Return(nil, internalErr)
				return repo
			},
			expect: func(t *testing.T, info *userModel.UserInfo, err error) {
				assert.Nil(t, info)
				assert.Equal(t, internalErr, err)
			},
		},
		{
			name: "success",
			setupRepo: func(ctx context.Context) *repo_mocks.UserRepository {
				repo := SetupRepoForInfo(t)
				existing := &entity.User{
					Base:        fixture.GetBaseEntity("u-123"),
					DisplayName: "Alice",
					Email:       "alice@example.com",
					UserName:    "alice",
				}
				repo.On("GetUserById", ctx, "u-123").Return(existing, nil)
				return repo
			},
			expect: func(t *testing.T, info *userModel.UserInfo, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, info)
				assert.Equal(t, "u-123", info.Id)
				assert.Equal(t, "Alice", info.DisplayName)
				assert.Equal(t, "alice@example.com", info.Email)
				assert.Equal(t, "alice", info.UserName)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo := tc.setupRepo(ctx)
			// hasher and jwt not needed for GetUserInfo
			service := NewUserService(repo, nil, nil)
			info, err := service.GetUserInfo(ctx, func() string {
				if tc.name == "not found" {
					return "missing-id"
				}
				if tc.name == "repo error" {
					return "some-id"
				}
				return "u-123"
			}())
			tc.expect(t, info, err)
		})
	}
}
