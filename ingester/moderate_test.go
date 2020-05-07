package ingester

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModeratePercentage(t *testing.T) {

	assert.Equal(t, 20.0, requiredPercent(50, 10, 10))
	assert.Equal(t, 10.0, requiredPercent(150, 10, 10))

}

func TestSelect(t *testing.T) {

	fm := make(map[string]bool)

	for _, char := range loweralpha() {
		fm[string(char)] = false
	}

	selectByPercent(&fm, 20)

	assert.Equal(t, 6, countSelected(fm))

	for _, char := range loweralpha() {
		fm[string(char)] = false
	}

	selectByPercent(&fm, 10)

	assert.Equal(t, 3, countSelected(fm))

	// this should just select all and return
	selectByPercent(&fm, 99)
	assert.Equal(t, 26, countSelected(fm))

}

//https://rosettacode.org/wiki/Generate_lower_case_ASCII_alphabet#Go
func loweralpha() string {
	p := make([]byte, 26)
	for i := range p {
		p[i] = 'a' + byte(i)
	}
	return string(p)
}

func countSelected(fileMap map[string]bool) int {

	selected := 0

	for _, v := range fileMap {
		if v {
			selected++
		}
	}

	return selected
}
