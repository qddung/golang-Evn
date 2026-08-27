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
	// table-driven test cases; each subtest gets its own isolated DB and runs in parallel
	tests := []struct {
		name           string
		initialUser    *entity.User
		updateUserName string
		updatePassword string
		expectUserName string
		expectPassword string
	}{
		{
			name: "update username only",
			initialUser: &entity.User{
				DisplayName: "bradtest",
				Email:       "bradtest@gmail.com",
				Password:    "123456",
				UserName:    "bradtest",
			},
			updateUserName: "newusername",
			updatePassword: "",
			expectUserName: "newusername",
			expectPassword: "123456",
		},
		{
			name: "update password only",
			initialUser: &entity.User{
				DisplayName: "bradtest",
				Email:       "bradtest@gmail.com",
				Password:    "123456",
				UserName:    "bradtest",
			},
			updateUserName: "",
			updatePassword: "new-secret",
			expectUserName: "bradtest",
			expectPassword: "new-secret",
		},
		{
			name: "update username and password",
			initialUser: &entity.User{
				DisplayName: "bradtest",
				Email:       "bradtest@gmail.com",
				Password:    "123456",
				UserName:    "bradtest",
			},
			updateUserName: "finalname",
			updatePassword: "final-pass",
			expectUserName: "finalname",
			expectPassword: "final-pass",
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fix := fixture.NewUserTestCase(t)
			db := fixture.NewFixture(t, fix)

			userRepository := NewUserRepository(db)
			ctx := context.Background()

			// create initial user for this case
			user := &entity.User{
				DisplayName: tc.initialUser.DisplayName,
				Email:       tc.initialUser.Email,
				Password:    tc.initialUser.Password,
				UserName:    tc.initialUser.UserName,
			}
			if err := userRepository.CreateUser(ctx, user); err != nil {
				t.Fatalf("failed to create user: %v", err)
			}

			// prepare update
			update := &domain_model.UpdateUser{
				Id:       user.Id,
				UserName: tc.updateUserName,
				Password: tc.updatePassword,
			}

			// perform update
			if err := userRepository.UpdateUser(ctx, update); err != nil {
				t.Fatalf("UpdateUser failed: %v", err)
			}

			// verify
			updated, err := userRepository.GetUserById(ctx, user.Id)
			if err != nil {
				t.Fatalf("GetUserById failed: %v", err)
			}
			assert.Equal(t, tc.expectUserName, updated.UserName)
			assert.Equal(t, tc.expectPassword, updated.Password)
		})
	}
}
