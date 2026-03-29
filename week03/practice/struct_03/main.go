package main

import "fmt"

type Book struct {
	Title  string
	Author string
	Year   int
}

func FindBooksByAuthor(author string, books []Book) []Book {
	var result []Book
	for _, book := range books {
		if book.Author == author {
			result = append(result, book)
		}
	}
	return result
}

func main() {
	books := []Book{
		{Title: "Go结构体入门", Author: "张三", Year: 2020},
		{Title: "Goweek03", Author: "张三", Year: 2022},
		{Title: "Goweek03/struct03", Author: "张三", Year: 2023},
	}

	zhangSanBooks := FindBooksByAuthor("张三", books)

	fmt.Println("作者'张三'的图书有：")
	for i, book := range zhangSanBooks {
		fmt.Printf("第%d本：书名《%s》，出版年份：%d\n", i+1, book.Title, book.Year)
	}
}
