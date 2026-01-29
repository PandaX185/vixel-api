package user

type RegisterDto struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (d RegisterDto) ToUser() *User {
	return &User{
		Username: d.Username,
		Email:    d.Email,
		Password: d.Password,
	}
}

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
