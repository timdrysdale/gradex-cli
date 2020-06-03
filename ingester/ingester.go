package ingester

import (
	"fmt"
	"os"
	"path/filepath"
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
	backgroundIsVanilla   bool
	opticalExpand         int
	SkipQuestionFile      bool
}

func New(path string, msgCh chan chmsg.MessageInfo, logger *zerolog.Logger) (*Ingester, error) {

	g := &Ingester{}

	g.msgCh = msgCh

	g.timeout = time.Millisecond //timeout on chmsg sending

	g.root = path

	g.overlayTemplatePath = "layout.svg"
	g.ingestTemplatePath = "layout-flatten-312pt.svg"
	g.backgroundIsVanilla = true
	g.opticalExpand = -10

	err := g.SetupGradexDirs()

	if logger != nil { //for testing
		g.logger = logger
	}
	return g, err
}

func (g *Ingester) SetBackgroundIsVanilla(vanilla bool) {
	g.logger.Info().Bool("vanilla", vanilla).Msg(fmt.Sprintf("Changing backgroundIsVanilla from %v to %v", g.backgroundIsVanilla, vanilla))
	g.backgroundIsVanilla = vanilla
}

func (g *Ingester) SetOpticalShrink(shrink int) {
	g.logger.Info().Int("shrink", shrink).Msg(fmt.Sprintf("Changing optical shrink from %d to %d", -1*g.opticalExpand, shrink))
	g.opticalExpand = -1 * shrink
}

func (g *Ingester) SetSkipQuestionFile(skip bool) {
	g.logger.Info().Bool("skip", skip).Msg(fmt.Sprintf("Changing skipQuestionFile from %v to %v", g.SkipQuestionFile, skip))
	g.SkipQuestionFile = skip
}

func (g *Ingester) SetOverlayTemplatePath(path string) error {

	_, err := os.Stat(filepath.Join(g.OverlayTemplate(), path))
	if err != nil {
		return err
	}

	g.overlayTemplatePath = path
	g.logger.Info().Str("path", path).Msg("Changed overlay template file")
	return nil
}

func (g *Ingester) SetIngestTemplatePath(path string) error {
	g.logger.Info().Str("path", path).Msg("Changed ingest template file")
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
