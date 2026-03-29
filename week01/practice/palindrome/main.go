package practice

import (
	"unicode"
)

// IsPalindrome 判断字符串是否为回文（忽略大小写、空格和标点）
func IsPalindrome(s string) bool {
	// 处理字符串：过滤空格和标点，转为小写
	var processed []rune
	for _, r := range s {
		if unicode.IsSpace(r) || unicode.IsPunct(r) {
			continue // 忽略空格和标点
		}
		processed = append(processed, unicode.ToLower(r)) // 统一转为小写
	}

	// 判断处理后的字符串是否对称
	length := len(processed)
	for i := 0; i < length/2; i++ {
		if processed[i] != processed[length-1-i] {
			return false
		}
	}
	return true
}

// CountCharacters 统计字符串中每个Unicode字符的出现次数
func CountCharacters(s string) map[rune]int {
	counts := make(map[rune]int)
	// 遍历字符串中的每个rune（支持Unicode字符）
	for _, r := range s {
		counts[r]++
	}
	return counts
}
