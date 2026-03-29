package main

import (
	"fmt"
	"sync" // 用于WaitGroup
)

func calculateSquareSum(nums []int, start, end int, resultChan chan<- int) {
	sum := 0
	for i := start; i < end; i++ {
		sum += nums[i] * nums[i]
	}
	resultChan <- sum
}

func main() {
	nums := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		nums[i] = i
	}

	partSize := 250
	parts := 4
	var wg sync.WaitGroup

	resultChan := make(chan int, parts)

	for i := 0; i < parts; i++ {
		start := i * partSize
		end := start + partSize
		if end > len(nums) {
			end = len(nums)
		}

		wg.Add(1)
		go func(s, e int) {
			defer wg.Done()
			calculateSquareSum(nums, s, e, resultChan)
		}(start, end)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	totalSum := 0

	for partSum := range resultChan {
		totalSum += partSum
	}

	fmt.Println("所有元素的平方和为：", totalSum)
}
