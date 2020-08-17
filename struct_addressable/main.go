package main

import (
	"fmt"
	"reflect"
)

type Animal interface {
	CoolName() string
}

type Cat struct {
	name string
}

func (c *Cat) CoolName() string {
	return c.name
}

type Man struct {
	ManPet Cat
}

type Woman struct {
	WomanPet Cat
}

func main() {
	m := Man{
		ManPet: Cat{
			name: "亞當的貓",
		},
	}
	WhatYourPetName(m)

	fmt.Println()

	w := &Woman{
		WomanPet: Cat{
			name: "夏娃的貓",
		},
	}
	WhatYourPetName(w)
}

func WhatYourPetName(person interface{}) {
	p := reflect.ValueOf(person)

	switch p.Kind() {
	case reflect.Struct:
		StructAction(p)

	case reflect.Ptr:
		elem := p.Elem()
		if elem.Kind() != reflect.Struct {
			fmt.Printf("I am %v, not struct", elem.Type().String())
			return
		}
		StructAction(elem)
	}
}

func StructAction(v reflect.Value) {
	fmt.Println("=====start=====")

	PersonFieldName := v.Type().Field(0).Name
	PersonFieldValue := v.Field(0)
	// 1
	IamStruct(v)

	// 2
	WhatISFieldType(PersonFieldName, PersonFieldValue)

	// ●*
	ExecuteAnimalInterface(PersonFieldValue)

	fmt.Println()

	// 3
	if CanFieldAddressable(v) {
		PersonFieldValue = PersonFieldValue.Addr() // 重點!! 可尋址就可以轉回指標值, Addr() returns a pointer value

		// 4
		fmt.Println("4. Use reflect.Value.Addr(): get pointer")

		// ●*
		ExecuteAnimalInterface(PersonFieldValue)
	}

	fmt.Println("=====end======")
}

func IamStruct(v reflect.Value) {
	structName := v.Type().String()
	fieldNames := make([]string, 0)
	for i := 0; i < v.NumField(); i++ {
		fieldNames = append(fieldNames, v.Type().Field(i).Name)
	}
	fmt.Printf("1. I am struct %v, having fields %#v\n", structName, fieldNames)
}

func WhatISFieldType(fieldName string, v reflect.Value) {
	fieldType := v.Type().String()
	fmt.Printf("2. %v type is %v\n", fieldName, fieldType)
}

func CanFieldAddressable(v reflect.Value) bool {
	fieldName := v.Type().Field(0).Name
	fieldValue := v.Field(0)
	fmt.Printf("3. %v can addressable? %v\n", fieldName, fieldValue.CanAddr())
	return fieldValue.CanAddr()
}

func ExecuteAnimalInterface(v reflect.Value) bool {
	animalInterfaceType := reflect.TypeOf((*Animal)(nil)).Elem()
	interfaceName := animalInterfaceType.String()

	typeName := v.Type().String()
	IsImplement := v.Type().Implements(animalInterfaceType)
	fmt.Printf("●*1. %v Implements %v? %v\n", typeName, interfaceName, IsImplement)

	if IsImplement {
		animal := v.Interface().(Animal)
		fmt.Printf("●*2. Execute %v interface method: %v\n", interfaceName, animalInterfaceType.Method(0).Name)
		fmt.Printf("●*3. WhatYourPetName: %v\n", animal.CoolName())
	}
	return IsImplement
}
