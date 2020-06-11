package ingester

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/timdrysdale/gradex-cli/comment"
	"github.com/timdrysdale/gradex-cli/pagedata"
)

type PageReport struct {
	Error        string `csv:"error"`
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

func (g *Ingester) ReportOnProcessedDir(exam, dir string, showOK bool, showMark bool, reconcile bool) ([]string, error) {

	tokens := []string{}

	files, err := g.GetFileList(dir)

	if err != nil {
		return []string{}, err
	}

	noLinkError := true

	errorPageReports := []PageReport{}

	destMap := make(map[string]map[int]PageReport)

	for _, file := range files {

		if !IsPDF(file) {
			continue
		}

		prMap, err := GetPageSummaryMapFromFile(file)

		destMap[file] = prMap

		if err != nil { // linkError
			noLinkError = false

			for _, pr := range prMap {
				if !pr.IsLinked {
					tokens = append(tokens, "BROKEN-LINK: "+pr.String())
					pr.Error = "BROKEN-LINK"
					errorPageReports = append(errorPageReports, pr)
				}
			}

		}

		for _, pr := range prMap {
			if pr.Status == statusBad {
				tokens = append(tokens, "BAD : "+pr.String())
				pr.Error = "BAD-PAGE"
				errorPageReports = append(errorPageReports, pr)

			} else if pr.Status == statusSkipped {
				tokens = append(tokens, "SKIP: "+pr.String())
				pr.Error = "SKIPPED"
				errorPageReports = append(errorPageReports, pr)

			} else if showMark && pr.Status == statusMarked {

				tokens = append(tokens, "QBOX: "+pr.String())
				pr.Error = "QBOX"
				errorPageReports = append(errorPageReports, pr)

			} else if showOK {
				tokens = append(tokens, "OK  : "+pr.String())
				pr.Error = "OK"
				errorPageReports = append(errorPageReports, pr)
			}

		}

	}

	if !reconcile {
		return tokens, nil
	}

	destPages := make(map[string]int)

	// populate the source map
	for _, reportMap := range destMap {
		for _, report := range reportMap {
			destPages[report.FirstLink] = 0
		}
	}

	fmt.Printf("Destination: Found %d files and %d unique pages in %s\n", len(destMap), len(destPages), dir)

	srcDir := g.GetExamDir(exam, anonPapers)
	sourceFiles, err := g.GetFileList(srcDir)

	if err != nil {
		return []string{}, err
	}

	srcMap := make(map[string]map[int]PageReport)

	for _, file := range sourceFiles {

		if !IsPDF(file) {
			continue
		}
		prMap, _ := GetPageSummaryMapFromFile(file) //ignore link errors
		srcMap[file] = prMap

	}

	srcPages := make(map[string]int)

	uuidMap := make(map[string]PageReport)

	// populate the source map
	for _, reportMap := range srcMap {
		for _, report := range reportMap {
			srcPages[report.FirstLink] = 0
			uuidMap[report.FirstLink] = report
		}
	}

	fmt.Printf("Source     : Found %d files and %d unique pages in %s\n", len(srcMap), len(srcPages), srcDir)

	if len(destMap) != len(srcMap) {
		fmt.Printf("WARNING: number of files differs in each location (%d != %d) ", len(destMap), len(srcMap))
		fmt.Printf("(This may not be an issue, because files may have been merged into batches or split)\n")
	}
	if len(destPages) != len(srcPages) {
		fmt.Printf("ERROR: number of unique pages differs in each location (%d != %d) ", len(destPages), len(srcPages))
		fmt.Printf("(This is a problem because it means actual pages have been lost or found!)\n")
	}

	for _, reportMap := range destMap {
		for _, report := range reportMap {
			if _, ok := srcPages[report.FirstLink]; ok {
				srcPages[report.FirstLink] = srcPages[report.FirstLink] + 1
			} else {
				noLinkError = false
				report.Error = "NO-ORIGINAL-PAGE-FOR"
				errorPageReports = append(errorPageReports, report)
				fmt.Printf("%s: %s-%s %s PAGE %d\n", report.Error, report.What, report.Who, report.When, report.PageNumber)
			}
		}
	}

	for key, count := range srcPages {
		if count < 1 {
			noLinkError = false
			pr := uuidMap[key]
			pr.Error = "MISSING-PROCESSED-PAGE-FOR"
			fmt.Printf("%s: %s-%s %s PAGE %d\n", pr.Error, pr.What, pr.Who, pr.When, pr.PageNumber)
			errorPageReports = append(errorPageReports, pr)
		}
	}

	if noLinkError {
		errorPageReports = append(errorPageReports, PageReport{
			Error: fmt.Sprintf("SUCCESS: %d unique original and processed pages are all correctly linked between %s and %s", len(srcPages), srcDir, dir),
		})
		fmt.Printf("SUCCESS: %d unique original and processed pages are all correctly linked between %s and %s\n", len(srcPages), srcDir, dir)
	}

	//write to reports folder
	reportPath := filepath.Join(g.GetExamDir(exam, reports),
		fmt.Sprintf("%s-%s-%d.csv",
			shortenAssignment(exam),
			filepath.Base(dir),
			time.Now().Unix()))

	err = WritePageReportsToCSV(errorPageReports, reportPath)
	if err != nil {
		fmt.Printf("Error writing report to %s\n", reportPath)
	}
	return tokens, nil

}

func WritePageReportsToCSV(reports []PageReport, outputPath string) error {
	// wrap the marshalling library in case we need converters etc later
	file, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	return gocsv.MarshalFile(&reports, file)
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
		cmts = append(cmts, "["+cmt.Label+"]: "+cmt.Text)
	}

	return strings.Join(cmts, "; ")
}
