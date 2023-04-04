package algorithm

import "sort"

// TopK  取前K个元素,基于切片实现
func TopK[T any](slice []T, k int, comparator Comparator) []T {
	if k > len(slice) {
		k = len(slice)
	}
	sort.Slice(slice, func(i, j int) bool {
		return comparator.Compare(slice[i], slice[j]) == -1
	})
	return slice[:k]
}
