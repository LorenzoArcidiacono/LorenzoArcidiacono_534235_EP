package main

import (
	"fmt"
	"strings"

	"./menu"
	"./test"
)

func main() {
	var sel string
	fmt.Print("Vuoi giocare o eseguire un test? [giocare/test]:")
	fmt.Scanf("%v", &sel)
	sel = strings.ToLower(sel)
	switch sel {
	case "giocare":
		menu.Start()
	case "test":
		test.Start()
	default:
		fmt.Println("Selenzione non valida.")
	}

}
