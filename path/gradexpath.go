package gradexpath

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
		reports,
	}
)

const (
	config = "00-config"

	tempImages = "03-temporary-images"
	tempPages  = "04-temporary-pages"

	markedCombined   = "23-marked-combined"
	markedMerged     = "24-marked-merged"
	markedPruned     = "25-marked-pruned" //whatever gets trimmed goes here for potential audit
	markedReady      = "26-marked-ready"
	moderateActive   = "30-moderate-active"
	moderateInActive = "31-moderate-inactive"

	moderateInActiveBack = "35-moderate-inactive-back"
	moderatedCombined    = "36-moderated-combined"
	moderatedMerged      = "37-moderated-merged"
	moderatedPruned      = "38-moderated-pruned"
	moderatedReady       = "39-moderated-ready"

	checkedCombined = "43-checked-combined"
	checkedMerged   = "44-checked-merged"
	checkedPruned   = "45-checked-pruned"
	checkedReady    = "46-checked-ready"
	reports         = "99-reports"

	acceptedReceipts = "02-accepted-receipts"
	acceptedPapers   = "03-accepted-papers"
	anonPapers       = "05-anonymous-papers"

	markerReady = "20-marker-ready"
	markerSent  = "21-marker-sent"
	markerBack  = "22-marker-back"

	moderatorReady = "32-moderator-ready"
	moderatorSent  = "33-moderator-sent"
	moderatorBack  = "34-moderator-back"

	checkerReady = "40-checker-ready"
	checkerSent  = "41-checker-sent"
	checkerBack  = "42-checker-back"

	remarkerReady = "50-remarker-ready"
	remarkerSent  = "51-remarker-sent"
	remarkerBack  = "52-remarker-back"

	recheckerReady = "60-rechecker-ready"
	recheckerSent  = "61-rechecker-sent"
	recheckerBack  = "62-rechecker-back"

	N = 3
)

type DirSum struct {
	Files int
	Pages int
	Size  int
}

