package slidesutil

import (
	slides "google.golang.org/api/slides/v1"
)

// CreateSlideMarkdown creates a slide using Markdown
// given a PresentationID, title, and markdown body.
func CreateSlideMainPoint(srv *slides.Service, psv *slides.PresentationsService, presentationID, titleText string) error {
	reqs1 := []*slides.Request{CreateSlideRequestLayout(LayoutMainPoint)}

	_, err := psv.BatchUpdate(
		presentationID,
		&slides.BatchUpdatePresentationRequest{Requests: reqs1}).Do()
	if err != nil {
		return err
	}

	presentation, err := srv.Presentations.Get(presentationID).Do()
	if err != nil {
		return err
	}
	newSlide := presentation.Slides[len(presentation.Slides)-1]

	newSlideTitleID := newSlide.PageElements[0].ObjectId

	req := InsertTextRequest(newSlideTitleID, titleText)

	_, err = psv.BatchUpdate(
		presentationID,
		&slides.BatchUpdatePresentationRequest{Requests: []*slides.Request{req}}).Do()
	return err
}
