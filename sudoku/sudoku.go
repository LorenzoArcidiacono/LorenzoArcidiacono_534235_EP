package sudoku

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"./sdkerror"
)

//Costanti di inizializzazione dello schema.
const (
	rows     = 9
	columns  = 9
	empty    = 0
	maxValue = 9
)

//Tipo che descrive il livello di difficoltà.
type difficulty int

//Enumerazione dei livelli di difficoltà e degli indizi relativi.
const (
	easy       difficulty = 31
	medium     difficulty = 26
	hard       difficulty = 22
	veryHard   difficulty = 20
	impossible difficulty = 17
)

//Cell descrive una singola casella della griglia.
type cell struct {
	val   int8 //numero nella cella; int8 occupa meno memoria.
	fixed bool //false se il valore è modificabile.
}

//Grid descrive una griglia composta da rows*columns cell (esportato).
type Grid [rows][columns]cell

//setDifficulty restituisce il livello di difficoltà in base alla stringa passata
func setDifficulty(str string) (difficulty, error) {
	str = strings.ToLower(str)
	fmt.Println("letto:", str)
	switch str {
	case "facile":
		return easy, nil
	case "medio":
		return medium, nil
	case "difficile":
		return hard, nil
	case "arduo":
		return veryHard, nil
	case "impossibile":
		return impossible, nil
	default:
		return -1, sdkerror.ErrDifficulty
	}
}

//CreateByEnum crea una griglia a partire dall'enumerazione completa scelta dall'utente (esportata)
func CreateByEnum(digits [rows][columns]int8) *Grid {
	var g Grid
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			d := digits[i][j]
			if d != empty {
				g[i][j].val = d
				g[i][j].fixed = true
			}
		}
	}
	return &g
}

//CreateFromString restituisce una griglia creata a partire da una singola strina (esportata).
func CreateFromString(digits []string, void string) (*Grid, error) {
	var g Grid
	count := 0
	for index := range digits {
		if count >= rows {
			break
		}
		single := strings.Split(digits[index], "")
		for item := range single {
			if single[item] != void {
				i := index % 9
				value, err := strconv.Atoi(single[item])
				if err != nil {
					return nil, err
				}
				g[i][item].val = int8(value)
				if value != empty {
					g[i][item].fixed = true
				}
			}
		}
		count++
	}
	return &g, nil
}

//CreateFromArray restituisce una griglia creata a partire da una matrice (esportata).
func CreateFromArray(digit [9][9]int) *Grid {
	var g Grid
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			if digit[i][j] != empty {
				g[i][j].val = int8(digit[i][j])
				g[i][j].fixed = true
			}
		}
	}
	return &g
}

//Generate restituisce una griglia creata randomicamente (esportata).
func Generate(lev string) (*Grid, error) {
	level, err := setDifficulty(lev)
	if err != nil {
		return nil, err
	}
	var tab Grid

	SolveBT(&tab)                 //risolve una griglia vuota
	randomStart(&tab, int(level)) //imposta alcune caselle a zero in base alla difficoltà scelta
	tab.setUnmodifiable()         //setta le caselle che non possono essere modificate
	return &tab, nil
}

//randomStart elimina dalla griglia i valori di alcune celle in modo da rimanere solo con clues celle impostate.
func randomStart(g *Grid, clues int) {
	rand.Seed(time.Now().UTC().UnixNano())
	for i := rows * columns; i > clues; i-- {
		elim := rand.Intn(maxValue) + 1
		r, c := g.find(int8(elim))
		if r != -1 && c != -1 {
			g.clear(r, c)
		} else {
			i++
		}
	}
}

//find trova il primo valore == digit a partire da una riga generata randomicamente
func (g *Grid) find(digit int8) (int, int) {
	tryAll := false
	count := 0
	rand.Seed(time.Now().UTC().UnixNano())
	i := rand.Intn(rows)
	for !tryAll { //se lo ho già cancellato dalla riga random vado a quella successiva
		for j := 0; j < columns; j++ {
			if g[i][j].val == digit {
				return i, j
			}
		}
		count++
		i = (i + 1) % rows
		if count == rows-1 { //se digit non appare nella griglia
			tryAll = true
		}
	}
	return -1, -1
}

