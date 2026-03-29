package main

import "testing"

// 接下来开始写单侧
// 首先是isPrime函数
func TestIsPrime(t *testing.T) {
	testCases := []struct {
		name string
		num  int
		want bool
	}{
		{"负数不是素数", -5, false},
		{"0不是素数", 0, false},
		{"1不是素数", 1, false},
		{"2是素数", 2, true},
		{"3是素数", 3, true},
		{"4不是素数", 4, false},
		{"5是素数", 5, true},
		{"9不是素数", 9, false},
		{"17是素数", 17, true},
		{"100不是素数", 100, false},
		{"997是素数", 997, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := isPrime(tc.num)
			if got != tc.want {
				t.Errorf("isPrime(%d) = %v, 期望 %v", tc.num, got, tc.want)
			}
		})
	}
}

// 接下来是测试findPrimes函数
func TestFindPrimes(t *testing.T) {
	testCases := []struct {
		name  string
		start int
		end   int
		want  []int
	}{
		{"1到10的素数", 1, 10, []int{2, 3, 5, 7}},
		{"10到20的素数", 10, 20, []int{11, 13, 17, 19}},
		{"20到30的素数", 20, 30, []int{23, 29}},
		{"开始值大于结束值", 5, 3, []int{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := findPrimes(tc.start, tc.end)
			if len(got) != len(tc.want) {
				t.Errorf("findPrimes(%d, %d) 长度不匹配: 实际 %d, 期望 %d",
					tc.start, tc.end, len(got), len(tc.want))
				return
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("findPrimes(%d, %d)[%d] = %d, 期望 %d",
						tc.start, tc.end, i, got[i], tc.want[i])
				}
			}
		})
	}
}
