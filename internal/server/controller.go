package server

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"project9/internal/db"
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

	id, err := db.CreateNote(note.Title, note.Content)
	if errors.Is(err, db.ErrInvalidFormatJson) {
		log.Printf("%s %s - 400 Bad Request: %v", r.Method, r.URL.Path, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		log.Printf("%s %s - 500 Internal Server Error: %v", r.Method, r.URL.Path, err)
		http.Error(w, "ошибка обработки данных", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(201)
	w.Write(fmt.Append(nil, id))
}

func GetNoteHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	note, err := db.GetNoteById(id)
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

	resNote, err := db.UpdateNote(id, note.Title, note.Content)
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
