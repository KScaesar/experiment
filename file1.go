package main

import (
	"fmt"

	"github.com/zouyx/agollo/v3"
)

func A() {
	v := agollo.GetValue("key")
	fmt.Print(v)
}
