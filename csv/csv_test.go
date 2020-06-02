package csv

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {

	s := New()

	assert.Equal(t, reflect.TypeOf(&Store{}), reflect.TypeOf(s))

	s.SetFixedHeader([]string{"what", "who", "when"})
	s.SetRequiredHeader([]string{"A1", "B1", "B2", "B3"})
	line1 := s.Add()

	line1.Add("what", "Foo")
	line1.Add("when", "May")
	line1.Add("who", "AAA")
	line1.Add("A1", "10")
	line1.Add("B1", "10")
	line1.Add("B2", "10")

	line2 := s.Add()
	line2.Add("what", "Foo")
	line2.Add("when", "May")
	line2.Add("who", "AAB")
	line2.Add("B3", "10")
	line2.Add("A1", "9")
	line2.Add("B2", "10")

	expected := `what,who,when,A1,B1,B2,B3
Foo,AAA,May,10,10,10,-
Foo,AAB,May,9,-,10,10`

	assert.Equal(t, expected, s.ToString())
	f, err := os.OpenFile("test.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	assert.NoError(t, err)
	defer f.Close()

	_, err = s.WriteCSV(f)

	assert.NoError(t, err)

}
