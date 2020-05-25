package ingester

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/timdrysdale/chmsg"
)

type Ingester struct {
	root                  string
	msgCh                 chan chmsg.MessageInfo
	timeout               time.Duration
	logger                *zerolog.Logger
	Redo                  bool
	UseFullAssignmentName bool
	overlayTemplatePath   string
	ingestTemplatePath    string
}

func New(path string, msgCh chan chmsg.MessageInfo, logger *zerolog.Logger) (*Ingester, error) {

	g := &Ingester{}

	g.msgCh = msgCh

	g.timeout = time.Millisecond //timeout on chmsg sending

	g.root = path

	g.overlayTemplatePath = "layout.svg"
	g.ingestTemplatePath = "layout-flatten-312pt.svg"

	err := g.SetupGradexPaths()

	if logger != nil { //for testing
		g.logger = logger
	}
	return g, err
}

func (g *Ingester) SetOverlayTemplatePath(path string) error {

	_, err := os.Stat(filepath.Join(g.OverlayTemplate(), path))
	if err != nil {
		return err
	}

	g.overlayTemplatePath = path

	return nil
}

func (g *Ingester) SetIngestTemplatePath(path string) error {

	_, err := os.Stat(filepath.Join(g.IngestTemplate(), path))

	if err != nil {
		return err
	}

	g.ingestTemplatePath = path

	return nil
}

func (g *Ingester) SetUseFullAssignmentName() {

	g.UseFullAssignmentName = true
}

/*func NewIngester(path string, msgCh chan chmsg.MessageInfo) (*Ingester, error) {

	g := &Ingester{}

	g.msgCh = msgCh

	g.timeout = time.Millisecond //timeout on chmsg sending

	g.root = path
	err := g.SetupGradexPaths()

	return g, err
}*/

func limit(initials string, N int) string {
	if len(initials) < 3 {
		N = len(initials)
	}
	return strings.ToUpper(initials[0:N])
}
