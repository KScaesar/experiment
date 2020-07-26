package main

type Animal interface {
	Call() string
}

type CatValueMethod struct {
	Name string
	v    [100]int
}

func (c CatValueMethod) Call() string {
	return c.Name
}

type DogPointerMethod struct {
	Name string
	v    [100]int
}

func (d *DogPointerMethod) Call() string {
	return d.Name
}
