package coverage

// Size returns a string characterizing the size of the given integer.
// it is used here to demonstrate test code coverage
func Size(i int) string {
	switch {
	case i < 0:
		return "negative"
	case i == 0:
		return "zero"
	case i < 10:
		return "small"
	case i < 100:
		return "big"
	case i < 1000:
		return "huge"
	}
	return "enormous"
}
