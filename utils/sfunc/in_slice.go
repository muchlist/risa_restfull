package sfunc

// InSlice seperti fungsi in, apakah target ada didalam slice.
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

// ValueInSliceIsAvailable memasukkan input request berupa slice dan
// membandingkan isi slicenya apakah tersedia untuk digunakan
func ValueInSliceIsAvailable(inputSlice []string, availableSlice []string) bool {
	for _, input := range inputSlice {
		if !InSlice(input, availableSlice) {
			return false
		}
	}
	return true
}
