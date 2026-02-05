package main

import (
	"fmt"
)

func main() {
	var arr1 = [5]int{1, 2, 3, 4, 5}
	arr2 := [3]int{1, 2, 3}
	var arr3 = [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	cars := [4]string{"Volvo", "BMW", "Ford", "Mazda"}

	fmt.Println(arr1)
	fmt.Println(arr2)
	fmt.Println(arr3)
	fmt.Println(cars)

	// Index start from 0
	fmt.Println(cars[0])
	fmt.Println(cars[1])
	fmt.Println(cars[2])
	fmt.Println(cars[3])
}
