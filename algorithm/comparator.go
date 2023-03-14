package algorithm

// Comparator 比较器用于比较两个值
type Comparator interface {
	// Compare v1 and v2
	// Ascending order: should return 1 -> v1 > v2, 0 -> v1 = v2, -1 -> v1 < v2
	// Descending order: should return 1 -> v1 < v2, 0 -> v1 = v2, -1 -> v1 > v2
	Compare(v1, v2 any) int
}
