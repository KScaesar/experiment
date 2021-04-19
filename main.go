package main

import (
	"experiment/build_dir_var"
)

func main() {
	build_dir_var.Version()
}

// https://henvic.dev/posts/my-go-mistakes/
// -X 的名稱 必須是 pkgName.varName
// pkgName 必須加上 moduleName的路徑
//
// go build -ldflags="-X 'experiment/build_dir_var.version=小寫變數也可以'" -o build_test
// ./build_test
