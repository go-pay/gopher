package util

import (
	"sort"
	"testing"
)

func TestIntDeduplicate(t *testing.T) {
	in := []int{3, 2, 1, 4, 3, 2, 1, 4, 1} // any item can be sorted
	merged := IntDeduplicate(in)
	t.Logf("in: %d", in)         // in: [3 2 1 4 3 2 1 4 1]
	t.Logf("merged: %d", merged) // merged: [3 2 1 4]
}

func BenchmarkIntDeduplicate(b *testing.B) {
	in := []int{3, 2, 1, 4, 3, 2, 1, 4, 1} // any item can be sorted
	for i := 0; i < b.N; i++ {
		IntDeduplicate(in)
	}
}

func TestIntSortDeduplicate(t *testing.T) {
	in := []int{3, 2, 1, 4, 3, 2, 1, 4, 1} // any item can be sorted
	merged := IntSortDeduplicate(in)
	t.Logf("in: %d", in)              // in: [3 2 1 4 3 2 1 4 1]
	t.Logf("sort merged: %d", merged) // sort merged: [1 2 3 4]
}

func BenchmarkIntSortDeduplicate(b *testing.B) {
	in := []int{3, 2, 1, 4, 3, 2, 1, 4, 1} // any item can be sorted
	for i := 0; i < b.N; i++ {
		IntSortDeduplicate(in)
	}
}

func TestIntMergeDeduplicate(t *testing.T) {
	in := []int{3, 3, 5, 7, 14, 11, 13, 15, 12}  // slice1
	in2 := []int{3, 4, 5, 7, 18, 11, 22, 15, 35} // slice2
	merged := IntMergeDeduplicate(in, in2)
	sort.Ints(merged)
	t.Logf("merged: %d", merged) // merged: [3 4 5 7 11 12 13 14 15 18 22 35]
}

func BenchmarkIntMergeDeduplicate(b *testing.B) {
	in := []int{3, 3, 5, 7, 14, 11, 13, 15, 12}  // slice1
	in2 := []int{3, 4, 5, 7, 18, 11, 22, 15, 35} // slice2
	for i := 0; i < b.N; i++ {
		IntMergeDeduplicate(in, in2)
	}
}

func TestMergeStringDuplicate(t *testing.T) {
	ins := []string{"abc", "hello", "fhgk", "jerry", "world", "jerry", "abc", "hello"}
	merged := StringDeduplicate(ins)
	t.Logf("in: %s", ins)        // in: [abc hello fhgk jerry world jerry abc hello]
	t.Logf("merged: %s", merged) // merged: [abc hello fhgk jerry world]
}

func BenchmarkStringDeduplicate(b *testing.B) {
	ins := []string{"abc", "hello", "fhgk", "jerry", "world", "jerry", "abc", "hello"}
	for i := 0; i < b.N; i++ {
		StringDeduplicate(ins)
	}
}

func TestStringSortDeduplicate(t *testing.T) {
	ins := []string{"abc", "hello", "fhgk", "jerry", "world", "jerry", "abc", "hello"}
	merged := StringSortDeduplicate(ins)
	t.Logf("in: %s", ins)             // in: [abc hello fhgk jerry world jerry abc hello]
	t.Logf("sort merged: %s", merged) // sort merged: [abc fhgk hello jerry world]
}

func BenchmarkStringSortDeduplicate(b *testing.B) {
	ins := []string{"abc", "hello", "fhgk", "jerry", "world", "jerry", "abc", "hello"}
	for i := 0; i < b.N; i++ {
		StringSortDeduplicate(ins)
	}
}

func TestStringMergeDeduplicate(t *testing.T) {
	ins := []string{"abc", "hello", "fhgk", "jerry", "world", "jerry", "abc", "hello"}
	ins2 := []string{"dfsf", "hello", "qwer", "jerry", "hello", "tom", "abc", "fuck"}
	merged := StringMergeDeduplicate(ins, ins2)
	t.Logf("merged: %s", merged) // merged: [abc hello fhgk jerry world dfsf qwer tom fuck]
}

func BenchmarkStringMergeDeduplicate(b *testing.B) {
	ins := []string{"abc", "hello", "fhgk", "jerry", "world", "jerry", "abc", "hello"}
	ins2 := []string{"dfsf", "hello", "qwer", "jerry", "hello", "tom", "abc", "fuck"}
	for i := 0; i < b.N; i++ {
		StringMergeDeduplicate(ins, ins2)
	}
}

func TestStringMergeSortDeduplicate(t *testing.T) {
	ins := []string{"abc", "hello", "fhgk", "jerry", "world", "jerry", "abc", "hello"}
	ins2 := []string{"dfsf", "hello", "qwer", "jerry", "hello", "tom", "abc", "fuck"}
	merged := StringMergeSortDeduplicate(ins, ins2)
	t.Logf("sort merged: %s", merged) // sort merged: [abc dfsf fhgk fuck hello jerry qwer tom world]
}

