package utils

import "golang.org/x/exp/constraints"

// InSlice 判断数据是否在切片存在
func InSlice[T comparable](needle T, hyStack []T) bool {
	for _, v := range hyStack {
		if needle == v {
			return true
		}
	}

	return false
}

func SliceMaximum[T constraints.Ordered](s []T) T {
	if len(s) == 0 {
		var zero T
		return zero
	}
	m := s[0]
	for _, v := range s {
		if m < v {
			m = v
		}
	}
	return m
}

func SliceMinimum[T constraints.Ordered](s []T) T {
	if len(s) == 0 {
		var zero T
		return zero
	}
	m := s[0]
	for _, v := range s {
		if m > v {
			m = v
		}
	}
	return m
}

// SliceUnique 切片去重
func SliceUnique[T comparable](data []T) []T {
	if len(data) > 1024 {
		return sliceUniqueByMap(data)
	} else {
		return sliceUniqueByLoop(data)
	}
}

func sliceUniqueByMap[T comparable](data []T) []T {
	result := make([]T, 0)
	tempMap := make(map[T]byte)

	for _, e := range data {
		l := len(tempMap)
		tempMap[e] = 0

		if len(tempMap) != l {
			result = append(result, e)
		}
	}

	return result
}

func sliceUniqueByLoop[T comparable](data []T) []T {
	result := make([]T, 0)

	for i := range data {
		flag := true
		for j := range result {
			if data[i] == result[j] {
				flag = false
				break
			}
		}
		if flag {
			result = append(result, data[i])
		}
	}

	return result
}

// SliceIndex2Map 通过主键把结构体slice放到map
func SliceIndex2Map[T any, K constraints.Ordered](data []*T, f func(data T) (index K)) map[K]*T {
	sm := make(map[K]*T)
	if len(data) < 1 {
		return sm
	}
	for _, v := range data {
		key := f(*v)
		sm[key] = v
	}
	return sm
}
