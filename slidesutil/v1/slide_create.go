package slidesutil

import (
	slides "google.golang.org/api/slides/v1"
)

// CreateSlideTitleOnly creates a slide using Markdown
// given a PresentationID and title.
func CreateSlideTitleOnly(srv *slides.Service, psv *slides.PresentationsService, presentationID, titleText string) (string, error) {
	reqs1 := []*slides.Request{CreateSlideRequestLayout(LayoutTitleOnly)}

	resp1, err := psv.BatchUpdate(
		presentationID,
		&slides.BatchUpdatePresentationRequest{Requests: reqs1}).Do()
	if err != nil {
		return "", err
	}

	slideID := resp1.Replies[0].CreateSlide.ObjectId
	presentation, err := srv.Presentations.Get(presentationID).Do()
	if err != nil {
		return slideID, err
	}

	newSlide := presentation.Slides[len(presentation.Slides)-1]
	newSlideTitleID := newSlide.PageElements[0].ObjectId

	_, err = psv.BatchUpdate(
		presentationID,
		&slides.BatchUpdatePresentationRequest{
			Requests: []*slides.Request{InsertTextRequest(
				newSlideTitleID, titleText)},
		}).Do()
	return slideID, err
}

// CreateSlideTitleAndBody creates a slide using Markdown
// given a PresentationID, title, and markdown body.
func CreateSlideTitleAndBody(srv *slides.Service, psv *slides.PresentationsService, presentationID, titleText string) (string, error) {
	reqs1 := []*slides.Request{CreateSlideRequestLayout(LayoutTitleAndBody)}

	resp1, err := psv.BatchUpdate(
		presentationID,
		&slides.BatchUpdatePresentationRequest{Requests: reqs1}).Do()
	if err != nil {
		return "", err
	}

	slideID := resp1.Replies[0].CreateSlide.ObjectId
	presentation, err := srv.Presentations.Get(presentationID).Do()
	if err != nil {
		return slideID, err
	}

	newSlide := presentation.Slides[len(presentation.Slides)-1]
	newSlideTitleID := newSlide.PageElements[0].ObjectId

	_, err = psv.BatchUpdate(
		presentationID,
		&slides.BatchUpdatePresentationRequest{
			Requests: []*slides.Request{InsertTextRequest(
				newSlideTitleID, titleText)},
		}).Do()

	return slideID, err
}
