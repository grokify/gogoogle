package slidesutil

import (
	"google.golang.org/api/slides/v1"
)

// CreatePresentation Creates a new presentation with
// filename, title and subtitle.
func CreatePresentation(srv *slides.Service, psv *slides.PresentationsService,
	filename, titleText, subtitleText string) (string, error) {
	pres := &slides.Presentation{Title: filename}
	rsp1, err := psv.Create(pres).Do()
	if err != nil {
		return "", err
	}

	presentationID := rsp1.PresentationId
	titleSlide := rsp1.Slides[0]
	titleID := titleSlide.PageElements[0].ObjectId
	subtitleID := titleSlide.PageElements[1].ObjectId

	reqs := []*slides.Request{}
	if len(titleText) > 0 {
		reqs = append(reqs, InsertTextRequest(titleID, titleText))
	}
	if len(subtitleText) > 0 {
		reqs = append(reqs, InsertTextRequest(subtitleID, subtitleText))
	}
	if len(reqs) > 0 {
		_, err := psv.BatchUpdate(
			presentationID,
			&slides.BatchUpdatePresentationRequest{Requests: reqs}).Do()
		if err != nil {
			return presentationID, err
		}
	}
	return presentationID, nil
}
