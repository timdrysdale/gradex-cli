package ingester

import (
	"github.com/timdrysdale/chmsg"
	"github.com/timdrysdale/gradex-cli/extract"
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
	InputPath                string
	CoverPath                string //no cover page if empty ""
	AncestorPath             string //no redo of ancestry if empty
	Msg                      *chmsg.Messager
	NewFieldMap              map[int][]pagedata.Field
	OldPageDataMap           map[int]pagedata.PageData //this has the individual bits filled in?
	OutputPath               string
	PageCount                int
	ProcessDetail            pagedata.ProcessDetail
	SpreadName               string
	Template                 string
	Who                      string
	OpticalBoxSpread         string
	ReadOpticalBoxes         bool
	TextFields               map[int]map[string]extract.TextField
	OmitPreviousComments     bool //this is for the checked stage, where we don't want earlier comments
	PropagateTextFieldValues bool // this is for enter active - copy textfield values out of pagedata into enter bar
}

// Overlay command struct - for backwards compatability
// ExamName: This is our internal system representation of the exam,
//and MAY NOT equal the value in PageData for cosmetic reasons - hence
//it is also not equivalent for functional reasons as this string MUST
//exactly match our internal representation

type OverlayCommand struct {
	CoverPath                string //get a cover page from here if it exists, and if cover page exists
	FromPath                 string
	ToPath                   string
	ExamName                 string
	TemplatePath             string
	SpreadName               string
	ProcessDetail            pagedata.ProcessDetail
	Msg                      *chmsg.Messager
	PathDecoration           string //this is the "-ma1" for marker1, "mo2" for moderator 2, "d" for done etc
	OpticalBoxSpread         string
	ReadOpticalBoxes         bool
	AncestorPath             string
	OmitPreviousComments     bool //this is for the checked stage, where we don't want earlier comments
	PropagateTextFieldValues bool // this is for enter active - copy textfield values out of pagedata into enter bar
}

type CoverPageCommand struct {
	Questions      []string
	FromPath       string
	ToPath         string
	ConfigPath     string
	ExamName       string
	TemplatePath   string
	SpreadName     string
	ProcessDetail  pagedata.ProcessDetail
	PathDecoration string
}

