package model

type Person struct {
	Name string
	age  int
}

func NewPerson(name string, age int) Person {
	return Person{
		Name: name,
		age:  age, // unexported field
	}
}

type Teacher struct {
	Name string
	Age  int // exported field
}

func NewTeacher(name string, age int) Teacher {
	return Teacher{
		Name: name,
		Age:  age,
	}
}
