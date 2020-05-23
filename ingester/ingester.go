package ingester

import (
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
}

func New(path string, msgCh chan chmsg.MessageInfo, logger *zerolog.Logger) (*Ingester, error) {

	g := &Ingester{}

	g.msgCh = msgCh

	g.timeout = time.Millisecond //timeout on chmsg sending

	g.root = path
	err := g.SetupGradexPaths()

	if logger != nil { //for testing
		g.logger = logger
	}
	return g, err
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
