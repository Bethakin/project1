package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Bethakin/project1/internal/model"
	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase() (*Database, error) {
	connStr := "user=postgres password=27Berfin72 dbname=users sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	log.Println("Successfully connected to database")
	return &Database{DB: db}, nil
}

func (d *Database) GetAllTodos() ([]*model.TodoUser, error) {
	rows, err := d.DB.Query("SELECT id, email, password FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*model.TodoUser
	for rows.Next() {
		var todo model.TodoUser
		if err := rows.Scan(&todo.ID, &todo.Email, &todo.Password); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}
	return todos, nil
}

func (d *Database) GetTodoByID(id int) (*model.TodoUser, error) {
	var todo model.TodoUser
	err := d.DB.QueryRow("SELECT id, email, password FROM todos WHERE id = $1", id).
		Scan(&todo.ID, &todo.Email, &todo.Password)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (d *Database) CreateTodo(todo *model.TodoUser) error {
	return d.DB.QueryRow(
		"INSERT INTO todos (email, password) VALUES ($1, $2) RETURNING id",
		todo.Email, todo.Password,
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

func (d *Database) UpdateTodo(id int, todo *model.TodoUser) error {
	result, err := d.DB.Exec(
		"UPDATE todos SET email = $1, password = $2 WHERE id = $3",
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
		return fmt.Errorf("todo with id %d not found", id)
	}
	return nil
}
