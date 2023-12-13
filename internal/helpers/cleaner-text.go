package helper

func CutLineBreak(s string) string {
	l, r := 0, len(s)-1
	for i := range s {
		if s[i] == '\n' || s[i] == ' ' || s[i] == '.' || s[i] == ',' {
			l++
		} else if s[r] == '\n' || s[r] == ' ' || s[r] == ',' || s[r] == '.' {
			r--
		} else {
			break
		}
	}
	return s[l : r+1]
}
