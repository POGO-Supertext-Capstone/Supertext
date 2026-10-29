package upload

import (
	"fmt"
	"os"
	"strings"

	"github.com/kenny-not-dead/gopptx"
)

// PresentationUpload holds data provided by the user in
// the form upload, as well as the temporay file location
// holding the copy for decoding
type PresentationUpload struct {
	PresentationName string
	Filename         string
	TmpFilePath      string
}

// runInfo holds info about a run of text, including its
// content and formatting
type runInfo struct {
	Text   string
	Bold   bool
	Italic bool
}

// SlideContent holds the info about an individual slide
// including its index (order in presentation), its unique slide ID
// and the text content in that slide
type SlideContent struct {
	Index   int
	ID      int
	Content string
}

// DecodePPTX goes throgh a powerpoint file and converts the slide content
// into plain text with appropiate HTML tags, include <p> to denote a new line,
// <i> to denote italics, and <b> to denote boldness.
func DecodePPTX(uploadData PresentationUpload) ([]SlideContent, error) {
	defer os.Remove(uploadData.TmpFilePath)

	f, err := gopptx.OpenFile(uploadData.TmpFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open temp file: %w", err)
	}
	defer f.Close()

	// read slide data
	slides := make([]SlideContent, len(f.Presentation.Slides.Slide))

	for i, slide := range f.Presentation.Slides.Slide {
		shapes, err := f.GetShapes(slide.SlideID)
		if err != nil {
			return nil, fmt.Errorf("Error reading shapes from slide %d: %v", slide.SlideID, err)
		}

		//fmt.Printf("Slide %d (index %d):\n", slide.SlideID, i)
		slides[i].ID = slide.SlideID
		slides[i].Index = i

		content := ""
		for _, shape := range shapes {
			//fmt.Printf("Type: %T, Value: %+v\n", shape.TextBody.Paragraph, shape.TextBody.Paragraph)
			for _, paragraph := range shape.TextBody.Paragraph {
				content += "<p>" + decodeParagraph(paragraph) + "</p>"
			}
		}
		slides[i].Content = content
	}

	return slides, nil
}

// decodeRun decodes the properties of a run and saves them to a
// run info struct for later use
func decodeRun(run gopptx.DecodeRuns) runInfo {
	text := run.Text
	props := run.RunProperties

	bold := props != nil && props.Bold != nil && *props.Bold == 1
	italic := props != nil && props.Italic != nil && *props.Italic == 1

	return runInfo{text, bold, italic}
}

// decodeParagraph decodes all the runs inside a paragraph, and combines them
// based on similar formatting to save space and reduce number of redundant tags
func decodeParagraph(paragraph gopptx.DecodeParagraph) string {
	runs := make([]runInfo, len(paragraph.Runs))
	for i, run := range paragraph.Runs {
		runs[i] = decodeRun(run)
	}

	var result strings.Builder

	for i := 0; i < len(runs); {
		current := runs[i]

		// Combine consecutive runs with identical formatting.
		j := i + 1
		for j < len(runs) && runs[j].Bold == current.Bold && runs[j].Italic == current.Italic {
			current.Text += runs[j].Text
			j++
		}

		// Apply formatting tags.
		if current.Bold && current.Italic {
			result.WriteString("<b><i>" + current.Text + "</i></b>")
		} else if current.Bold {
			result.WriteString("<b>" + current.Text + "</b>")
		} else if current.Italic {
			result.WriteString("<i>" + current.Text + "</i>")
		} else {
			result.WriteString(current.Text)
		}

		i = j
	}

	return result.String()
}

func DebugPrintSlides(slides []SlideContent) {
	for _, slide := range slides {
		fmt.Printf("Slide %d (ID: %d)", slide.Index, slide.ID)
		fmt.Println(slide.Content)
	}
}
