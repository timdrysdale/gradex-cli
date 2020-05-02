package ingester

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoneFilePath(t *testing.T) {
	assert.Equal(t, ".blah.done", doneFilePath("blah.pdf"))
	assert.Equal(t, "/foo/bar/.dab.done", doneFilePath("/foo/bar/dab.txt"))
}
func TestIsZIP(t *testing.T) {

	assert.True(t, IsZIP("foo.ZIP"))
	assert.True(t, IsZIP("foo.ZIp"))
	assert.True(t, IsZIP("foo.zip"))
	assert.True(t, IsZIP("bar.foo.zip"))
	assert.False(t, IsZIP("bar.zip.foo"))
}
func TestIsTXT(t *testing.T) {
	assert.True(t, IsTXT("foo.TXT"))
	assert.True(t, IsTXT("foo.TXt"))
	assert.True(t, IsTXT("foo.txt"))
	assert.False(t, IsTXT("bar.foo.text"))
	assert.False(t, IsTXT("bar.foo.teXt"))
	assert.False(t, IsTXT("bar.TEXT.zip"))

}
func TestIsPDF(t *testing.T) {

	assert.True(t, IsPDF("foo.PDF"))
	assert.True(t, IsPDF("foo.PDf"))
	assert.True(t, IsPDF("foo.pdf"))
	assert.True(t, IsPDF("bar.foo.pDF"))
	assert.False(t, IsPDF("bar.PDF.foo"))
}

func TestBareFile(t *testing.T) {

	assert.Equal(t, BareFile("./test/some/foo.pdf"), "foo")
	assert.Equal(t, BareFile("./test/some.foo.pdf"), "some.foo")

}
