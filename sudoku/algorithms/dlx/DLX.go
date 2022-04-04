package dlx

// Autore: Jakub Rozanski
// Progetto originale: https://github.com/golangchallenge/GCSolutions/tree/master/nov15/normal/jakub-rozanski/jrozansk-go-challenge8-d4d18058ef2c

type Solver struct {
	solution []int
	header   *Node
}

//InitSolver genera la matrice da usare per trovare la soluzione
func InitSolver(input string) (*Solver, error) {
	if len(input) != NumOfCells {
		return nil, InputTooShortError
	}
	solver := Solver{header: generateDlxMatrix()}
	fixedCells := asciiToInts(input)
	if err := solver.addToSolution(fixedCells); err != nil {
		return nil, err
	}
	return &solver, nil
}

//Solve cerca una soluzione ricororsivamente applicando l'algoritmo 'Dancing Links'
func (s *Solver) Solve() bool {
	col := s.header.Right
	if s.header.Right == s.header { //tutte le possibilità già controllate
		return true
	}
	if col.Down == col {
		return false
	}
	cover(col)
	for row := col.Down; row != col; row = row.Down {
		for cell := row.Right; cell != row; cell = cell.Right {
			cover(cell.Head)
		}
		if s.Solve() { //cerca ricorsimante nell nuovo header più a destra
			s.solution = append(s.solution, row.Row) //se trova una soluzione la salva
			return true
		}
		for cell := row.Left; cell != row; cell = cell.Left {
			uncover(cell.Head) //se non trova una soluzione fa la uncover dell'header delle righe
		}
	}
	uncover(col) //se non trova una soluzione a partire da questa colonna fa la uncover della colonna
	return false
}

//restituisce la soluzione del problema sotto forma di stringa
func (s *Solver) GetSolution() string {
	result := make([]byte, NumOfCells)
	for _, v := range s.solution {
		row, col, val := decodePosibility(v)
		result[row*GridSize+col] = intToAscii(val) + 1
	}
	return string(result)
}

func (s *Solver) addToSolution(fixedCells []int) error {
	for i, val := range fixedCells {
		if val == 0 {
			continue
		}

		row := i / GridSize
		col := i % GridSize

		constraints := constraintPositions(row, col, val-1)
		for _, columnIdx := range constraints {
			col := getColumnHeader(s.header, columnIdx)
			if col == nil {
				return InputMalformedError
			}
			cover(col) //fa la cover di tutta la colonna a partire dal'header
		}
		s.solution = append(s.solution, encodePossibility(row, col, val-1))
	}
	return nil
}

func cover(col *Node) {
	col.Cover() //cover dell'header
	coverRows(col)
}

func uncover(header *Node) {
	uncoverRows(header)
	uncoverCol(header)
}

func coverRows(header *Node) {
	for ptr := header.Down; ptr != header; ptr = ptr.Down { //cover di tutta la colonna
		for cell := ptr.Right; cell != ptr; cell = cell.Right { //cover della riga associata
			cell.Cover()
		}
	}
}

func uncoverRows(header *Node) {
	for ptr := header.Up; ptr != header; ptr = ptr.Up {
		for cell := ptr.Left; cell != ptr; cell = cell.Left {
			uncoverRow(cell)
		}
	}
}

func uncoverCol(header *Node) {
	header.Left.Right = header
	header.Right.Left = header
}

func uncoverRow(cell *Node) {
	cell.Up.Down = cell
	cell.Down.Up = cell
}
