package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// 全局变量（避免复杂的依赖注入）
var (
	reportURL  string        // 远程上报地址
	resultsDir = "./results" // 结果文件存储目录
)

// FileResult 存储单个文件的计算结果
type FileResult struct {
	OriginalFilename string   // 原始文件名（如a.txt，方便用户识别）
	Filename         string   // 唯一文件名（含时间戳）
	Expressions      []string // 原始表达式列表
	ValidResults     []int    // 仅正确的计算结果
	ValidExpressions []string // 仅正确的表达式（用于生成结果文件）
	ErrorExpressions []string // 错误的表达式（格式："表达式=错误原因"）
	Sum              int      // 仅正确结果的总和
	HasError         bool     // 是否存在错误表达式（而非整个文件失败）
}

// 初始化函数：读取.env配置、创建results目录
func init() {
	// 读取.env文件
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("警告：读取.env文件失败，%v\n", err)
	}
	// 获取上报地址
	reportURL = os.Getenv("REPORT_URL")
	if reportURL == "" {
		fmt.Println("警告：REPORT_URL未配置，远程上报功能将失效")
	}
	// 创建results目录（不存在则创建）
	err = os.MkdirAll(resultsDir, 0755)
	if err != nil {
		fmt.Printf("错误：创建results目录失败，%v\n", err)
		os.Exit(1) // 目录创建失败，程序无法运行，直接退出
	}
}

// 辅助函数：获取当前时间戳（毫秒）
func getTimestampMs() int64 {
	return time.Now().UnixNano() / 1e6 // 纳秒转毫秒
}

// 辅助函数：计算字符串的MD5（32位小写）
func calculateMD5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// 辅助函数：排序文件名（升序）
func sortFilenames(filenames []string) []string {
	sort.Strings(filenames)
	return filenames
}

// 核心函数：解析并计算单个算术表达式
// 支持的运算符：+、-、*；支持负数（如-5+3、10*-4）、单个字符开头的运算符（如+123、-456）
func calculateExpression(expr string) (int, error) {
	// 步骤1：找到运算符的位置（处理负数/单个字符运算符情况）
	opIndex := -1
	// 先判断第一个字符是否为运算符（处理+123、-456这种情况）
	if len(expr) > 0 {
		switch expr[0] {
		case '+', '-', '*':
			opIndex = 0
		}
	}
	// 若第一个字符不是运算符，从索引1开始找（处理-5+3这种负数情况）
	if opIndex == -1 {
		for i := 1; i < len(expr); i++ {
			switch expr[i] {
			case '+', '-', '*':
				opIndex = i
				break
			}
		}
	}
	if opIndex == -1 {
		return 0, fmt.Errorf("表达式%s无有效运算符（仅支持+、-、*）", expr)
	}

	// 步骤2：分割操作数和运算符
	num1Str := expr[:opIndex]
	operator := expr[opIndex]
	num2Str := expr[opIndex+1:]

	// 处理num1Str为空的情况（如+123 → num1Str为空，视为0）
	if num1Str == "" {
		num1Str = "0"
	}
	// 处理num2Str为空的情况（如123+ → num2Str为空，视为0）
	if num2Str == "" {
		num2Str = "0"
	}

	// 步骤3：转换为整数
	num1, err := strconv.Atoi(num1Str)
	if err != nil {
		return 0, fmt.Errorf("表达式%s中操作数1[%s]转换失败：%v", expr, num1Str, err)
	}
	num2, err := strconv.Atoi(num2Str)
	if err != nil {
		return 0, fmt.Errorf("表达式%s中操作数2[%s]转换失败：%v", expr, num2Str, err)
	}

	// 步骤4：计算结果
	var result int
	switch operator {
	case '+':
		result = num1 + num2
	case '-':
		result = num1 - num2
	case '*':
		result = num1 * num2
	default:
		return 0, fmt.Errorf("不支持的运算符%c", operator)
	}

	return result, nil
}

