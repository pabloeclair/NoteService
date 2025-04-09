package db

import (
	"encoding/json"
	"errors"
	"io"
)

var ErrInvalidFormatJson = errors.New(`поля "title" и "content" являются обязательными`)

type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func ParseToNote(inputJson io.Reader, noteStruct *Note) error {

	decoder := json.NewDecoder(inputJson)
	err := decoder.Decode(noteStruct)
	if err != nil {
		return err
	}

	if noteStruct.Title == "" || noteStruct.Content == "" {
		return ErrInvalidFormatJson
	}

	return nil

}

func ParseToJson(inputNote *Note) ([]byte, error) {

	jsonByte, err := json.Marshal(inputNote)
	if err != nil {
		return nil, err
	}

	return jsonByte, nil
}
