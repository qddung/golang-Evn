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
	errCreateUser := userRepository.CreateUser(ctx, user)
	assert.NoError(t, errCreateUser)
	assert.NoError(t, nil)

	userId := user.Id
	userWithID, err := userRepository.GetUserById(ctx, userId)
	assert.NoError(t, err)
	assert.Equal(t, userWithID.Id, user.Id)

	userWithName, err := userRepository.GetUserByUserName(ctx, user.UserName)

	assert.Equal(t, userWithName.UserName, user.UserName)

	userWithMail, err := userRepository.GetUserByEmail(ctx, user.Email)

	assert.Equal(t, userWithMail.Email, user.Email)

}
