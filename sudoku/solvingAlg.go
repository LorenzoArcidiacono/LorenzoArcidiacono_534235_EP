package sudoku

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"../analyzetime"
	"../setting"
	"./algorithms/dlx"
)

//Solve cerca la soluzione con entrambi gli algoritmi e ne calcola il tempo di esecuzione
func Solve(g *Grid) (bool, string, string) {
	outputBt, err := os.OpenFile(setting.OutPathBT, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	outputDlx, err := os.OpenFile(setting.OutPathDLX, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)

	defer outputBt.Close()
	defer outputDlx.Close()

	if err != nil {
		fmt.Println(err)
		return false, "", ""
	}
	start := time.Now() //setta il tempo di inizio dell'elaborazione
	check, sol1 := SolveDLX(g)
	analyzetime.Track(start, "DLX", outputDlx) //stampa il tempo impiegato
	if check == false {
		return false, "", ""
	}
	start = time.Now()
	check = SolveBT(g)
	analyzetime.Track(start, "BT", outputBt)
	return check, sol1, g.toString()

}

//SolveDLX cerca la soluzione tramite l'algoritmo Dancing Links
func SolveDLX(g *Grid) (bool, string) {
	var str string
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			val := int(g[i][j].val)
			str += strconv.Itoa(val)
		}
	}
	solver, err := dlx.InitSolver(str)
	if err != nil || !solver.Solve() {
		fmt.Println("Couldn't solve given sudoku: ", err)
		return false, ""
	}
	return true, solver.GetSolution()
}

//SolveBT cerca la soluzione tramite l'algoritmo di Backtracking
func SolveBT(g *Grid) bool {
	i, j := findUnassignedLoc(*g) //cerca una locazione libera
	if i == -1 && j == -1 {
		return true
	}

	values := randValues()
	for index := range values {
		err := g.checkAll(i, j, values[index])
		if err == nil {
			g[i][j].val = values[index]
			if SolveBT(g) { //cerca di risolverlo tramite ricorsione
				return true
			}
			g[i][j].val = empty //se non ci riesce cambia il valore
		}
	}
	return false
}

//randValues restituisce un array con una permutazione di tutti i valori da 1 a maxValue
func randValues() [maxValue]int8 {
	var values [maxValue]int8
	rand.Seed(time.Now().UTC().UnixNano())
	perm := rand.Perm(maxValue)
	for i := range perm {
		perm[i]++
		values[i] = int8(perm[i])
	}
	return values
}

//findUnassignedLoc restituisce una coppia che indica una posizione vuota nella griglia
func findUnassignedLoc(g Grid) (int, int) {
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			if g[i][j].val == empty {
				return i, j
			}
		}
	}
	return -1, -1
}
