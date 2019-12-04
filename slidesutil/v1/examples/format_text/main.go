// Formatting text with the Google Slides API
// Video: https://www.youtube.com/watch?v=_O2aUCJyCoQ
// This intentially does not use any of the helpers
// in https://github.com/grokify/googleutil/slidesutil/v1/
// See the `examples/format_markdown` for how to use
// the slidesutil request helpers.
package main

import (
	"github.com/grokify/gotilla/fmt/fmtutil"
	log "github.com/sirupsen/logrus"

	"google.golang.org/api/slides/v1"

	slidesutilexamples "github.com/grokify/googleutil/slidesutil/v1/examples"
)

func main() {
	gss, err := slidesutilexamples.Setup()
	if err != nil {
		log.Fatal(err)
	}

	srv := gss.SlidesSerivce
	psv := gss.PresentationsService

	pres := &slides.Presentation{Title: "Slides text formatting DEMO"}
	rsp1, err := psv.Create(pres).Do()
	if err != nil {
		log.Fatal(err)
	}

	deckID := rsp1.PresentationId
	titleSlide := rsp1.Slides[0]
	titleID := titleSlide.PageElements[0].ObjectId
	subtitleID := titleSlide.PageElements[1].ObjectId
	log.Infof("PresentationID: %v\nTitleID: %v\nSubtitleID: %v\n", deckID, titleID, subtitleID)

	log.Info(`== Create "main point" layout slide & add titles `)
	reqs := []*slides.Request{
		{
			CreateSlide: &slides.CreateSlideRequest{
				//ObjectId: newPageId,
				SlideLayoutReference: &slides.LayoutReference{
					PredefinedLayout: "MAIN_POINT",
				},
			},
		},
		{
			InsertText: &slides.InsertTextRequest{
				ObjectId: titleID,
				Text:     "Formatting text",
			},
		},
		{
			InsertText: &slides.InsertTextRequest{
				ObjectId: subtitleID,
				Text:     "via the Google Slides API",
			},
		},
	}

	rsp2, err := psv.BatchUpdate(
		deckID,
		&slides.BatchUpdatePresentationRequest{Requests: reqs}).Do()
	if err != nil {
		log.Fatal(err)
	}

	fmtutil.PrintJSON(rsp2)

	slideID := rsp2.Replies[0].CreateSlide.ObjectId
	log.Infof("Created SlideID: %v\n", slideID)

	log.Info(`== Fetch "main point" slide title (textbox) ID`)
	presentation, err := srv.Presentations.Get(deckID).Do()
	fmtutil.PrintJSON(presentation)
	if err != nil {
		log.Fatal(err)
	}
	newSlide := presentation.Slides[len(presentation.Slides)-1]
	textboxID := newSlide.PageElements[0].ObjectId

	log.Info("== Insert text & perform various formatting operations")

	reqs = []*slides.Request{
		{
			InsertText: &slides.InsertTextRequest{
				ObjectId: textboxID,
				Text:     "\nBold 1\nItal 2\n\tfoo\n\tbar\n\tbaz\n\t\tquz\nMono 3",
			},
		},
		{
			UpdateTextStyle: &slides.UpdateTextStyleRequest{
				ObjectId: textboxID,
				Style: &slides.TextStyle{
					FontSize: &slides.Dimension{
						Magnitude: float64(32),
						Unit:      "PT",
					},
				},
				Fields: "fontSize",
			},
		},
		{ // change word 1 in para 1 ("Bold") to bold
			UpdateTextStyle: &slides.UpdateTextStyleRequest{
				ObjectId: textboxID,
				Style: &slides.TextStyle{
					Bold: true,
				},
				TextRange: &slides.Range{
					Type:       "FIXED_RANGE",
					StartIndex: int64(1),
					EndIndex:   int64(5),
				},
				Fields: "bold",
			},
		},
		{ // change word 1 in para 2 ("Ital") to italics
			UpdateTextStyle: &slides.UpdateTextStyleRequest{
				ObjectId: textboxID,
				Style: &slides.TextStyle{
					Italic: true,
				},
				TextRange: &slides.Range{
					Type:       "FIXED_RANGE",
					StartIndex: int64(8),
					EndIndex:   int64(12),
				},
				Fields: "italic",
			},
		},
		{ // change word 1 in para 6 ("Mono") to Courier New
			UpdateTextStyle: &slides.UpdateTextStyleRequest{
				ObjectId: textboxID,
				Style: &slides.TextStyle{
					FontFamily: "Courier New",
				},
				TextRange: &slides.Range{
					Type:       "FIXED_RANGE",
					StartIndex: int64(36),
					EndIndex:   int64(40),
				},
				Fields: "fontFamily",
			},
		},
		{ // bulletize everything
			CreateParagraphBullets: &slides.CreateParagraphBulletsRequest{
				ObjectId: textboxID,
				TextRange: &slides.Range{
					Type:       "FIXED_RANGE",
					StartIndex: int64(1),
					EndIndex:   int64(42),
				},
			},
		},
	}

	_, err = psv.BatchUpdate(
		deckID,
		&slides.BatchUpdatePresentationRequest{Requests: reqs}).Do()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("DONE")
}
