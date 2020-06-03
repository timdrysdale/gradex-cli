package ingester

import (
	"fmt"
	"math"
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

func createTestFileMap(size int) map[string]bool {

	testSet := make(map[string]bool)

	for i := 0; i < size; i++ {

		testSet[safeUUID()] = false
	}

	return testSet

}

func TestCreateTestFileMap(t *testing.T) {

	for i := 0; i < 200; i++ {
		assert.Equal(t, i, len(createTestFileMap(i)))
	}

}

func testSelectSetsAbove100(t *testing.T) {
	//not implemented yet
}

func howManyToMakeActive(i int) int {

	reqdPercent := requiredPercent(i, 10, 10)

	activeCount := int(math.Round(float64(i) * reqdPercent / 100))

	if activeCount > i {
		activeCount = i
	}

	return activeCount
}

func countActive(fm map[string]bool) int {

	count := 0

	for _, f := range fm {
		if f {
			count++
		}
	}

	return count
}

func TestSelectSetsBelow10(t *testing.T) {
	// since we want 10percent or 10, whichever is greater
	// in this band it is always 10  -easy!

	failedSizes := []int{}

	for setSize := 1; setSize <= 10; setSize++ {

		fm := createTestFileMap(setSize)

		assert.Equal(t, setSize, len(fm))

		assert.Equal(t, 0, countActive(fm))

		selectByPercent(&fm, requiredPercent(setSize, 10, 10))

		expectedActiveCount := setSize

		countOK := expectedActiveCount == countActive(fm) || (expectedActiveCount+1) == countActive(fm)

		if !countOK {
			failedSizes = append(failedSizes, countActive(fm))
		}
	}

	assert.Equal(t, 0, len(failedSizes))

	if len(failedSizes) > 0 {
		fmt.Println(failedSizes)
	}
}

func TestSelectSets11To100(t *testing.T) {
	// since we want 10percent or 10, whichever is greater
	// in this band it is always 10  -easy!

	expectedActiveCount := 10

	failedSizes := []int{}

	for setSize := 11; setSize <= 100; setSize++ {

		fm := createTestFileMap(setSize)

		assert.Equal(t, setSize, len(fm))

		assert.Equal(t, 0, countActive(fm))

		selectByPercent(&fm, requiredPercent(setSize, 10, 10))

		countOK := expectedActiveCount == countActive(fm) || (expectedActiveCount+1) == countActive(fm)

		if !countOK {
			failedSizes = append(failedSizes, countActive(fm))
		}
	}

	assert.Equal(t, 0, len(failedSizes))

	if len(failedSizes) > 0 {
		fmt.Println(failedSizes)
	}
}

func TestSelectSetsAbove100(t *testing.T) {
	// since we want 10percent or 10, whichever is greater
	// in this band it is always 10  -easy!

	failedSizes := []int{}

	for setSize := 100; setSize <= 300; setSize++ {

		fm := createTestFileMap(setSize)

		assert.Equal(t, setSize, len(fm))

		assert.Equal(t, 0, countActive(fm))

		selectByPercent(&fm, requiredPercent(setSize, 10, 10))

		expectedActiveCount := int(0.1 * float64(setSize))

		countOK := expectedActiveCount == countActive(fm) || (expectedActiveCount+1) == countActive(fm)

		if !countOK {
			failedSizes = append(failedSizes, countActive(fm))
		}
	}

	assert.Equal(t, 0, len(failedSizes))

	if len(failedSizes) > 0 {
		fmt.Println(failedSizes)
	}
}

//
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
