/*
 * funtions to get PDF comments, and flatten them
 *
 *
 */

package comment

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/timdrysdale/gradex-cli/geo"
	pdfcore "github.com/timdrysdale/unipdf/v3/core"
	creator "github.com/timdrysdale/unipdf/v3/creator"
	pdf "github.com/timdrysdale/unipdf/v3/model"
)

type Comment struct {
	Pos  geo.Point
	Text string
	Page int
}

type Comments map[int][]Comment

func GetComments(reader *pdf.PdfReader) (Comments, error) {

	comments := make(map[int][]Comment)

	for p, page := range reader.PageList {

		if annotations, err := page.GetAnnotations(); err == nil {

			for _, annot := range annotations {

				if reflect.TypeOf(annot.GetContext()).String() == "*model.PdfAnnotationText" {

					if rect, is := annot.Rect.(*pdfcore.PdfObjectArray); is {

						x, err := strconv.ParseFloat(rect.Get(0).String(), 64)
						if err != nil {
							return comments, err
						}
						y, err := strconv.ParseFloat(rect.Get(1).String(), 64)
						if err != nil {
							return comments, err
						}

						newComment := Comment{
							Pos:  geo.Point{X: x, Y: y},
							Text: annot.Contents.String(),
							Page: p,
						}

						comments[p] = append(comments[p], newComment)

					}

				}

			}

		}

	}

	return comments, nil

}

func (c Comments) GetByPage(page int) []Comment {

	return c[page]

}

func DrawMarker(c *creator.Creator, comment Comment, label string) {

	r := c.NewRectangle(comment.Pos.X, comment.Pos.Y, 5*creator.PPMM, 5*creator.PPMM)
	r.SetBorderColor(creator.ColorYellow)
	r.SetFillColor(creator.ColorYellow)
	c.Draw(r)
	p := c.NewParagraph(fmt.Sprintf("[%s]", label))
	p.SetPos(comment.Pos.X, comment.Pos.Y)
	c.Draw(p)

}

func DrawText(c *creator.Creator, comment Comment, label string, X, Y float64) {

	p := c.NewParagraph(fmt.Sprintf("[%s] %s", label, comment.Text))
	p.SetPos(X, Y)
	c.Draw(p)

}

func DrawComment(c *creator.Creator, comment Comment, label string, X, Y float64) {
	comment.Pos.Y = c.Height() - comment.Pos.Y
	DrawText(c, comment, label, X, Y)
	DrawMarker(c, comment, label)
	comment.Pos.X = X
	comment.Pos.Y = Y
	DrawMarker(c, comment, label)
}
