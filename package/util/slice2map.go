package util

func Slice2Map(s []int64) map[int64]struct{} {
	res := make(map[int64]struct{}, len(s))
	for _, v := range s {
		res[v] = struct{}{}
	}
	return res
}
