package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func main() {
	fmt.Println("Start main function")

	var Cuong Person
	Cuong.Age = 27
	Cuong.Name = "Cuong"

	var Hau Person
	Hau.Age = 27
	Hau.Name = "Hau"

	fmt.Println(Cuong.Age)
	fmt.Println(Hau.Name)

	fmt.Println("End main function")
}
