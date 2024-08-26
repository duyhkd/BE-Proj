package main

import (
	"fmt"
)

func main() {
	calculateRect(2, 3)
	evenLengthString("")
	evenLengthString("1")
	sliceStats([]int{28, 32, 1, 16, 8, 2, 64, 4})
}

func calculateRect(w, h int) {
	perimeter := 2 * (w + h)
	fmt.Printf("Perimeter: %v\n", perimeter)
	area := w * h
	fmt.Printf("Area: %v\n", area)
}

func evenLengthString(str string) {
	isEven := len(str)%2 == 0
	fmt.Printf("Is even: %v\n", isEven)
}

func sliceStats(slice []int) {
	var sum, max, min int
	for _, v := range slice {
		sum += v
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	average := sum / len(slice)
	sortedSlice := mergeSort(slice)
	fmt.Printf("Sum: %v\nMax: %v\nMin: %v\nAverage: %v\nSorted slice: %v\n", sum, max, min, average, sortedSlice)
}

func mergeSort(slice []int) []int {
	length := len(slice)
	var slice1, slice2 []int
	if length < 2 {
		return slice
	} else {
		slice1 = mergeSort(slice[0 : length/2])
		slice2 = mergeSort(slice[length/2 : length])
	}
	mergedSlice := merge(slice1, slice2)

	return mergedSlice
}

func merge(slice1, slice2 []int) []int {
	length1 := len(slice1)
	length2 := len(slice2)
	var mergedSlice []int
	var i, j int
	for i < length1 || j < length2 {
		if i == length1 {
			mergedSlice = append(mergedSlice, slice2[j:]...)
			j = length2
		} else if j == length2 {
			mergedSlice = append(mergedSlice, slice1[i:]...)
			i = length1
		} else if slice1[i] < slice2[j] {
			mergedSlice = append(mergedSlice, slice1[i])
			i++
		} else {
			mergedSlice = append(mergedSlice, slice2[j])
			j++
		}
	}
	return mergedSlice
}

func twoSum(nums []int, target int) []int {
	var result []int
	m := make(map[int][]int)
	for i, v := range nums {
		// Slice of index so we can handle dupes
		m[v] = append(m[v], i)
	}
	for i, v := range nums {
		term2 := target - v
		val, ok := m[term2]
		if ok && len(val) > 0 {
			// if duplicate means we have 2 val with different index
			if i == val[len(val)-1] && len(val) > 1 {
				result = append(result, i, val[len(val)-2])
				break
			} else if i != val[len(val)-1] {
				result = append(result, i, val[len(val)-1])
				break
			}

		}
	}
	return result
}
