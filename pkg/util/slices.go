package util

func InSlice[t comparable](a t, b []t) bool {
	for _, v := range b {
		if a == v {
			return true
		}
	}
	return false
}
