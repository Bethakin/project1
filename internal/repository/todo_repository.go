package repository

import (
	"fmt"

	"github.com/Bethakin/project1/internal/database"
	"github.com/Bethakin/project1/model"
	"github.com/jmoiron/sqlx"
)

type TodoRepository struct {
	db *sqlx.DB
}

func NewTodoRepository(db *database.Database) *TodoRepository {
	return &TodoRepository{
		db: db.DB,
	}
}

func (r *TodoRepository) GetAll() ([]*model.Todo, error) {
	rows, err := r.db.Query("SELECT id, title, description, user_id FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*model.Todo
	for rows.Next() {
		var todo model.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.UsersID); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}
	return todos, nil
}

func (r *TodoRepository) GetByID(id int) (*model.Todo, error) {
	var todo model.Todo
	err := r.db.QueryRow("SELECT id, title, description, user_id FROM todos WHERE id = $1", id).
		Scan(&todo.ID, &todo.Title, &todo.Description, &todo.UsersID)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *TodoRepository) GetByUserID(userID int) ([]*model.Todo, error) {
	rows, err := r.db.Query("SELECT id, title, description, user_id FROM todos WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*model.Todo
	for rows.Next() {
		var todo model.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.UsersID); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	if len(todos) == 0 {
		return nil, fmt.Errorf("no todos found for user with ID %d", userID)
	}

	return todos, nil
}

func (r *TodoRepository) Create(todo *model.Todo) error {
	return r.db.QueryRow(
		"INSERT INTO todos (title, description, user_id) VALUES ($1, $2, $3) RETURNING id",
		todo.Title, todo.Description, todo.UsersID,
	).Scan(&todo.ID)
}

func (r *TodoRepository) Update(id int, userID int, todo *model.Todo) error {
	result, err := r.db.Exec(
		"UPDATE todos SET title = $1, description = $2 WHERE id = $3 AND user_id = $4",
		todo.Title, todo.Description, id, userID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo with id %d not found", id)
	}
	return nil
}

func (r *TodoRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo with id %d not found", id)
	}
	return nil
}
