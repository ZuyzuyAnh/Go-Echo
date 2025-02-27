package repository

import (
	"context"
	"echo-demo/internal/model"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type RoleRepository struct {
	DB *sqlx.DB
}

func NewRoleRepository(db *sqlx.DB) *RoleRepository {
	return &RoleRepository{
		DB: db,
	}
}

func (r *RoleRepository) getRoleByUserId(userID int64) (*model.Role, error) {
	query := `
		SELECT r.id, r.name, r.description 
		FROM roles r JOIN userroles ur ON r.id = userroles.role_id
		WHERE ur.user_id = $1 
	`

	var role model.Role

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, userID).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
	)

	return &role, err
}

func (r *RoleRepository) InsertUserRole(exec Executor, userId int64, roleID int64) error {
	query := `
		INSERT INTO userroles (user_id, role_id)
		VALUES ($1, $2)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := exec.ExecContext(ctx, query, userId, roleID)
	if err != nil {
		if table, column, ok := isForeignKeyError(err); ok {
			return &ForeignKeyError{Message: fmt.Sprintf("No such %s with %s", table, column)}
		}

		return err
	}

	return nil
}
