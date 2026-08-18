package user

import (
	"context"
	"testing"

	domain_model "github.com/homework/lab/internal/models/domain"
	"github.com/homework/lab/internal/models/entity"
	"github.com/homework/lab/internal/test/data/fixture"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUser(t *testing.T) {
	fix := fixture.NewUserTestCase(t)
	db := fixture.NewFixture(t, fix)

	userRepository := NewUserRepository(db)
	ctx := context.Background()

	// // create a base user
	user := &entity.User{
		DisplayName: "bradtest",
		Email:       "bradtest@gmail.com",
		Password:    "123456",
		UserName:    "bradtest",
	}
	errCreateUser := userRepository.CreateUser(ctx, user)
	assert.NoError(t, errCreateUser)

	origPassword := user.Password

	// Subtest: update username only
	t.Run("update username only", func(t *testing.T) {
		update := &domain_model.UpdateUser{
			Id:       user.Id,
			UserName: "newusername",
			Password: "",
		}
		err := userRepository.UpdateUser(ctx, update)
		assert.NoError(t, err)

		updated, err := userRepository.GetUserById(ctx, user.Id)
		assert.NoError(t, err)
		assert.Equal(t, "newusername", updated.UserName)
		// password should remain unchanged
		assert.Equal(t, origPassword, updated.Password)
	})

	// Subtest: update password only
	t.Run("update password only", func(t *testing.T) {
		update := &domain_model.UpdateUser{
			Id:       user.Id,
			UserName: "",
			Password: "new-secret",
		}
		err := userRepository.UpdateUser(ctx, update)
		assert.NoError(t, err)

		updated, err := userRepository.GetUserById(ctx, user.Id)
		assert.NoError(t, err)
		// username should remain as previously set
		assert.Equal(t, "newusername", updated.UserName)
		// password should be updated
		assert.Equal(t, "new-secret", updated.Password)
	})

	// Subtest: update both username and password
	t.Run("update username and password", func(t *testing.T) {
		update := &domain_model.UpdateUser{
			Id:       user.Id,
			UserName: "finalname",
			Password: "final-pass",
		}
		err := userRepository.UpdateUser(ctx, update)
		assert.NoError(t, err)

		updated, err := userRepository.GetUserById(ctx, user.Id)
		assert.NoError(t, err)
		assert.Equal(t, "finalname", updated.UserName)
		assert.Equal(t, "final-pass", updated.Password)
	})
}
