package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var ErrInvalidFormatJson = errors.New(`поля "title" и "content" являются обязательными`)

type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type lastId struct {
	ID int `json:"id"`
}

func ParseToNote(inputJson io.Reader, noteStruct *Note) error {

	decoder := json.NewDecoder(inputJson)
	err := decoder.Decode(noteStruct)
	if err != nil {
		return fmt.Errorf("ParseToNote: %w", err)
	}

	if noteStruct.Title == "" || noteStruct.Content == "" {
		return ErrInvalidFormatJson
	}

	return nil

}

func ParseToJson(inputNote *Note) ([]byte, error) {

	jsonByte, err := json.Marshal(inputNote)
	if err != nil {
		return nil, fmt.Errorf("ParseToJson: %w", err)
	}

	return jsonByte, nil
}
