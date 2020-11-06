package main

func NewHome() Home {
	oV := Order{}
	oP := &Order{}
	list_oV := []Order{oV}
	list_oP := []*Order{oP}
	return Home{
		ItemValue:   list_oV,
		ItemPointer: list_oP,
	}
}

type Home struct {
	ItemValue   []Order
	ItemPointer []*Order
}

func (h *Home) Counter() {}

type Order struct {
	Name    string
	Amount1 int64
}

func main() {
	h := NewHome()
	h.Counter()
}
