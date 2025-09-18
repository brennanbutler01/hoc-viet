package vocabulary

import "time"

type Word struct {
	Word        string    `json:"word"`
	Translation string    `json:"translation"`
	DateCreated time.Time `json:"dateCreated"`
}

type AddWordRequest struct {
	Word        string `json:"word" doc:"The word to add"`
	Translation string `json:"translation" doc:"The translation of the word"`
}

type AddWordResponse struct {
	Body Word
}
