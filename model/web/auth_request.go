package web

type RegisterRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	RoleID   string `json:"role_id"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type ForgetPasswordRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResetPassword struct {
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
}
