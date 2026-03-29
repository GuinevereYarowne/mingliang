package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	filename := "calculations.txt"
	if len(os.Args) > 1 { // 判断是否有命令行参数（os.Args是命令行参数切片）,有就用这个作为文件名
		filename = os.Args[1]
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("读取文件失败：", err)
		return
	}
	lines := strings.Split(string(content), "\n")

	var results []string
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}
		res := processOneLine(trimmedLine)
		if res != "" {
			results = append(results, res)
		}
	}

	os.MkdirAll("results", 0755)
	outputFilename := strings.TrimSuffix(filename, ".txt") + "_result.txt"
	outputPath := "results/" + outputFilename
	outputContent := strings.Join(results, "\n")
	err = os.WriteFile(outputPath, []byte(outputContent), 0644)
	if err != nil {
		fmt.Println("写入结果失败：", err)
		return
	}

	fmt.Printf("处理完成，共计算正确 %d 道题\n", len(results))
	fmt.Printf("结果在：%s\n", outputPath)
}

func processOneLine(line string) string {
	cleanLine := strings.ReplaceAll(line, " ", "")
	if cleanLine == "" {
		return ""
	}

	opPos := -1
	var op string
	for i := 0; i < len(cleanLine); i++ {
		c := cleanLine[i]
		if i == 0 && c == '-' {
			continue
		}
		if c == '+' || c == '-' || c == '*' || c == '/' {
			opPos = i
			op = string(c)
			break
		}
	}

	if opPos == -1 {
		return ""
	}

	num1Str := cleanLine[:opPos]
	num2Str := cleanLine[opPos+1:]

	num1, err1 := strconv.ParseFloat(num1Str, 64)
	num2, err2 := strconv.ParseFloat(num2Str, 64)
	if err1 != nil || err2 != nil {
		return ""
	}

	var result float64
	switch op {
	case "+":
		result = num1 + num2
	case "-":
		result = num1 - num2
	case "*":
		result = num1 * num2
	case "/":
		if num2 == 0 {
			return ""
		}
		result = num1 / num2
	default:
		return ""
	}

	if result == float64(int(result)) {
		return fmt.Sprintf("%s=%d", cleanLine, int(result))
	} else {
		return fmt.Sprintf("%s=%.2f", cleanLine, result)
	}
}
