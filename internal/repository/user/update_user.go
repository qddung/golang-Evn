package user

import "github.com/homework/lab/internal/models/entity"

func (u *userRepository) UpdateUser(user entity.User) (*entity.User, error) {
	var userEntity entity.User
	u.db.First(&userEntity, "id = ?", user.Id)

	return &userEntity, nil
}
