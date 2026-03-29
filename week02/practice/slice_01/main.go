package main

import "fmt"

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println("步骤1：初始切片 =", nums)
	subNums := nums[2:7]
	fmt.Println("步骤2：第3到第7个元素 =", subNums)
	subNums = append(subNums, 11, 12, 13)
	fmt.Println("步骤3：添加后 =", subNums)
	subNums = append(subNums[:4], subNums[5:]...)
	fmt.Println("步骤4：删除第5个元素后 =", subNums)
	for i := 0; i < len(subNums); i++ {
		subNums[i] *= 2
	}
	fmt.Println("步骤5：所有元素乘以2后 =", subNums)
	fmt.Println("\n最终切片内容：", subNums)
	fmt.Println("最终切片长度：", len(subNums))
	fmt.Println("最终切片容量：", cap(subNums))
}
