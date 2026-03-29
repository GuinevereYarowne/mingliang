package main // 改成 main 包（可执行包）

import (
	"errors"
	"fmt" // 加个 fmt 用来打印测试结果
)

// 原来的三个函数不变，直接保留
func BubbleSort(numbers []int) []int {
	newNums := make([]int, len(numbers))
	copy(newNums, numbers)
	length := len(newNums)
	for i := 0; i < length-1; i++ {
		for j := 0; j < length-1-i; j++ {
			if newNums[j] > newNums[j+1] {
				newNums[j], newNums[j+1] = newNums[j+1], newNums[j]
			}
		}
	}
	return newNums
}

func IsSorted(numbers []int) bool {
	for i := 0; i < len(numbers)-1; i++ {
		if numbers[i] > numbers[i+1] {
			return false
		}
	}
	return true
}

func FindMedian(numbers []int) (float64, error) {
	if len(numbers) == 0 {
		return 0, errors.New("数组为空")
	}
	sorted := BubbleSort(numbers)
	length := len(sorted)
	if length%2 == 1 {
		return float64(sorted[length/2]), nil
	} else {
		mid1 := sorted[length/2-1]
		mid2 := sorted[length/2]
		return float64(mid1+mid2) / 2.0, nil
	}
}

// 加个 main 函数（程序入口），里面调用上面的函数测试
func main() {
	// 测试数据
	testNums := []int{7, 3, 5, 2, 9, 1, 4}

	// 调用排序函数
	sortedNums := BubbleSort(testNums)
	fmt.Println("原始数组：", testNums)
	fmt.Println("排序后数组：", sortedNums)

	// 调用是否排序判断
	fmt.Println("数组是否升序：", IsSorted(sortedNums)) // 输出 true

	// 调用中位数函数
	median, err := FindMedian(testNums)
	if err != nil {
		fmt.Println("获取中位数失败：", err)
	} else {
		fmt.Println("数组中位数：", median) // 测试数据排序后是 [1,2,3,4,5,7,9]，中位数是 4.0
	}
}
