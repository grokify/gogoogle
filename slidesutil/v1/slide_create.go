package slidesutil

import (
	"github.com/grokify/gotilla/fmt/fmtutil"
	"google.golang.org/api/slides/v1"
)

// CreateSlideTitleOnly creates a slide using Markdown
// given a PresentationID, title, and markdown body.
func CreateSlideTitleOnly(srv *slides.Service, psv *slides.PresentationsService, presentationID, titleText string) (string, error) {
	reqs1 := []*slides.Request{CreateSlideRequestLayout(LayoutTitleOnly)}

	resp1, err := psv.BatchUpdate(
		presentationID,
		&slides.BatchUpdatePresentationRequest{Requests: reqs1}).Do()
	if err != nil {
		return "", err
	}

	slideID := resp1.Replies[0].CreateSlide.ObjectId
	//log.Infof("CREATED SLIDE [%v]\n", slideID)
	//log.Info(`== Fetch "main point" slide title (textbox) ID`)
	presentation, err := srv.Presentations.Get(presentationID).Do()
	fmtutil.PrintJSON(presentation)
	if err != nil {
		return slideID, err
	}
	newSlide := presentation.Slides[len(presentation.Slides)-1]
	fmtutil.PrintJSON(presentation.Slides)

	newSlideTitleID := newSlide.PageElements[0].ObjectId

	_, err = psv.BatchUpdate(
		presentationID,
		&slides.BatchUpdatePresentationRequest{
			Requests: []*slides.Request{InsertTextRequest(
				newSlideTitleID, titleText)},
		}).Do()
	return slideID, err
}
