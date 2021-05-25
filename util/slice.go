package util

// int数组合并，去重复
func MergeSliceRemoveDuplicate(slice1, slice2 []int) (merged []int) {
	var dupMap = make(map[int]int)
	slice1 = append(slice1, slice2...)
	for _, v := range slice1 {
		length := len(dupMap)
		dupMap[v] = 1
		if len(dupMap) != length {
			merged = append(merged, v)
		}
	}
	return merged
}

// 过滤数组 去除src中item，在dst中存在的item
// src[1,2,3,4,5]   dst[2,4,6,8]	result[1,3,5]
func FilterSlice(src []int, dst []int) (result []int) {
	aMap := make(map[int]struct{})
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
