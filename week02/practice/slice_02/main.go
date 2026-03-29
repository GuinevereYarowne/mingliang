package main

import "fmt"

func main() {
	slice1 := []int{1, 2, 3, 4}
	slice2 := []int{3, 4, 5, 6}
	fmt.Println("原始切片 slice1：", slice1)
	fmt.Println("原始切片 slice2：", slice2)
	combinedslice := append(slice1, slice2...)
	fmt.Println("拼接后的切片：", combinedslice)
	tempMap := make(map[int]bool)
	uniqueslice := []int{}
	for _, num := range combinedslice {
		if !tempMap[num] {
			tempMap[num] = true
			uniqueslice = append(uniqueslice, num)
		}
	}
	fmt.Println("去重后的切片 uniqueslice：", uniqueslice)
}
