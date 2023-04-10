package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"gitlab.com/learn-micorservices/auth-service/model/domain"
)

type RoleRepository interface {
	GetRoleByID(c context.Context, role_id string) domain.Role

	//kafka
	Create(c context.Context, role domain.Role) error
	Update(c context.Context, id string, role domain.Role) error
	Delete(c context.Context, id string) error
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

// function for kafka
func (repository *roleRepository) Create(c context.Context, role domain.Role) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := `INSERT INTO roles (
		id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4)`

	if _, err := db.Prepare(c, "create_role", query); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, "create_role", role.ID, role.Name, role.CreatedAt, role.UpdatedAt); err != nil {
		return err
	}

	return nil
}

func (repository *roleRepository) Update(c context.Context, id string, role domain.Role) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := `UPDATE roles SET
		name = $1,
		created_at = $2,
		updated_at = $3
		WHERE id = $4`

	if _, err := db.Prepare(c, "update_role", query); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, "update_role", role.Name, role.CreatedAt, role.UpdatedAt, role.ID); err != nil {
		return err
	}

	return nil
}

func (repository *roleRepository) Delete(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := `DELETE FROM roles WHERE id = $1`

	if _, err := db.Prepare(c, "data", query); err != nil {
		return err
	}

	if _, err := db.Exec(c, "data", id); err != nil {
		return err
	}

	return nil
}
