package line_breaker

import "strings"

type SimpleLineBreaker struct{}

func (s SimpleLineBreaker) GetLines(text string, maxCharacters int) (lines []string) {
	words := strings.Split(text, " ")
	var line string

	for idx, word := range words {

		if wordLength := len(word); wordLength+1 > spaceAvailable(maxCharacters, line) {
			appendAndFlush(&line, &lines, word)
		} else {
			line = line + " " + word
		}

		if idx == len(words)-1 {
			appendAndFlush(&line, &lines, word)
		}
	}

	return lines
}

func appendAndFlush(line *string, lines *[]string, word string) {
	*lines = append(*lines, strings.Trim(*line, " "))
	*line = word
}

func spaceAvailable(maxCharacters int, line string) int {
	return (maxCharacters - len(strings.Trim(line, " ")))
}
