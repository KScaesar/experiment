package config

import (
	"fmt"

	"github.com/zouyx/agollo/v3"
)

func B() {
	v := Get()
	fmt.Print(v)
}

func Get() string {
	// de something

	v := agollo.GetValue("key")
	return v
}