// 核心函数：处理单个文件的所有表达式（供goroutine调用）
func processFile(fileHeader *multipart.FileHeader, originalName string, resultChan chan<- FileResult) {
	// 初始化文件结果
	fileResult := FileResult{
		Filename:         fileHeader.Filename, // 唯一文件名（含时间戳）
		OriginalFilename: originalName,        // 原始文件名（不修改，避免下划线分割错误）
		ErrorExpressions: make([]string, 0),   // 初始化错误列表
	}

	// 步骤1：打开上传的文件
	file, err := fileHeader.Open()
	if err != nil {
		// 这里是「文件无法打开」，属于题目要求的“文件解析失败”
		fileResult.HasError = true
		fileResult.ErrorExpressions = append(fileResult.ErrorExpressions, fmt.Sprintf("文件打开失败：%v", err))
		resultChan <- fileResult
		return
	}
	defer file.Close() // 延迟关闭文件

	// 步骤2：读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		// 这里是「文件无法读取」，属于题目要求的“文件解析失败”
		fileResult.HasError = true
		fileResult.ErrorExpressions = append(fileResult.ErrorExpressions, fmt.Sprintf("文件读取失败：%v", err))
		resultChan <- fileResult
		return
	}

	// 步骤3：按行分割内容（处理不同换行符：\n、\r\n）
	lines := strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")

	// 步骤4：遍历每行表达式并计算
	total := 0
	for _, line := range lines {
		line = strings.TrimSpace(line) // 去除空格和换行
		if line == "" {
			continue // 跳过空行
		}
		// 记录原始表达式
		fileResult.Expressions = append(fileResult.Expressions, line)
		// 计算表达式
		res, err := calculateExpression(line)
		if err != nil {
			// 记录错误表达式，不中断循环
			fileResult.ErrorExpressions = append(fileResult.ErrorExpressions, fmt.Sprintf("%s=%v", line, err))
			fileResult.HasError = true
			continue
		}
		// 收集正确的表达式和结果
		fileResult.ValidExpressions = append(fileResult.ValidExpressions, line)
		fileResult.ValidResults = append(fileResult.ValidResults, res)
		total += res
	}

	fileResult.Sum = total
	// 发送结果到通道（仅发送一次）
	resultChan <- fileResult
}

// 辅助函数：远程上报计算结果（异步执行，失败仅记录日志）
func reportResult(reportData map[string]interface{}) {
	// 异步上报：用goroutine不阻塞主流程
	go func() {
		// 1. 先检查上报地址是否配置
		if reportURL == "" {
			fmt.Println("上报失败：REPORT_URL未配置")
			return
		}

		// 2. 创建带5秒超时的context（控制整个请求的生命周期）
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel() // 确保函数结束时取消context，避免资源泄漏

		// 3. 配置http.Client（细粒度超时，避免堵塞）
		client := &http.Client{
			Timeout: 5 * time.Second, // 整体超时（兜底，防止底层超时配置漏了）
			Transport: &http.Transport{
				// 连接超时：2秒（建立TCP连接的时间限制，快速判断网络是否通）
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					conn, err := net.DialTimeout(network, addr, 2*time.Second)
					if err != nil {
						return nil, err
					}
					// 读写超时：3秒（连接建立后，收发数据的时间限制）
					if err := conn.SetDeadline(time.Now().Add(3 * time.Second)); err != nil {
						return nil, err
					}
					return conn, nil
				},
				TLSHandshakeTimeout: 2 * time.Second, // HTTPS握手超时（如果上报地址是HTTPS）
			},
		}

		// 4. 将上报数据序列化为JSON（失败则记日志返回）
		jsonData, err := json.Marshal(reportData)
		if err != nil {
			fmt.Printf("上报失败：JSON序列化失败，%v\n", err)
			return
		}

		// 5. 创建HTTP POST请求
		req, err := http.NewRequestWithContext(
			ctx,                                 // 带超时的context
			"POST",                              // 请求方法
			reportURL,                           // 上报地址
			strings.NewReader(string(jsonData)), // 请求体（JSON字符串）
		)
		if err != nil {
			fmt.Printf("上报失败：创建请求对象失败，%v\n", err)
			return
		}

		// 6. 设置请求头（告诉远程服务：请求体是JSON格式）
		req.Header.Set("Content-Type", "application/json")

		// 7. 发送请求
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("上报失败：请求发送失败（超时/网络错误），%v\n", err)
			return
		}
		defer resp.Body.Close() // 确保响应体关闭，避免资源泄漏

		// 8. 检查远程响应状态码（非2xx视为失败）
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			fmt.Printf("上报失败：远程服务返回错误状态码，code=%d\n", resp.StatusCode)
			return
		}

		// 9. 上报成功（仅日志提示）
		fmt.Println("上报成功：远程服务接收数据")
	}()
}

