package gcfg

import "math"

type IntW struct {
	// 整数原值
	IntOri int
	// 原值除以10000之后向上取整的值
	IntCalCeil int
	// 原值除以10000之后向下取整的值
	IntCalFloor int
	// 原值除以10000之后四舍五入取整的值
	IntCalRound int
	// 原值除以10000后的浮点数值
	FloatCal float32
}

func NewIntW(val int32) IntW {
	floatVal := float64(val) / 10000

	return IntW{
		int(val),
		int(math.Ceil(floatVal)),
		int(math.Floor(floatVal)),
		int(math.Round(floatVal)),
		float32(floatVal),
	}
}
