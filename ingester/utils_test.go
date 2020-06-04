package ingester

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAnonymousFromPath(t *testing.T) {
	assert.Equal(t, "B999999", GetAnonymousFromPath("Practice-B999999-maTDD-marked-comments.pdf"))
}

func TestShortenBaseFileName(t *testing.T) {

	assert.Equal(t, "s0000000_attempt_2020-05-01-02-00-00", shortenBaseFileName("PGEEnnnn A Super Long Exam Name - Exam Dropbox_s0000000_attempt_2020-05-01-02-00-00_PGEEnnnn-B000000"))
	assert.Equal(t, "s0000000_attempt_2020-05-01-02-00-00", shortenBaseFileName("PGEEnnnn A Super Long Exam Name - Exam Dropbox_s0000000_attempt_2020-05-01-02-00-00_EXTRA_STUFF_WITH_UNDERSCORES"))
	assert.Equal(t, "AB", shortenBaseFileName("AB"))
	assert.Equal(t, "ABkjahsdfjkhasdjkhfkjashdjfkhaskjdhfkjasdhfkjahs", shortenBaseFileName("ABkjahsdfjkhasdjkhfkjashdjfkhaskjdhfkjasdhfkjahs"))
}

func TestShortenAssignment(t *testing.T) {

	assert.Equal(t, "PGEE11120", shortenAssignment("PGEE11120 Advanced Wireless Communications (MSc) - Exam Dropbox"))
	assert.Equal(t, "FOO", shortenAssignment("FOO"))
	assert.Equal(t, "0123456789AB", shortenAssignment("0123456789ABCDEFGHIJKLMNOPQRSTUVW"))
	assert.Equal(t, "0123456789AB", shortenAssignment("0123456789ABCDEFGHIJKLMNOPQRSTUVW asdfa  asdf asdf as asdf asdf asdf asdf asdf "))
}

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
