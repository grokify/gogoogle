package slidesutil

import (
	"github.com/grokify/gotilla/fmt/fmtutil"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/slides/v1"
)

// CreateSlideImage creates a slide using a main image
// as the body given a PresentationID, title, and imageURL.
func CreateSlideImage(srv *slides.Service, psv *slides.PresentationsService, presentationID, titleText, bodyMarkdown string, underlineLinks bool) error {
	reqs1 := []*slides.Request{CreateSlideRequestLayout(LayoutTitleOnly)}

	resp1, err := psv.BatchUpdate(
		presentationID,
		&slides.BatchUpdatePresentationRequest{Requests: reqs1}).Do()
	if err != nil {
		return err
	}

	if 1 == 0 {
		slideID := resp1.Replies[0].CreateSlide.ObjectId
		log.Infof("CREATED SLIDE [%v]\n", slideID)
	}
	//log.Info(`== Fetch "main point" slide title (textbox) ID`)
	presentation, err := srv.Presentations.Get(presentationID).Do()
	fmtutil.PrintJSON(presentation)
	if err != nil {
		return err
	}
	newSlide := presentation.Slides[len(presentation.Slides)-1]
	fmtutil.PrintJSON(presentation.Slides)

	newSlideTitleID := newSlide.PageElements[0].ObjectId
	//newSlideBodyTextboxID := newSlide.PageElements[1].ObjectId

	cm := NewCommonMarkData(bodyMarkdown)
	cm.Inflate()
	//fmtutil.PrintJSON(cm.Lines())

	reqs2 := []*slides.Request{InsertTextRequest(newSlideTitleID, titleText)}

	_, err = psv.BatchUpdate(
		presentationID,
		&slides.BatchUpdatePresentationRequest{Requests: reqs2}).Do()
	return err
}
