package main

import "fmt"

func main() {
	/*
		- Slices are similar to arrays, but are more powerful and flexiable
		- Unlike arrays, the length of a slice can grow and shrink as needed
	*/

	var slice1 = []int{1, 2, 3, 4, 5}
	slice2 := []int{1, 2, 3, 4, 5}

	fmt.Println(slice1)
	fmt.Println(slice2)

	slice2 = append(slice2, 6)
	fmt.Println(slice2)
}
