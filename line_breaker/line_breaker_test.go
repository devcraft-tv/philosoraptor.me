package line_breaker

import "testing"

func TestGetLines(t *testing.T) {
	slb := SimpleLineBreaker{}
	line := "This is an awesome sentence that does not make any sense."
	lines := slb.GetLines(line, 22)

	expected := "This is an awesome"
	result := lines[0]

	if result != expected {
		t.Errorf("Broken line should be '%s', got: '%s'", expected, result)
	}
}
