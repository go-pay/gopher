package util

import "sort"

// recommend use IntDeduplicate
// Deprecated
func MergeIntDuplicate(slice []int) (merged []int) {
	return IntDeduplicate(slice)
}

// int 数组去重
func IntDeduplicate(slice []int) (result []int) {
	var dupMap = make(map[int]struct{})
	for _, v := range slice {
		length := len(dupMap)
		dupMap[v] = struct{}{}
		if len(dupMap) != length {
			result = append(result, v)
		}
	}
	return result
}

// int 数组排序+去重
func IntSortDeduplicate(slice []int) (result []int) {
	tmp := make([]int, len(slice))
	copy(tmp, slice)
	sort.Ints(tmp)
	j := 0
	for i := 1; i < len(tmp); i++ {
		if tmp[j] == tmp[i] {
			continue
		}
		j++
		tmp[j] = tmp[i]
	}
	return tmp[:j+1]
}

// recommend use IntMergeDeduplicate
// Deprecated
func MergeSliceRemoveDuplicate(slice1, slice2 []int) (result []int) {
	slice1 = append(slice1, slice2...)
	return IntDeduplicate(slice1)
}

// int 数组合并+去重
func IntMergeDeduplicate(slice1, slice2 []int) (result []int) {
	slice1 = append(slice1, slice2...)
	return IntDeduplicate(slice1)
}

// int 数组合并排序+去重
func IntMergeSortDeduplicate(slice1, slice2 []int) (result []int) {
	slice1 = append(slice1, slice2...)
	return IntSortDeduplicate(slice1)
}

// recommend use StringDeduplicate
// Deprecated
func MergeStringDuplicate(slice []string) (result []string) {
	return StringDeduplicate(slice)
}

// string 数组去重
func StringDeduplicate(slice []string) (result []string) {
	var dupMap = make(map[string]struct{})
	for _, v := range slice {
		length := len(dupMap)
		dupMap[v] = struct{}{}
		if len(dupMap) != length {
			result = append(result, v)
		}
	}
	return result
}

// string 数组排序+去重
func StringSortDeduplicate(slice []string) (result []string) {
	tmp := make([]string, len(slice))
	copy(tmp, slice)
	sort.Strings(tmp)
	j := 0
	for i := 1; i < len(tmp); i++ {
		if tmp[j] == tmp[i] {
			continue
		}
		j++
		tmp[j] = tmp[i]
	}
	return tmp[:j+1]
}

// recommend use StringMergeDeduplicate
// Deprecated
func MergeStringSliceRemoveDuplicate(slice1, slice2 []string) (result []string) {
	return StringMergeDeduplicate(slice1, slice2)
}

// string 数组合并+去重
func StringMergeDeduplicate(slice1, slice2 []string) (result []string) {
	slice1 = append(slice1, slice2...)
	return StringDeduplicate(slice1)
}

// string 数组合并排序+去重
func StringMergeSortDeduplicate(slice1, slice2 []string) (result []string) {
	slice1 = append(slice1, slice2...)
	return StringSortDeduplicate(slice1)
}

// 过滤数组，去除src在dst中存在的item
// src[1,2,3,4,5]   dst[2,4,6,8]	result[1,3,5]
func FilterIntSlice(src []int, dst []int) (result []int) {
	aMap := make(map[int]struct{})
	result = make([]int, 0)
	for _, v := range dst {
		aMap[v] = struct{}{}
	}
	for _, v := range src {
		if _, has := aMap[v]; !has {
			result = append(result, v)
		}
	}
	return result
}

// 过滤数组，去除src在dst中存在的item
// src["a","b","c","d","e"]   dst["b","d","f","h"]	result["a","c","e"]
func FilterStringSlice(src []string, dst []string) (result []string) {
	aMap := make(map[string]struct{})
	result = make([]string, 0)
	for _, v := range dst {
		aMap[v] = struct{}{}
	}
	for _, v := range src {
		if _, has := aMap[v]; !has {
			result = append(result, v)
		}
	}
	return result
}
