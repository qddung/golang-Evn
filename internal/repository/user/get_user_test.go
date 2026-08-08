package user

import (
	"context"
	"testing"

	"github.com/homework/lab/internal/models/entity"
	"github.com/homework/lab/pkg/sqldb"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	db, err := sqldb.NewMiniPostgres(t)

	assert.NoError(t, err)
	if err != nil {
		t.Fatal(err)
		return
	}

	err_migrate := db.AutoMigrate(entity.User{})
	assert.NoError(t, err_migrate)
	if err_migrate != nil {
		t.Fatal(err_migrate.Error())
		return
	}

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
	assert.NoError(t, nil)

	userWithName, err := userRepository.GetUserByUserName(ctx, user.UserName)

	assert.Equal(t, userWithName.UserName, user.UserName)

	userWithMail, err := userRepository.GetUserByEmail(ctx, user.Email)

	assert.Equal(t, userWithMail.Email, user.Email)

}