//setUnmodifiable setta il campo fixed di ogni cella dove val != empty a true.
func (g *Grid) setUnmodifiable() {
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			if g[i][j].val != empty {
				g[i][j].fixed = true
			}
		}
	}
}

//inBound controlla che la coppia row,column appartenga alla griglia.
func inBound(row, column int) bool {
	if row < 0 || row >= rows {
		return false
	}

	if column < 0 || column >= columns {
		return false
	}
	return true
}

//validDigit controlla che il valore passato come argomento sia accettabile
func validDigit(digit int8) bool {
	return digit >= 1 && digit <= 9
}

func (g *Grid) isFixed(row, column int) bool {
	return g[row][column].fixed
}

//inRow restituisce true se nella riga esiste già una cella con val == digit, false altrimenti
func (g *Grid) inRow(row int, digit int8) bool {
	for c := 0; c < columns; c++ {
		if g[row][c].val == digit {
			return true
		}
	}
	return false
}

//inColumn restituisce true se nella colonna esiste già una cella con val == digit, false altrimenti
func (g *Grid) inColumn(column int, digit int8) bool {
	for r := 0; r < rows; r++ {
		if g[r][column].val == digit {
			return true
		}
	}
	return false
}

//inRegion restituisce true se nella regione esiste già una cella con val == digit, false altrimenti
func (g *Grid) inRegion(row, column int, digit int8) bool {
	startRow, startCol := (row/3)*3, (column/3)*3
	for r := startRow; r < startRow+3; r++ {
		for c := startCol; c < startCol+3; c++ {
			if g[r][c].val == digit {
				return true
			}
		}
	}
	return false
}

//checkAll controlla tutte le condizioni
func (g Grid) checkAll(row, column int, digit int8) error {
	switch {
	case !inBound(row, column):
		return sdkerror.ErrBuonds
	case !validDigit(digit):
		return sdkerror.ErrValue
	case g.isFixed(row, column):
		return sdkerror.ErrFixed
	case g.inRow(row, digit):
		return sdkerror.ErrRow
	case g.inColumn(column, digit):
		return sdkerror.ErrCol
	case g.inRegion(row, column, digit):
		return sdkerror.ErrRegion
	}
	return nil
}

//Set se le condizioni sono verificate setta la cella (row,column) a digit (esportata).
func (g *Grid) Set(row, column int, digit int8) error {
	err := g.checkAll(row, column, digit)
	if err != nil {
		return err
	}
	g[row][column].val = digit
	return nil
}

//clear setta a empty la cella (row,column).
func (g *Grid) clear(row, column int) error {
	switch {
	case !inBound(row, column):
		return sdkerror.ErrBuonds
	case g.isFixed(row, column):
		return sdkerror.ErrFixed
	}

	g[row][column].val = empty
	return nil
}

//Show stampa una versione semplice della griglia (esportata).
func (g Grid) Show() {
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			fmt.Printf("%v ", g[i][j].val)
		}
		fmt.Print("\n")
	}
	fmt.Print("\n\n")
}

//ShowScheme stampa una versione elaborata della griglia (esportata).
func (g Grid) ShowScheme(f *os.File) {
	fmt.Print("\n--------SUDOKU------\n\n")
	for i := 0; i < rows; i++ {
		if i == 3 || i == 6 {
			fmt.Fprint(f, "---|---|---\n")
		}
		for j := 0; j < columns; j++ {
			if j == 3 || j == 6 {
				fmt.Fprint(f, "|")
			}
			fmt.Fprint(f, g[i][j].val)
		}
		fmt.Fprint(f, "\n")
	}
}

//toString restituisce la griglia sotto forma di stringa.
func (g Grid) toString() string {
	var str string
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			val := int(g[i][j].val)
			str += strconv.Itoa(val)
		}
	}
	return str
}

//GetVal restituisce il valore della cella (i,j)
func (g Grid) GetVal(i, j int) int {
	return int(g[i][j].val)
}

//IsSolved restituisce true se tutte le celle sono state riempite
func IsSolved(g Grid) bool {
	i, j := findUnassignedLoc(g)
	if i == -1 && j == -1 {
		return true
	}
	return false
}
