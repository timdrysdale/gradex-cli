package gradexpath

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetupPaths(t *testing.T) {

	gp := New("./tmp-delete-me")

	if gp.Root() != "./tmp-delete-me" {
		t.Errorf("test root set up wrong %s", root)
	}

	// don't use Root() here
	// JUST in case we kill a whole working installation
	os.RemoveAll("./tmp-delete-me")

	err := gp.SetupGradexPaths()
	assert.NoError(t, err)

	err = gp.SetupExamPaths("sample")
	assert.NoError(t, err)

}

//check we can move files without adjusting the modification time
func TestStructFileMod(t *testing.T) {

	gp := New("./tmp-delete-me")

	d1 := []byte("Gradex Testing\n")
	basepath := filepath.Join(gp.Root(), "tmp")
	err := gp.EnsureDir(basepath)
	assert.NoError(t, err)
	testPath := filepath.Join(basepath, "test.txt")
	err = ioutil.WriteFile(testPath, d1, 0755)
	assert.NoError(t, err)
	err = os.Chmod(testPath, 0755)
	assert.NoError(t, err)

	info, err := os.Stat(testPath)
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	assert.NotEqual(t, info.ModTime(), time.Now())

	newPath := filepath.Join(gp.Root(), "tmp", "new.txt")
	err = os.Rename(testPath, newPath)
	infoNew, err := os.Stat(newPath)

	assert.NoError(t, err)
	assert.Equal(t, info.ModTime(), infoNew.ModTime())
}

func TestNewFileMove(t *testing.T) {

	gp := New("./tmp-delete-me")
	d0 := []byte("Gradex Testing\n")
	basepath := filepath.Join(gp.Root(), "tmp")
	err := gp.EnsureDir(basepath)
	assert.NoError(t, err)

	test0 := filepath.Join(basepath, "test0.txt")
	err = ioutil.WriteFile(test0, d0, 0755)
	assert.NoError(t, err)
	err = os.Chmod(test0, 0755)
	assert.NoError(t, err)
	info0, err := os.Stat(test0)
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	test1 := filepath.Join(basepath, "test1.txt")
	d1 := []byte("XXXX\n")
	err = ioutil.WriteFile(test1, d1, 0755)
	assert.NoError(t, err)
	err = os.Chmod(test1, 0755)
	assert.NoError(t, err)
	info1, err := os.Stat(test1)
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	test2 := filepath.Join(basepath, "test2.txt")
	d2 := []byte("YYYY\n")
	err = ioutil.WriteFile(test2, d2, 0755)
	assert.NoError(t, err)
	err = os.Chmod(test2, 0755)
	assert.NoError(t, err)
	info2, err := os.Stat(test2)
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	test3 := filepath.Join(basepath, "test3.txt")
	d3 := []byte("ZZZZ\n")
	err = ioutil.WriteFile(test3, d3, 0755)
	assert.NoError(t, err)
	err = os.Chmod(test3, 0755)
	assert.NoError(t, err)
	info3, err := os.Stat(test3)
	assert.NoError(t, err)

	// check file modtimes

	assert.True(t, info3.ModTime().After(info2.ModTime()))
	assert.True(t, info2.ModTime().After(info1.ModTime()))
	assert.True(t, info1.ModTime().After(info0.ModTime()))

	//should move
	err = gp.MoveIfNewerThanDestination(test1, test0)
	assert.NoError(t, err)

	//should NOT move - but throw no error
	err = gp.MoveIfNewerThanDestination(test2, test3)
	assert.NoError(t, err)

	info0new, err := os.Stat(test0)
	assert.NoError(t, err)
	_, err = os.Stat(test1)
	assert.Error(t, err) // ERROR should have moved!
	_, err = os.Stat(test2)
	assert.NoError(t, err) // no error - should NOT have moved
	info3new, err := os.Stat(test3)
	assert.NoError(t, err)

	if !info0new.ModTime().After(info0.ModTime()) {
		t.Error("first file mod time should have changed")
	}

	if !info3new.ModTime().Equal(info3.ModTime()) {
		t.Error("last file mod time should NOT have changed")
	}

	c0, err := ioutil.ReadFile(test0)
	assert.NoError(t, err)
	c2, err := ioutil.ReadFile(test2)
	assert.NoError(t, err)
	c3, err := ioutil.ReadFile(test3)
	assert.NoError(t, err)

	assert.Equal(t, c0, d1) //content changed
	assert.Equal(t, c2, d2)
	assert.Equal(t, c3, d3) //content not changed

}
