package config

import (
	"github.com/zouyx/agollo/v3"
)

func Get() string {
	// de something

	v := agollo.GetValue("key")
	return v
}
