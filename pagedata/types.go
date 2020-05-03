package pagedata

type PageData struct {
	Exam        ExamDetails         `json:"exam"`
	Author      AuthorDetails       `json:"author"`
	Page        PageDetails         `json:"page"`
	Contact     ContactDetails      `json:"contact"`
	Submission  SubmissionDetails   `json:"submission"`
	Questions   []QuestionDetails   `json:"questions"`
	Processing  []ProcessingDetails `json:"processing"`
	Custom      []CustomDetails     `json:"custom"`
	Revision    int                 `json:"revision"`
	PreparedFor string              `json:"preparedfor"`
	ToDo        string              `json:"todo"`
}

// don't use this in anonymous pages
type SubmissionDetails struct {
	FilePrefix       string `json:"filePrefix"`
	OriginalFilename string `json:"originalFilename"`
	OriginalFormat   string `json:"originalFormat"`
	NewFilename      string `json:"newFilename"`
	NewFormat        string `json:"newFormat"`
}

type ExamDetails struct {
	CourseCode string `json:"courseCode"`
	Diet       string `json:"diet"`
	Date       string `json:"date"`
	UUID       string `json:"UUID"`
}

type AuthorDetails struct {
	Anonymous string `json:"Anonymous"`
	Identity  string `json:"Identity"`
}

type PageDetails struct {
	UUID     string `json:"UUID"`
	Number   int    `json:"number"`
	Of       int    `json:"of"`
	Filename string `json:"filename"`
}

type ContactDetails struct {
	Name    string `json:"name"`
	UUID    string `json:"UUID"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

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

type MarkingAction struct {
	Actor    string         `json:"actor"`
	Contact  ContactDetails `json:"contact"`
	Mark     MarkDetails    `json:"mark"`
	Done     bool           `json:"done"`
	UnixTime int64          `json:"unixTime"`
	Custom   CustomDetails  `json:"custom"`
}

type MarkDetails struct {
	Given     float64 `json:"given"`
	Available float64 `json:"available"`
	Comment   float64 `json:"comment"`
}

type CustomDetails struct {
	Key   string `json:"name"`
	Value string `json:"value"`
}

type ProcessingDetails struct {
	UUID       string             `json:"UUID"`
	Previous   string             `json:"previous"`
	UnixTime   int64              `json:"unixTime"`
	Name       string             `json:"name"`
	Parameters []ParameterDetails `json:"parameters"`
	By         ContactDetails     `json:"by"`
	Sequence   int                `json:"sequence"`
}

type ParameterDetails struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Sequence int    `json:"sequence"`
}

const (
	StartTag       = "<gradex-pagedata>"
	EndTag         = "</gradex-pagedata>"
	StartTagOffset = len(StartTag)
	EndTagOffset   = len(EndTag)
)
