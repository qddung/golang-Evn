package user

import (
	"context"

	"github.com/homework/lab/internal/models/entity"
)

// UserRepository interface
//
//go:generate mockery --name UserRepository --filename user_repository.go
type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserById(ctx context.Context, id string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByUserName(ctx context.Context, userName string) (*entity.User, error)
}
