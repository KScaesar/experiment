package build_dir_var

import "fmt"

var version = "Nothing"

func Version() {
	fmt.Println("version=", version)
}
