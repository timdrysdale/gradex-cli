package ingester

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/chmsg"
	"github.com/timdrysdale/gradex-cli/pagedata"
)

func TestQNumber(t *testing.T) {

	n, what, is := isQ("tf-q1-mark")
	assert.True(t, is)
	assert.Equal(t, "1", n)
	assert.Equal(t, "mark", what)

	n, what, is = isQ("tf-q02-section")
	assert.True(t, is)
	assert.Equal(t, "02", n)
	assert.Equal(t, "section", what)

	n, what, is = isQ("q1-asdfasd")
	assert.False(t, is)
	assert.Equal(t, "", n)
	assert.Equal(t, "", what)

}

func TestGetNum(t *testing.T) {

	v, err := getNum("0.5/23")
	assert.NoError(t, err)
	assert.Equal(t, 0.5, v)

	v, err = getNum(".5/23")
	assert.NoError(t, err)
	assert.Equal(t, 0.5, v)

	v, err = getNum(" .5/23")
	assert.NoError(t, err)
	assert.Equal(t, 0.5, v)

	v, err = getNum("8 ")
	assert.NoError(t, err)
	assert.Equal(t, 8.0, v) //test needs a float to match return type

	v, err = getNum(" 3.5\\12")
	assert.NoError(t, err)
	assert.Equal(t, 3.5, v)

	v, err = getNum(" 1.1-20/3") //not sure what this would really mean,
	assert.NoError(t, err)
	assert.Equal(t, 1.1, v)

}

func TestSelectPageDetailsWithMarks(t *testing.T) {

	winnerA := pagedata.PageDetail{
		Process: pagedata.ProcessDetail{
			Name: "enter-active-bar",
			For:  "A",
		},
	}
	winnerB := pagedata.PageDetail{
		Process: pagedata.ProcessDetail{
			Name: "merge-marked",
			For:  "E",
		},
	}
	pageData1 := pagedata.PageData{
		Current: winnerA,
		Previous: []pagedata.PageDetail{
			pagedata.PageDetail{
				Process: pagedata.ProcessDetail{
					Name: "merge-marked",
					For:  "B",
				},
			},
			pagedata.PageDetail{
				Process: pagedata.ProcessDetail{
					Name: "enter-inactive-bar",
					For:  "C",
				},
			},
		},
	}

	pageData2 := pagedata.PageData{
		Current: pagedata.PageDetail{
			Process: pagedata.ProcessDetail{
				Name: "enter-inactive-bar",
				For:  "D",
			},
		},
		Previous: []pagedata.PageDetail{
			winnerB,
			pagedata.PageDetail{
				Process: pagedata.ProcessDetail{
					Name: "flattened-marked",
					For:  "F",
				},
			},
		},
	}
	pageData3 := pagedata.PageData{
		Current: pagedata.PageDetail{
			Process: pagedata.ProcessDetail{
				Name: "enter-inactive-bar-foo",
				For:  "D",
			},
		},
		Previous: []pagedata.PageDetail{
			pagedata.PageDetail{
				Process: pagedata.ProcessDetail{
					Name: "merge-marked-foo",
					For:  "E",
				},
			},
			pagedata.PageDetail{
				Process: pagedata.ProcessDetail{
					Name: "flattened-marked",
					For:  "F",
				},
			},
		},
	}

	pdMap := make(map[int]pagedata.PageData)

	pdMap[1] = pageData1
	pdMap[2] = pageData2
	pdMap[3] = pageData3

	pageDetails := selectPageDetailsWithMarks(pdMap)

	assert.Equal(t, winnerA, pageDetails[0])
	assert.Equal(t, winnerB, pageDetails[1])
	assert.Equal(t, 2, len(pageDetails))
}

