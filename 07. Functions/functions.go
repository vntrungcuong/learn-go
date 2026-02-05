package main

import (
	"fmt"
)

func sayHello() {
	fmt.Println("Hello World")
}

func sayName(name string, old int) {
	fmt.Println("Hello", name, "you are", old, "years old")
}

func add(a int, b int) int {
	return a + b
}

func stringReturn(name string) string {
	return "Hello " + name
}

func myFunc(x int, y int) (result int) {
	result = x + y
	return
}

func multipleReturn(name string, old int) (name2 string, old2 int) {
	name2 = "Hello " + name
	old2 = old + 10
	return
}

func main() {
	fmt.Println("Start main function")

	sayHello()
	sayName("John", 30)

	fmt.Println(add(5, 2))
	fmt.Println(stringReturn("John"))
	fmt.Println(myFunc(5, 2))
	fmt.Println(multipleReturn("John", 30))

	fmt.Println("End main function")
}
