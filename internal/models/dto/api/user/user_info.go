package user

import (
	"github.com/homework/lab/internal/models/entity"
)

type UserInfo struct {
	Id          string `json:"id"`
	UpdateAt    string `json:"updated_at"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	UserName    string `json:"username"`
	CreateAt    string `json:"created_at"`
}

func (u *UserInfo) PopulateInfoFromUserEntity(entityUser *entity.User) {
	u.Id = entityUser.Id
	u.DisplayName = entityUser.DisplayName
	u.Email = entityUser.Email
	u.UserName = entityUser.UserName
	u.CreateAt = entityUser.CreatedAt.String()
	u.UpdateAt = entityUser.UpdatedAt.String()
}
