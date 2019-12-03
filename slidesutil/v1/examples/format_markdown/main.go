// Formatting text with the Google Slides API
// Video: https://www.youtube.com/watch?v=_O2aUCJyCoQ
package main

import (
	su "github.com/grokify/googleutil/slidesutil/v1"
	"github.com/grokify/gotilla/fmt/fmtutil"
	log "github.com/sirupsen/logrus"

	"google.golang.org/api/slides/v1"

	slidesutilexamples "github.com/grokify/googleutil/slidesutil/v1/examples"
)

//"github.com/google/google-api-go-client/slides/v1"

const Markdown = "Foo\n* [**Foo**](https://example.com/foo)\n* [**Bar**](http://example.com/bar)\nBar\n* **Foo**\n* **Bar**\n    * Baz"

func main() {
	gss, err := slidesutilexamples.Setup()
	if err != nil {
		log.Fatal(err)
	}

	srv := gss.SlidesSerivce
	psv := gss.PresentationsService

	pres := &slides.Presentation{Title: "Slides markdown formatting DEMO"}
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
					PredefinedLayout: su.LayoutTitleAndBody,
				},
			},
		},
		{
			InsertText: &slides.InsertTextRequest{
				ObjectId: titleID,
				Text:     "Formatting Markdown",
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
	fmtutil.PrintJSON(presentation.Slides)
	//panic("Z")
	newSlideTitleID := newSlide.PageElements[0].ObjectId
	textboxID := newSlide.PageElements[1].ObjectId

	log.Info("== Insert text & perform various formatting operations")

	reqs = []*slides.Request{
		{
			InsertText: &slides.InsertTextRequest{
				ObjectId: newSlideTitleID,
				Text:     "Test Slide",
			},
		},
	}

	cm := su.NewCommonMarkData(Markdown)
	cm.Inflate()
	fmtutil.PrintJSON(cm.Lines())

	reqs = su.CommonMarkDataToRequests(textboxID, cm)
	reqs = append(reqs, &slides.Request{
		InsertText: &slides.InsertTextRequest{
			ObjectId: newSlideTitleID,
			Text:     "Test Slide",
		},
	})
	fmtutil.PrintJSON(reqs)

	_, err = psv.BatchUpdate(
		deckID,
		&slides.BatchUpdatePresentationRequest{Requests: reqs}).Do()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("DONE")
}

/*
func main() {
	cm := su.NewCommonMarkData(Markdown)
	cm.Inflate()
	fmtutil.PrintJSON(cm.Lines())

	reqs := su.CommonMarkDataToRequests("abc", cm)
	fmtutil.PrintJSON(reqs)

	fmt.Println("DONE")
}
*/
