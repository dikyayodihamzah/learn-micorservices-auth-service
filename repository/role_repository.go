package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"gitlab.com/learn-micorservices/auth-service/model/domain"
)

type RoleRepository interface {
	GetRoleByID(c context.Context, role_id string) domain.Role
}

type roleRepository struct {
	Database func(dbName string) *pgx.Conn
}

func NewRoleRepository(database func(dbName string) *pgx.Conn) RoleRepository {
	return &roleRepository{
		Database: database,
	}
}

func (repository *roleRepository) GetRoleByID(c context.Context, role_id string) domain.Role {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := `SELECT * FROM roles WHERE id = $1"`

	roles := db.QueryRow(ctx, query, role_id)

	var role domain.Role
	roles.Scan(&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt)

	return role
}
