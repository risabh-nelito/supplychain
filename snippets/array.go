package main

import (
	"fmt"
)

func main() {
	a := []string{"A", "B", "C", "D", "E"}
	i := 2

	a[i] = a[len(a)-1] // Copy last element to index i
	a[len(a)-1] = ""   // Erase last element (write zero value)
	a = a[:len(a)-1]   // Truncate slice

	fmt.Println(a)
}