func TestQMap(t *testing.T) {

	winnerA := pagedata.PageDetail{
		Data: []pagedata.Field{
			pagedata.Field{
				Key:   "not a textfield",
				Value: "happy days",
			},
			pagedata.Field{
				Key:   "tf-q1-mark",
				Value: ".5/23",
			},
			pagedata.Field{
				Key:   "tf-q1-section",
				Value: "A",
			},
			pagedata.Field{
				Key:   "tf-q1-number",
				Value: "1",
			},
			pagedata.Field{
				Key:   "tf-q2-mark",
				Value: "7",
			},
			pagedata.Field{
				Key:   "tf-q2-section",
				Value: "A",
			},
			pagedata.Field{
				Key:   "tf-q2-number",
				Value: "2",
			},
		},
	}

	winnerB := pagedata.PageDetail{
		Data: []pagedata.Field{
			pagedata.Field{
				Key:   "not a textfield",
				Value: "happy days",
			},
			pagedata.Field{
				Key:   "tf-q1-mark",
				Value: "13-20",
			},
			pagedata.Field{
				Key:   "tf-q1-section",
				Value: "", //deliberately blank
			},
			pagedata.Field{
				Key:   "tf-q1-number",
				Value: "B1",
			},
		},
	}
	winnerC := pagedata.PageDetail{
		Data: []pagedata.Field{
			pagedata.Field{
				Key:   "not a textfield",
				Value: "happy days",
			},
			pagedata.Field{
				Key:   "tf-q1-mark",
				Value: "5-20",
			},
			pagedata.Field{
				Key:   "tf-q1-section",
				Value: "B",
			},
			pagedata.Field{
				Key:   "tf-q1-number",
				Value: "1",
			},
		},
	}

	pageDetails := []pagedata.PageDetail{
		winnerA,
		winnerB,
		winnerC,
	}

	expectedQMap := map[string]string{
		"A1": "0.5",
		"A2": "7",
		"B1": "18", //involves adding up two part marks
	}

	QMap := getQMap(pageDetails)

	assert.Equal(t, expectedQMap, QMap)

}

func TestCoverPage(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	verbose := false

	mch := make(chan chmsg.MessageInfo)

	closed := make(chan struct{})
	defer close(closed)
	go func() {
		for {
			select {
			case <-closed:
				break
			case msg := <-mch:
				if verbose {
					fmt.Printf("MC:%s\n", msg.Message)
				}
			}

		}
	}()

	logFile := "./cover-testing.log"

	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	logger := zerolog.New(f).With().Timestamp().Logger()

	//logger := zerolog.Nop()

	g, err := New("./tmp-cover", mch, &logger)
	assert.NoError(t, err)

	assert.Equal(t, "./tmp-cover", g.Root())

	//>>>>>>>>>>>>>>>>>>>>>>>>> SETUP >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	// don't use GetRoot() here
	// JUST in case we kill a whole working installation

	os.RemoveAll("./tmp-cover")

	g.EnsureDirectoryStructure()

	templateFiles, err := g.GetFileList("./test-fs/etc/overlay/template")
	assert.NoError(t, err)

	for _, file := range templateFiles {
		destination := filepath.Join(g.OverlayTemplate(), filepath.Base(file))
		err := Copy(file, destination)
		assert.NoError(t, err)
	}
	procDetail := pagedata.ProcessDetail{
		UUID:     safeUUID(),
		UnixTime: time.Now().UnixNano(),
		Name:     "check-cover",
		By:       "gradex-cli",
		ToDo:     "checking",
		For:      "X",
	}

	questions := []string{"A1", "A2", "a3", "b1", " B3"} //deliberate case and white space

	cp := CoverPageCommand{
		Questions:      questions,
		FromPath:       "test-cover",
		ToPath:         "test-cover-out",
		ExamName:       "Exam",
		TemplatePath:   g.OverlayLayoutSVG(),
		SpreadName:     "addition",
		ProcessDetail:  procDetail,
		PathDecoration: "-cover",
	}

	err = g.CoverPage(cp, &logger)

	assert.NoError(t, err)
}
