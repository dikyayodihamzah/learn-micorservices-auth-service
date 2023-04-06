package web

import (
	"time"

	"gitlab.com/learn-micorservices/auth-service/model/domain"
)

type RoleResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RegisterResponse struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Username  string       `json:"username"`
	Email     string       `json:"email"`
	Phone     string       `json:"phone"`
	Role      RoleResponse `json:"role"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type LoginResponse struct {
	ID     string `json:"id"`
	RoleID string `json:"role_id"`
}

func NewRegisterResponse(user domain.User) RegisterResponse {
	return RegisterResponse{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Role: RoleResponse{
			ID:   user.RoleID,
			Name: user.RoleName,
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
