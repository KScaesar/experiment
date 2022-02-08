package main

import (
	"fmt"
	"reflect"
)

func main() {
	t1 := &testType{Name: "t1"}
	action(t1)

	fmt.Println()

	t2 := testType{Name: "t2"}
	action(t2)

	// https://stackoverflow.com/questions/18570391/check-if-struct-implements-a-given-interface
}

// 如果定義為以下形式, refType = nil
// 拿到的 interface = nil, 會造成 panic
// var refType = reflect.TypeOf((Valuer)(nil))

// 應該 使用 pointer of interface, 再使用 elem()
// 如此 拿到 interface != nil
var refType = reflect.TypeOf((*Valuer)(nil)).Elem()

func action(obj interface{}) {
	rv := reflect.ValueOf(obj)
	tv := rv.Type()

	ok := tv.Implements(refType)
	if ok {
		fmt.Println("Implements", tv.String())
		obj.(Valuer).M1()
	} else {
		fmt.Println("no Implements", tv.Name())
	}
}

type testType struct {
	Name string
}

func (v *testType) M1() {
	fmt.Println("M1", v.Name)
}

type Valuer interface {
	M1()
}
