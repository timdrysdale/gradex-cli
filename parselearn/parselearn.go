package parselearn

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
)

type Submission struct {
	Revision           int     `csv:"Revision"`
	Action             string  `csv:"Action"`
	FirstName          string  `csv:"FirstName"`
	LastName           string  `csv:"LastName"`
	Matriculation      string  `csv:"Matriculation"`
	Assignment         string  `csv:"Assignment"`
	DateSubmitted      string  `csv:"DateSubmitted"`
	SubmissionField    string  `csv:"SubmissionField"`
	Comments           string  `csv:"Comments"`
	OriginalFilename   string  `csv:"OriginalFilename"`
	Filename           string  `csv:"Filename"`
	ExamNumber         string  `csv:"ExamNumber"`
	MatriculationError string  `csv:"MatriculationError"`
	ExamNumberError    string  `csv:"ExamNumberError"`
	FiletypeError      string  `csv:"FiletypeError"`
	FilenameError      string  `csv:"FilenameError"`
	NumberOfPages      string  `csv:"NumberOfPages"`
	FilesizeMB         float64 `csv:"FilesizeMB"`
	NumberOfFiles      int     `csv:"NumberOfFiles"`
	OwnPath            string  `csv:"OwnPath"`
}

func CheckFilename(receiptPath string) error {

	files, err := GetFilePaths(receiptPath)

	if err != nil {
		return fmt.Errorf("Could not get file list from receipt because %s", err.Error())
	}

	for _, file := range files {
		_, err := os.Stat(file)

		if os.IsNotExist(err) {
			return fmt.Errorf("Can't find file %s", file)

		}

	}
	return nil
}

func HandleIgnoreReceipts(receiptMap *map[string]Submission) {

	// ignore receipts where requested
	ignoreKeys := []string{}

	for k, v := range *receiptMap {
		if v.Action == "ignore" {
			ignoreKeys = append(ignoreKeys, k)
		}
	}

	for _, k := range ignoreKeys {
		delete(*receiptMap, k)
	}

}

func GetFilePaths(inputPath string) ([]string, error) {

	files := []string{}

	file, err := os.Open(inputPath)

	if err != nil {
		return files, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lline := strings.ToLower(line) //used for case insensitive field identification
		if strings.HasPrefix(lline, "filename:") {
			line = strings.TrimSpace(line)
			line = strings.TrimPrefix(line, "Filename:")
			line = strings.TrimSpace(line)
			files = append(files, line)
		}
	}

	return files, nil
}

func ParseLearnReceipt(inputPath string) (Submission, error) {

	sub := Submission{}

	sub.OwnPath = inputPath

	file, err := os.Open(inputPath)

	if err != nil {
		return sub, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

SCAN:
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		//use for case insenstive field identification only
		//need to keep capitalisation in filenames
		lline := strings.ToLower(line)

		switch {
		case strings.HasPrefix(lline, "revision:"):
			processRevision(line, &sub)
		case strings.HasPrefix(lline, "action:"):
			processAction(line, &sub)
		case strings.HasPrefix(lline, "name:"):
			processName(line, &sub)
		case strings.HasPrefix(lline, "assignment:"):
			processAssignment(line, &sub)
		case strings.HasPrefix(lline, "date submitted:"):
			processDateSubmitted(line, &sub)
		case strings.HasPrefix(lline, "submission field:"):
			scanner.Scan()
			processSubmission(scanner.Text(), &sub)
		case strings.HasPrefix(lline, "comments:"):
			scanner.Scan()
			processComments(scanner.Text(), &sub)
		case strings.HasPrefix(lline, "files:"):
			break SCAN
		default:
			continue
		}
	}

	// now read in the files ....
	// TODO figure out nested csv so we can record multiple files
	// meanwhile for safety, count the number of original files

	sub.NumberOfFiles = 0

	// we will take the first file as the one we expect to see renamed
	gotOriginal := false
	gotNew := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lline := strings.ToLower(line) //used for case insensitive field identification
		switch {
		case strings.HasPrefix(lline, "original filename:"):
			if !gotOriginal {
				processOriginalFilename(line, &sub)
				gotOriginal = true
			}
			sub.NumberOfFiles++
		case strings.HasPrefix(lline, "filename:"):
			if !gotNew {
				processFilename(line, &sub)
				gotNew = true
			}
		default:
			continue
		}

	}

	return sub, scanner.Err()
}

//Name: First Last (sxxxxxxx)
func processRevision(line string, sub *Submission) {

	m := strings.Index(line, ":")
	n := len(line)

	revision := strings.TrimSpace(line[m+1 : n])

	var rev int64

	rev, err := strconv.ParseInt(strings.TrimSpace(revision), 10, 64)

	if err != nil {

		rev = 0

	}

	sub.Revision = int(rev)
}

// we lower case this one since made by hand
func processAction(line string, sub *Submission) {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(strings.ToLower(line), "action:")
	sub.Action = strings.TrimSpace(line)
}

//Name: First Last (sxxxxxxx)
func processName(line string, sub *Submission) {

	m := strings.Index(line, ":")
	n := strings.Index(line, "(")
	p := strings.Index(line, ")")

	name := strings.TrimSpace(line[m+1 : n])
	matric := strings.TrimSpace(line[n+1 : p])

	sub.FirstName = "-"
	sub.LastName = name
	sub.Matriculation = matric
}

//Assignment: Practice Exam Drop Box
func processAssignment(line string, sub *Submission) {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "Assignment:")
	sub.Assignment = strings.TrimSpace(line)
}

//Date Submitted: Monday, dd April yyyy hh:mm:ss o'clock BST
func processDateSubmitted(line string, sub *Submission) {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "Date Submitted:")
	sub.DateSubmitted = strings.TrimSpace(line)
}

//Submission Field:
//There is no student submission text data for this assignment.
func processSubmission(line string, sub *Submission) {
	sub.SubmissionField = strings.TrimSpace(line)

}

//Comments:
//There are no student comments for this assignment
func processComments(line string, sub *Submission) {
	sub.Comments = strings.TrimSpace(line)
}

//Files:
//	Original filename: OnlineExam-Bxxxxxx.pdf
//	Filename: Practice Exam Drop Box_sxxxxxxx_attempt_yyyy-mm-dd-hh-mm-ss_OnlineExam-Bxxxxxx.pdf
func processOriginalFilename(line string, sub *Submission) {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "Original filename:")
	sub.OriginalFilename = strings.TrimSpace(line)
}
func processFilename(line string, sub *Submission) {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "Filename:")
	sub.Filename = strings.TrimSpace(line)
}

func WriteSubmissionsToCSV(subs []Submission, outputPath string) error {
	// wrap the marshalling library in case we need converters etc later
	file, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	return gocsv.MarshalFile(&subs, file)
}

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
