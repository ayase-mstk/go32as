package parse

func isAlpha(b byte) bool {
	return ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z')
}

func isAlphaStr(s string) bool {
	for i := 0; i < len(s); i++ {
		if !isAlpha(s[i]) {
			return false
		}
	}
	return true
}

func isNumeric(b byte) bool {
	return '0' <= b && b <= '9'
}

func isNumericStr(s string) bool {
	for i := 0; i < len(s); i++ {
		if !isNumeric(s[i]) {
			return false
		}
	}
	return true
}