// 接口1：POST /api/calculate - 处理计算请求
func handleCalculate(c *gin.Context) {
	// 步骤1：校验必填参数
	username := c.PostForm("username")
	uuid := c.PostForm("uuid")
	if username == "" || uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and uuid are required"})
		return
	}

	// 步骤2：获取上传的文件
	files, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get files"})
		return
	}
	fileHeaders := files.File["files"]
	if len(fileHeaders) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no files uploaded"})
		return
	}

	// 步骤3：记录请求开始时间
	reqStart := getTimestampMs()

	// 步骤4：并发处理所有文件
	resultChan := make(chan FileResult, len(fileHeaders)) // 带缓冲通道，避免阻塞
	var wg sync.WaitGroup
	for _, fh := range fileHeaders {
		wg.Add(1)
		// 复制文件头（避免循环变量复用问题）
		originalFilename := fh.Filename // 保存原始文件名
		// 生成唯一文件名（原文件名_纳秒时间戳，确保绝对唯一）
		uniqueFilename := fmt.Sprintf("%s_%d", originalFilename, time.Now().UnixNano())
		fileHeaderCopy := *fh
		fileHeaderCopy.Filename = uniqueFilename

		// 启动独立协程处理文件，传入复制后的文件头和原始文件名
		go func(fc *multipart.FileHeader, originalName string) {
			defer wg.Done()
			processFile(fc, originalName, resultChan)
		}(&fileHeaderCopy, originalFilename)
	}

	// 等待所有goroutine完成后关闭通道
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 步骤5：收集所有文件的计算结果
	fileResults := make(map[string]FileResult)
	var filenames []string
	var allErrors []string // 收集所有文件的错误表达式（最终返回给用户）
	for fr := range resultChan {
		// 收集该文件的所有错误（如果有的话）
		if len(fr.ErrorExpressions) > 0 {
			errMsg := fmt.Sprintf("%s：%s", fr.OriginalFilename, strings.Join(fr.ErrorExpressions, "；"))
			allErrors = append(allErrors, errMsg)
		}

		// 即使文件有部分错误，也收集其正确结果（如果有的话）
		if len(fr.ValidResults) > 0 {
			fileResults[fr.Filename] = fr
			filenames = append(filenames, fr.Filename)
		} else if fr.HasError {
			// 只有「文件完全无法读取/打开」（无任何有效结果），才返回400
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("文件解析失败：%s", strings.Join(fr.ErrorExpressions, "；"))})
			return
		}
	}

	// 步骤6：整理响应数据
	sortedFilenames := sortFilenames(filenames)
	details := make(map[string]string)
	var sumList []string
	var totalCount int
	var resultFileContent []string

	for _, fn := range sortedFilenames {
		fr := fileResults[fn]
		resultStrs := make([]string, len(fr.ValidResults))
		for i, res := range fr.ValidResults {
			resultStrs[i] = strconv.Itoa(res)
			resultFileContent = append(resultFileContent, fmt.Sprintf("%s=%d", fr.ValidExpressions[i], res))
		}
		details[fr.OriginalFilename] = strings.Join(resultStrs, ",")
		sumList = append(sumList, strconv.Itoa(fr.Sum))
		totalCount += len(fr.ValidResults)
	}

	// 步骤7：定义关键变量（请求结束时间、结果文件名）
	reqEnd := getTimestampMs()
	resultFilename := fmt.Sprintf("result_%d.txt", reqStart)
	sumListStr := strings.Join(sumList, ",")

	// 步骤8：写入结果文件（处理空结果场景）
	if len(resultFileContent) == 0 {
		resultFileContent = append(resultFileContent, "无有效计算结果（所有表达式均错误）")
	}
	resultFilePath := fmt.Sprintf("%s/%s", resultsDir, resultFilename)
	err = os.WriteFile(resultFilePath, []byte(strings.Join(resultFileContent, "\n")), 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to write result file: %v", err)})
		return
	}

	// 步骤9：准备远程上报数据
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		fmt.Printf("序列化details失败：%v\n", err)
	}
	detailsStr := string(detailsJSON)
	md5Str := calculateMD5(detailsStr)
	reportData := map[string]interface{}{
		"username":    username,
		"uuid":        uuid,
		"req_start":   reqStart,
		"req_end":     reqEnd,
		"details":     detailsStr,
		"sum_list":    sumListStr,
		"total_count": totalCount,
		"md5":         md5Str,
	}
	// 异步上报（不阻塞响应）
	reportResult(reportData)

	// 步骤10：构建响应并统一返回
	response := gin.H{
		"req_start":   reqStart,
		"req_end":     reqEnd,
		"details":     details,
		"sum_list":    sumListStr,
		"total_count": totalCount,
		"result_file": resultFilename,
	}
	// 若有错误表达式，添加提示（不影响有效结果）
	if len(allErrors) > 0 {
		response["error_hint"] = "部分表达式计算失败：" + strings.Join(allErrors, "；")
	}
	c.JSON(http.StatusOK, response)
}

