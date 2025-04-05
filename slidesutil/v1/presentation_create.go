package slidesutil

import (
	"fmt"
	"net/http"

	"github.com/grokify/mogo/errors/errorsutil"
	slides "google.golang.org/api/slides/v1"
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

func CreateEmptyPresentation(psv *slides.PresentationsService, filename string) (string, error) {
	pres := &slides.Presentation{Title: filename}
	res, err := psv.Create(pres).Do()
	if err != nil {
		return "", errorsutil.Wrap(err, "CreateRoadmapSlide - psv.Create(pres).Do()")
	}

	fmt.Printf("CREATED Presentation with Id %v\n", res.PresentationId)

	if 1 == 0 {
		for i, slide := range res.Slides {
			fmt.Printf("- Slide #%d id %v contains %d elements.\n", (i + 1),
				slide.ObjectId,
				len(slide.PageElements))
		}
	}

	pageID := res.Slides[0].ObjectId

	requests := []*slides.Request{
		{
			DeleteObject: &slides.DeleteObjectRequest{ObjectId: pageID},
		},
	}
	breq := &slides.BatchUpdatePresentationRequest{
		Requests: requests,
	}
	_, err = psv.BatchUpdate(res.PresentationId, breq).Do() // resu
	if err != nil {
		return "", errorsutil.Wrap(err, "CreateRoadmapSlide - psv.BatchUpdate(res.PresentationId, breq).Do()")
	}
	return res.PresentationId, nil
}

func CreateEmptyPresentationGSS(gss *GoogleSlidesService, filename string) (string, error) {
	return CreateEmptyPresentation(gss.PresentationsService, filename)
}

func CreateEmptyPresentationHTTP(googleClient *http.Client, filename string) (string, error) {
	gss, err := NewGoogleSlidesService(googleClient)
	if err != nil {
		return "", errorsutil.Wrap(err, "CreateRoadmapSlide - slidesutil.NewGoogleSlidesService()")
	}

	return CreateEmptyPresentationGSS(gss, filename)
}
