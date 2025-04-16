package repository

import (
	"fmt"

	"github.com/Bethakin/project1/internal/database"
	"github.com/Bethakin/project1/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *database.Database) *UserRepository {
	return &UserRepository{
		db: db.DB,
	}
}

func (r *UserRepository) GetAll() ([]*model.Todousers, error) {
	rows, err := r.db.Query("SELECT id, email, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.Todousers
	for rows.Next() {
		var user model.Todousers
		if err := rows.Scan(&user.ID, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *UserRepository) (id int) (*model.Todousers, error) {
	var user model.Todousers
	err := r.db.QueryRow("SELECT id, email, password FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *model.Todousers) error {
	return r.db.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id",
		user.Email, user.Password,
	).Scan(&user.ID)
}

func (r *UserRepository) Update(id int, user *model.Todousers) error {
	result, err := r.db.Exec(
		"UPDATE users SET email = $1, password = $2 WHERE id = $3",
		user.Email, user.Password, id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}
	return nil
}

func (r *UserRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}
	return nil
}

func (r *UserRepository) DeleteUserTodos(userID int) error {
	_, err := r.db.Exec("DELETE FROM todos WHERE user_id = $1", userID)
	if err != nil {
		return err
	}
	return nil
}
