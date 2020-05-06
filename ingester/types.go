package ingester

import (
	"github.com/timdrysdale/chmsg"
	"github.com/timdrysdale/gradex-cli/pagedata"
)

type PDFSummary struct {
	CourseCode  string
	PreparedFor string
	ToDo        string
}

type FlattenTask struct {
	InputPath   string
	PageCount   int
	PageDataMap map[int]pagedata.PageData
	OutputPath  string
	PreparedFor string
	ToDo        string
}

type OverlayTask struct {
	InputPath string
	PageCount int
	//PreparedFor   string
	//ToDo          string
	//NewProcessing pdfpagedata.ProcessingDetails
	//NewQuestion   pdfpagedata.QuestionDetails
	ProcessDetail  pagedata.ProcessDetail
	OldPageDataMap map[int]pagedata.PageData //this has the individual bits filled in?
	OutputPath     string
	SpreadName     string
	Template       string
	Msg            *chmsg.Messager
	Who            string
}

// Overlay command struct - for backwards compatability
// ExamName: This is our internal system representation of the exam,
//and MAY NOT equal the value in PageData for cosmetic reasons - hence
//it is also not equivalent for functional reasons as this string MUST
//exactly match our internal representation

type OverlayCommand struct {
	FromPath      string
	ToPath        string
	ExamName      string
	TemplatePath  string
	SpreadName    string
	ProcessDetail pagedata.ProcessDetail
	//PreparedFor       string
	//ToDo              string
	//ProcessingDetails pdfpagedata.ProcessingDetails
	//QuestionDetails   pdfpagedata.QuestionDetails
	Msg            *chmsg.Messager
	PathDecoration string //this is the "-ma1" for marker1, "mo2" for moderator 2, "d" for done etc
}

var (
	isTesting bool
	testroot  = "./tmp-delete-me"
	ExamStage = []string{
		config,
		acceptedPapers,
		acceptedReceipts,
		tempImages,
		tempPages,
		anonPapers,
		qImages,
		qPages,
		qReady,
		qBack,
		markerReady,
		markerSent,
		markerBack,
		markedCombined,
		markedMerged,
		markedPruned,
		markedReady,
		moderateActive,
		moderatorReady,
		moderatorSent,
		moderatorBack,
		moderatedCombined,
		moderatedMerged,
		moderatedPruned,
		moderatedReady,
		moderateInActive,
		moderateInActiveBack,
		checkerReady,
		checkerSent,
		checkerBack,
		checkedCombined,
		checkedMerged,
		checkedPruned,
		checkedReady,

		remarkerReady,
		remarkerSent,
		remarkerBack,
		remarkedCombined,
		remarkedMerged,
		remarkedPruned,
		remarkedReady,
		recheckerReady,
		recheckerSent,
		recheckerBack,
		recheckedCombined,
		recheckedMerged,
		recheckedPruned,
		recheckedReady,
		reports,
	}
)

const (
	config = "00-config"

	tempImages = "03-temporary-images"
	tempPages  = "04-temporary-pages"

	acceptedReceipts = "02-accepted-receipts"
	acceptedPapers   = "03-accepted-papers"
	anonPapers       = "05-anonymous-papers"

	qImages = "06-question-images"
	qPages  = "07-question-pages"
	qReady  = "08-question-ready"
	qBack   = "09-question-back"

	markerReady          = "20-marker-ready"
	markerSent           = "21-marker-sent"
	markerBack           = "22-marker-back"
	markedCombined       = "23-marked-combined"
	markedMerged         = "24-marked-merged"
	markedPruned         = "25-marked-pruned" //whatever gets trimmed goes here for potential audit
	markedReady          = "26-marked-ready"
	moderateActive       = "30-moderate-active"
	moderateInActive     = "31-moderate-inactive"
	moderatorReady       = "32-moderator-ready"
	moderatorSent        = "33-moderator-sent"
	moderatorBack        = "34-moderator-back"
	moderateInActiveBack = "35-moderate-inactive-back"
	moderatedCombined    = "36-moderated-combined"
	moderatedMerged      = "37-moderated-merged"
	moderatedPruned      = "38-moderated-pruned"
	moderatedReady       = "39-moderated-ready"

	checkerReady    = "40-checker-ready"
	checkerSent     = "41-checker-sent"
	checkerBack     = "42-checker-back"
	checkedCombined = "43-checked-combined"
	checkedMerged   = "44-checked-merged"
	checkedPruned   = "45-checked-pruned"
	checkedReady    = "46-checked-ready"
	reports         = "99-reports"

	remarkerReady    = "50-remarker-ready"
	remarkerSent     = "51-remarker-sent"
	remarkerBack     = "52-remarker-back"
	remarkedCombined = "53-marked-combined"
	remarkedMerged   = "54-marked-merged"
	remarkedPruned   = "55-marked-pruned" //whatever gets trimmed goes here for potential audit
	remarkedReady    = "56-marked-ready"

	recheckerReady    = "60-rechecker-ready"
	recheckerSent     = "61-rechecker-sent"
	recheckerBack     = "62-rechecker-back"
	recheckedCombined = "63-rechecked-combined"
	recheckedMerged   = "64-rechecked-merged"
	recheckedPruned   = "65-rechecked-pruned"
	recheckedReady    = "66-rechecked-ready"

	N = 3
)
