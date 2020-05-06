package ingester

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/timdrysdale/chmsg"
	"github.com/timdrysdale/gradex-cli/merge"
	"github.com/timdrysdale/gradex-cli/parsesvg"
	"github.com/timdrysdale/pool"
)

func (g *Ingester) Annotate(exam string) error {

	logger := g.logger.With().Str("process", "annotate").Logger()

	mc := chmsg.MessagerConf{
		ExamName:     exam,
		FunctionName: "annotate",
		TaskName:     "question",
	}

	cm := chmsg.New(mc, g.msgCh, g.timeout)

	err := g.SetupExamPaths(exam)
	if err != nil {
		return err
	}

	inputPath, err := g.GetFileList(g.AcceptedPaperImages(exam))

	if err != nil {
		return err
	}

	N := len(inputPath)

	tasks := []*pool.Task{}

	for i := 0; i < N; i++ {

		pageImage := inputPath[i]
		exam := exam
		newtask := pool.NewTask(func() error {
			err := g.annotateOneImage(pageImage, exam)
			if err != nil {
				fmt.Println(err)
				g.logger.Error().
					Str("error", err.Error()).
					Msg("Error annotating")
			}
			return err
		})
		tasks = append(tasks, newtask)
	}

	p := pool.NewPool(tasks, runtime.GOMAXPROCS(-1))

	p.Run()

	var numErrors int
	for _, task := range p.Tasks {
		if task.Err != nil {
			logger.Error().
				Str("error", task.Err.Error()).
				Msg("Annotating task error")
			numErrors++
		}
	}

	pageFiles, err := g.GetFileList(g.QuestionPages(exam))

	var divided [][]string

	chunkSize := 100 //len(pageFiles)

	for i := 0; i < len(pageFiles); i += chunkSize {
		end := i + chunkSize

		if end > len(pageFiles) {
			end = len(pageFiles)
		}

		divided = append(divided, pageFiles[i:end])

		outputPath := filepath.Join(g.QuestionReady(exam), fmt.Sprintf("batch-%02d.pdf", i+1))
		err = merge.PDF(pageFiles[i:end], outputPath)
		if err != nil {
			logger.Error().
				Str("file", outputPath).
				Str("error", err.Error()).
				Msg(fmt.Sprintf("Error merging processed pages for (%s) because %v\n", outputPath, err))

			return err
		}
	}
	if err == nil {
		cm.Send("Finished Processing annotations")
		logger.Info().
			Str("exam", exam).
			Msg("Finished annotating questions")
	} else {
		logger.Error().
			Str("exam", exam).
			Str("error", err.Error()).
			Msg("Finished annotating questions")

	}
	return err
}

func (g *Ingester) annotateOneImage(imagePath, exam string) error {

	cropImage := filepath.Join(g.QuestionImages(exam), filepath.Base(imagePath)) //keep extension

	err := CropToQuestion(imagePath, cropImage)
	if err != nil {
		g.logger.Error().
			Str("error", err.Error()).
			Msg("Error cropping question image")
	}

	if strings.ToLower(filepath.Ext(imagePath)) != ".jpg" {
		return errors.New(fmt.Sprintf("%s does not appear to be a jpg", imagePath))
	}

	// construct image name
	svgLayoutPath := filepath.Join(g.Root(), "etc/annotate/template/layout-annotate.svg")

	baseName := strings.TrimSuffix(filepath.Base(cropImage), ".jpg")

	outputPath := filepath.Join(g.QuestionPages(exam), baseName+".pdf")
	fmt.Println(g.QuestionPages(exam))
	fmt.Println(outputPath)
	docTextPrefills := make(map[int]parsesvg.PagePrefills)

	docTextPrefills[0] = make(parsesvg.PagePrefills)

	docTextPrefills[0]["filename"] = baseName

	prefillImagePaths := make(map[string]string)

	prefillImagePaths["annotate-question"] = strings.TrimSuffix(cropImage, ".jpg")

	contents := parsesvg.SpreadContents{
		SvgLayoutPath:         svgLayoutPath,
		SpreadName:            "annotate",
		PreviousImagePath:     "",
		PageNumber:            0,
		PdfOutputPath:         outputPath,
		Prefills:              docTextPrefills,
		PrefillImagePaths:     prefillImagePaths,
		TemplatePathsRelative: true,
	}

	err = parsesvg.RenderSpreadExtra(contents)
	if err != nil {
		fmt.Println(err)

	}
	return err

}
