//nolint:gosec
package algorithm

import (
	"math/rand"
	"time"
)

// DoubleAverage 二倍均值算法 抢红包
func DoubleAverage(amount float64, num int) []float64 {
	result := make([]float64, num)
	if num == 0 {
		return result
	}
	if num == 1 {
		result[0] = amount
		return result
	}
	remainAmount := int(amount * 100)
	for i := 0; i < num; i++ {
		remainNum := num - i
		if remainNum == 1 {
			result[i] = float64(remainAmount) / 100.0
		} else {
			m := remainAmount / remainNum * 2
			money := 1 + rand.New(rand.NewSource(time.Now().UnixNano())).Intn(m-1)
			remainAmount -= money
			result[i] = float64(money) / 100.0
		}
	}
	return result
}
