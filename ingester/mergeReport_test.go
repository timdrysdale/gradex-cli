package ingester

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/gradex-cli/comment"
	"github.com/timdrysdale/gradex-cli/pagedata"
)

func TestPageReport(t *testing.T) {

	originalPath := "a/file/some/where.pdf"
	ownPath := "a/b/c.pdf"
	cmt0 := "So this is a pretty long commment but we will probably see longer"
	cmt1 := "And here is another commment"
	cmtCombined := strings.Join([]string{"[0]: " + cmt0, "[1]: " + cmt1}, "; ")
	pdMap := map[int]pagedata.PageData{
		1: pagedata.PageData{
			Current: pagedata.PageDetail{
				Item: pagedata.ItemDetail{
					What: "EL03",
					Who:  "B00",
				},
				Own: pagedata.FileDetail{
					Path: ownPath,
				},
				Original: pagedata.FileDetail{
					Path:   originalPath,
					Number: 1,
				},
				Comments: []comment.Comment{
					comment.Comment{
						Label: "0",
						Text:  cmt0,
					},
					comment.Comment{
						Label: "1",
						Text:  cmt1,
					},
				},
				Data: []pagedata.Field{
					pagedata.Field{
						Key:   "not a textfield",
						Value: "happy days",
					},
					pagedata.Field{
						Key:   "tf-page-ok-optical",
						Value: "",
					},
					pagedata.Field{
						Key:   "tf-question-01-section-optical",
						Value: markDetected,
					},
					pagedata.Field{
						Key:   "tf-page-bad-optical",
						Value: markDetected,
					},
				},
			},
		},
		2: pagedata.PageData{
			Current: pagedata.PageDetail{
				Item: pagedata.ItemDetail{
					What: "EL03",
					Who:  "B00",
				},
				Own: pagedata.FileDetail{
					Path: ownPath,
				},
				Original: pagedata.FileDetail{
					Path:   originalPath,
					Number: 2,
				},
				Data: []pagedata.Field{
					pagedata.Field{
						Key:   "not a textfield",
						Value: "happy days",
					},
					pagedata.Field{
						Key:   "tf-question-01-section-optical",
						Value: "",
					},
					pagedata.Field{
						Key:   "tf-page-bad-optical",
						Value: "",
					},
				},
			},
		},

		3: pagedata.PageData{
			Current: pagedata.PageDetail{
				Item: pagedata.ItemDetail{
					What: "EL03",
					Who:  "B00",
				},
				Own: pagedata.FileDetail{
					Path: ownPath,
				},
				Original: pagedata.FileDetail{
					Path:   originalPath,
					Number: 3,
				},
				Data: []pagedata.Field{
					pagedata.Field{
						Key:   "not a textfield",
						Value: "happy days",
					},
					pagedata.Field{
						Key:   "tf-question-01-section-optical",
						Value: markDetected,
					},
					pagedata.Field{
						Key:   "tf-page-ok-optical",
						Value: markDetected,
					},
				},
			},
		},
	}

	prMap, err := GetPageSummaryMap(pdMap)

	assert.NoError(t, err) //no linking in the test map

	assert.Equal(t, 3, len(prMap))

	assert.Equal(t, statusBad, prMap[1].Status)

	assert.Equal(t, cmtCombined, prMap[1].Comments)

	assert.Equal(t, statusSkipped, prMap[2].Status)

	assert.Equal(t, statusMarked, prMap[3].Status)

	assert.Equal(t, markDetected, GetField(pdMap[3].Current.Data, "tf-page-ok-optical"))

}
