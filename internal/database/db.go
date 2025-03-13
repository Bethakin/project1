package database

import (
	_ "database/sql"
	"fmt"
	"log"

	"github.com/Bethakin/project1/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	DB *sqlx.DB
}

func NewDatabase() (*Database, error) {
	connStr := "user=postgres password=27Berfin72 dbname=users sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	log.Println("Successfully connected to database")
	return &Database{DB: db}, nil
}

func (d *Database) GetAlluserss() ([]*model.Todousers, error) {
	rows, err := d.DB.Query("SELECT id, email, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*model.Todousers
	for rows.Next() {
		var todo model.Todousers
		if err := rows.Scan(&todo.ID, &todo.Email, &todo.Password); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}
	return todos, nil
}

func (d *Database) GetAllTodos() ([]*model.Todousers, error) {
	rows, err := d.DB.Query("SELECT id, title, description FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*model.Todousers
	for rows.Next() {
		var todo model.Todousers
		if err := rows.Scan(&todo.ID, &todo.Email, &todo.Password); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}
	return todos, nil
}

func (d *Database) GetTodoByID(id int) (*model.Todo, error) {
	var todo model.Todo
	err := d.DB.QueryRow("SELECT id, title, description FROM todos WHERE id = $1", id).
		Scan(&todo.ID, &todo.Title, &todo.Description)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (d *Database) GetusersByID(id int) (*model.Todousers, error) {
	var todo model.Todousers
	err := d.DB.QueryRow("SELECT id, email, password FROM users WHERE id = $1", id).
		Scan(&todo.ID, &todo.Email, &todo.Password)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (d *Database) Createusers(todo *model.Todousers) error {
	return d.DB.QueryRow(
		"INSERT INTO userss (email, password) VALUES ($1, $2) RETURNING id",
		todo.Email, todo.Password,
	).Scan(&todo.ID)
}

func (d *Database) CreateTodo(todo *model.Todo) error {
	return d.DB.QueryRow(
		"INSERT INTO todos (title, description, user_id) VALUES ($1, $2, $3) RETURNING id",
		todo.Title, todo.Description, todo.UsersID,
	).Scan(&todo.ID)
}

func (d *Database) DeleteTodo(id int) error {
	result, err := d.DB.Exec("DELETE FROM todos WHERE id = $1", id)
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

func (d *Database) DeleteTheusers(id int) error {
	result, err := d.DB.Exec("DELETE FROM users WHERE id = $1", id)
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

func (d *Database) DeleteuserssTodo(user_id int) error {
	result, err := d.DB.Exec("DELETE FROM todos WHERE user_id = $1", user_id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo with id %d not found", user_id)
	}
	return nil
}

func (d *Database) Updateusers(id int, todo *model.Todousers) error {
	result, err := d.DB.Exec(
		"UPDATE users SET email = $1, password = $2 WHERE id = $3",
		todo.Email, todo.Password, id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("users with id %d not found", id)
	}
	return nil
}

func (d *Database) UpdateTodo(id int, users_id int, todo *model.Todo) error {
	result, err := d.DB.Exec(
		"UPDATE todos SET Title = $1, Description = $2 WHERE id = $3 and user_id = $4",
		todo.Title, todo.Description, id, users_id,
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
