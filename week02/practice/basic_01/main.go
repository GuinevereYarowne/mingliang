package main

import "fmt"

func Len(s string) int {
	runeArr := []rune(s)
	return len(runeArr) //返回字符，非字节
}

func Huiwen(x int) bool {
	if x < 0 || (x != 0 && x%10 == 0) {
		return false
	}
	rev := 0
	for x > rev {
		rev = rev*10 + x%10
		x /= 10
	}
	return x == rev || x == rev/10
}

func main() {
	name := "小明"
	age := 23
	gender := true // 男为true

	fmt.Println("姓名：", name)
	fmt.Println("年龄：", age)
	genderStr := "男"
	if !gender {
		genderStr = "女"
	}
	fmt.Println("性别：", genderStr)

	fmt.Println("------------------------")

	testStr1 := "小明今年23岁"
	fmt.Printf("字符串「%s」的字符个数：%d\n", testStr1, Len(testStr1))
	fmt.Printf("字符串「%s」的字节数：%d\n", testStr1, len(testStr1))

	fmt.Println("------------------------")

	testNums := []int{121, -121, 0}
	for _, num := range testNums {
		result := Huiwen(num)
		fmt.Printf("整数 %d 是回文数吗？%t\n", num, result)
	}
}
