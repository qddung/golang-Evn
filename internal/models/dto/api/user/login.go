package user

type UserLogin struct {
	UserName string `json:"username" validate:"required"`
	Password string `json:"password" validate:"min=8"`
}
