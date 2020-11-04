package main

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

func NewHome() Home {
	name := [1]string{"caesar"}
	return Home{Person: name}
}

type Home struct {
	Person [1]string
	time   time.Time
	db     *sqlx.DB
}

func (h *Home) Counter() {}

func test(h *Home) Home {
	h.Counter()
	return *h
}

func main() {
	h := NewHome()
	test(&h)

	a := A{
		Name: "Car",
	}
	b := B(a)
	b.Name = "lala"
	fmt.Printf("%v %v\n", &a, &b)
}

type A struct {
	Name string
}

type B A
