package upload

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kenny-not-dead/gopptx"
)

type PresentationUpload struct {
	PresentationName	string;
	Filename					string;
	TmpFilePath				string;
}

type runInfo struct {
	Text		string;
	Bold		bool;
	Italic	bool;
}

type slideContent struct {
	Index			int;
	ID				int;
	Content		string;
}

// DecodePPTX goes throgh a powerpoint file and converts the slide content
// into plain text with appropiate HTML tags, include <p> to denote a new line,
// <i> to denote italics, and <b> to denote boldness. 
func DecodePPTX(uploadData PresentationUpload) {
  defer os.Remove(uploadData.TmpFilePath)
	
	f, err := gopptx.OpenFile(uploadData.TmpFilePath)
	if err != nil {
    log.Fatal("Failed to open temp:",err)
		return;
  }
  defer f.Close()

	// read slide data
	slides := make([]slideContent, len(f.Presentation.Slides.Slide))

	for i, slide := range f.Presentation.Slides.Slide {
    shapes, err := f.GetShapes(slide.SlideID)
    if err != nil {
      log.Printf("Error reading shapes from slide %d: %v", slide.SlideID, err)
      continue
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
	debugPrintSlides(slides)
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

func debugPrintSlides(slides []slideContent) {
	for _, slide := range slides {
		fmt.Printf("Slide %d (ID: %d)", slide.Index, slide.ID)
		fmt.Println(slide.Content)
	}
}