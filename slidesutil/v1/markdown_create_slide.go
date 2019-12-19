package slidesutil

import (
	"fmt"

	"github.com/grokify/gotilla/fmt/fmtutil"
	"google.golang.org/api/slides/v1"
)

// CreateSlideMarkdown creates a slide using Markdown
// given a PresentationID, title, and markdown body.
func CreateSlideMarkdown(srv *slides.Service, psv *slides.PresentationsService, presentationID, titleText, bodyMarkdown string, underline bool) error {
	reqs1 := []*slides.Request{
		CreateSlideRequestLayout(LayoutTitleAndBody)}

	resp1, err := psv.BatchUpdate(
		presentationID,
		&slides.BatchUpdatePresentationRequest{Requests: reqs1}).Do()
	if err != nil {
		return err
	}

	if 1 == 0 {
		slideID := resp1.Replies[0].CreateSlide.ObjectId
		fmt.Printf("CREATED SLIDE [%v]\n", slideID)
	}
	//log.Infof("Created SlideID: %v\n", slideID)
	//log.Info(`== Fetch "main point" slide title (textbox) ID`)
	presentation, err := srv.Presentations.Get(presentationID).Do()
	fmtutil.PrintJSON(presentation)
	if err != nil {
		return err
	}
	newSlide := presentation.Slides[len(presentation.Slides)-1]
	fmtutil.PrintJSON(presentation.Slides)

	newSlideTitleID := newSlide.PageElements[0].ObjectId
	newSlideBodyTextboxID := newSlide.PageElements[1].ObjectId

	cm := NewCommonMarkData(bodyMarkdown)
	cm.Inflate()
	//fmtutil.PrintJSON(cm.Lines())

	reqs2 := CommonMarkDataToRequests(newSlideBodyTextboxID, cm, underline)
	reqs2 = append(
		reqs2,
		InsertTextRequest(newSlideTitleID, titleText))
	lineCount := cm.LineCount()
	fmt.Printf("LINE_COUNT [%v]\n", lineCount)
	// 22
	if lineCount > 15 {
		reqs2 = append(
			reqs2,
			UpdateTextStyleRequestFontSizePT(newSlideBodyTextboxID, float64(8)),
			UpdateParagraphStyleRequestLineSpacing(newSlideBodyTextboxID, float64(100)),
		)
	}
	//panic("Z")
	//fmtutil.PrintJSON(reqs2)

	_, err = psv.BatchUpdate(
		presentationID,
		&slides.BatchUpdatePresentationRequest{Requests: reqs2}).Do()
	return err
}
