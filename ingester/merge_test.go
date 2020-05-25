package ingester

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/gradex-cli/pagedata"
)

// check state machine won't go back to a lower priority state
// e.g. once we're marked, we shouldn't go back to bad
// even though that seems odd, if merely "bad" we might
// drop the page from the merge, whereas if marked, we need
// to keep it so we can see what was marked, despite it being
// bad (something partly marked, AND bad, is not a problem
// for this system to handle - that's for a human, we just
// need to make sure human work is not skipped, hence
// showing marked work even if bad is marked

func TestPageFSM(t *testing.T) {

	pageFSM := newPageFSM()

	assert.Equal(t, statusSkipped, pageFSM.Current())

	pageFSM.Event(statusSeen)

	assert.Equal(t, statusSeen, pageFSM.Current())

	pageFSM.Event(statusBad)

	assert.Equal(t, statusBad, pageFSM.Current())

	pageFSM.Event(statusSeen)

	assert.Equal(t, statusBad, pageFSM.Current())

	pageFSM.Event(statusMarked)

	assert.Equal(t, statusMarked, pageFSM.Current())

	pageFSM.Event(statusBad)

	assert.Equal(t, statusMarked, pageFSM.Current())

	pageFSM.Event(statusSeen)

	assert.Equal(t, statusMarked, pageFSM.Current())

}

func TestSummarisePageSkipped(t *testing.T) {

	originalPath := "some/original/path with spaces/file.pdf"
	ownPath := "a/b/c.pdf"
	pageNumber := 3
	wasFor := "DEF"

	pageData := pagedata.PageData{
		Current: pagedata.PageDetail{
			Own: pagedata.FileDetail{
				Path: ownPath,
			},
			Original: pagedata.FileDetail{
				Path:   originalPath,
				Number: pageNumber,
			},
			Data: []pagedata.Field{
				pagedata.Field{ //include to make sure we are checking for actual text-fields
					Key:   "not a textfield",
					Value: "happy days",
				},
				pagedata.Field{
					Key:   "tf-page-bad",
					Value: "",
				},
				pagedata.Field{
					Key:   "tf-page-ok",
					Value: "",
				},
				pagedata.Field{
					Key:   "tf-question-01-section",
					Value: "",
				},
				pagedata.Field{
					Key:   "tf-page-bad-optical",
					Value: "",
				},
				pagedata.Field{
					Key:   "tf-page-ok-optical",
					Value: "",
				},
				pagedata.Field{
					Key:   "tf-question-01-section-optical",
					Value: "",
				},
			},
		},
		Previous: []pagedata.PageDetail{
			pagedata.PageDetail{
				Process: pagedata.ProcessDetail{
					ToDo: "FirstProcess",
					For:  "ABC",
				},
			},
			pagedata.PageDetail{
				Process: pagedata.ProcessDetail{
					For:  wasFor,
					ToDo: "SecondProcess",
				},
			},
		},
	}

	summary := summarisePage(pageData)

	assert.Equal(t, wasFor, summary.WasFor)
	assert.Equal(t, statusSkipped, summary.Status)
	assert.Equal(t, originalPath, summary.Original)
	assert.Equal(t, ownPath, summary.OwnPath)
	assert.Equal(t, pageNumber, summary.PageNumber)

}

