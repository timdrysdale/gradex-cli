package report

/*
// things that did not sit in pagedata, which is not specifically tied to marking questions
// but also supports question annotation and scan checking

// use section for (a), (b) and number for (i)
type QuestionDetails struct {
	UUID           string            `json:"UUID"`
	Name           string            `json:"name"` //what to call it in a dropbox etc
	Section        string            `json:"section"`
	Number         int               `json:"number"` //No Harry Potter Platform 9&3/4 questions
	Parts          []QuestionDetails `json:"parts"`
	MarksAvailable float64           `json:"marksAvailable"`
	MarksAwarded   float64           `json:"marksAwarded"`
	Marking        []MarkingAction   `json:"markers"`
	Moderating     []MarkingAction   `json:"moderators"`
	Checking       []MarkingAction   `json:"checkers"`
	Sequence       int               `json:"sequence"`
	UnixTime       int64             `json:"unixTime"`
	Previous       string            `json:"previous"`
}

type MarkDetails struct {
	Given     float64 `json:"given"`
	Available float64 `json:"available"`
	Comment   float64 `json:"comment"`
}

type MarkingAction struct {
	Actor    string         `json:"actor"`
	Contact  ContactDetails `json:"contact"`
	Mark     MarkDetails    `json:"mark"`
	Done     bool           `json:"done"`
	UnixTime int64          `json:"unixTime"`
	Custom   CustomDetails  `json:"custom"`
}

type ParameterDetails struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Sequence int    `json:"sequence"`
}


type cmdOptions struct {
	pdfPassword string
}

type ScanResult struct {
	ScanPerfect            bool   `csv:"ScanPerfect"`
	ScanRotated            bool   `csv:"ScanRotated"`
	ScanContrast           bool   `csv:"ScanContrast"`
	ScanFaint              bool   `csv:"ScanFaint"`
	ScanIncomplete         bool   `csv:"ScanIncomplete"`
	ScanBroken             bool   `csv:"ScanBroken"`
	ScanComment1           string `csv:"ScanComment1"`
	ScanComment2           string `csv:"ScanComment2"`
	HeadingPerfect         bool   `csv:"HeadingPerfect"`
	HeadingVerbose         bool   `csv:"HeadingVerbose"`
	HeadingNoLine          bool   `csv:"HeadingNoLine"`
	HeadingNoQuestion      bool   `csv:"HeadingNoQuestion"`
	HeadingNoExamNumber    bool   `csv:"HeadingNoExamNumber"`
	HeadingAnonymityBroken bool   `csv:"HeadingAnonymityBroken"`
	HeadingComment1        string `csv:"HeadingComment1"`
	HeadingComment2        string `csv:"HeadingComment2"`
	FilenamePerfect        bool   `csv:"FilenamePerfect"`
	FilenameVerbose        bool   `csv:"FilenameVerbose"`
	FilenameNoCourse       bool   `csv:"FilenameNoCourse"`
	FilenameNoID           bool   `csv:"FilenameNoID"`
	InputFile              string `csv:"InputFile"`
	BatchFile              string `csv:"BatchFile"`
	BatchPage              int    `csv:"BatchPage"`
	Submission             parselearn.Submission
}

//type Submission struct {
//	FirstName          string  `csv:"FirstName"`
//	LastName           string  `csv:"LastName"`
//	Matriculation      string  `csv:"Matriculation"`
//	Assignment         string  `csv:"Assignment"`
//	DateSubmitted      string  `csv:"DateSubmitted"`
//	SubmissionField    string  `csv:"SubmissionField"`
//	Comments           string  `csv:"Comments"`
//	OriginalFilename   string  `csv:"OriginalFilename"`
//	Filename           string  `csv:"Filename"`
//	ExamNumber         string  `csv:"ExamNumber"`
//	MatriculationError string  `csv:"MatriculationError"`
//	ExamNumberError    string  `csv:"ExamNumberError"`
//	FiletypeError      string  `csv:"FiletypeError"`
//	FilenameError      string  `csv:"FilenameError"`
//	NumberOfPages      string  `csv:"NumberOfPages"`
//	FilesizeMB         float64 `csv:"FilesizeMB"`
//	NumberOfFiles      int     `csv:"NumberOfFiles"`
//}
*/
