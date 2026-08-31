package user

import (
	"context"

	"github.com/homework/lab/internal/models/base"
	domain_model "github.com/homework/lab/internal/models/domain"
	"github.com/homework/lab/internal/models/entity"
)

// UpdateUser
func (u *userRepository) UpdateUser(ctx context.Context, user *domain_model.UpdateUser) error {
	usr := &entity.User{
		Base: base.Base{
			Id: user.Id,
		},
	}
	if err := u.db.WithContext(ctx).Where(usr).First(usr).Error; err != nil {
		return err
	}

	// Apply only provided fields
	if user.Password != "" {
		usr.Password = user.Password
	}
	if user.UserName != "" {
		usr.UserName = user.UserName
	}

	// Save the updated entity
	return u.db.WithContext(ctx).Save(usr).Error
}
