package pagedata

import "github.com/timdrysdale/gradex-cli/comment"

const (
	IsPage    = "page"
	IsRegion  = "region"
	IsCover   = "cover"
	IsMontage = "montage"

	IsAnonymous = "anonymous"
	IsIdentity  = "identity"
)

// Used in triaging files at ingest/staging
type Summary struct {
	Is   string //page, region, cover-page etc
	What string //item
	For  string //proc
	ToDo string //proc
}

// >>>>>>>>>>>>>>>>>>>>> Actual page data >>>>>>>>>>>>>>>>>>>>>>>>

// use an array so that we can decide later that keys to do the mapping of history on
// e.g. we will smash together ALL markers' work, and key on ToDo-Marker-Revision etc.

type PageData struct {
	Current  PageDetail   `json:"current"`
	Previous []PageDetail `json:"previous"`
}

// use custom data for group authorship, if individual authorship must be tracked here
// else use a group id e.g. group-<uuid> which has the individual authors recorded
// elsewhere, along with the original submission.
type PageDetail struct {
	Is       string            `json:"is"` //page, region
	Own      FileDetail        `json:"own"`
	Original FileDetail        `json:"original"`
	Current  FileDetail        `json:"current"`
	Item     ItemDetail        `json:"item"`
	Process  ProcessDetail     `json:"process"`
	UUID     string            `json:"UUID"` //for mapping the previous page datas later
	Follows  string            `json:"follows"`
	Revision int               `json:"revision"` //if we want to rewrite history ....
	Data     []Field           `json:"data"`
	Comments []comment.Comment `json:"comments"`
}

type FileDetail struct {
	Path   string `json:"path"`
	UUID   string `json:"UUID"`
	Number int    `json:"number"`
	Of     int    `json:"of"`
}

//whotype exam number:EN matriculation number:UUN etc
type ItemDetail struct {
	What    string `json:"what"`
	When    string `json:"when"`
	Who     string `json:"who"`
	UUID    string `json:"UUID"`
	WhoType string `json:"whoType"`
}

type ProcessDetail struct {
	Name     string  `json:"name"`
	UUID     string  `json:"UUID"` // process batch UUID
	UnixTime int64   `json:"unixTime"`
	For      string  `json:"for"`
	ToDo     string  `json:"toDo"`
	By       string  `json:"by"`
	Data     []Field `json:"data"`
}

type Field struct { //for clarity in code, and brevity in pagedata
	Key   string `json:"k"`
	Value string `json:"v"`
}

const (
	StartTag        = "<gradex-pagedata>"
	EndTag          = "</gradex-pagedata>"
	StartTagOffset  = len(StartTag)
	EndTagOffset    = len(EndTag)
	StartHash       = "<hash>"
	EndHash         = "</hash>"
	StartHashOffset = len(StartHash)
	EndHashOffset   = len(EndHash)
)
