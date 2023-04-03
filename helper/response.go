package helper

import (
	"gitlab.com/learn-micorservices/auth-service/model/domain"
	"gitlab.com/learn-micorservices/auth-service/model/web"
)

// Profile Responses
func ToRegisterResponse(user domain.User) web.RegisterResponse {
	return web.RegisterResponse{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Role: web.RoleResponse{
			ID:   user.Role.ID,
			Name: user.Role.Name,
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