/*
func dirsummary(path string) []int {
	return []int{}

}

// Path is in the package name, so anything without another noun IS a path
// non path things need a noun

func limit(initials string, N int) string {
	if len(initials) < 3 {
		N = len(initials)
	}
	return strings.ToUpper(initials[0:N])
}

func markedcombined(exam string) string {
	return filepath.Join(Exam(), exam, markedCombined)
}
func markedmerged(exam string) string {
	return filepath.Join(Exam(), exam, markedMerged)
}
func markedpruned(exam string) string {
	return filepath.Join(Exam(), exam, markedPruned)
}
func markedReady(exam string) string {
	return filepath.Join(Exam(), exam, markedReady)
}
func moderateActive(exam string) string {
	return filepath.Join(Exam(), exam, moderateActive)
}

func moderatedCombined(exam string) string {
	return filepath.Join(Exam(), exam, moderatedCombined)
}
func moderatedMerged(exam string) string {
	return filepath.Join(Exam(), exam, moderatedMerged)
}
func moderatedPruned(exam string) string {
	return filepath.Join(Exam(), exam, moderatedPruned)
}

func moderatedReady(exam string) string {
	return filepath.Join(Exam(), exam, moderatedReady)
}

func moderateInActive(exam string) string {
	return filepath.Join(Exam(), exam, moderateInActive)
}

func moderatedInActiveBack(exam string) string {
	return filepath.Join(Exam(), exam, moderateInActiveBack)
}

func checkedCombined(exam string) string {
	return filepath.Join(Exam(), exam, checkedCombined)
}
func checkedMerged(exam string) string {
	return filepath.Join(Exam(), exam, checkedMerged)
}
func checkedPruned(exam string) string {
	return filepath.Join(Exam(), exam, checkedPruned)
}
func checkedReady(exam string) string {
	return filepath.Join(Exam(), exam, checkedReady)
}

func doneDecoration() string {
	return "d"
}

func markerABCDecoration(initials string) string {
	return fmt.Sprintf("-ma%s", limit(initials, N))
}

func markerABCDirName(initials string) string {
	return limit(initials, N)
}

func moderatorABCDecoration(initials string) string {
	return fmt.Sprintf("-mo%s", limit(initials, N))
}

func moderatorABCDirName(initials string) string {
	return limit(initials, N)
}

func checkerABCDecoration(initials string) string {
	return fmt.Sprintf("-c%s", limit(initials, N))
}

func checkerABCDirName(initials string) string {
	return limit(initials, N)
}

func markerNDecoration(number int) string {
	return fmt.Sprintf("-ma%d", number)
}

func markerNDirName(number int) string {
	return fmt.Sprintf("marker%d", number)
}

func moderatorNDecoration(number int) string {
	return fmt.Sprintf("-mo%d", number)
}

func moderatorNDirName(number int) string {
	return fmt.Sprintf("moderator%d", number)
}

func checkerNDecoration(number int) string {
	return fmt.Sprintf("-c%d", number)
}

func checkerNDirName(number int) string {
	return fmt.Sprintf("checker%d", number)
}

func markerReady(exam, marker string) string {
	path := filepath.Join(Exam(), exam, markerReady, limit(marker, N))
	ensureDirAll(path)
	return path
}

func markerSent(exam, marker string) string {
	path := filepath.Join(Exam(), exam, markerSent, limit(marker, N))
	ensureDirAll(path)
	return path
}

func markerBack(exam, marker string) string {
	path := filepath.Join(Exam(), exam, markerBack, limit(marker, N))
	ensureDirAll(path)
	return path
}

func moderatorReady(exam, moderator string) string {
	path := filepath.Join(Exam(), exam, moderatorReady, limit(moderator, N))
	ensureDirAll(path)
	return path
}

func moderatorSent(exam, moderator string) string {
	path := filepath.Join(Exam(), exam, moderatorSent, limit(moderator, N))
	ensureDirAll(path)
	return path
}

func moderatorBack(exam, moderator string) string {
	path := filepath.Join(Exam(), exam, moderatorBack, limit(moderator, N))
	ensureDirAll(path)
	return path
}

func checkerReady(exam, checker string) string {
	path := filepath.Join(Exam(), exam, checkerReady, limit(checker, N))
	ensureDirAll(path)
	return path
}

func checkerSent(exam, checker string) string {
	path := filepath.Join(Exam(), exam, checkerSent, limit(checker, N))
	ensureDirAll(path)
	return path
}

func checkerBack(exam, checker string) string {
	path := filepath.Join(Exam(), exam, checkerBack, limit(checker, N))
	ensureDirAll(path)
	return path
}

func reMarkerReady(exam, marker string) string {
	path := filepath.Join(Exam(), exam, remarkerReady, limit(marker, N))
	ensureDirAll(path)
	return path
}

func reMarkerSent(exam, marker string) string {
	path := filepath.Join(Exam(), exam, remarkerSent, limit(marker, N))
	ensureDirAll(path)
	return path
}

func reMarkerBack(exam, marker string) string {
	path := filepath.Join(Exam(), exam, remarkerBack, limit(marker, N))
	ensureDirAll(path)
	return path
}

func reCheckerReady(exam, checker string) string {
	path := filepath.Join(Exam(), exam, recheckerReady, limit(checker, N))
	ensureDirAll(path)
	return path
}

func reCheckerSent(exam, checker string) string {
	path := filepath.Join(Exam(), exam, recheckerSent, limit(checker, N))
	ensureDirAll(path)
	return path
}

func reCheckerBack(exam, checker string) string {
	path := filepath.Join(Exam(), exam, recheckerBack, limit(checker, N))
	ensureDirAll(path)
	return path
}

func flattenLayoutSVG() string {
	return filepath.Join(IngestTemplate(), "layout-flatten-312pt.svg")
}

func overlayLayoutSVG() string {
	return filepath.Join(OverlayTemplate(), "layout.svg")
}

func acceptedPapers(exam string) string {
	return filepath.Join(Exam(), exam, acceptedPapers)
}

func acceptedReceipts(exam string) string {
	return filepath.Join(Exam(), exam, acceptedReceipts)
}

//TODO in flatten, swap these paths for the general named ones below
func acceptedPaperImages(exam string) string {
	return filepath.Join(Exam(), exam, tempImages)
}
func acceptedPaperPages(exam string) string {
	return filepath.Join(Exam(), exam, tempPages)
}
func paperImages(exam string) string {
	return filepath.Join(Exam(), exam, tempImages)
}
func paperPages(exam string) string {
	return filepath.Join(Exam(), exam, tempPages)
}

func anonymousPapers(exam string) string {
	return filepath.Join(Exam(), exam, anonPapers)
}

func identity() string {
	return filepath.Join(Etc(), "identity")
}

func identityCSV() string {
	return filepath.Join(Identity(), "identity.csv")
}

func ingest() string {
	return filepath.Join(Root(), "ingest")
}

func ingestTemplate() string {
	return filepath.Join(IngestConf(), "template")
}

func overlayTemplate() string {
	return filepath.Join(OverlayConf(), "template")

}
func tempPdf() string {
	return filepath.Join(Root(), "temp-pdf")
}

func tempTxt() string {
	return filepath.Join(Root(), "temp-txt")
}

func export() string {
	return filepath.Join(Root(), "export")
}

func etc() string {
	return filepath.Join(Root(), "etc")
}

func variable() string {
	return filepath.Join(Root(), "var")
}

func usr() string {
	return filepath.Join(Root(), "usr")
}

func exam() string {
	return filepath.Join(Usr(), "exam")
}

func ingestConf() string {
	return filepath.Join(Etc(), "ingest")
}

func overlayConf() string {
	return filepath.Join(Etc(), "overlay")
}

func extractConf() string {
	return filepath.Join(Etc(), "extract")
}

func setupConf() string {
	return filepath.Join(Etc(), "setup")
}

func setTesting() { //need this when testing other tools
	isTesting = true
}

func root() string {
	if isTesting {
		return testroot
	}
	return root
}

func getExamPath(name string) string {
	return filepath.Join(Exam(), name)
}

func getExamStagePath(name, stage string) string {
	return filepath.Join(Exam(), name, stage)
}

func setupGradexPaths() error {

	paths := []string{
		root(),
		ingest(),
		identity(),
		export(),
		variable(),
		usr(),
		exam(),
		tempPdf(),
		tempTxt(),
		etc(),
		ingestconf(),
		overlayconf(),
		ongestTemplate(),
		overlayTemplate(),
		extractConf(),
		setupConf(),
	}

	for _, path := range paths {

		err := ensureDirAll(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func setupExamPaths(exam string) error {
	// don't use ensureDirAll so it flags if we are not otherwise setup
	err := ensureDir(GetExamPath(exam))
	if err != nil {
		return err
	}

	for _, stage := range ExamStage {
		err := ensureDir(GetExamStagePath(exam, stage))
		if err != nil {
			return err
		}
	}
	return nil
}

// if the source file is not newer, it's not an error
// we just won't move it - anything left we deal with later
func moveIfNewerThanDestination(source, destination string) error {

	//check both exist
	sourceInfo, err := os.Stat(source)

	if err != nil {
		return err
	}

	destinationInfo, err := os.Stat(destination)

	// source newer by definition if destination does not exist
	if os.IsNotExist(err) {
		err = os.Rename(source, destination)
		return err
	}
	if err != nil {
		return err
	}
	if sourceInfo.ModTime().After(destinationInfo.ModTime()) {
		err = os.Rename(source, destination)
		return err
	}

	return nil

}

func isSameAsSelfInDir(source, destinationDir string) bool {

	//check both exist
	sourceInfo, err := os.Stat(source)

	if err != nil {
		return false
	}

	destination := filepath.Join(destinationDir, filepath.Base(source))

	destinationInfo, err := os.Stat(destination)

	if err != nil {
		return false
	}

	timeEqual := sourceInfo.ModTime().Equal(destinationInfo.ModTime())
	sizeEqual := sourceInfo.Size() == destinationInfo.Size()
	return timeEqual && sizeEqual

}

func moveIfNewerThanDestinationInDir(source, destinationDir string) error {

	//check both exist
	sourceInfo, err := os.Stat(source)

	if err != nil {
		return err
	}

	destination := filepath.Join(destinationDir, filepath.Base(source))

	destinationInfo, err := os.Stat(destination)

	// source newer by definition if destination does not exist
	if os.IsNotExist(err) {
		err = os.Rename(source, destination)
		return err
	}
	if err != nil {
		return err
	}
	if sourceInfo.ModTime().After(destinationInfo.ModTime()) {
		err = os.Rename(source, destination)
		return err
	}

	return nil

}

func examDiet(exam string) string {

	m := int(time.Now().Month())

	switch {
	case m > 4 && m < 6:
		return fmt.Sprintf("May-%d", time.Now().Year())
	case m > 6 && m < 10:
		return fmt.Sprintf("Aug-%d", time.Now().Year())
	case m > 10 || m < 3:
		return fmt.Sprintf("Dec-%d", time.Now().Year())
	default:
		return fmt.Sprintf("%d", time.Now().Year())
	}

}
*/
