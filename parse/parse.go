package parse

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

//ErrPath errore nel caso il path passato sia una stringa vuota
var ErrPath = errors.New("Il path passato non è valido")

//Parse restituisce un array di stringhe lette dal file per riga
//ha come argomenti il path del file da cui leggere e una stringa in cui scrivere la stringa con cui iniziano le frasi da saltare
func Parse(path string, delim string) ([]string, error) {
	var red []string
	if path == "" {
		return nil, ErrPath
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		str := scanner.Text()
		if strings.Index(str, delim) != 0 && str != "" {
			red = append(red, str)
		}
	}
	file.Close()
	return red, nil
}
