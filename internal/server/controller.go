package server

import (
	"errors"
	"log"
	"net/http"
	"project9/internal/db"
)

func AddNoteHandler(w http.ResponseWriter, r *http.Request) {
	var note db.Note

	err := db.ParseToNote(r.Body, &note)
	if errors.Is(err, db.ErrInvalidFormatJson) {
		log.Printf("%s %s - %v", r.Method, r.URL.Path, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		log.Printf("%s %s - %v", r.Method, r.URL.Path, err)
		http.Error(w, "ошибка обработки данных", http.StatusInternalServerError)
		return
	}

	err = db.CreateNote(note.Title, note.Content)
	if err != nil {
		log.Printf("%s %s - %v", r.Method, r.URL.Path, err)
		http.Error(w, "ошибка обработки данных", http.StatusInternalServerError)
	}
}

func GetNoteHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	note, err := db.GetNoteById(id)
	if err != nil {
		log.Printf("%s %s - %v", r.Method, r.URL.Path, err)
		http.Error(w, "ошибка обработки данных", http.StatusInternalServerError)
		return
	}

	jsonByte, err := db.ParseToJson(&note)
	if err != nil {
		log.Printf("%s %s - %v", r.Method, r.URL.Path, err)
		http.Error(w, "ошибка обработки данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonByte)
}
