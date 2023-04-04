package algorithm

// Search algorithms see https://github.com/TheAlgorithms/Go/tree/master/search

// LinearSearch 线性搜索
// 如果找不到返回 -1
func LinearSearch[T any](slice []T, target T, equal func(a, b T) bool) int {
	for i, v := range slice {
		if equal(v, target) {
			return i
		}
	}
	return -1
}

// BinarySearch  二分查找
// 如果找不到返回 -1
func BinarySearch[T any](sortedSlice []T, target T, lowIndex, highIndex int, comparator Comparator) int {
	if highIndex < lowIndex || len(sortedSlice) == 0 {
		return -1
	}

	midIndex := lowIndex + (highIndex-lowIndex)/2
	isMidValGreatTarget := comparator.Compare(sortedSlice[midIndex], target) == 1
	isMidValLessTarget := comparator.Compare(sortedSlice[midIndex], target) == -1

	if isMidValGreatTarget {
		return BinarySearch(sortedSlice, target, lowIndex, midIndex-1, comparator)
	} else if isMidValLessTarget {
		return BinarySearch(sortedSlice, target, midIndex+1, highIndex, comparator)
	}

	return midIndex
}

// BinaryIterativeSearch 二分迭代查找
// 如果找不到返回 -1
func BinaryIterativeSearch[T any](sortedSlice []T, target T, lowIndex, highIndex int, comparator Comparator) int {
	startIndex := lowIndex
	endIndex := highIndex

	var midIndex int
	for startIndex <= endIndex {
		midIndex = startIndex + (endIndex-startIndex)/2
		isMidValGreatTarget := comparator.Compare(sortedSlice[midIndex], target) == 1
		isMidValLessTarget := comparator.Compare(sortedSlice[midIndex], target) == -1

		if isMidValGreatTarget {
			endIndex = midIndex - 1
		} else if isMidValLessTarget {
			startIndex = midIndex + 1
		} else {
			return midIndex
		}
	}
	return -1
}
