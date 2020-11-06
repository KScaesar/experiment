# 逃逸分析實驗

會產生逃逸的  
只有該指標  
不會整個物件一起逃  
method 賦值 給自己, 無回傳值, 也會逃逸  

## 回傳指標物件
回傳一個結構 by value  
內部生成一個指標一起回傳  
會產生逃逸  

```
func NewHome() Home {
	var db sqlx.DB
	return Home{db: &db}
}

type Home struct {
	Person []string
	time   time.Time
	db     *sqlx.DB
}

func (h *Home) Counter() {}

func main() {
	h := NewHome()
	h.Counter()
}

```
```
 go run -gcflags "-m -l -N" .

# experiment/escape
./escape.go:10:6: moved to heap: db
./escape.go:20:7: (*Home).Counter h does not escape
```

## 回傳一個結構的空值
```
func NewHome() Home {
	return Home{}
}

type Home struct {
	Person []string
	time   time.Time
	db     *sqlx.DB
}

func (h *Home) Counter() {}

func main() {
	h := NewHome()
	h.Counter()
}

```
```
 go run -gcflags "-m -l -N" .

# experiment/escape
./escape.go:19:7: (*Home).Counter h does not escape
```

## 回傳時間值
```
func NewHome() Home {
	return Home{time: time.Now()}
}

type Home struct {
	Person []string
	time   time.Time
	db     *sqlx.DB
}

func (h *Home) Counter() {}

func main() {
	h := NewHome()
	h.Counter()
}
```
```
 go run -gcflags "-m -l -N" .

# experiment/escape
./escape.go:19:7: (*Home).Counter h does not escape
```

## 回傳 slice

指針 跟 值  
都會發生逃逸  
但 gc 時, 在 heap 的 指針, 需要多追蹤一次  
增加 gc mark 記憶體的時間  
所以用值傳遞比較好  

```
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

```
```
 go run -gcflags "-m -l -N" .

# experiment/escape
./escape.go:5:8: &Order literal escapes to heap
./escape.go:6:20: []Order literal escapes to heap
./escape.go:7:21: []*Order literal escapes to heap
./escape.go:19:7: (*Home).Counter h does not escape
```

## 二次回傳
```
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

func test(h *Home) *Home {
	h.Counter()
	var db sqlx.DB
	h.db = &db
	return h
}

func main() {
	h := NewHome()
	test(&h)
}

```
```
 go run -gcflags "-m -l -N" .

# experiment/escape
./escape.go:20:7: (*Home).Counter h does not escape
./escape.go:22:11: leaking param: h to result ~r1 level=0
./escape.go:24:6: moved to heap: db


```

## method 賦值 給自己, 無回傳值, 逃逸
```
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

```
```
 go run -gcflags "-m -l -N" .

# experiment/escape
./escape.go:10:6: moved to heap: db
./escape.go:20:7: leaking param: h to result ~r0 level=0
./escape.go:21:15: []string literal escapes to heap

```