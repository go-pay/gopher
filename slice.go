package goutil

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