func TestSummarisePageSeen(t *testing.T) {

	originalPath := "some/original/path with spaces/file.pdf"
	ownPath := "a/b/c.pdf"
	pageNumber := 3
	wasFor := "DEF"

	pageData := pagedata.PageData{
		Current: pagedata.PageDetail{
			Own: pagedata.FileDetail{
				Path: ownPath,
			},
			Original: pagedata.FileDetail{
				Path:   originalPath,
				Number: pageNumber,
			},
			Data: []pagedata.Field{
				pagedata.Field{
					Key:   "not a textfield",
					Value: "happy days",
				},
				pagedata.Field{
					Key:   "tf-page-bad",
					Value: "",
				},
				pagedata.Field{
					Key:   "tf-page-ok",
					Value: "",
				},
				pagedata.Field{
					Key:   "tf-question-01-section",
					Value: "",
				},
				pagedata.Field{
					Key:   "tf-page-bad-optical",
					Value: "",
				},
				pagedata.Field{
					Key:   "tf-page-ok-optical",
					Value: markDetected,
				},
				pagedata.Field{
					Key:   "tf-question-01-section-optical",
					Value: "",
				},
			},
		},
		Previous: []pagedata.PageDetail{
			pagedata.PageDetail{
				Process: pagedata.ProcessDetail{
					ToDo: "FirstProcess",
					For:  "ABC",
				},
			},
			pagedata.PageDetail{
				Process: pagedata.ProcessDetail{
					For:  wasFor,
					ToDo: "SecondProcess",
				},
			},
		},
	}

	summary := summarisePage(pageData)

	assert.Equal(t, wasFor, summary.WasFor)
	assert.Equal(t, statusSeen, summary.Status)
	assert.Equal(t, originalPath, summary.Original)
	assert.Equal(t, ownPath, summary.OwnPath)
	assert.Equal(t, pageNumber, summary.PageNumber)

}

func TestSummarisePageBad(t *testing.T) {

	originalPath := "some/original/path with spaces/file.pdf"
	ownPath := "a/b/c.pdf"
	pageNumber := 3
	wasFor := "DEF"

	pageData := pagedata.PageData{
		Current: pagedata.PageDetail{
			Own: pagedata.FileDetail{
				Path: ownPath,
			},
			Original: pagedata.FileDetail{
				Path:   originalPath,
				Number: pageNumber,
			},
			Data: []pagedata.Field{
				pagedata.Field{
					Key:   "not a textfield",
					Value: "happy days",
				},
				pagedata.Field{
					Key:   "tf-page-bad",
					Value: "",
				},
				pagedata.Field{
					Key:   "tf-page-ok",
					Value: "",
				},
				pagedata.Field{
					Key:   "tf-question-01-section",
					Value: "",
				},
				pagedata.Field{
					Key:   "tf-page-bad-optical",
					Value: markDetected,
				},
				pagedata.Field{
					Key:   "tf-page-ok-optical",
					Value: "",
				},
				pagedata.Field{
					Key:   "tf-question-01-section-optical",
					Value: "",
				},
			},
		},
		Previous: []pagedata.PageDetail{
			pagedata.PageDetail{
				Process: pagedata.ProcessDetail{
					ToDo: "FirstProcess",
					For:  "ABC",
				},
			},
			pagedata.PageDetail{
				Process: pagedata.ProcessDetail{
					For:  wasFor,
					ToDo: "SecondProcess",
				},
			},
		},
	}

	summary := summarisePage(pageData)

	assert.Equal(t, wasFor, summary.WasFor)
	assert.Equal(t, statusBad, summary.Status)
	assert.Equal(t, originalPath, summary.Original)
	assert.Equal(t, ownPath, summary.OwnPath)
	assert.Equal(t, pageNumber, summary.PageNumber)

}

