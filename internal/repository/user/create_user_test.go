package user

import (
	"context"
	"testing"

	"github.com/homework/lab/internal/models/entity"
	"github.com/homework/lab/internal/test/data/fixture"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	fix := fixture.NewUserTestCase(t)
	db := fixture.NewFixture(t, fix)

	userRepository := NewUserRepository(db)
	ctx := context.Background()
	user := &entity.User{
		DisplayName: "brad",
		Email:       "abc@gmail.com",
		Password:    "123456",
		UserName:    "acd",
	}
	errCreateUser := userRepository.CreateUser(ctx, user)
	assert.NoError(t, errCreateUser)

}
