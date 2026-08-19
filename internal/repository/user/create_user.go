package user

import (
	"context"

	"github.com/homework/lab/internal/models/entity"
)

// CreateUser Repository
func (u *userRepository) CreateUser(ctx context.Context, user *entity.User) error {
	return u.db.WithContext(ctx).Create(user).Error
}
