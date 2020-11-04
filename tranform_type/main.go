package main

import "github.com/davecgh/go-spew/spew"

type A struct {
	Name string
	Age  int
}

type V struct {
	Name string
	Age  int
}

func main() {
	a := A{
		Name: "asdfas",
		Age:  13,
	}

	b := V(a)

	spew.Dump(b)
}
