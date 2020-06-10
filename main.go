package main

import (
	"fmt"

	"github.com/zouyx/agollo/v3"
)

func main() {
	A()
	B()
	v := agollo.GetValue("key")
	fmt.Print(v)
}
