package dlx

// Autore: Jakub Rozanski
// Progetto originale: https://github.com/golangchallenge/GCSolutions/tree/master/nov15/normal/jakub-rozanski/jrozansk-go-challenge8-d4d18058ef2c

//Node struttura che descrive un singolo nodo
type Node struct {
	Left  *Node
	Right *Node
	Up    *Node
	Down  *Node

	Head *Node

	Col int
	Row int
}

//nodeInit inizializza un nodo
func nodeInit(row int, col int) *Node {
	n := Node{}
	n.Left = &n
	n.Right = &n
	n.Up = &n
	n.Down = &n
	n.Row = row
	n.Col = col
	return &n
}

//NodeRegular inizializza un nodo qualsiasi
func NodeRegular(row int, col int) *Node {
	return nodeInit(row, col)
}

//NodeHeader inizializza un nodo header
func NodeHeader(col int) *Node {
	return nodeInit(0, col)
}

//Cover elimina il nodo dal grafo
func (n *Node) Cover() {
	if n.Head != nil {
		n.Up.Down = n.Down
		n.Down.Up = n.Up
		return
	}
	n.Left.Right = n.Right
	n.Right.Left = n.Left
}
