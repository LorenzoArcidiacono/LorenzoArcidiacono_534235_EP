package sdkerror

import "errors"

var (
	ErrBuonds     = errors.New("Fuori dai limiti")
	ErrValue      = errors.New("Valore errato")
	ErrCol        = errors.New("Cifra già presente in questa colonna")
	ErrRow        = errors.New("Cifra già presente in questa riga")
	ErrRegion     = errors.New("Cifra già presente in questo settore")
	ErrFixed      = errors.New("Casella non modificabile")
	ErrDifficulty = errors.New("Livello di difficoltà sbagliato")
)