// 接口2：GET /api/result/list - 获取结果文件列表
func handleResultList(c *gin.Context) {
	// 读取results目录下的所有文件
	entries, err := os.ReadDir(resultsDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to read results dir: %v", err)})
		return
	}

	// 收集文件名
	var filenames []string
	for _, entry := range entries {
		if !entry.IsDir() {
			filenames = append(filenames, entry.Name())
		}
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{"files": filenames})
}

// 接口3：GET /api/result/detail - 查看/下载结果文件
func handleResultDetail(c *gin.Context) {
	// 步骤1：校验fileid参数
	fileid := c.Query("fileid")
	if fileid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fileid is required"})
		return
	}

	// 步骤2：拼接文件路径
	filePath := fmt.Sprintf("%s/%s", resultsDir, fileid)
	// 检查文件是否存在
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to check file: %v", err)})
		return
	}

	// 步骤3：判断是否下载
	download := c.Query("download")
	if download == "1" {
		// 下载文件
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileid))
		c.Header("Content-Type", "application/octet-stream")
		c.File(filePath)
		return
	}

	// 步骤4：返回文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to read file: %v", err)})
		return
	}
	c.Data(http.StatusOK, "text/plain; charset=utf-8", content)
}

func main() {
	// 初始化Gin引擎
	r := gin.Default()

	// 注册路由
	r.POST("/api/calculate", handleCalculate)
	r.GET("/api/result/list", handleResultList)
	r.GET("/api/result/detail", handleResultDetail)

	// 启动服务
	fmt.Println("服务器启动在 :8080 端口")
	err := r.Run(":8080")
	if err != nil {
		fmt.Printf("服务器启动失败：%v\n", err)
	}
}
