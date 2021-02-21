package utils

func InSlice(target string, slice []string) bool {
	if len(slice) == 0 {
		return false
	}

	for _, value := range slice {
		if target == value {
			return true
		}
	}

	return false
}
