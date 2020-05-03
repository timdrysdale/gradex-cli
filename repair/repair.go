package repair

import (
	"fmt"

	"github.com/rs/zerolog"
)

func RepairPDF(path string, textMap map[int]string, logger *zerolog.Logger) error {

	//TODO as a repair / catch up tool
	//Allow someone to do future steps of processing using the current tool
	//likely insertion points :- flatten,

	// anonPaper,
	//	markDone,
	//	moderateDone,
	//	checkDone,
	//	ReMarkDone,
	//	ReCheckDone

	// LOAD PDF
	// EXTRACT ALL PAGETEXT, RAW to Data dic using k:"pagetext","v":<pagetext>
	// EXTRACT ALL TEXTFIELDS, to Data dic using k:<textfield-name>,v:<textfield-value>>
	// Extract comments
	// FLATTEN PAGE TO IMAGE
	// Render a plain spread of the same size as the original image
	// insert the new page data into current on each page
	// insert the extracted data into previous, the Data fields
	// on a per page basis

	return fmt.Errorf("Not implemented")
}
