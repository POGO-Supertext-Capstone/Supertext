package upload

import (
	"fmt"
	"io"
	"strings"

	"github.com/kenny-not-dead/gopptx"
)

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
func DecodePPTX(r io.Reader) ([]SlideContent, error) {
	f, err := gopptx.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to open PPTX: %w", err)
	}
	defer f.Close()

	// read slide data
	slides := make([]SlideContent, len(f.Presentation.Slides.Slide))

	for i, slide := range f.Presentation.Slides.Slide {
		shapes, err := f.GetShapes(slide.SlideID)
		if err != nil {
			return nil, fmt.Errorf("Error reading shapes from slide %d: %v", slide.SlideID, err)
		}

		slides[i].ID = slide.SlideID
		slides[i].Index = i

		content := ""
		for _, shape := range shapes {
			for _, paragraph := range shape.TextBody.Paragraph {
				if len(paragraph.Breaks) > 0 {
					content += decodeLineBreakParagraph(paragraph)
				} else {
					content += decodeParagraph(paragraph)
				}
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

	return "<p>" + encodeRunTags(runs) + "</p>"
}

// decodeParagraph decodes all the runs inside a paragraph, seperated by line breaks
// and combines them based on similar formatting to save space and reduce
// number of redundant tags
func decodeLineBreakParagraph(paragraph gopptx.DecodeParagraph) string {
	runs := make([]runInfo, len(paragraph.Runs))
	run_indices := make([]int, len(paragraph.Runs))
	for i, run := range paragraph.Runs {
		runs[i] = decodeRun(run)
		run_indices[i] = strings.Index(paragraph.InnerXML, runs[i].Text)
	}

	break_indices := []int{}
	start := 0
	for {
		idx := strings.Index(paragraph.InnerXML[start:], "<a:br")
		if idx == -1 {
			break
		}

		start += idx
		break_indices = append(break_indices, start)
		start += len("<a:br")
	}

	// Split the runs into groups based on the breaks.
	result := []gopptx.DecodeParagraph{}
	current := gopptx.DecodeParagraph{ParagraphProperties: paragraph.ParagraphProperties}

	breakIdx := 0
	for i, run := range paragraph.Runs {
		// If a break occurs before this run, start a new paragraph.
		for breakIdx < len(break_indices) && break_indices[breakIdx] < run_indices[i] {
			result = append(result, current)
			current = gopptx.DecodeParagraph{ParagraphProperties: paragraph.ParagraphProperties}
			breakIdx++
		}

		current.Runs = append(current.Runs, run)
	}
	result = append(result, current)

	result_string := ""
	for _, par := range result {
		result_string += decodeParagraph(par)
	}

	return result_string

}

// encodeRunTags encodes all the runs in a paragraph with their
// limited HTML tags
func encodeRunTags(runs []runInfo) string {
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
