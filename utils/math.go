package utils

type IntX interface {
	int8 | uint8 | int | uint | int32 | uint32 | int64 | uint64 | int16 | uint16
}

type FloatX interface {
	float32 | float64
}

type Number interface {
	IntX | FloatX
}

// Max 返回大的数
func Max[T Number](p1, p2 T) T {
	if p1 > p2 {
		return p1
	}

	return p2
}

// Min 返回小的数
func Min[T Number](p1, p2 T) T {
	if p1 < p2 {
		return p1
	}

	return p2
}
