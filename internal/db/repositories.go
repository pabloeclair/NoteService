package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var (
	ErrInvalidFormatJson           = errors.New(`поля "title" и "content" являются обязательными`)
	ErrInvalifFotmatJsonPutRequset = errors.New(`одно из полей "title" или "content" должно обязательно присутствовать`)
)

type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type NotesList struct {
	SearchResult []Note `json:"search_result"`
}

type LastId struct {
	LastId int `json:"id"`
}

func ParseToNote(inputJson io.Reader, noteStruct *Note) error {

	decoder := json.NewDecoder(inputJson)
	err := decoder.Decode(noteStruct)
	if err != nil {
		return fmt.Errorf("ParseToNote: %w", err)
	}

	return nil

}

func ParseToJson(inputNote any) ([]byte, error) {

	jsonByte, err := json.Marshal(inputNote)
	if err != nil {
		return nil, fmt.Errorf("ParseToJson: %w", err)
	}

	return jsonByte, nil
}
