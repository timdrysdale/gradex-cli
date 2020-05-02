package ingester

import (
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/timdrysdale/chmsg"
)

type Ingester struct {
	root    string
	msgCh   chan chmsg.MessageInfo
	timeout time.Duration
	logger  *zerolog.Logger
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
