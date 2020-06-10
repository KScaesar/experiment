package main

import (
	"fmt"

	"github.com/zouyx/agollo/v3"
)

func B() {
	v := agollo.GetValue("key")
	fmt.Print(v)
}
