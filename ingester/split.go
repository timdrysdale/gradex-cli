package ingester

import "github.com/timdrysdale/gradex-cli/image"

// simplified https://github.com/catherinelu/evangelist/blob/master/server.go

func ConvertPDFToJPEGs(pdfPath string, jpegPath string, outputFile string) error {

	return image.ConvertPDFToJPEGs(pdfPath, jpegPath, outputFile)
}

func CropToQuestion(inputPath, outputPath string) error {

	return image.CropToQuestion(inputPath, outputPath)

}
