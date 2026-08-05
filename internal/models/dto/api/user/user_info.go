package user

import (
	"github.com/homework/lab/internal/models/entity"
)

type UserInfo struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	UserName    string `json:"userName"`
	CreateAt    string `json:"created_at"`
}

func (u *UserInfo) PopulateInfoFromUserEntity(entityUser *entity.User) {
	u.DisplayName = entityUser.DisplayName
	u.Email = entityUser.Email
	u.UserName = entityUser.UserName
	u.CreateAt = entityUser.CreateAt.String()
}
