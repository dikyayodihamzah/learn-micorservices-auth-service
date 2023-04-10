package service

import (
	"context"
	"encoding/base64"
	"log"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/thanhpk/randstr"
	"gitlab.com/learn-micorservices/auth-service/exception"
	"gitlab.com/learn-micorservices/auth-service/helper"
	"gitlab.com/learn-micorservices/auth-service/model/domain"
	"gitlab.com/learn-micorservices/auth-service/model/web"
	"gitlab.com/learn-micorservices/auth-service/repository"
)

type AuthService interface {
	Register(c context.Context, request web.RegisterRequest) (web.RegisterResponse, error)
	Login(c context.Context, request web.LoginRequest) (fiber.Cookie, web.LoginResponse, error)
	Logout(c context.Context) fiber.Cookie
	ForgetPassword(c context.Context, email string) error
	ResetPassword(c context.Context, email, token string, request web.ResetPassword) error
	CheckToken(c context.Context, token, email string) error
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
		Name:      request.Name,
		Username:  request.Username,
		Email:     request.Email,
		Password:  request.Password,
		Phone:     request.Phone,
		RoleID:    request.RoleID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	user.GenerateID()
	user.SetPassword(request.Password)

	if err := service.AuthRepository.Register(c, user); err != nil {
		return web.RegisterResponse{}, err
	}

	user, err := service.AuthRepository.GetUsersByQuery(c, "id", user.ID)
	if err != nil {
		return web.RegisterResponse{}, err
	}

	// KAFKA
	helper.ProduceToKafka(user, "POST.USER", helper.KafkaTopic)

	return helper.ToRegisterResponse(user), nil
}

func (service *authService) Login(c context.Context, request web.LoginRequest) (fiber.Cookie, web.LoginResponse, error) {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	if request.Username == "" && request.Email == "" && request.Phone == "" {
		return fiber.Cookie{}, web.LoginResponse{}, exception.ErrBadRequest("username, email, or phone are missing")
	}

	var user domain.User

	// Login using username
	if request.Username != "" {
		user, err = service.AuthRepository.GetUsersByQuery(c, "username", request.Username)
		if err != nil || user.ID == "" {
			return fiber.Cookie{}, web.LoginResponse{}, exception.ErrNotFound("user not found")
		}
	}

	// Login using email
	if request.Email != "" {
		user, err = service.AuthRepository.GetUsersByQuery(c, "email", request.Email)
		if err != nil || user.ID == "" {
			return fiber.Cookie{}, web.LoginResponse{}, exception.ErrNotFound("user not found")
		}
	}

	// Login using phone
	if request.Phone != "" {
		user, err = service.AuthRepository.GetUsersByQuery(c, "phone", request.Phone)
		if err != nil || user.ID == "" {
			return fiber.Cookie{}, web.LoginResponse{}, exception.ErrNotFound("user not found")
		}
	}

	if err := user.ComparePassword(user.Password, request.Password); err != nil {
		return fiber.Cookie{}, web.LoginResponse{}, exception.ErrBadRequest("wrong password")
	}

	claims := helper.UserClaimsData{
		ID:     user.ID,
		RoleID: user.RoleID,
	}

	token, err := helper.GenerateJWT(user.ID, claims)
	helper.PanicIfError(err)

	cookie := fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	response := web.LoginResponse{
		ID:     user.ID,
		RoleID: user.RoleID,
	}

	log.Println(token)
	return cookie, response, nil
}

func (service *authService) Logout(c context.Context) fiber.Cookie {
	cookie := fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	return cookie
}

func (service *authService) ForgetPassword(c context.Context, email string) error {
	user, err := service.AuthRepository.GetUsersByQuery(c, "email", email)

	if err != nil || user.Email == "" {
		return exception.ErrBadRequest("email not match for the record")
	}

	tokens, _ := service.AuthRepository.CheckTokenWithQuery(c, "email", email)
	if tokens.Tokens != "" {
		service.AuthRepository.DeleteToken(c, tokens.Tokens)
	}

	token := strings.ToLower(randstr.String(30))

	data := domain.ResetPasswordToken{
		Tokens:    token,
		Email:     email,
		CreatedAt: time.Now(),
	}

	if err := service.AuthRepository.CreateToken(c, data); err != nil {
		return err
	}

	if err := helper.EmailSender(email, token); err != nil {
		return fiber.NewError(fiber.StatusBadGateway, "request error")
	}

	return nil
}

func (service *authService) ResetPassword(c context.Context, email, token string, request web.ResetPassword) error {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	var decodedByte, _ = base64.StdEncoding.DecodeString(token)
	var resetToken = string(decodedByte)

	if request.Password != request.PasswordConfirm {
		panic(exception.ErrBadRequest("password didn't match"))
	}

	checkToken, err := service.AuthRepository.CheckToken(c, resetToken)

	if err != nil {
		return exception.ErrBadRequest("token invalid")
	}

	if checkToken.Tokens != resetToken || checkToken.Email != email {
		return exception.ErrBadRequest("token invalid")
	}

	var user domain.User

	if user, err = service.AuthRepository.GetUsersByQuery(c, "email", email); err != nil {
		return exception.ErrNotFound(err.Error())
	}

	user.SetPassword(request.Password)
	user.UpdatedAt = time.Now()

	service.AuthRepository.UpdatePassword(c, user)

	// KAFKA
	helper.ProduceToKafka(user, "PUT.USER", helper.KafkaTopic)

	err = service.AuthRepository.DeleteToken(c, resetToken)
	if err != nil {
		return exception.ErrInternalServer(err.Error())
	}
	return nil
}

func (service *authService) CheckToken(c context.Context, token, email string) error {

	var decodedByte, _ = base64.StdEncoding.DecodeString(token)
	var resetToken = string(decodedByte)

	checkToken, err := service.AuthRepository.CheckToken(c, resetToken)

	if err != nil {
		return exception.ErrBadRequest("token invalid")
	}

	if checkToken.Tokens != resetToken || checkToken.Email != email {
		return exception.ErrBadRequest("token invalid")
	}
	return nil
}
