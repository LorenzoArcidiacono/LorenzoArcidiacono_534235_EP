package game

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"../../sudoku"
	"../../terminal"
)

//Game avvia una singola partita
func Game(g *sudoku.Grid) {
	var str string
	exit := false
	rule()
	for !exit {
		g.ShowScheme(os.Stdout)
		fmt.Print("\n»  ")
		fmt.Scanf("%s", &str)
		str = strings.ToLower(str) //si aspetta di leggere una tripla i,j,v dove v->(i,j)
		sub := strings.Split(str, ",")
		_, err := strconv.Atoi(sub[0]) //prova a trasformare i singoli valori in int
		if err == nil {                //se non c'è stato un errore prova a inserire v in (i,j)
			insert(sub, g)
		} else { //altrimenti controlla la stringa
			switch str {
			case "soluzione": //è stata richiesta la soluzione
				solve(g)
				exit = true
			case "exit": //si vuole uscire dal gioco
				terminal.Clear()
				exit = true
			case "regole": //si chiede di stampare a schermo le regole
				terminal.Clear()
				rule()
			default:
				terminal.Clear()
				fmt.Println("Errore nella selezione")
			}
		}
		if sudoku.IsSolved(*g) && !exit { //prima di ricominciare il ciclo controlla che il sudoku non sia risolto
			fmt.Println("Hai risolto correttamente il sudoku!")
			exit = true
		}
	}
	g.ShowScheme(os.Stdout)
}

//rule stampa a schermo la spiegazione delle possibili scelte
func rule() {
	fmt.Println("------- REGOLE -------")
	fmt.Println("- Scrivendo: soluzione -> viene stampata la soluzione del sudoku.")
	fmt.Println("- Scrivendo: exit -> viene chiuso lo schema attuale.")
	fmt.Println("- Scrivendo: regole -> vengono stampate le regole del gioco.")
	fmt.Println("- Scrivendo una i,j,v -> prova ad inserire v nella cella (i,j).")
	fmt.Println("- La cifra 0 indica la casella vuota.")
	fmt.Println("- Il gioco non permette di scrivere una cifra in una casella se è già presente nel relativo: blocco, riga o colonna.")
	fmt.Println("- Le caselle vanno da (1,1) a (9,9).")
}

//insert riceve una slice con tre valori numerici che indicano la posizione e il valore da inserire
func insert(str []string, g *sudoku.Grid) {
	r, err := strconv.Atoi(str[0]) //indice di riga
	if err != nil {
		panic(err)
	}
	c, err := strconv.Atoi(str[1]) //indice di colonna
	if err != nil {
		panic(err)
	}
	v, err := strconv.Atoi(str[2]) //valore da inserire
	if err != nil {
		panic(err)
	}
	err = g.Set(r-1, c-1, int8(v)) //vengono passati r-1 e c-1 poichè la griglia è indirizzata a partire da 0
	terminal.Clear()               // mentre i valori passati da 1
	if err != nil {
		fmt.Println(err)
	}
}

//solve scrive nella griglia passata la soluzione del sudoku
func solve(g *sudoku.Grid) {
	check, _, _ := sudoku.Solve(g)
	if !check {
		fmt.Println("Soluzione non trovata")
	}
}
