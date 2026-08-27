package user

type UpdateUserInput struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
