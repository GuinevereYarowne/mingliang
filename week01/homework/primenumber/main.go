package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

//单测内容放在该目录另一个01_test.go文件w下，终端go test -v就行

// isPrime 判断一个数是否为素数
func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// findPrimes 在指定范围内查找所有素数
func findPrimes(start, end int) []int {
	var primes []int
	for i := start; i <= end; i++ {
		if isPrime(i) {
			primes = append(primes, i)
		}
	}
	return primes
}

// writePrimesToFile 将素数写入文件
func writePrimesToFile(primes []int, start, end int) error {
	var strPrimes []string // 声明字符串切片strPrimes，用于存储素数的字符串形式（因为文件写入的是字节
	for _, p := range primes {
		strPrimes = append(strPrimes, strconv.Itoa(p)) //// strconv.Itoa(p)：把整数p转成字符串（比如3→"3"）
	}
	// 把字符串切片用空格连接成一个字符串（比如["2","3","5"]→"2 3 5"）
	content := strings.Join(strPrimes, " ") // strings.Join(切片, 分隔符)：拼接切片元素

	filename := fmt.Sprintf("mingliang_primes_%d_%d.txt", start, end)

	// 写入文件：os.WriteFile(文件名, 字节内容, 文件权限)
	return os.WriteFile(filename, []byte(content), 0644)
	// 1. filename：要写入的文件名；
	// 2. []byte(content)：把字符串转成字节切片（文件写入的是字节流）；
	// 3. 0644：文件权限（八进制），表示“所有者可读可写，其他用户只读”（Go中文件权限必须用0开头表示八进制）；
	// 返回值：error类型，如果写入成功返回nil，失败返回具体错误（比如权限不够）
}

func main() {
	// 检查命令行参数是否正确
	if len(os.Args) != 3 {
		// os.Args：命令行参数切片（索引0是程序名，1是第一个输入参数，2是第二个输入参数）
		// 比如运行命令：go run main.go 2 10 → os.Args = ["./main", "2", "10"]，长度是3
		fmt.Println("用法: go run main.go 开始值 结束值")
		os.Exit(1)
	}

	// 解析命令行参数
	start, err1 := strconv.Atoi(os.Args[1]) // os.Args[1]是第一个输入参数（字符串），转成int类型的start
	end, err2 := strconv.Atoi(os.Args[2])   // os.Args[2]是第二个输入参数，转成int类型的end
	if err1 != nil || err2 != nil || start > end {
		fmt.Println("无效的参数，请输入有效的整数范围")
		os.Exit(1)
	}

	// 记录开始时间
	startTime := time.Now()

	// 查找素数
	primes := findPrimes(start, end)

	// 计算耗时
	duration := time.Since(startTime) // time.Since(t)：返回从t到现在的时间差

	// 写入文件
	err := writePrimesToFile(primes, start, end)
	if err != nil {
		fmt.Printf("写入文件失败: %v\n", err)
		os.Exit(1)
	}

	// 打印结果
	fmt.Printf("计算时间: %v\n", duration)
	fmt.Printf("找到的素数个数: %d\n", len(primes))
}
