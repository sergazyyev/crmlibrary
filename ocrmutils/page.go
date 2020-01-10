package ocrmutils

func TotalPageNumber(total, offset int) int {
	if total == 0 || offset == 0 {
		return 1
	}
	i := total / offset
	r := total % offset
	if r > 0 {
		i++
	}
	return i
}
