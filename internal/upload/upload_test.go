package upload

import (
	"testing"
	
	"github.com/kenny-not-dead/gopptx"
)

func TestDecodeRun(t *testing.T) {
	bold := 1
	italic := 1

	tests := []struct {
		name	string
		run		gopptx.DecodeRuns
		want	runInfo
	}{
		{
			name:	"empty text",
			run:	gopptx.DecodeRuns{Text: ""},
			want:	runInfo{Text: "", Bold: false, Italic: false},
		},
		{
			name:	"plain text",
			run:	gopptx.DecodeRuns{Text: "Hello"},
			want:	runInfo{Text: "Hello", Bold: false, Italic: false},
		},
		{
			name:	"bold text",
			run:	gopptx.DecodeRuns{
				Text: "Hello", 
				RunProperties: &gopptx.DecodeRunProperties{Bold: &bold},
			},
			want:	runInfo{Text: "Hello", Bold: true, Italic: false},
		},
		{
			name:	"italic text",
			run:	gopptx.DecodeRuns{
				Text: "Hello", 
				RunProperties: &gopptx.DecodeRunProperties{Italic: &italic},
			},
			want:	runInfo{Text: "Hello", Bold: false, Italic: true},
		},
		{
			name:	"bold and italic text",
			run:	gopptx.DecodeRuns{
				Text: "Hello", 
				RunProperties: &gopptx.DecodeRunProperties{Bold: &bold, Italic: &italic},
			},
			want:	runInfo{Text: "Hello", Bold: true, Italic: true},
		},
		{
			name:	"nil properties",
			run:	gopptx.DecodeRuns{
				Text: "Hello", 
				RunProperties: nil,
			},
			want:	runInfo{Text: "Hello", Bold: false, Italic: false},
		},
		{
			name:	"special characters",
			run:	gopptx.DecodeRuns{Text: "áéíóúüñ¿¡"},
			want:	runInfo{Text: "áéíóúüñ¿¡", Bold: false, Italic: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodeRun(tt.run)
			if got != tt.want {
				t.Errorf("decodeRun() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestDecodeParagraph(t *testing.T) {
	true_int := 1

	tests := []struct {
		name string
		paragraph gopptx.DecodeParagraph
		want string
	}{
		{
			name: "empty paragraph",
			paragraph: gopptx.DecodeParagraph{},
			want: "",
		},
		{
			name: "simple paragraph",
    	paragraph: gopptx.DecodeParagraph{
        Runs: []gopptx.DecodeRuns{
        	{ Text: "Hello, world!", },
        },
    	},
    	want: "Hello, world!",
		},
		{
			name: "combined paragraph 1",
    	paragraph: gopptx.DecodeParagraph{
        Runs: []gopptx.DecodeRuns{
        	{ Text: "Hello, world! ", },
					{ Text: "¡Hola Mundo!", RunProperties: &gopptx.DecodeRunProperties{Italic: &true_int}, },
        },
    	},
    	want: "Hello, world! <i>¡Hola Mundo!</i>",
		},
		{
			name: "combined paragraph 2",
    	paragraph: gopptx.DecodeParagraph{
        Runs: []gopptx.DecodeRuns{
        	{ Text: "What is", },
					{ Text: " up", RunProperties: &gopptx.DecodeRunProperties{Italic: &true_int}, },
					{ Text: " my", },
					{ Text: " dudes", RunProperties: &gopptx.DecodeRunProperties{Bold: &true_int}, },
        },
    	},
    	want: "What is<i> up</i> my<b> dudes</b>",
		},
		{
			name: "messy paragraph 1",
    	paragraph: gopptx.DecodeParagraph{
        Runs: []gopptx.DecodeRuns{
					{ Text: "likewise ", RunProperties: &gopptx.DecodeRunProperties{Italic: &true_int}, },
					{ Text: "punctuation ", RunProperties: &gopptx.DecodeRunProperties{Italic: &true_int}, },
					{ Text: "should ", },
					{ Text: "be ", RunProperties: &gopptx.DecodeRunProperties{Italic: &true_int}, },
					{ Text: "combined!", RunProperties: &gopptx.DecodeRunProperties{Italic: &true_int}, },
        },
    	},
    	want: "<i>likewise punctuation </i>should <i>be combined!</i>",
		},
		{
			name: "messy paragraph 2",
    	paragraph: gopptx.DecodeParagraph{
        Runs: []gopptx.DecodeRuns{
					{ Text: "likewise ", RunProperties: &gopptx.DecodeRunProperties{Bold: &true_int}, },
					{ Text: "punctuation ", RunProperties: &gopptx.DecodeRunProperties{Italic: &true_int, Bold: &true_int}, },
					{ Text: "should ", },
					{ Text: "be ", RunProperties: &gopptx.DecodeRunProperties{Italic: &true_int, Bold: &true_int}, },
					{ Text: "combined!", RunProperties: &gopptx.DecodeRunProperties{Italic: &true_int, Bold: &true_int}, },
        },
    	},
    	want: "<b>likewise </b><b><i>punctuation </i></b>should <b><i>be combined!</i></b>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodeParagraph(tt.paragraph)

			if got != tt.want {
				t.Errorf("decodeParagraph() = %q, want %q", got, tt.want)
			}
		})
	}
}