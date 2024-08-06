package main

import "fmt"

type man struct {
	manID string
	name  string
	phone string
}

func main() {
	man1 := man{
		manID: "7045500",
		name:  "Palm",
		phone: "4156",
	}
	fmt.Println("man1 =", man1)
}
