package main

import (
	"fmt"
	"sync"
)

type LogEntry struct {
	ID      int
	Content string
}

func generateLogs(logChannel chan<- LogEntry, stopChan <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 100; i++ {
		log := LogEntry{
			ID:      i,
			Content: fmt.Sprintf("日志内容 %d", i),
		}
		select {
		case <-stopChan:
			fmt.Println("生成器: 收到停止信号，退出")
			return
		case logChannel <- log:
			fmt.Printf("生成器: 生成日志 ID=%d\n", i)
		}
	}
	close(logChannel)
	fmt.Println("生成器: 完成所有日志生成")
}

func filterLogs(logChannel <-chan LogEntry, filteredChannel chan<- LogEntry,
	stopChan chan struct{}, wg *sync.WaitGroup, id int) {
	defer wg.Done()
	for log := range logChannel {
		if log.ID%2 == 0 {
			if log.ID == 50 {
				fmt.Printf("过滤器 %d: 发现ID=50的日志，发送停止信号\n", id)
				close(stopChan)
				return
			}
			select {
			case <-stopChan:
				fmt.Printf("过滤器 %d: 收到停止信号，退出\n", id)
				return
			case filteredChannel <- log:
				fmt.Printf("过滤器 %d: 过滤日志 ID=%d (偶数)\n", id, log.ID)
			}
		} else {
			fmt.Printf("过滤器 %d: 忽略日志 ID=%d (奇数)\n", id, log.ID)
		}
	}
	fmt.Printf("过滤器 %d: logChannel已关闭，退出\n", id)
}

func storeLogs(filteredChannel <-chan LogEntry, stopChan <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-stopChan:
			fmt.Println("存储器: 收到停止信号，退出")
			return
		case log, ok := <-filteredChannel:
			if !ok {
				fmt.Println("存储器: filteredChannel已关闭，退出")
				return
			}
			fmt.Printf("存储器: 存储日志 [ID=%d, 内容=\"%s\"]\n", log.ID, log.Content)
		}
	}
}

func main() {
	fmt.Println("开始日志处理系统...")

	logChannel := make(chan LogEntry, 10)      // 传递原始日志
	filteredChannel := make(chan LogEntry, 10) // 传递过滤后的日志
	stopChan := make(chan struct{})            // 发送停止信号
	var wg sync.WaitGroup
	fmt.Println("\n--- 阶段一：启动日志生成器 ---")
	wg.Add(1)
	go generateLogs(logChannel, stopChan, &wg)
	fmt.Println("\n--- 阶段三：启动日志存储器 ---")
	wg.Add(1)
	go storeLogs(filteredChannel, stopChan, &wg)
	fmt.Println("\n--- 阶段二：启动日志过滤器 (3个) ---")
	var filterWg sync.WaitGroup
	filterWg.Add(3)
	// 启动3个过滤goroutine
	for i := 1; i <= 3; i++ {
		go func(filterId int) {
			fmt.Printf("过滤器 %d: 启动\n", filterId)
			filterLogs(logChannel, filteredChannel, stopChan, &filterWg, filterId)
			fmt.Printf("过滤器 %d: 退出\n", filterId)
		}(i)
	}
	// 等待所有过滤goroutine完成
	go func() {
		filterWg.Wait()
		close(filteredChannel)
		fmt.Println("所有过滤器已完成，关闭filteredChannel")
	}()
	// 等待所有goroutine完成
	wg.Wait()
	fmt.Println("\n日志处理系统已停止")
}
