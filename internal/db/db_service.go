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

	notesTable := `CREATE TABLE IF NOT EXISTS notes (
		id SERIAL PRIMARY KEY,
		title VARCHAR NOT NULL,
		content TEXT
	);`

	_, err := db.ExecContext(ctx, notesTable)

	return err
}
