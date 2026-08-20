package user

// -- Models
type UserRegister struct {
	DisplayName string `json:"display_name" binding:"required"`
	Email       string `json:"email" binding:"email"`
	Password    string `json:"password" binding:"min=8"`
	UserName    string `json:"username" binding:"required"`
}
