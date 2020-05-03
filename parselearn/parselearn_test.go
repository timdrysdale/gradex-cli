package parselearn

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Example receipt (anonymised)

//Name: First Last (sxxxxxxx)
//Assignment: Practice Exam Drop Box
//Date Submitted: Monday, dd April yyyy hh:mm:ss o'clock BST
//Current Mark: Needs Marking
//
//Submission Field:
//There is no student submission text data for this assignment.
//
//Comments:
//There are no student comments for this assignment.
//
//Files:
//	Original filename: OnlineExam-Bxxxxxx.pdf
//	Filename: Practice Exam Drop Box_sxxxxxxx_attempt_yyyy-mm-dd-hh-mm-ss_OnlineExam-Bxxxxxx.pdf

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func TestProcessRevision(t *testing.T) {

	sub := Submission{}

	processRevision("Revision: 2", &sub)
	assertEqual(t, sub.Revision, 2)

}
func TestProcessName(t *testing.T) {

	sub := Submission{}

	processName("Name: Donald The Duck (sxxxxxxx)", &sub)
	assertEqual(t, sub.FirstName, "-")

	assertEqual(t, sub.LastName, "Donald The Duck")
	assertEqual(t, sub.Matriculation, "sxxxxxxx")
}

func TestProcessAction(t *testing.T) {

	sub := Submission{}
	processAction("Action: ignore", &sub)
	assertEqual(t, sub.Action, "ignore")
}

func TestProcessAssignment(t *testing.T) {

	sub := Submission{}
	processAssignment("Assignment: Practice Exam Drop Box", &sub)
	assertEqual(t, sub.Assignment, "Practice Exam Drop Box")
}

func TestDateSubmitted(t *testing.T) {

	sub := Submission{}
	processDateSubmitted("Date Submitted: Monday, dd April yyyy hh:mm:ss o'clock BST", &sub)
	assertEqual(t, sub.DateSubmitted, "Monday, dd April yyyy hh:mm:ss o'clock BST")
}

func TestSubmissionField(t *testing.T) {

	sub := Submission{}
	processSubmission("There is no student submission text data for this assignment.", &sub)
	assertEqual(t, sub.SubmissionField, "There is no student submission text data for this assignment.")
}

func TestComments(t *testing.T) {

	sub := Submission{}
	processSubmission("There are no student comments for this assignment", &sub)
	assertEqual(t, sub.SubmissionField, "There are no student comments for this assignment")
}

func TestOriginalFilename(t *testing.T) {

	sub := Submission{}
	processOriginalFilename("Original filename: OnlineExam-Bxxxxxx.pdf", &sub)
	assertEqual(t, sub.OriginalFilename, "OnlineExam-Bxxxxxx.pdf")
}

func TestFilename(t *testing.T) {

	sub := Submission{}
	processFilename("Filename: Practice Exam Drop Box_sxxxxxxx_attempt_yyyy-mm-dd-hh-mm-ss_OnlineExam-Bxxxxxx.pdf", &sub)
	assertEqual(t, sub.Filename, "Practice Exam Drop Box_sxxxxxxx_attempt_yyyy-mm-dd-hh-mm-ss_OnlineExam-Bxxxxxx.pdf")
}

func TestParseFile(t *testing.T) {

	sub, err := ParseLearnReceipt("./test/receipt2.txt")
	if err != nil {
		t.Error(err)
	}

	assertEqual(t, sub.FirstName, "-")
	assertEqual(t, sub.LastName, "John Smith")
	assertEqual(t, sub.Matriculation, "s00000000")
	assertEqual(t, sub.Assignment, "Some Exam Or Other")
	assertEqual(t, sub.DateSubmitted, "Tuesday, dd April yyyy hh:mm:ss o'clock BST")
	assertEqual(t, sub.SubmissionField, "There is no student submission text data for this assignment.")
	assertEqual(t, sub.Comments, "There are no student comments for this assignment.")
	assertEqual(t, sub.OriginalFilename, "ENGI1234-Bxxxxxx.pdf")
	assertEqual(t, sub.Filename, "Practice Exam Drop Box_sxxxxxxx_attempt_yyyy-mm-dd-hh-mm-ss_ENGI1234-Bxxxxxx.pdf")

}