func BenchmarkStringMergeSortDeduplicate(b *testing.B) {
	ins := []string{"abc", "hello", "fhgk", "jerry", "world", "jerry", "abc", "hello"}
	ins2 := []string{"dfsf", "hello", "qwer", "jerry", "hello", "tom", "abc", "fuck"}
	for i := 0; i < b.N; i++ {
		StringMergeSortDeduplicate(ins, ins2)
	}
}

func TestStringIntersect(t *testing.T) {
	ins := []string{"abc", "hello", "fhgk", "jerry", "world", "jerry", "abc", "hello"}
	ins2 := []string{"dfsf", "hello", "qwer", "jerry", "hello", "tom", "abc", "fuck"}
	result := StringIntersect(ins, ins2)
	t.Logf("result: %s", result) // result: [abc hello jerry]

}

func TestStringSortIntersect(t *testing.T) {
	ins := []string{"abc", "hello", "fhgk", "jerry", "world", "jerry", "abc", "hello"}
	ins2 := []string{"dfsf", "hello", "qwer", "jerry", "hello", "tom", "abc", "fuck"}
	result := StringSortIntersect(ins, ins2)
	t.Logf("result: %s", result) // result: [abc hello jerry]
}

func TestStringUnion(t *testing.T) {
	ins := []string{"abc", "hello", "fhgk", "jerry", "world", "jerry", "abc", "hello"}
	ins2 := []string{"dfsf", "hello", "qwer", "jerry", "hello", "tom", "abc", "fuck"}
	result := StringUnion(ins, ins2)
	t.Logf("result: %s", result) // result: [hello dfsf tom abc jerry world qwer fuck fhgk]

}

func TestStringSortUnion(t *testing.T) {
	ins := []string{"abc", "hello", "fhgk", "jerry", "world", "jerry", "abc", "hello"}
	ins2 := []string{"dfsf", "hello", "qwer", "jerry", "hello", "tom", "abc", "fuck"}
	result := StringSortUnion(ins, ins2)
	t.Logf("result: %s", result) // result: [abc dfsf fhgk fuck hello jerry qwer tom world]
}

func TestIntIntersect(t *testing.T) {
	in := []int{3, 3, 5, 7, 14, 11, 13, 15, 12}  // slice1
	in2 := []int{3, 4, 5, 7, 18, 11, 22, 15, 35} // slice2
	result := IntIntersect(in, in2)
	t.Logf("result: %d", result) // result: [3 5 7 11 15]
}

func TestIntSortIntersect(t *testing.T) {
	in := []int{3, 3, 5, 7, 14, 11, 13, 15, 12}  // slice1
	in2 := []int{3, 4, 5, 7, 18, 11, 22, 15, 35} // slice2
	result := IntSortIntersect(in, in2)
	t.Logf("result: %d", result) // result: [3 5 7 11 15]
}

func TestIntUnion(t *testing.T) {
	in := []int{1, 2, 3, 4}  // slice1
	in2 := []int{2, 4, 6, 8} // slice2
	result := IntUnion(in, in2)
	t.Logf("result: %d", result) // result: [6 8 1 2 3 4]
}

func TestIntSortUnion(t *testing.T) {
	in := []int{1, 2, 3, 4}  // slice1
	in2 := []int{2, 4, 6, 8} // slice2
	result := IntSortUnion(in, in2)
	t.Logf("result: %d", result) // result: [1 2 3 4 6 8]
}

func TestIntRemoveElementByIndex(t *testing.T) {
	in := []int{1, 2, 3, 4} // slice1
	result := IntRemoveElementByIndex(in, 1)
	t.Logf("in: %d", in)         // result: [1 2 3 4]
	t.Logf("result: %d", result) // result: [1 3 4]
}

func TestIntRemoveElement(t *testing.T) {
	in := []int{1, 2, 3, 4, 4, 5, 4, 4} // slice1
	result := IntRemoveElement(in, 4, 2)
	t.Logf("in: %d", in)         // result: [1 2 3 4 4 5 4 4]
	t.Logf("result: %d", result) // result: [1 2 3 5 4 4]
}

func TestStringRemoveElementByIndex(t *testing.T) {
	ins := []string{"abc", "hello", "fhgk", "jerry", "hello"}
	result := StringRemoveElementByIndex(ins, 3)
	t.Logf("ins: %s", ins)       // result: [abc hello fhgk jerry hello]
	t.Logf("result: %s", result) // result: [abc hello fhgk hello]
}

func TestStringRemoveElement(t *testing.T) {
	ins := []string{"abc", "hello", "abc", "abc", "hello", "abc"}
	result := StringRemoveElement(ins, "abc", 3)
	t.Logf("ins: %s", ins)       // ins: [abc hello abc abc hello abc]
	t.Logf("result: %s", result) // result: [hello hello abc]
}
