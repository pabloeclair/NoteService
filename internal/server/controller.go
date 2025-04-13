package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"project9/internal/db"
	"strconv"
	"time"
)

type loggingReponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingReponseWriter {
	return &loggingReponseWriter{w, http.StatusOK}
}

func (lrw *loggingReponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		lrw := NewLoggingResponseWriter(w)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
		defer cancel()
		r = r.WithContext(ctx)
		next.ServeHTTP(lrw, r)

		if lrw.statusCode == http.StatusOK {
			log.Printf("%s %s - 200 OK", r.Method, r.URL.Path)
		} else if lrw.statusCode == http.StatusCreated {
			log.Printf("%s %s - 201 Created", r.Method, r.URL.Path)
		}
	})
}

func AddNoteHandler(w http.ResponseWriter, r *http.Request) {
	var note db.Note

	err := db.ParseToNote(r.Body, &note)
	if err != nil {
		log.Printf("%s %s - 400 Bad Request: %v", r.Method, r.URL.Path, err)
		http.Error(w, "неверный тип json", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	idStruct, err := db.CreateNote(note.Title, note.Content, ctx)
	if errors.Is(err, db.ErrInvalidFormatJson) {
		log.Printf("%s %s - 400 Bad Request: %v", r.Method, r.URL.Path, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		log.Printf("%s %s - 500 Internal Server Error: %v", r.Method, r.URL.Path, err)
		http.Error(w, "ошибка обработки данных", http.StatusInternalServerError)
		return
	}

	idByte, err := json.Marshal(&idStruct)
	if err != nil {
		log.Printf("%s %s - 500 Internal Server Error: %v", r.Method, r.URL.Path, err)
		http.Error(w, "ошибка обработки данных", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(201)
	w.Write(idByte)

}

func GetNoteHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	_, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("%s %s - 400 Bad Request: в строку пути необходимо ввести число – id заметки", r.Method, r.URL.Path)
		http.Error(w, "в строку пути необходимо ввести число – id заметки", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	note, err := db.GetNoteById(id, ctx)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("%s %s - 404 Not Found: %v", r.Method, r.URL.Path, fmt.Sprintf("заметки с id = %s не существует", id))
		http.Error(w, fmt.Sprintf("заметки с id = %s не существует", id), http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("%s %s - 500 Internal Server Error: %v", r.Method, r.URL.Path, err)
		http.Error(w, "ошибка обработки данных", http.StatusInternalServerError)
		return
	}

	jsonByte, err := db.ParseToJson(&note)
	if err != nil {
		log.Printf("%s %s - 500 Internal Server Error: %v", r.Method, r.URL.Path, err)
		http.Error(w, "ошибка обработки данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonByte)
}

func PutNoteHandler(w http.ResponseWriter, r *http.Request) {
	var note db.Note

	id := r.PathValue("id")
	err := db.ParseToNote(r.Body, &note)
	if err != nil {
		log.Printf("%s %s - 400 Bad Request: %v", r.Method, r.URL.Path, err)
		http.Error(w, "неверный тип json", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	resNote, err := db.UpdateNote(id, note.Title, note.Content, ctx)
	if errors.Is(err, db.ErrInvalifFotmatJsonPutRequset) {
		log.Printf("%s %s - 400 Bad Request: %v", r.Method, r.URL.Path, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if errors.Is(err, sql.ErrNoRows) {
		log.Printf("%s %s - 404 Not Found: %v", r.Method, r.URL.Path, fmt.Sprintf("заметки с id = %s не существует", id))
		http.Error(w, fmt.Sprintf("заметки с id = %s не существует", id), http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("%s %s - 500 Internal Server Error: %v", r.Method, r.URL.Path, err)
		http.Error(w, "ошибка обработки данных", http.StatusInternalServerError)
		return
	}

	resJson, err := db.ParseToJson(&resNote)
	if err != nil {
		log.Printf("%s %s - 500 Internal Server Error: %v", r.Method, r.URL.Path, err)
		http.Error(w, "ошибка обработки данных", http.StatusInternalServerError)
		return
	}
	w.Write(resJson)
}

func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	ctx := r.Context()
	err := db.DropNote(id, ctx)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("%s %s - 404 Not Found: %v", r.Method, r.URL.Path, fmt.Sprintf("заметки с id = %s не существует", id))
		http.Error(w, fmt.Sprintf("заметки с id = %s не существует", id), http.StatusNotFound)
	} else if err != nil {
		log.Printf("%s %s - 500 Internal Server Error: %v", r.Method, r.URL.Path, err)
		http.Error(w, "ошибка обработки данных", http.StatusInternalServerError)
	}
}

func GetSearchNotesHandlerv(w http.ResponseWriter, r *http.Request) {

	filter := r.URL.Query().Get("q")
	ctx := r.Context()
	notes, err := db.GetNotesByContent(filter, ctx)
	if err != nil {
		log.Printf("%s %s - 500 Internal Server Error: %v", r.Method, r.URL.Path, err)
		http.Error(w, "ошибка обработки данных", http.StatusInternalServerError)
		return
	}

	notesByte, err := db.ParseToJson(&notes)
	if err != nil {
		log.Printf("%s %s - 500 Internal Server Error: %v", r.Method, r.URL.Path, err)
		http.Error(w, "ошибка обработки данных", http.StatusInternalServerError)
		return
	}

	w.Write(notesByte)
}
