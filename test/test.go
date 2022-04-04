package test

import (
	"fmt"
	"log"
	"os"
	"strings"

	"../analyzetime"
	"../parse"
	"../setting"
	"../sudoku"
)

//Start avvia tutti i test e permette di aggiungerne di nuovi
func Start() {
	os.Remove(setting.OutPathBT)
	os.Remove(setting.OutPathDLX)
	os.Remove(setting.OutPathClues)
	var sel string
	fmt.Print("Vuoi eseguire i test standard o uno nuovo? [standard/nuovo]:")
	_, err := fmt.Scanf("%v", &sel)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	sel = strings.ToLower(sel)
	switch sel {
	case "standard":
		for _, path := range setting.InputTest {
			testByLine(path, "#", ".")
			if len(setting.InputTest) == 1 {
				avgBt, err := analyzetime.Average(setting.OutPathBT)
				maxBt, err := analyzetime.Max(setting.OutPathBT)
				minBt, err := analyzetime.Min(setting.OutPathBT)
				avgDlx, err := analyzetime.Average(setting.OutPathDLX)
				maxDlx, err := analyzetime.Max(setting.OutPathDLX)
				minDlx, err := analyzetime.Min(setting.OutPathDLX)
				if err != nil {
					fmt.Println(err)
					log.Fatal(err)
				}
				fmt.Printf("%v backtracking: average %.4f, max value %.4f,min value %.4f\n", path, avgBt, maxBt, minBt)
				fmt.Printf("%v dancing links: average %.4f, max value %.4f,min value %.4f\n", path, avgDlx, maxDlx, minDlx)
			} else {
				f, err := os.OpenFile(setting.OutPathBT, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					log.Fatal(err)
				}
				f2, err := os.OpenFile(setting.OutPathDLX, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					log.Fatal(err)
				}
				fc, err := os.OpenFile(setting.OutPathClues, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					log.Fatal(err)
				}
				_, err = fmt.Fprintf(f, "-------------------------\n")
				if err != nil {
					log.Fatal(err)
				}
				_, err = fmt.Fprintf(f2, "-------------------------\n")
				if err != nil {
					log.Fatal(err)
				}
				_, err = fmt.Fprintf(fc, "-------------------------\n")
				if err != nil {
					log.Fatal(err)
				}

				f.Close()
				f2.Close()
				fc.Close()
			}
		}

	case "nuovo":
		var pathin, pathout, delim, void string
		fmt.Print("Indicare il path del file da cui leggere i test: ")
		fmt.Scanf("%v", &pathin)
		fmt.Print("Indicare il path del file da cui leggere le soluzioni: ")
		fmt.Scanf("%v", &pathout)
		fmt.Print("Indicare il carattere che precede le righe da non leggere (solitamente #): ")
		fmt.Scanf("%v", &delim)
		fmt.Print("Indicare il valore che indica il numero mancante (di solito 0): ")
		fmt.Scanf("%v", &void)
		test(pathin, pathout, delim, void)
	default:
		fmt.Println("Errore nella selezione")
	}
}

//test legge i file di input ed output e esegue i test
func test(pathin, pathout, delim, void string) {
	gridIn := readTest(pathin, delim)
	gridOut := readTest(pathout, delim)
	for g := range gridIn {
		fmt.Println("test #:", g)
		test, err := sudoku.CreateFromString(gridIn[g], void)
		if err != nil {
			panic(err)
		}
		check, s1, s2 := sudoku.Solve(test)
		if !check {
			fmt.Println("Schema non risolto")
		}

		if strings.Compare(s1, s2) != 0 {
			fmt.Println("Due soluzioni trovate")
			fmt.Printf("%v\n%v\n", s1, s2)
		}
		compare(gridOut[g], s1)
	}
}

//readTest restituisce tutte le griglie lette dal file
func readTest(path string, delim string) [][]string {
	allGrid, err := parse.Parse(path, delim)
	var grid [][]string
	var single []string
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for row := range allGrid {
		single = append(single, allGrid[row])
		if (row+1)%9 == 0 && row != 0 {
			grid = append(grid, single)
			single = nil
		}
	}
	return grid
}

//testByLine legge ed esegue tutte le griglie del file scritte come singole righe
func testByLine(pathin, delim, void string) {
	gridIn := readTestByLine(pathin, delim)
	for g := range gridIn {
		fmt.Printf("test #:%d ", g)
		test, err := sudoku.CreateFromString(gridIn[g], void)
		countClues(gridIn[g])
		if err != nil {
			panic(err)
		}
		check, s1, s2 := sudoku.Solve(test)
		if !check {
			fmt.Println("Schema non risolto")
		}

		if strings.Compare(s1, s2) != 0 {
			fmt.Println("Due soluzioni trovate")
			fmt.Printf("%v\n%v\n", s1, s2)
		}
		fmt.Println("risolto correttamente: ", s1)
	}
}

//readTestByLine legge e restituisce tutte le griglie di un file scritte come singole righe
func readTestByLine(path string, delim string) [][]string {
	allGrid, err := parse.Parse(path, delim)
	var grid [][]string
	var single []string
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for row := range allGrid {
		str := strings.Split(allGrid[row], "")
		s := ""
		for i := range str {
			if i%9 == 0 && i != 0 {
				single = append(single, s)
				s = ""
			}
			s += str[i]
		}
		single = append(single, s) //append dell'ultima riga
		grid = append(grid, single)
		single = nil
	}
	return grid
}

//compare controllora se una griglia su più righe coincide con una su una singola riga
func compare(sol []string, str string) {
	var output string
	for i := range sol {
		output += sol[i]
	}
	if strings.Compare(output, str) == 0 {
		fmt.Println("Soluzione corretta")
	} else {
		fmt.Println("Soluzione errata")
		fmt.Println(str, "\n", output)
	}

}

//countClues conta il numero di indizi iniziali
func countClues(str []string) {
	c, _ := os.OpenFile(setting.OutPathClues, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	var empty int
	for i := range str {
		empty += strings.Count(str[i], ".")
	}
	fmt.Fprintf(c, "%d\n", 81-empty)
}