var (
	textFieldPrefix = "tf-"
	markDetected    = "mark-detected"
	opticalSuffix   = "-optical"
	isTesting       bool
	testroot        = "./tmp-delete-me"
	ExamStage       = []string{
		config,
		pageBad,
		acceptedReceipts,
		acceptedPapers,
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
		markerFlattened,
		markerProcessed,
		moderatorInactive,
		moderatorActive,
		moderatorReady,
		moderatorSent,
		moderatorBack,
		moderatorFlattened,
		moderatorProcessed,
		enterInactive,
		enterActive,
		enterReady,
		enterSent,
		enterBack,
		enterFlattened,
		enterProcessed,
		checkerReady,
		checkerSent,
		checkerBack,
		checkerFlattened,
		checkerProcessed,
		finalCover,
		reMarkerInactive,
		reMarkerActive,
		reMarkerReady,
		reMarkerSent,
		reMarkerBack,
		reMarkerFlattened,
		reMarkerProcessed,
		reModeratorInactive,
		reModeratorActive,
		reModeratorReady,
		reModeratorSent,
		reModeratorBack,
		reModeratorFlattened,
		reModeratorProcessed,
		reCheckerReady,
		reCheckerSent,
		reCheckerBack,
		reCheckerFlattened,
		reCheckerProcessed,
		finalPapers,
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
	config               = "00-config"
	pageBad              = "01-page-bad"
	acceptedReceipts     = "02-accepted-receipts"
	acceptedPapers       = "03-accepted-papers"
	tempImages           = "04-temporary-images"
	tempPages            = "04-temporary-pages"
	anonPapers           = "05-anonymous-papers"
	questionImages       = "06-question-images"
	questionPages        = "07-question-pages"
	questionReady        = "08-question-ready"
	questionSent         = "09-question-sent"
	questionBack         = "10-question-back"
	questionSplit        = "11-question-split"
	markerReady          = "20-marker-ready"
	markerSent           = "21-marker-sent"
	markerBack           = "22-marker-back"
	markerFlattened      = "23-marker-flattened"
	markerProcessed      = "24-marker-processed"
	moderatorInactive    = "28-moderator-inactive"
	moderatorActive      = "29-moderator-active"
	moderatorReady       = "30-moderator-ready"
	moderatorSent        = "31-moderator-sent"
	moderatorBack        = "32-moderator-back"
	moderatorFlattened   = "33-moderator-flattened"
	moderatorProcessed   = "34-moderator-processed"
	moderatorCover       = "35-moderator-cover"
	enterInactive        = "38-enter-inactive"
	enterActive          = "39-enter-active"
	enterReady           = "40-enter-ready"
	enterSent            = "41-enter-sent"
	enterBack            = "42-enter-back"
	enterFlattened       = "43-enter-flattened"
	enterProcessed       = "44-enter-processed"
	checkerCover         = "49-checker-cover"
	checkerReady         = "50-checker-ready"
	checkerSent          = "51-checker-sent"
	checkerBack          = "52-checker-back"
	checkerFlattened     = "53-checker-flattened"
	checkerProcessed     = "54-checker-processed"
	finalCover           = "55-final-cover"
	reMarkerInactive     = "58-remarker-inactive"
	reMarkerActive       = "59-remarker-active"
	reMarkerReady        = "60-remarker-ready"
	reMarkerSent         = "61-remarker-sent"
	reMarkerBack         = "62-remarker-back"
	reMarkerFlattened    = "63-remarker-flattened"
	reMarkerProcessed    = "64-remarker-processed"
	reModeratorInactive  = "68-remoderator-inactive"
	reModeratorActive    = "69-remoderator-active"
	reModeratorReady     = "70-remoderator-ready"
	reModeratorSent      = "71-remoderator-sent"
	reModeratorBack      = "72-remoderator-back"
	reModeratorFlattened = "73-remoderator-flattened"
	reModeratorProcessed = "74-remoderator-processed"
	reModeratorCover     = "75-remoderator-cover"
	reEnterInactive      = "78-reenter-inactive"
	reEnterActive        = "79-reenter-active"
	reEnterReady         = "80-reenter-ready"
	reEnterSent          = "81-reenter-sent"
	reEnterBack          = "82-reenter-back"
	reEnterFlattened     = "83-reenter-flattened"
	reEnterProcessed     = "84-reenter-processed"
	reCheckerCover       = "89-rechecker-cover"
	reCheckerReady       = "90-rechecker-ready"
	reCheckerSent        = "91-rechecker-sent"
	reCheckerBack        = "92-rechecker-back"
	reCheckerFlattened   = "93-rechecker-flattened"
	reCheckerProcessed   = "94-rechecker-processed"
	finalPapers          = "95-final-papers"
	reports              = "99-reports"

	inactive       = "inactive"
	doneDecoration = "d"
	labelling      = "labelling"
	marking        = "marking"
	remarking      = "remarking"
	moderating     = "moderating"
	remoderating   = "remoderating"
	checking       = "checking"
	rechecking     = "rechecking"
	entering       = "entering"
	reentering     = "rentering"

	marked      = "marked"
	moderated   = "moderated"
	entered     = "entered"
	checked     = "checked"
	remarked    = "remarked"
	remoderated = "remoderated"
	reentered   = "reentered"
	rechecked   = "rechecked"

	//External directories

	PageBad           = pageBad
	MarkerProcessed   = markerProcessed
	ModeratorActive   = moderatorActive
	ModeratorInactive = moderatorInactive

	//External stages
	Checked = checked
)
