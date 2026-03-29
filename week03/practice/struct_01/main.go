package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name  string
	Age   int
	Email string
}

func NewPerson(name string, age int, email string) Person {
	return Person{
		Name:  name,
		Age:   age,
		Email: email,
	}
}

func PrintPerson(p Person) {
	fmt.Println("Person信息：")
	fmt.Printf("姓名：%s\n", p.Name)
	fmt.Printf("年龄：%d\n", p.Age)
	fmt.Printf("邮箱：%s\n", p.Email)

	jsonData, err := json.Marshal(p)
	if err != nil {
		fmt.Println("JSON转换失败：", err)
		return
	}

	fmt.Println("JSON格式：", string(jsonData))
}

func main() {
	person := NewPerson("小明", 20, "xiaoming@example.com")

	PrintPerson(person)
}
