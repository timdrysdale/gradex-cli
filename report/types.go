package report

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
