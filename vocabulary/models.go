package vocabulary

import "time"

type Word struct {
	Word        string    `json:"word"`
	Translation string    `json:"translation"`
	DateCreated time.Time `json:"dateCreated"`
}

type AddWordRequest struct {
	Word        string `json:"word" minLength:"1" maxLength:"200" doc:"The word to add"`
	Translation string `json:"translation" minLength:"1" maxLength:"500" doc:"The translation of the word"`
}

type AddWordResponse struct {
	Body Word
}