func TestSummarisePageMarked(t *testing.T) {

	originalPath := "some/original/path with spaces/file.pdf"
	ownPath := "a/b/c.pdf"
	pageNumber := 3
	wasFor := "DEF"

	pageData := pagedata.PageData{
		Current: pagedata.PageDetail{
			Own: pagedata.FileDetail{
				Path: ownPath,
			},
			Original: pagedata.FileDetail{
				Path:   originalPath,
				Number: pageNumber,
			},
			Data: []pagedata.Field{
				pagedata.Field{
					Key:   "not a textfield",
					Value: "happy days",
				},
				pagedata.Field{
					Key:   "tf-page-bad",
					Value: "",
				},
				pagedata.Field{
					Key:   "tf-page-ok",
					Value: "",
				},
				pagedata.Field{
					Key:   "tf-question-01-section",
					Value: markDetected,
				},
				pagedata.Field{
					Key:   "tf-page-bad-optical",
					Value: "",
				},
				pagedata.Field{
					Key:   "tf-page-ok-optical",
					Value: markDetected,
				},
				pagedata.Field{
					Key:   "tf-question-01-section-optical",
					Value: "",
				},
			},
		},
		Previous: []pagedata.PageDetail{
			pagedata.PageDetail{
				Process: pagedata.ProcessDetail{
					ToDo: "FirstProcess",
					For:  "ABC",
				},
			},
			pagedata.PageDetail{
				Process: pagedata.ProcessDetail{
					For:  wasFor,
					ToDo: "SecondProcess",
				},
			},
		},
	}

	summary := summarisePage(pageData)

	assert.Equal(t, wasFor, summary.WasFor)
	assert.Equal(t, statusMarked, summary.Status)
	assert.Equal(t, originalPath, summary.Original)
	assert.Equal(t, ownPath, summary.OwnPath)
	assert.Equal(t, pageNumber, summary.PageNumber)

}

func makePaperMap1() map[string]map[int]PageCollection {

	summaries := []PageSummary{
		PageSummary{
			Original:   "A",
			PageNumber: 1,
			OwnPath:    "A1-ABC.pdf",
			Status:     statusMarked,
			WasFor:     "ABC",
		},
		PageSummary{
			Original:   "A",
			PageNumber: 1,
			OwnPath:    "A1-DEF.pdf",
			Status:     statusSeen,
			WasFor:     "DEF",
		},
		PageSummary{
			Original:   "B",
			PageNumber: 1,
			OwnPath:    "B1-ABC.pdf",
			Status:     statusMarked,
			WasFor:     "ABC",
		},
		PageSummary{
			Original:   "B",
			PageNumber: 1,
			OwnPath:    "B1-DEF.pdf",
			Status:     statusSeen,
			WasFor:     "DEF",
		},
		PageSummary{
			Original:   "B",
			PageNumber: 2,
			OwnPath:    "B2-ABC.pdf",
			Status:     statusMarked,
			WasFor:     "ABC",
		},
		PageSummary{
			Original:   "B",
			PageNumber: 2,
			OwnPath:    "B2-DEF.pdf",
			Status:     statusMarked,
			WasFor:     "DEF",
		},
	}

	return createPaperMap(summaries)

}

func TestCreatePaperMap(t *testing.T) {

	paperMap := makePaperMap1()

	assert.Equal(t, 1, len(paperMap["A"][1].Seen))
	assert.Equal(t, 1, len(paperMap["A"][1].Marked))
	assert.Equal(t, 1, len(paperMap["B"][1].Seen))
	assert.Equal(t, 1, len(paperMap["B"][1].Marked))
	assert.Equal(t, 2, len(paperMap["B"][2].Marked))
	assert.Equal(t, 0, len(paperMap["B"][2].Seen))

}

func TestCreatePageItem(t *testing.T) {

	paperMap := makePaperMap1()

	page := createPageItem(paperMap["B"][1], paperMap["B"][1].Seen[0])

	message := `This page seen by DEF
Marked: ABC
Bad:
Seen: DEF
Skipped:`

	assert.Equal(t, "B1-DEF.pdf", page.Path)
	assert.Equal(t, message, page.Message)
}

func TestCreatePageList(t *testing.T) {

	paperMap := makePaperMap1()

	pageList1 := createPageList(paperMap["B"][1])
	assert.Equal(t, 1, len(pageList1))
	assert.Equal(t, "B1-ABC.pdf", pageList1[0].Path)

	pageList2 := createPageList(paperMap["B"][2])
	assert.Equal(t, 2, len(pageList2))
	assert.Equal(t, "B2-ABC.pdf", pageList2[0].Path)
	assert.Equal(t, "B2-DEF.pdf", pageList2[1].Path)

}
