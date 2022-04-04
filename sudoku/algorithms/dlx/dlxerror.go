package dlx

// Autore: Jakub Rozanski
// Progetto originale: https://github.com/golangchallenge/GCSolutions/tree/master/nov15/normal/jakub-rozanski/jrozansk-go-challenge8-d4d18058ef2c

import "errors"

var (
	InputTooShortError  = errors.New("Input is too short!")
	InputMalformedError = errors.New("Input is malformed!")
)