var expected1 = `Revision,Action,FirstName,LastName,Matriculation,Assignment,DateSubmitted,SubmissionField,Comments,OriginalFilename,Filename,ExamNumber,MatriculationError,ExamNumberError,FiletypeError,FilenameError,NumberOfPages,FilesizeMB,NumberOfFiles,OwnPath`

var expected2 = `0,,-,First Last,sxxxxxxx,Practice Exam Drop Box,"Monday, dd April yyyy hh:mm:ss o'clock BST",There is no student submission text data for this assignment.,There are no student comments for this assignment.,OnlineExam-Bxxxxxx.pdf,Practice Exam Drop Box_sxxxxxxx_attempt_yyyy-mm-dd-hh-mm-ss_OnlineExam-Bxxxxxx.pdf,,,,,,,0,1,./test/receipt1.txt`

var expected3 = `99,ignore,-,John Smith,s00000000,Some Exam Or Other,"Tuesday, dd April yyyy hh:mm:ss o'clock BST",There is no student submission text data for this assignment.,There are no student comments for this assignment.,ENGI1234-Bxxxxxx.pdf,Practice Exam Drop Box_sxxxxxxx_attempt_yyyy-mm-dd-hh-mm-ss_ENGI1234-Bxxxxxx.pdf,,,,,,,0,1,./test/receipt2.txt`

var expected4 = `1,,-,First Last,sxxxxxxx,Practice Exam Drop Box,"Monday, dd April yyyy hh:mm:ss o'clock BST",There is no student submission text data for this assignment.,There are no student comments for this assignment.,OnlineExam-Bxxxxxx.pdf,Practice Exam Drop Box_sxxxxxxx_attempt_yyyy-mm-dd-hh-mm-ss_OnlineExam-Bxxxxxx.pdf,,,,,,,0,1,./test/receipt3.txt`

func TestMarshallToFile(t *testing.T) {

	sub1, err := ParseLearnReceipt("./test/receipt1.txt")
	if err != nil {
		t.Error(err)
	}

	sub2, err := ParseLearnReceipt("./test/receipt2.txt")
	if err != nil {
		t.Error(err)
	}
	sub3, err := ParseLearnReceipt("./test/receipt3.txt")
	if err != nil {
		t.Error(err)
	}
	subs := []Submission{sub1, sub2, sub3}

	WriteSubmissionsToCSV(subs, "./test/out.csv")

	file, err := os.Open("./test/out.csv")

	if err != nil {
		t.Error(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	assertEqual(t, scanner.Text(), expected1)
	scanner.Scan()
	assertEqual(t, scanner.Text(), expected2)
	scanner.Scan()
	assertEqual(t, scanner.Text(), expected3)
	scanner.Scan()
	assertEqual(t, scanner.Text(), expected4)
}

func TestGetFiles(t *testing.T) {

	list, err := GetFilePaths("./test/receiptmulti.txt")

	assert.NoError(t, err)
	assert.Equal(t, []string{"A B-C.pdf", "D-D F.pdf", "G H_I J.pdf", "G H_I J Copy.pdf", "G H_I J Copy (2).pdf"}, list)

}

func TestHandleIgnoreReceipts(t *testing.T) {

	rm := make(map[string]Submission)

	rm["keep"] = Submission{Action: ""}
	rm["goaway"] = Submission{Action: "ignore"}

	assert.Equal(t, 2, len(rm))

	HandleIgnoreReceipts(&rm)

	assert.Equal(t, 1, len(rm))

	for k, _ := range rm {
		assert.Equal(t, "keep", k)
	}
}
