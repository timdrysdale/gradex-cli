package ingester

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/chmsg"
)

func TestIngesterOverlayTemplatePaths(t *testing.T) {

	mch := make(chan chmsg.MessageInfo)
	logger := zerolog.Nop()

	g, err := New("./tmp-delete-me", mch, &logger)
	assert.NoError(t, err)
	assert.Equal(t, "./tmp-delete-me", g.Root())

	os.RemoveAll("./tmp-delete-me")
	g.EnsureDirectoryStructure()

	templateFiles, err := g.GetFileList("./test-fs/etc/overlay/template")
	assert.NoError(t, err)

	for _, file := range templateFiles {
		err := g.CopyToDir(file, g.OverlayTemplate())
		assert.NoError(t, err)
	}

	assert.Equal(t, "tmp-delete-me/etc/overlay/template/layout.svg", g.OverlayLayoutSVG())

	err = g.SetOverlayTemplatePath("layout-q5.svg")
	assert.NoError(t, err)
	assert.Equal(t, "tmp-delete-me/etc/overlay/template/layout-q5.svg", g.OverlayLayoutSVG())

}
