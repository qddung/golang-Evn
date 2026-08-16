package user

type UserLogin struct {
	UserName string `json:"username" validate:"email"`
	Password string `json:"password" validate:"min=8"`
}
