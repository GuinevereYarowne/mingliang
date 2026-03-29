package practice

import "errors"

// FindMinMax 查找数组中的最大值和最小值
// 若数组为空，返回错误
func FindMinMax(numbers []int) (min int, max int, err error) {
	if len(numbers) == 0 {
		return 0, 0, errors.New("数组为空")
	}

	// 初始化最小值和最大值为第一个元素
	min = numbers[0]
	max = numbers[0]

	// 遍历数组更新最小值和最大值
	for _, num := range numbers[1:] {
		if num < min {
			min = num
		}
		if num > max {
			max = num
		}
	}

	return min, max, nil
}
