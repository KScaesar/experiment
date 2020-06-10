package main

import (
	"fmt"

	"experiment/config"
)

func B() {
	v := config.Get()
	fmt.Print(v)
}
