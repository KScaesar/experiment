package main

import (
	"time"

	"github.com/jmoiron/sqlx"
)

func NewHome() Home {
	var db sqlx.DB
	return Home{db: &db}
}

type Home struct {
	Person []string
	time   time.Time
	db     *sqlx.DB
}

func (h *Home) Counter() *Home {
	p := []string{""}
	h.Person = p
	return h
}

func main() {
	h := NewHome()
	h2 := h.Counter()
	h2.Counter()
}
