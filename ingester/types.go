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
		pageBad,
		acceptedPapers,
		acceptedReceipts,
		tempImages,
		tempPages,
		anonPapers,
		questionImages,
		questionPages,
		questionReady,
		questionSent,
		questionBack,
		questionSplit,
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

		reMarkerReady,
		reMarkerSent,
		reMarkerBack,
		reMarkedCombined,
		reMarkedMerged,
		reMarkedPruned,
		reMarkedReady,
		reCheckerReady,
		reCheckerSent,
		reCheckerBack,
		reCheckedCombined,
		reCheckedMerged,
		reCheckedPruned,
		reCheckedReady,
		reports,
	}
)

const (

	// EXTERNAL, e.g. command args, so no numbers

	QuestionReady  = "questionReady"
	QuestionSent   = "questionSent"
	MarkerReady    = "markerReady"
	MarkerSent     = "markerSent"
	ModeratorReady = "moderatorReady"
	ModeratorSent  = "moderatorSent"
	CheckerReady   = "checkerReady"
	CheckerSent    = "checkerSent"
	ReMarkerReady  = "remarkerReady"
	ReMarkerSent   = "remarkerSent"
	ReCheckerReady = "recheckerReady"
	ReCheckerSent  = "recheckerSent"

	//>>>>>>>>>>>> INTERNAL >>>>>>>>>>>>>>>>>>
	config = "00-config"

	pageBad = "01-page-bad"

	tempImages = "03-temporary-images"
	tempPages  = "04-temporary-pages"

	acceptedReceipts = "02-accepted-receipts"
	acceptedPapers   = "03-accepted-papers"
	anonPapers       = "05-anonymous-papers"

	questionImages = "06-question-images"
	questionPages  = "07-question-pages"
	questionReady  = "08-question-ready"
	questionSent   = "09-question-sent"
	questionBack   = "10-question-back"
	questionSplit  = "11-question-split"

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

	reMarkerReady    = "50-remarker-ready"
	reMarkerSent     = "51-remarker-sent"
	reMarkerBack     = "52-remarker-back"
	reMarkedCombined = "53-marked-combined"
	reMarkedMerged   = "54-marked-merged"
	reMarkedPruned   = "55-marked-pruned" //whatever gets trimmed goes here for potential audit
	reMarkedReady    = "56-marked-ready"

	reCheckerReady    = "60-rechecker-ready"
	reCheckerSent     = "61-rechecker-sent"
	reCheckerBack     = "62-rechecker-back"
	reCheckedCombined = "63-rechecked-combined"
	reCheckedMerged   = "64-rechecked-merged"
	reCheckedPruned   = "65-rechecked-pruned"
	reCheckedReady    = "66-rechecked-ready"

	N = 3
)
