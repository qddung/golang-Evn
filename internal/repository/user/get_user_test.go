package user

import (
	"context"
	"testing"

	"github.com/homework/lab/internal/models/entity"
	"github.com/homework/lab/internal/test/data/fixture"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	fix := fixture.NewUserTestCase(t)
	db := fixture.NewFixture(t, fix)

	userRepository := NewUserRepository(db)
	ctx := context.Background()
	user := &entity.User{
		DisplayName: "bradtest",
		Email:       "bradtest@gmail.com",
		Password:    "123456",
		UserName:    "bradtest",
	}
	if err := userRepository.CreateUser(ctx, user); !assert.NoError(t, err) {
		t.Fatal(err)
	}

	testCases := []struct {
		name      string
		action    func() (interface{}, error)
		expects   func(t *testing.T, res interface{}, err error)
	}{
		{
			name: "Get by ID",
			action: func() (interface{}, error) {
				return userRepository.GetUserById(ctx, user.Id)
			},
			expects: func(t *testing.T, res interface{}, err error) {
				assert.NoError(t, err)
				u, ok := res.(*entity.User)
				assert.True(t, ok)
				assert.Equal(t, user.Id, u.Id)
			},
		},
		{
			name: "Get by UserName",
			action: func() (interface{}, error) {
				return userRepository.GetUserByUserName(ctx, user.UserName)
			},
			expects: func(t *testing.T, res interface{}, err error) {
				assert.NoError(t, err)
				u, ok := res.(*entity.User)
				assert.True(t, ok)
				assert.Equal(t, user.UserName, u.UserName)
			},
		},
		{
			name: "Get by Email",
			action: func() (interface{}, error) {
				return userRepository.GetUserByEmail(ctx, user.Email)
			},
			expects: func(t *testing.T, res interface{}, err error) {
				assert.NoError(t, err)
				u, ok := res.(*entity.User)
				assert.True(t, ok)
				assert.Equal(t, user.Email, u.Email)
			},
		},
		{
			name: "Get by unknown ID returns nil",
			action: func() (interface{}, error) {
				return userRepository.GetUserById(ctx, "non-existent-id")
			},
			expects: func(t *testing.T, res interface{}, err error) {
				// expect nil result and non-nil error from repo (gorm.ErrRecordNotFound or nil result)
				assert.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.action()
			tc.expects(t, res, err)
		})
	}
}
