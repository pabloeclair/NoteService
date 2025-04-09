package db

import (
	"context"
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

	return err
}

func CreateNote(title string, content string) error {

	query := `INSERT INTO notes (title, content) VALUES ($1, $2);`

	dsn := os.Getenv("DSN")
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()
	db := sqlx.MustConnect("pgx", dsn)
	defer db.Close()

	_, err := db.QueryContext(ctx, query, title, content)

	return err
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
		return note, err
	}

	return note, nil
}
