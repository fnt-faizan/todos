package db

import (
	"context"
	"database/sql"
	"todos/models"
)

func GetAllTodos(ctx context.Context, db *sql.DB) ([]*models.Todo, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, title, status, deleted FROM todos WHERE deleted = false ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*models.Todo
	for rows.Next() {
		t := new(models.Todo)
		if err := rows.Scan(&t.Id, &t.Title, &t.Status, &t.Deleted); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}

func GetTodoByID(ctx context.Context, db *sql.DB, id int) (*models.Todo, error) {
	t := new(models.Todo)
	err := db.QueryRowContext(ctx,
		"SELECT title, status, deleted FROM todos WHERE id = $1 AND deleted = false", id).
		Scan(&t.Title, &t.Status, &t.Deleted)
	if err != nil {
		return nil, err
	}
	t.Id = id
	return t, nil
}

func InsertTodo(ctx context.Context, db *sql.DB, t *models.Todo) error {
	return db.QueryRowContext(ctx,
		"INSERT INTO todos (title, status, deleted) VALUES ($1, $2, $3) RETURNING id",
		t.Title, t.Status, t.Deleted).Scan(&t.Id)
}

func UpdateTodo(ctx context.Context, db *sql.DB, t *models.Todo) (bool, error) {
	res, err := db.ExecContext(ctx,
		"UPDATE todos SET title = $1, status = $2, deleted = $3 WHERE id = $4 and deleted = false",
		t.Title, t.Status, t.Deleted, t.Id)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows == 1, nil
}

func DeleteTodo(ctx context.Context, db *sql.DB, id int) (bool, error) {
	res, err := db.ExecContext(ctx,
		"UPDATE todos SET deleted = true WHERE id = $1 and deleted = false", id)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows == 1, nil
}
