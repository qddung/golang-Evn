package user

import (
	"context"

	"github.com/homework/lab/internal/models/entity"
	"gorm.io/gorm"
)

func (u *userRepository) GetUserById(ctx context.Context, id string) (*entity.User, error) {
	user, err := gorm.G[entity.User](u.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := gorm.G[entity.User](u.db).Where("email = ?", email).First(ctx)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &user, nil
}

func (u *userRepository) GetUserByUserName(ctx context.Context, userName string) (*entity.User, error) {
	user, err := gorm.G[entity.User](u.db).Where("user_name = ?", userName).First(ctx)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, nil
}
