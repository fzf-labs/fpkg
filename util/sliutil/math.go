package sliutil

import (
	"math"

	"golang.org/x/exp/constraints"
)

// Sum 求和是所有元素的和。
func Sum[T constraints.Integer | constraints.Float](ss []T) (sum T) {
	for _, s := range ss {
		sum += s
	}

	return
}

// Average 元素的平均值，如果没有则为零
func Average[T constraints.Integer | constraints.Float](ss []T) float64 {
	if l := len(ss); l > 0 {
		return float64(Sum(ss)) / float64(l)
	}

	return 0
}

// Abs 返回绝对值。
func Abs[T constraints.Integer | constraints.Float](val T) T {
	if val < 0 {
		return -val
	}

	return val
}

// Max 是最大值，或零。
func Max[T constraints.Ordered](ss []T) (min T) {
	if len(ss) == 0 {
		return
	}

	min = ss[0]
	for _, s := range ss {
		if s > min {
			min = s
		}
	}

	return
}

// Min 是最小值，或零。
func Min[T constraints.Ordered](ss []T) (min T) {
	if len(ss) == 0 {
		return
	}

	min = ss[0]
	for _, s := range ss {
		if s < min {
			min = s
		}
	}

	return
}

// Median  中位数
// 中位数返回数据样本的上半部分与下半部分之间的分隔值。如果切片中没有元素，则返回0。如果元素数量为偶数，则返回两个“中值”的ElementType平均值。
func Median[T constraints.Integer | constraints.Float](ss []T) T {
	n := len(ss)
	if n == 0 {
		return 0
	}

	if n == 1 {
		return ss[0]
	}

	// This implementation aims at linear time O(n) on average.
	// It uses the same idea as QuickSort, but makes only 1 recursive
	// call instead of 2. See also Quickselect.

	work := make([]T, len(ss))
	copy(work, ss)

	limit1, limit2 := n/2, n/2+1
	if n%2 == 0 {
		limit1, limit2 = n/2-1, n/2+1
	}

	var rec func(a, b int)
	rec = func(a, b int) {
		if b-a <= 1 {
			return
		}
		ipivot := (a + b) / 2
		pivot := work[ipivot]
		work[a], work[ipivot] = work[ipivot], work[a]
		j := a
		k := b
		for j+1 < k {
			if work[j+1] < pivot {
				work[j+1], work[j] = work[j], work[j+1]
				j++
			} else {
				work[j+1], work[k-1] = work[k-1], work[j+1]
				k--
			}
		}
		// 1 or 0 recursive calls
		if j > limit1 {
			rec(a, j)
		}
		if j+1 < limit2 {
			rec(j+1, b)
		}
	}

	rec(0, len(work))

	if n%2 == 1 {
		return work[n/2]
	} else {
		return (work[n/2-1] + work[n/2]) / 2
	}
}

// Product 所有元素的乘积。
func Product[T constraints.Integer | constraints.Float](ss []T) (product T) {
	if len(ss) == 0 {
		return
	}

	product = ss[0]
	for _, s := range ss[1:] {
		product *= s
	}

	return
}

// StandardDeviation 标准差
func StandardDeviation[T constraints.Integer | constraints.Float](ss []T) float64 {
	if len(ss) == 0 {
		return 0.0
	}

	avg := Average(ss)

	var sd float64
	for i := range ss {
		sd += math.Pow(float64(ss[i])-avg, 2)
	}
	sd = math.Sqrt(sd / float64(len(ss)))

	return sd
}
