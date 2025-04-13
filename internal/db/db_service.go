package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
)

func CreateNoteTable() error {

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	dsn := os.Getenv("DSN")
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

func CreateNote(title string, content string, ctx context.Context) (LastId, error) {
	var id int
	var idStruct LastId

	if title == "" || content == "" {
		return idStruct, ErrInvalidFormatJson
	}

	dsn := os.Getenv("DSN")
	db := sqlx.MustConnect("pgx", dsn)
	defer db.Close()

	_, err := db.ExecContext(ctx, `INSERT INTO notes (title, content) VALUES ($1, $2);`, title, content)
	if err != nil {
		return idStruct, fmt.Errorf("CreateNote: %w", err)
	}

	err = db.GetContext(ctx, &id, "SELECT id FROM notes ORDER BY id DESC LIMIT 1")
	if err != nil {
		return idStruct, fmt.Errorf("CreateNote: %w", err)
	}

	idStruct.LastId = id

	return idStruct, nil

}

func GetNoteById(id string, ctx context.Context) (Note, error) {
	var note Note

	dsn := os.Getenv("DSN")
	db := sqlx.MustConnect("pgx", dsn)
	defer db.Close()

	err := db.GetContext(ctx, &note, `SELECT * FROM notes WHERE id = $1`, id)
	if err != nil {
		return note, fmt.Errorf("GetNoteById: %w", err)
	}

	return note, nil
}

func UpdateNote(id string, title string, content string, ctx context.Context) (Note, error) {
	var note Note

	dsn := os.Getenv("DSN")
	db := sqlx.MustConnect("pgx", dsn)
	defer db.Close()

	if content == "" && title == "" {
		return note, ErrInvalifFotmatJsonPutRequset
	}

	if title != "" {
		_, err := db.ExecContext(ctx, `UPDATE notes SET title = $1 WHERE id = $2`, title, id)
		if err != nil {
			return note, fmt.Errorf("UpdateNote: %w", err)
		}
	}

	if content != "" {
		_, err := db.ExecContext(ctx, `UPDATE notes SET content = $1 WHERE id = $2`, content, id)
		if err != nil {
			return note, fmt.Errorf("UpdateNote: %w", err)
		}
	}

	err := db.GetContext(ctx, &note, `SELECT * FROM notes WHERE id = $1`, id)
	if err != nil {
		return note, fmt.Errorf("UpdateNote: %w", err)
	}
	return note, nil
}

func DropNote(id string, ctx context.Context) error {
	var check int

	dsn := os.Getenv("DSN")
	db := sqlx.MustConnect("pgx", dsn)
	defer db.Close()

	err := db.GetContext(ctx, &check, `SELECT id FROM notes WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("DropNote: %w", err)
	}

	_, err = db.ExecContext(ctx, `DELETE FROM notes WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("DropNote: %w", err)
	}

	return nil
}

func GetNotesByContent(filter string, ctx context.Context) (NotesList, error) {
	var notes []Note
	var notesList NotesList

	dsn := os.Getenv("DSN")
	db := sqlx.MustConnect("pgx", dsn)
	defer db.Close()

	err := db.SelectContext(ctx, &notes, `SELECT * FROM notes WHERE content LIKE '%' || $1 || '%'`, filter)
	notesList.SearchResult = notes
	if err != nil {
		return notesList, fmt.Errorf("GetNotesByContent: %w", err)
	}
	return notesList, nil
}
