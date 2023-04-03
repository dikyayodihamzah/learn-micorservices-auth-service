package service

import (
	"context"
	"time"

	"github.com/go-playground/validator"
	"gitlab.com/learn-micorservices/auth-service/exception"
	"gitlab.com/learn-micorservices/auth-service/helper"
	"gitlab.com/learn-micorservices/auth-service/model/domain"
	"gitlab.com/learn-micorservices/auth-service/model/web"
	"gitlab.com/learn-micorservices/auth-service/repository"
)

type AuthService interface {
	GetCurrentProfile(c context.Context, claims helper.JWTClaims) (web.ProfileResponse, error)
	UpdateProfile(c context.Context, claims helper.JWTClaims, request web.UpdateProfileRequest) (web.ProfileResponse, error)
}

type authService struct {
	AuthRepository repository.AuthRepository
	RoleRepository repository.RoleRepository
	Validate       *validator.Validate
}

func NewAuthService(authRepository repository.AuthRepository, roleRepository repository.RoleRepository, validate *validator.Validate) AuthService {
	return &authService{
		AuthRepository: authRepository,
		RoleRepository: roleRepository,
		Validate:       validate,
	}
}

func (service *authService) Register(c context.Context, request web.RegisterRequest) (web.RegisterResponse, error) {
	if err := service.Validate.Struct(request); err != nil {
		return web.RegisterResponse{}, exception.ErrBadRequest(err.Error())
	}

	if userByEmail, _ := service.AuthRepository.GetUsersByQuery(c, "email", request.Email); userByEmail.ID != "" {
		exception.ErrBadRequest("email already registered")
	}

	if userByUsername, _ := service.AuthRepository.GetUsersByQuery(c, "username", request.Username); userByUsername.ID != "" {
		exception.ErrBadRequest("username already registered")
	}

	if request.Phone != "" {
		if !helper.IsNumeric(request.Phone) {
			panic(exception.ErrBadRequest("Phone should numeric"))
		}

		if len([]rune(request.Phone)) < 10 || len([]rune(request.Phone)) > 13 {
			panic(exception.ErrBadRequest("Phone should 10-13 digit"))
		}

		if userByPhone, _ := service.AuthRepository.GetUsersByQuery(c, "phone", request.Phone); userByPhone.ID != "" {
			exception.ErrBadRequest("phone already registered")
		}
	}

	if role := service.RoleRepository.GetRoleByID(c, request.RoleID); role.ID == "" {
		exception.ErrBadRequest("role not found")
	}

	user := domain.User{
		Name:     request.Name,
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
		Role: domain.Role{
			ID: request.RoleID,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	user.GenerateID()
	user.SetPassword(request.Password)

	err := service.AuthRepository.Register(c, user)
	if err != nil {
		return web.RegisterResponse{}, err
	}

	user, err = service.AuthRepository.GetUsersByQuery(c, "id", user.ID)
	if err != nil {
		return web.RegisterResponse{}, err
	}

	// KAFKA

	return helper.ToRegisterResponse(user), nil
}

func (service *profileService) GetCurrentProfile(c context.Context, claims helper.JWTClaims) (web.ProfileResponse, error) {
	user, err := service.ProfileRepository.GetProfileByID(c, claims.Profile.ID)
	if err != nil {
		return web.ProfileResponse{}, err
	}

	if user.ID == "" {
		return web.ProfileResponse{}, exception.ErrNotFound("user not found")
	}
	return helper.ToProfileResponse(user), nil
}

func (service *profileService) UpdateProfile(c context.Context, claims helper.JWTClaims, request web.UpdateProfileRequest) (web.ProfileResponse, error) {
	if err := service.Validate.Struct(request); err != nil {
		return web.ProfileResponse{}, exception.ErrBadRequest(err.Error())
	}

	user, err := service.ProfileRepository.GetProfileByID(c, claims.Profile.ID)
	if err != nil {
		return web.ProfileResponse{}, exception.ErrNotFound(err.Error())
	}

	if request.Name != "" {
		user.Name = request.Name
	}

	if request.Email != "" {
		if userByEmail, _ := service.ProfileRepository.GetProfilesByQuery(c, "email", request.Email); userByEmail.ID != "" && userByEmail.ID != claims.Profile.ID {
			exception.ErrBadRequest("email already registered")
		}
		user.Email = request.Email
	}

	if request.Phone != "" {
		if !helper.IsNumeric(request.Phone) {
			panic(exception.ErrBadRequest("Phone should numeric"))
		}

		if len([]rune(request.Phone)) < 10 || len([]rune(request.Phone)) > 13 {
			panic(exception.ErrBadRequest("Phone should 10-13 digit"))
		}

		if userByPhone, _ := service.ProfileRepository.GetProfilesByQuery(c, "phone", request.Phone); userByPhone.ID != "" && userByPhone.ID != claims.Profile.ID {
			exception.ErrBadRequest("phone already registered")
		}
		user.Phone = request.Phone
	}

	user.UpdatedAt = time.Now()

	if err := service.ProfileRepository.UpdateProfile(c, user); err != nil {
		return web.ProfileResponse{}, exception.ErrInternalServer(err.Error())
	}

	// KAFKA

	user, _ = service.ProfileRepository.GetProfileByID(c, claims.Profile.ID)
	return helper.ToProfileResponse(user), nil
}
