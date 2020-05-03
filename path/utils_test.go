package gradexpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsZip(t *testing.T) {

	assert.True(t, IsZip("foo.ZIP"))
	assert.True(t, IsZip("foo.ZIp"))
	assert.True(t, IsZip("foo.zip"))
	assert.True(t, IsZip("bar.foo.zip"))
	assert.False(t, IsZip("bar.zip.foo"))
}
func TestIsTxt(t *testing.T) {

	assert.True(t, IsTxt("foo.TXT"))
	assert.True(t, IsTxt("foo.TXt"))
	assert.True(t, IsTxt("foo.txt"))
	assert.False(t, IsTxt("bar.foo.text"))
	assert.False(t, IsTxt("bar.foo.teXt"))
	assert.False(t, IsTxt("bar.TEXT.zip"))

}
func TestIsPdf(t *testing.T) {

	assert.True(t, IsPdf("foo.PDF"))
	assert.True(t, IsPdf("foo.PDf"))
	assert.True(t, IsPdf("foo.pdf"))
	assert.True(t, IsPdf("bar.foo.pDF"))
	assert.False(t, IsPdf("bar.PDF.foo"))
}

func TestBareFile(t *testing.T) {

	assert.Equal(t, BareFile("./test/some/foo.pdf"), "foo")
	assert.Equal(t, BareFile("./test/some.foo.pdf"), "some.foo")

}
