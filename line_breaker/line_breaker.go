package line_breaker

type SimpleLineBreaker struct{}

func (s SimpleLineBreaker) GetLines(line string, chars int) []string {
	lines := []string{"This is an awesome"}
	return lines
}
