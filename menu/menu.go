package menu

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"../parse"
	"../sudoku"
	"./game"
)

var errExit = errors.New("Chiusura del gioco")

//recovery controlla se è stato lanciato un errore e nel caso lo stampa a schermo
func recovery() {
	if e := recover(); e != nil {
		fmt.Println(e)
	}
	os.Exit(1)
}

//getValue legge e restituisce un numero da tastiera
func getValue() (int, error) {
	var i int

	fmt.Print("Enter number: ")
	_, err := fmt.Scanf("%d", &i) //legge un numero da tastiera
	if err != nil {
		return -1, err
	}
	return i, err
}

//getValueString legge una stringa di 9 numeri separati da virgole da tastiera
//restituisce un array dei singoli numeri
func getValueString() ([9]int, error) {
	var val [9]int
	var str string
	_, err := fmt.Scanf("%v", &str) //legge una stringa da tastiera
	if err != nil {
		return val, err
	}
	single := strings.Split(str, ",") //crea una slice contenente solo i numeri
	for ind := range single {
		val[ind], err = strconv.Atoi(single[ind]) //converti i numeri in int
		if err != nil {
			return val, err
		}
	}
	return val, nil
}

//menu stampa le possibili scelte che l'utente può fare
func menu() (*sudoku.Grid, error) {
	var sel string
	var err error
	var grid *sudoku.Grid
	exit := false
	for exit == false {
		fmt.Println("\n------- MENU -------")
		fmt.Println("Come si desidera generare la griglia? [file/tastiera/random] dgigitare 'exit' per terminare, 'regole' per leggere le regole del gioco, 'help' per maggiori informazioni")
		fmt.Print("Selezionare l'opzione:")
		fmt.Scanf("%v", &sel)
		sel = strings.ToLower(sel)
		switch sel {
		case "file":
			grid, err = generateFromFile()
			if grid != nil {
				return grid, nil
			}
		case "tastiera":
			grid, err = generateFromInput()
			if grid != nil {
				return grid, nil
			}
		case "random":
			grid, err = generateRandom()
			if grid != nil {
				return grid, nil
			}
		case "help":
			help()
		case "regole":
			rule()
		case "exit":
			exit = true
		default:
			fmt.Println("Errore nella selezione, scegliere tra 'file/tastiera/random'")
		}
		if err != nil {
			fmt.Println(err)
		}
	}
	return grid, errExit
}

//help stampa a schermo la spiegazione delle varie funzioni
func help() {
	fmt.Println("\n---------- HELP -----------")
	fmt.Println("- Selezionare l'opzione desiderata scrivendo uno tra: file, tastiera, random")
	fmt.Println("- L'opzione file richiede la posizione del file, il carattere usato per le righe di commento, e il carattere usato come casella vuota.")
	fmt.Println("	In base a queste impostazioni legge la griglia dal file (se il file contiene più di una griglia viene letta solo la prima)")
	fmt.Println("- L'opzione tastiera richiede di scrivere manualmente una riga alla volta tutta la griglia.")
	fmt.Println("- L'opzione random richiede di scegliere un livello di difficoltà e genera una griglia in base alla scelta.")
}

//rule stampa a schermo le regole del sudoku
func rule() {
	fmt.Println("------ REGOLE SUDOKU ------")
	fmt.Println("- Scopo del gioco è riempire tutta la griglia.")
	fmt.Println("- Su ogni riga devono esserci tutti i valori da 1 a 9 senza ripetizioni")
	fmt.Println("- Su ogni colonna devono esserci tutti i valori da 1 a 9 senza ripetizioni")
	fmt.Println("- In ogni blocco devono esserci tutti i valori da 1 a 9 senza ripetizioni")
}

//generateFromFile legge una griglia da un file
func generateFromFile() (*sudoku.Grid, error) {
	var path, delim, void string
	fmt.Print("Indicare il path del file da cui leggere lo schema: ")
	fmt.Scanf("%s", &path)
	fmt.Print("Indicare il carattere che precede le righe da non leggere: ")
	fmt.Scanf("%s", &delim)
	fmt.Print("Indicare il valore che indica il numero mancante (di solito 0): ")
	fmt.Scanf("%s", &void)

	red, err := parse.Parse(path, delim)
	if err != nil {
		return nil, err
	}
	grid, err := sudoku.CreateFromString(red, void)
	if err != nil {
		return nil, err
	}
	return grid, nil
}

//generateFromInput permette all'utente di scrivere la griglia una riga alla volta
func generateFromInput() (*sudoku.Grid, error) {
	var values [9][9]int
	var err error
	fmt.Println("Scrivere una riga alla volta separando i valori con una virgola. Indicare le caselle vuote con 0")
	for i := 0; i < 9; i++ {
		values[i], err = getValueString()
		if err != nil {
			return nil, err
		}
	}
	grid := sudoku.CreateFromArray(values)
	return grid, nil
}

//generateRandom genera una griglia nuova
func generateRandom() (*sudoku.Grid, error) {
	var lev string
	fmt.Println("Selezionare un livello tra: facile, medio, difficile, arduo , impossibile")
	fmt.Scanf("%s", &lev)
	grid, err := sudoku.Generate(lev)
	if err != nil {
		return nil, err
	}
	return grid, nil
}

//Start avvia il gioco e stampa, all' inizio di ogni sessione di gioco, il menu
func Start() {
	defer recovery() //chiama la funzione recovery al momento dell'uscita dalla funzione Start

	exit := false
	for !exit {
		grid, err := menu()
		if err == errExit { //nel menu viene selezionata l'opzione di uscita
			os.Exit(0)
		}
		if err != nil {
			panic(err)
		}

		game.Game(grid)
	}
}
