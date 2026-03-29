package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// 8.html到48.html
func generateUrls() []string { //返回一个字符串类型的切片（存储所有 URL）
	var urls []string
	for i := 8; i <= 48; i++ {
		url := fmt.Sprintf("https://study-test.sixue.work/html/%d.html", i)
		urls = append(urls, url)
	}
	return urls
}

// 爬取单个网页，返回所有图片链接
func fetchImageUrls(url string) ([]string, error) {
	//返回两个值,图片链接切片([]string)和错误(error),error返回执行失败的原因(比如请求超时、HTML 解析失败)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()
	//defer：延迟执行后面的函数，直到当前函数（fetchImageUrls）退出前才执行。
	// resp.Body是一个 IO 资源（占用系统文件描述符），如果不关闭会导致资源泄露（程序运行久了会耗尽资源）。
	// 必须写在err判断之后：如果http.Get失败，resp是nil，调用resp.Body.Close()会 panic（空指针错误）。

	if resp.StatusCode != http.StatusOK {
		// resp.StatusCode：HTTP 响应状态码（比如 200 = 成功，404 = 页面不存在，500 = 服务器错误）。
		// http.StatusOK：net/http包的常量，值为 200，用常量比直接写 200 更易读。
		return nil, fmt.Errorf("响应错误: 状态码%d", resp.StatusCode)
	}
	doc, err := html.Parse(resp.Body)
	// 	html.Parse(resp.Body)：golang.org/x/net/html包的核心函数，接收一个io.Reader类型(resp.Body正好实现了该接口)，解析 HTML 文档,然后返回两个值：
	// doc *html.Node：HTML 文档的根节点（HTML 是树形结构，所有元素都是节点的子节点）。
	// err error：解析错误（比如 HTML 格式不合法）
	if err != nil {
		return nil, fmt.Errorf("解析HTML失败: %w", err)
	}
	var images []string
	var parseNode func(node *html.Node)
	parseNode = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "img" {
			// node.Type：节点类型，html.ElementNode是 “元素节点”（比如<img>、<div>标签），还有html.TextNode（文本节点）等。
			// node.Data：元素节点的标签名（比如 img、div、a），这里判断是否为img标签。
			for _, attr := range node.Attr {
				if attr.Key == "src" { //判断是否为src属性（图片链接存放在src中）
					images = append(images, attr.Val)
					break
				}
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			//HTML 是树形结构（比如<html>包含<body>，<body>包含<div>，<div>包含<img>）
			parseNode(child)
		}
	}
	parseNode(doc)
	return images, nil
}

func main() {
	startTime := time.Now()
	urls := generateUrls()
	// 并发控制，5个
	semaphore := make(chan struct{}, 5)
	// 	make(chan struct{}, 5)：创建一个 “带缓冲的通道”（chan），类型是struct{}（空结构体），缓冲容量是 5。
	// 核心作用：控制并发数不超过 5。
	// 为什么用struct{}？空结构体不占用任何内存（sizeof(struct{}) == 0），适合作为 “信号”（只需要通道的 “满/空” 状态，不需要传递实际数据）。
	//等待所有goroutine完成
	var wg sync.WaitGroup

	var (
		totalImages int                         //图片总数（不去重）
		uniqueMap   = make(map[string]struct{}) // 去重map（key=图片链接，value=空结构体)
		//map 的 key 不能重复，重复插入会覆盖，所以最终 map 的长度就是去重后的图片数。
		mutex sync.Mutex // 互斥锁
	)

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			// 获取、释放信号量
			semaphore <- struct{}{}
			// 往信号量通道发送一个空结构体，获取信号。如果通道已满（已经有 5 个 Goroutine 在执行），会阻塞在这里，直到有 Goroutine 释放信号。
			defer func() { <-semaphore }()
			//延迟执行 “从通道接收数据”，释放信号量。确保 Goroutine 无论是否出错，都会释放信号，避免通道永远满了导致其他 Goroutine 阻塞。

			//URL的图片链接
			images, err := fetchImageUrls(u)
			if err != nil {
				fmt.Printf("处理%s失败：%v\n", u, err)
				return
			}
			mutex.Lock()
			totalImages += len(images)
			for _, img := range images {
				uniqueMap[img] = struct{}{}
			}
			mutex.Unlock()
		}(url)
	}
	wg.Wait()
	//计算耗时
	elapsedTime := time.Since(startTime).Milliseconds()
	fmt.Printf("执行任务总耗时（毫秒）：%v，图片总数（不去重）：%d，图片总数（去重）：%d\n", elapsedTime, totalImages, len(uniqueMap))
}
