package image

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

//https://stackoverflow.com/questions/31663229/how-can-stdout-be-captured-or-suppressed-for-golang-testing
func TestCheck(t *testing.T) {
	defer quiet()()
	pdf1 := "./test/test.pdf"
	pdf2 := "./test/test-mod.pdf"

	//same as self should be true
	result, err := VisuallyIdenticalMultiPagePDF(pdf1, pdf1)
	assert.NoError(t, err)
	assert.True(t, result)

	//same as self should be true
	result, err = VisuallyIdenticalMultiPagePDF(pdf2, pdf2)
	assert.NoError(t, err)
	assert.True(t, result)

	//different, due to extra . on page one, so should be false
	result, err = VisuallyIdenticalMultiPagePDF(pdf1, pdf2)
	assert.NoError(t, err)
	assert.False(t, result)

}
func quiet() func() {
	null, _ := os.Open(os.DevNull)
	sout := os.Stdout
	serr := os.Stderr
	os.Stdout = null
	os.Stderr = null
	log.SetOutput(null)
	return func() {
		defer null.Close()
		os.Stdout = sout
		os.Stderr = serr
		log.SetOutput(os.Stderr)
	}
}
