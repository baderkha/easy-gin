package easygin

// SliceContains : iterates over a slice of something compareable and checks if the needle exists in the haystack
func SliceContains[T comparable](haystack []T, needle T) bool {
	for _, v := range haystack {
		if needle == v {
			return true
		}
	}
	return false
}
