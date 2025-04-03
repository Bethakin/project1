package database

import (
	"fmt"
	"log"
	"os"

	"github.com/Bethakin/project1/model"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Database struct {
	DB *sqlx.DB
}

func NewDatabase() (*Database, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSLMODE")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, sslMode)

	log.Println("Connecting to database with:", connStr)

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

func (d *Database) GetAllTodos() ([]*model.Todo, error) {
	rows, err := d.DB.Query("SELECT id, title, description FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*model.Todo
	for rows.Next() {
		var todo model.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description); err != nil {
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

func (d *Database) GetTodosByUserID(id int) ([]*model.Todo, error) {
	rows, err := d.DB.Query("SELECT id, title, description FROM todos WHERE user_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*model.Todo
	for rows.Next() {
		var todo model.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	if len(todos) == 0 {
		return nil, fmt.Errorf("no todos found for user with ID %d", id)
	}

	return todos, nil
}

func (d *Database) GetUserByID(id int) (*model.Todousers, error) {
	var user model.Todousers
	err := d.DB.QueryRow("SELECT id, email, password FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *Database) CreateUser(user *model.Todousers) error {
	return d.DB.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id",
		user.Email, user.Password,
	).Scan(&user.ID)
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

func (d *Database) DeleteUser(id int) error {
	result, err := d.DB.Exec("DELETE FROM users WHERE id = $1", id)
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

func (d *Database) DeleteUserTodos(user_id int) error {
	_, err := d.DB.Exec("DELETE FROM todos WHERE user_id = $1", user_id)
	if err != nil {
		return err
	}
	/*
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return fmt.Errorf("user todos with user_id %d not found", user_id)
		}*/
	return nil
}

func (d *Database) UpdateUser(id int, user *model.Todousers) error {
	result, err := d.DB.Exec(
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

func (d *Database) UpdateTodo(id int, user_id int, todo *model.Todo) error {
	result, err := d.DB.Exec(
		"UPDATE todos SET title = $1, description = $2 WHERE id = $3 AND user_id = $4",
		todo.Title, todo.Description, id, user_id,
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
