package database

import (
	"testing"
	"uuid"
	"encoding/json"
	"strings"
	"reflect"

	"github.com/POGO-Supertext-Capstone/Supertext/internal/upload"
)

// TestUploadShow tests the upload proces on a dummy show, and ensures
// all fields are saved in the database properly, then deletes the show
func TestUploadShow(t *testing.T) {
	testShow := uuid.New().String()
	slides := []upload.SlideContent{
		{Index: 1, ID: 1, Content: "<p>This is a test presentation</p>"},
		{Index: 2, ID: 2, Content: "<p>for TestUploadShow()</p>"},
	}
	testShowFilename := "test.pptx"

	err := UploadShow(slides, testShow, testShowFilename)
	if err != nil {
		t.Fatalf("UploadPPTX failed: %s", err)
	}

	filename, err := GetShowFilename(testShow)
	if err != nil {
		t.Fatalf("GetShowFilename failed: %s", err)
	}

	if filename != "test.pptx" {
		t.Errorf("expected test.pptx, got %s", filename)
	}

	err = DeleteShow(testShow)
	if err != nil {
		t.Fatalf("DeleteShow failed: %s", err)
	}
}

// TestShowEncode tests the uploading process formated the JSON properly
func TestShowEncode(t *testing.T) {
	tests := []struct {
		name string
		run  []upload.SlideContent
		want string
	}{
		{
			name: "Simple Slideshow",
			run: []upload.SlideContent{
				{Index: 1, ID: 1, Content: "<p>This is a test presentation</p>"},
				{Index: 2, ID: 2, Content: "<p>for TestShowEncode()</p>"},
			},
			want: `[{"index":1,"id":1,"text":"<p>This is a test presentation</p>"},{"index":2,"id":2,"text":"<p>for TestShowEncode()</p>"}]`,
		},
		{
			name: "Empty Slideshow",
			run: []upload.SlideContent{},
			want: `[]`,
		},
		{
			name: "Missing Fields Slideshow",
			run: []upload.SlideContent{
				{ID: 1, Content: "<p>1</p>"},
				{Index: 2, Content: "<p>2</p>"},
				{Index: 3, ID: 3},
			},
			want: `[{"index":0,"id":1,"text":"<p>1</p>"},{"index":2,"id":0,"text":"<p>2</p>"},{"index":3,"id":3,"text":""}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			showName := uuid.New().String()

			// Clean up test bucket
			defer func() {
				if err := DeleteShow(showName); err != nil {
					t.Errorf("failed to delete test show: %s", err)
				}
			}()

			// Upload test slides
			err := UploadShow(tt.run, showName, "dummy.pptx")
			if err != nil {
				t.Fatalf("UploadShow() failed: %s", err)
			}

			// Retrieve raw JSON
			got, err := GetShowJSON(showName)
			if err != nil {
				t.Fatalf("GetShowJSON() failed: %s", err)
			}

			if !json.Valid(got) {
				t.Fatal("GetShowJSON() returned invalid JSON")
			}
			
			// Compare JSON
			if strings.TrimSpace(string(got)) != tt.want {
				t.Errorf("GetShowJSON() = %s, want %s", got, tt.want)
			}
		})
	}
}

// TestShowEncodeDecode verifies that an uploaded show will decode
// into the same slide content as when it was uploaded
func TestShowEncodeDecode(t *testing.T) {
	tests := []struct {
		name string
		run  []upload.SlideContent
	}{
		{
			name: "Simple Slideshow",
			run: []upload.SlideContent{
				{Index: 1, ID: 1, Content: "<p>This is a test presentation</p>"},
				{Index: 2, ID: 2, Content: "<p>for TestShowEncode()</p>"},
			},
		},
		{
			name: "Empty Slideshow",
			run: []upload.SlideContent{},
		},
		{
			name: "Long Slideshow",
			run: []upload.SlideContent{
				{Index: 1, ID: 1, Content: "<p><b>THE NIGHT BEGINS</b></p>"},
				{Index: 2, ID: 2, Content: "<p>The stars are fading</p>"},
				{Index: 3, ID: 3, Content: "<p>And darkness falls</p>"},
				{Index: 4, ID: 4, Content: "<p><i>Listen...</i></p>"},
				{Index: 5, ID: 5, Content: "<p>Can you hear it?</p>"},
				{Index: 6, ID: 6, Content: "<p>A distant voice</p>"},
				{Index: 7, ID: 7, Content: "<p>Calling my name</p>"},
				{Index: 8, ID: 8, Content: "<p>Across the water</p>"},
				{Index: 9, ID: 9, Content: "<p>Beyond the shore</p>"},
				{Index: 10, ID: 10, Content: "<p>I cannot stay</p>"},
				{Index: 11, ID: 11, Content: "<p>I cannot turn back</p>"},
				{Index: 12, ID: 12, Content: "<p><b>TAKE ME AWAY</b></p>"},
				{Index: 13, ID: 13, Content: "<p>Through the silence</p>"},
				{Index: 14, ID: 14, Content: "<p>Through the storm</p>"},
				{Index: 15, ID: 15, Content: "<p>Through the shadows</p>"},
				{Index: 16, ID: 16, Content: "<p>I will go</p>"},
				{Index: 17, ID: 17, Content: "<p>Though I am frightened</p>"},
				{Index: 18, ID: 18, Content: "<p>Though I am alone</p>"},
				{Index: 19, ID: 19, Content: "<p>There is a promise</p>"},
				{Index: 20, ID: 20, Content: "<p>There is a home</p>"},
				{Index: 21, ID: 21, Content: "<p><i>Remember...</i></p>"},
				{Index: 22, ID: 22, Content: "<p>The words you whispered</p>"},
				{Index: 23, ID: 23, Content: "<p>The vows we made</p>"},
				{Index: 24, ID: 24, Content: "<p>Before the morning</p>"},
				{Index: 25, ID: 25, Content: "<p>Before the rain</p>"},
				{Index: 26, ID: 26, Content: "<p><b>NEVER FORGET</b></p>"},
				{Index: 27, ID: 27, Content: "<p>The road is broken</p>"},
				{Index: 28, ID: 28, Content: "<p>The wind is cold</p>"},
				{Index: 29, ID: 29, Content: "<p>But still I wander</p>"},
				{Index: 30, ID: 30, Content: "<p>Still I hope</p>"},
				{Index: 31, ID: 31, Content: "<p>And then I see</p>"},
				{Index: 32, ID: 32, Content: "<p>A light before me</p>"},
				{Index: 33, ID: 33, Content: "<p>A hand extended</p>"},
				{Index: 34, ID: 34, Content: "<p>A heart returned</p>"},
				{Index: 35, ID: 35, Content: "<p><b>AT LAST, THE DAWN</b></p>"},
				{Index: 36, ID: 36, Content: "<p>The night is over</p>"},
				{Index: 37, ID: 37, Content: "<p>The stars have gone</p>"},
				{Index: 38, ID: 38, Content: "<p>But we remain</p>"},
				{Index: 39, ID: 39, Content: "<p>Together, still</p>"},
				{Index: 40, ID: 40, Content: "<p><b>THE END</b></p>"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			showName := uuid.New().String()

			// Clean up test bucket
			defer func() {
				if err := DeleteShow(showName); err != nil {
					t.Errorf("failed to delete test show: %s", err)
				}
			}()

			// Upload test slides
			err := UploadShow(tt.run, showName, "dummy.pptx")
			if err != nil {
				t.Fatalf("UploadShow() failed: %s", err)
			}

			// Retrieve slides
			got, err := GetShowSlides(showName)
			if err != nil {
				t.Fatalf("GetShowSlides() failed: %s", err)
			}
			
			// ensure slides going in match slides going out
			if !reflect.DeepEqual(got, tt.run) {
				t.Errorf("GetShowSlides() = %v, want %v", got, tt.run)
			}
		})
	}
}

