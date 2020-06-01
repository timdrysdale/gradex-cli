package ingester

import (
	"fmt"
	"strings"

	"github.com/timdrysdale/gradex-cli/comment"
	"github.com/timdrysdale/gradex-cli/pagedata"
)

type PageReport struct {
	What         string `csv:"what"`
	Who          string `csv:"who"`
	When         string `csv:"when"`
	Original     string `csv:"original"`
	OwnPath      string `csv:"path"`
	WasFor       string `csv:"wasfor"`
	PageNumber   int    `csv:"page"`
	Status       string `csv:"status"`
	MergeMessage string `csv:"message"`
	Comments     string `csv:"comments"`
	IsLinked     bool   `csv:"linked"`
	FirstLink    string `csv:"firstlink"`
	LastLink     string `csv:"lastlink"`
}

func (p *PageReport) String() string {
	linkStatus := "LINK-ERROR"
	if p.IsLinked {
		linkStatus = "link-ok"
	}
	return p.What + "-" + p.Who + "(" + p.When + ") p" +
		fmt.Sprintf("%02d", p.PageNumber) + " is " + p.Status + "/" + linkStatus +
		":" + p.MergeMessage +
		" Comments: " + p.Comments

}

func (g *Ingester) ReportOnProcessedDir(dir string, showOK bool) ([]string, error) {

	tokens := []string{}

	files, err := g.GetFileList(dir)

	if err != nil {
		return []string{}, err
	}

	for _, file := range files {

		prMap, err := GetPageSummaryMapFromFile(file)

		if err != nil { // linkError

			for _, pr := range prMap {

				tokens = append(tokens, "LINK: "+pr.String())
			}
		} else {

			for _, pr := range prMap {
				if pr.Status == statusBad {
					tokens = append(tokens, "BAD : "+pr.String())
				} else if pr.Status == statusSkipped {
					tokens = append(tokens, "SKIP: "+pr.String())
				} else if showOK {
					tokens = append(tokens, "OK  : "+pr.String())
				}
			}

		}

	}

	return tokens, nil
}

func GetPageSummaryMapFromFile(path string) (map[int]PageReport, error) {

	pdMap, err := pagedata.UnMarshalAllFromFile(path)

	if err != nil {
		return map[int]PageReport{}, err
	}

	return GetPageSummaryMap(pdMap)

}

func GetPageSummaryMap(pdMap map[int]pagedata.PageData) (map[int]PageReport, error) {

	linkMap, linkErr := pagedata.GetLinkMap(pdMap)

	reportMap := make(map[int]PageReport)

	for page, pd := range pdMap {

		summary := summarisePage(pd)

		reportMap[page] = PageReport{
			What:         pd.Current.Item.What,
			Who:          pd.Current.Item.Who,
			When:         pd.Current.Item.When,
			Original:     summary.Original,
			OwnPath:      summary.OwnPath,
			WasFor:       summary.WasFor,
			PageNumber:   summary.PageNumber,
			Status:       summary.Status,
			MergeMessage: GetField(pd.Current.Data, "merge-message"),
			Comments:     CommentsToString(pd.Current.Comments),
			IsLinked:     linkMap[page].IsLinked,
			FirstLink:    linkMap[page].First,
			LastLink:     linkMap[page].Last,
		}

	}

	return reportMap, linkErr

}

func GetField(fields []pagedata.Field, key string) string {

	for _, item := range fields {
		if item.Key == key {
			return item.Value
		}
	}
	return ""
}

func CommentsToString(comments []comment.Comment) string {

	cmts := []string{}

	for _, cmt := range comments {
		cmts = append(cmts, cmt.Label+cmt.Text)
	}

	return strings.Join(cmts, "; ")
}
