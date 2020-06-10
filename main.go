package main

import (
	"fmt"

	"experiment/config"
)

func main() {
	A()
	B()
	v := config.Get()
	fmt.Print(v)
}
