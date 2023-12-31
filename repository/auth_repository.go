package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"gitlab.com/learn-micorservices/auth-service/exception"
	"gitlab.com/learn-micorservices/auth-service/model/domain"
)

var dbName = os.Getenv("DB_NAME")

type AuthRepository interface {
	Register(c context.Context, user domain.User) error
	GetUsersByQuery(c context.Context, params, value string) (domain.User, error)
	UpdatePassword(c context.Context, user domain.User) error

	CreateToken(c context.Context, tokens domain.ResetPasswordToken) error
	CheckToken(c context.Context, token string) (domain.ResetPasswordToken, error)
	CheckTokenWithQuery(c context.Context, params, value string) (domain.ResetPasswordToken, error)
	DeleteToken(c context.Context, token string) error

	//kafka
	UpdateUser(c context.Context, user domain.User) error
	DeleteUser(c context.Context, user_id string) error
}

type authRepository struct {
	Database func(dbName string) *pgx.Conn
}

func NewAuthRepository(database func(dbName string) *pgx.Conn) AuthRepository {
	return &authRepository{
		Database: database,
	}
}

func (repository *authRepository) Register(c context.Context, user domain.User) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := `INSERT INTO users (
		id,
		name,
		username,
		email,
		password,
		phone,
		role_id,
		created_at,
		updated_at
	)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	if _, err := db.Prepare(ctx, "data", query); err != nil {
		return exception.ErrInternalServer(err.Error())
	}

	if _, err := db.Exec(ctx, "data",
		user.ID,
		user.Name,
		user.Username,
		user.Email,
		user.Password,
		user.Phone,
		user.RoleID,
		user.CreatedAt,
		user.UpdatedAt); err != nil {
		return exception.ErrUnprocessableEntity(err.Error())
	}

	return nil
}

func (repository *authRepository) GetUsersByQuery(c context.Context, params, value string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := fmt.Sprintf(`
		SELECT users.*, roles.name 
		FROM users 
		LEFT JOIN roles ON roles.id = users.role_id
		WHERE %s = $1`, params)

	user := db.QueryRow(ctx, query, value)

	var data domain.User
	user.Scan(&data.ID, &data.Name, &data.Username, &data.Email, &data.Password, &data.Phone, &data.RoleID, &data.CreatedAt, &data.UpdatedAt, &data.RoleName)

	return data, nil
}

func (repository *authRepository) UpdatePassword(c context.Context, user domain.User) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := "UPDATE users SET password = $1 WHERE id = $2"

	if _, err := db.Prepare(c, "data", query); err != nil {
		return exception.ErrInternalServer(err.Error())
	}

	if _, err := db.Exec(c, "data", user.Password, user.ID); err != nil {
		return exception.ErrInternalServer(err.Error())
	}

	return nil
}

func (repository *authRepository) CreateToken(c context.Context, tokens domain.ResetPasswordToken) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := "INSERT INTO reset_token (tokens, email, created_at) VALUES($1,$2,$3)"

	if _, err := db.Prepare(ctx, "data", query); err != nil {
		return exception.ErrInternalServer(err.Error())
	}

	if _, err := db.Exec(ctx, "data", tokens.Tokens, tokens.Email, tokens.CreatedAt); err != nil {
		return exception.ErrInternalServer(err.Error())
	}

	return nil
}

func (repository *authRepository) CheckToken(c context.Context, token string) (domain.ResetPasswordToken, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := fmt.Sprintf("SELECT * FROM reset_token WHERE tokens = '%s';", token)
	user := db.QueryRow(ctx, query)

	data := new(domain.ResetPasswordToken)

	if err := user.Scan(data.Tokens, data.Email, data.CreatedAt); err != nil {
		return domain.ResetPasswordToken{}, err
	}

	return *data, nil
}

func (repository *authRepository) CheckTokenWithQuery(c context.Context, params, value string) (domain.ResetPasswordToken, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := fmt.Sprintf("SELECT * FROM reset_token WHERE %s = '%s';", params, value)
	user := db.QueryRow(ctx, query)

	data := new(domain.ResetPasswordToken)

	if err := user.Scan(data.Tokens, data.Email, data.CreatedAt); err != nil {
		return domain.ResetPasswordToken{}, err
	}

	return *data, nil
}

func (repository *authRepository) DeleteToken(c context.Context, token string) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := `DELETE FROM reset_token WHERE tokens = $1`

	if _, err := db.Prepare(ctx, "data", query); err != nil {
		return exception.ErrInternalServer(err.Error())
	}

	if _, err := db.Exec(ctx, "data", token); err != nil {
		return exception.ErrInternalServer(err.Error())
	}

	return nil
}

// function for kafka
func (repository *authRepository) UpdateUser(c context.Context, user domain.User) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := `UPDATE users SET 
		name = $1, 
		username = $2, 
		email = $3, 
		password = $4, 
		phone = $5, 
		role_id = $6, 
		updated_at = $7
		WHERE id = $8`

	if _, err := db.Prepare(ctx, "data", query); err != nil {
		return exception.ErrInternalServer(err.Error())
	}

	if _, err := db.Exec(ctx, "data",
		user.Name,
		user.Username,
		user.Email,
		user.Password,
		user.Phone,
		user.RoleID,
		user.UpdatedAt,
		user.ID); err != nil {
		return exception.ErrUnprocessableEntity(err.Error())
	}

	return nil
}

func (repository *authRepository) DeleteUser(c context.Context, user_id string) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := `DELETE FROM users WHERE id = $1`

	if _, err := db.Prepare(ctx, "data", query); err != nil {
		return exception.ErrInternalServer(err.Error())
	}

	if _, err := db.Exec(ctx, "data", user_id); err != nil {
		return exception.ErrUnprocessableEntity(err.Error())
	}

	return nil
}
