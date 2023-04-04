package algorithm

// 排序算法
// https://www.runoob.com/w3cnote/ten-sorting-algorithm.html

// BubbleSort 冒泡排序
func BubbleSort[T any](slice []T, comparator Comparator) {
	for i := 0; i < len(slice); i++ {
		for j := 0; j < len(slice)-1-i; j++ {
			isCurrGreatThanNext := comparator.Compare(slice[j], slice[j+1]) == 1
			if isCurrGreatThanNext {
				swap(slice, j, j+1)
			}
		}
	}
}

// InsertionSort 插入排序算法
func InsertionSort[T any](slice []T, comparator Comparator) {
	for i := 0; i < len(slice); i++ {
		for j := i; j > 0; j-- {
			isPreLessThanCurrent := comparator.Compare(slice[j], slice[j-1]) == -1
			if isPreLessThanCurrent {
				swap(slice, j, j-1)
			} else {
				break
			}
		}
	}
}

// SelectionSort 选择排序算法
func SelectionSort[T any](slice []T, comparator Comparator) {
	for i := 0; i < len(slice); i++ {
		min := i
		for j := i + 1; j < len(slice); j++ {
			if comparator.Compare(slice[j], slice[min]) == -1 {
				min = j
			}
		}
		swap(slice, i, min)
	}
}

// ShellSort 希尔排序
func ShellSort[T any](slice []T, comparator Comparator) {
	size := len(slice)

	gap := 1
	for gap < size/3 {
		gap = 3*gap + 1
	}

	for gap >= 1 {
		for i := gap; i < size; i++ {
			for j := i; j >= gap && comparator.Compare(slice[j], slice[j-gap]) == -1; j -= gap {
				swap(slice, j, j-gap)
			}
		}
		gap = gap / 3
	}
}

// QuickSort 快速排序
func QuickSort[T any](slice []T, comparator Comparator) {
	quickSort(slice, 0, len(slice)-1, comparator)
}

func quickSort[T any](slice []T, lowIndex, highIndex int, comparator Comparator) {
	if lowIndex < highIndex {
		p := partition(slice, lowIndex, highIndex, comparator)
		quickSort(slice, lowIndex, p-1, comparator)
		quickSort(slice, p+1, highIndex, comparator)
	}
}

// 将切片分成两部分
func partition[T any](slice []T, lowIndex, highIndex int, comparator Comparator) int {
	p := slice[highIndex]
	i := lowIndex
	for j := lowIndex; j < highIndex; j++ {
		if comparator.Compare(slice[j], p) == -1 { //slice[j] < p
			swap(slice, i, j)
			i++
		}
	}

	swap(slice, i, highIndex)

	return i
}

// HeapSort 堆排序
func HeapSort[T any](slice []T, comparator Comparator) {
	size := len(slice)

	for i := size/2 - 1; i >= 0; i-- {
		sift(slice, i, size-1, comparator)
	}
	for j := size - 1; j > 0; j-- {
		swap(slice, 0, j)
		sift(slice, 0, j-1, comparator)
	}
}

func sift[T any](slice []T, lowIndex, highIndex int, comparator Comparator) {
	i := lowIndex
	j := 2*i + 1

	temp := slice[i]
	for j <= highIndex {
		if j < highIndex && comparator.Compare(slice[j], slice[j+1]) == -1 { //slice[j] < slice[j+1]
			j++
		}
		if comparator.Compare(temp, slice[j]) == -1 { //tmp < slice[j]
			slice[i] = slice[j]
			i = j
			j = 2*i + 1
		} else {
			break
		}
	}
	slice[i] = temp
}

// MergeSort 归并排序
func MergeSort[T any](slice []T, comparator Comparator) {
	mergeSort(slice, 0, len(slice)-1, comparator)
}

func mergeSort[T any](slice []T, lowIndex, highIndex int, comparator Comparator) {
	if lowIndex < highIndex {
		mid := (lowIndex + highIndex) / 2
		mergeSort(slice, lowIndex, mid, comparator)
		mergeSort(slice, mid+1, highIndex, comparator)
		merge(slice, lowIndex, mid, highIndex, comparator)
	}
}

func merge[T any](slice []T, lowIndex, midIndex, highIndex int, comparator Comparator) {
	i := lowIndex
	j := midIndex + 1
	temp := []T{}

	for i <= midIndex && j <= highIndex {
		//slice[i] < slice[j]
		if comparator.Compare(slice[i], slice[j]) == -1 {
			temp = append(temp, slice[i])
			i++
		} else {
			temp = append(temp, slice[j])
			j++
		}
	}

	if i <= midIndex {
		temp = append(temp, slice[i:midIndex+1]...)
	} else {
		temp = append(temp, slice[j:highIndex+1]...)
	}

	for k := 0; k < len(temp); k++ {
		slice[lowIndex+k] = temp[k]
	}
}

// CountSort 计数排序
func CountSort[T any](slice []T, comparator Comparator) []T {
	size := len(slice)
	out := make([]T, size)

	for i := 0; i < size; i++ {
		count := 0
		for j := 0; j < size; j++ {
			//slice[i] > slice[j]
			if comparator.Compare(slice[i], slice[j]) == 1 {
				count++
			}
		}
		out[count] = slice[i]
	}

	return out
}

// 在索引i和j处交换两个切片值
func swap[T any](slice []T, i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
