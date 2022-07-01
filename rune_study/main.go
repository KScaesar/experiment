package main

import "fmt"

func main() {
	// "a" is string, 'a' is rune(int32)
	fmt.Printf(`
%[1]T:%#[1]v
%[2]T:%#[2]v
%[3]T:%#[3]v
`, []byte("a"), "a", 'a')
	// []uint8:[]byte{0x61}
	// string:"a"
	// int32:97

}
