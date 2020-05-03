# pdfcomment
get pdf comments out of doc / construct flattened versions

## Why 
Since editors [don't respect read-only settings](https://stackoverflow.com/questions/49941644/ios-pdfkit-make-text-widget-pdfannotation-readonly), we need to flatten text annotations

## Formatting

for a proper sticky note ...
```
// NewParagraph creates a new text paragraph.
// Default attributes:
// Font: Helvetica,
// Font size: 10
// Encoding: WinAnsiEncoding
// Wrap: enabled
// Text color: black
func (c *Creator) NewParagraph(text string) *Paragraph {
	return newParagraph(text, c.NewTextStyle())
}
```