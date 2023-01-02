package main

// Returns the lesser of two integers
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Returns the greater of two integers
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func sliceContains(s []coords, v coords) bool {
	for _, i := range s {
		if i == v {
			return true
		}
	}
	return false
}
