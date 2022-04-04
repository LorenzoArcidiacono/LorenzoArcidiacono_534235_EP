package analyzetime

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"time"
)

//Track stampa sul file i secondi passati da start
func Track(start time.Time, str string, f io.Writer) {
	t1 := time.Now()
	fmt.Fprintf(f, "%.4v\n", t1.Sub(start).Seconds())
}

//Average restituisce la media tra i valori letti dal file relativo a path
func Average(path string) (float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return -1, err
	}
	file.Seek(0, 0)
	defer file.Close()

	var total, count float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		red := scanner.Text()
		last, err := strconv.ParseFloat(red, 64)
		if err != nil {
			return -1, err
		}
		if err == nil {
			count++
			total += last
		}
	}
	if err := scanner.Err(); err != nil {
		return -1, err
	}
	return total / count, nil
}

//Max restituisce il valore più grande tra i numeri letti dal file relativo a path
func Max(path string) (float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return -1, err
	}
	file.Seek(0, 0)
	defer file.Close()

	var max, last float64
	scanner := bufio.NewScanner(file)
	max = math.Inf(-1)
	if err != nil {
		return -1, err
	}
	for scanner.Scan() {
		red := scanner.Text()
		last, err = strconv.ParseFloat(red, 64)
		if err != nil {
			return -1, err
		}
		if last > max {
			max = last
		}
	}
	if err := scanner.Err(); err != nil {
		return -1, err
	}
	return max, nil
}

//Min restituisce il valore più piccolo tra i numeri letti dal file relativo a path
func Min(path string) (float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return -1, err
	}
	file.Seek(0, 0)
	defer file.Close()

	var min, last float64
	scanner := bufio.NewScanner(file)
	min = math.Inf(1)
	for scanner.Scan() {
		red := scanner.Text()
		last, err = strconv.ParseFloat(red, 64)
		if err != nil {
			return -1, err
		}
		if min > last {
			min = last
		}
	}
	if err := scanner.Err(); err != nil {
		return -1, err
	}
	return min, nil
}
