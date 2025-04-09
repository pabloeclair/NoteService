package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
)

func CreateNoteTable() error {

	dsn := os.Getenv("DSN")

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	db := sqlx.MustConnect("pgx", dsn)
	defer db.Close()

	notesTable := `CREATE TABLE IF NOT EXISTS notes (
		id SERIAL PRIMARY KEY,
		title VARCHAR NOT NULL,
		content TEXT
	);`

	_, err := db.ExecContext(ctx, notesTable)
	if err != nil {
		err = fmt.Errorf("CreateNoteTable: %w", err)
	}

	return err
}

func CreateNote(title string, content string) (int, error) {

	var lastId lastId

	query := `INSERT INTO notes (title, content) VALUES ($1, $2);`

	dsn := os.Getenv("DSN")
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()
	db := sqlx.MustConnect("pgx", dsn)
	defer db.Close()

	_, err := db.ExecContext(ctx, query, title, content)
	if err != nil {
		return -1, fmt.Errorf("CreateNote: %w", err)
	}

	err = db.GetContext(ctx, &lastId, "SELECT id FROM notes ORDER BY id DESC LIMIT 1")
	if err != nil {
		return -1, fmt.Errorf("CreateNote: %w", err)
	}

	return lastId.ID, nil

}

func GetNoteById(id string) (Note, error) {

	var note Note
	query := `SELECT * FROM notes WHERE id = $1`

	dsn := os.Getenv("DSN")
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()
	db := sqlx.MustConnect("pgx", dsn)
	defer db.Close()

	err := db.GetContext(ctx, &note, query, id)
	if err != nil {
		return note, fmt.Errorf("GetNoteById: %w", err)
	}

	return note, nil
}
