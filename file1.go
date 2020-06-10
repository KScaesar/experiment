package main

import (
	"fmt"

	"experiment/config"
)

func A() {
	v := config.Get()
	fmt.Print(v)
}
