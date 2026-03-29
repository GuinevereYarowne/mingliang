package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

func main() {

	jsonStr := `{"name":"Jane Smith","age":25,"email":"janesmith@example.com"}`

	var person Person

	err := json.Unmarshal([]byte(jsonStr), &person)
	if err != nil {
		fmt.Println("反序列化失败：", err)
		return
	}

	fmt.Println("反序列化后的Person信息：")
	fmt.Printf("姓名：%s\n", person.Name)
	fmt.Printf("年龄：%d\n", person.Age)
	fmt.Printf("邮箱：%s\n", person.Email)
}
