package main

import "testing"

func BenchmarkValueMethod_AssignValue(b *testing.B) {
	var animal Animal = CatValueMethod{Name: "ball"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		animal.Call()
	}
}

func BenchmarkValueMethod_AssignPointer(b *testing.B) {
	var animal Animal = &CatValueMethod{Name: "ball"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		animal.Call()
	}
}

func BenchmarkPointerMethod_AssignPointer(b *testing.B) {
	var animal Animal = &DogPointerMethod{Name: "ball"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		animal.Call()
	}
}
