package domain

type UserParam struct {
	Age  int
	Role string `enums:"admin,ops,dev"`
}

type User struct {
	ID   string
	Name string
	Age  int
	Role string `enums:"admin,ops,dev"`
}
