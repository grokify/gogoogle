package slidesutil

import (
	"net/http"

	"github.com/grokify/mogo/errors/errorsutil"
	slides "google.golang.org/api/slides/v1"
)

type PresentationCreator struct {
	SlidesClient   *SlidesClient
	Filename       string
	Title          string
	Subtitle       string
	PresentationID string
}

func NewPresentationCreator(googHTTPClient *http.Client) (*PresentationCreator, error) {
	pc := &PresentationCreator{}
	slidesClient, err := NewSlidesClient(googHTTPClient)
	if err != nil {
		return pc, err
	}
	pc.SlidesClient = slidesClient
	return pc, nil
}

func (pc *PresentationCreator) Create(filename, title, subtitle string) (string, error) {
	presentationID, err := pc.SlidesClient.CreatePresentation(
		filename, title, subtitle)
	if err != nil {
		return presentationID, err
	}
	pc.Filename = filename
	pc.Title = title
	pc.Subtitle = subtitle
	pc.PresentationID = presentationID
	return presentationID, nil
}

func (pc *PresentationCreator) CreateEmpty(filename string) (string, error) {
	presentationID, err := CreateEmptyPresentation(
		pc.SlidesClient.GoogleSlidesService.PresentationsService, filename)
	if err != nil {
		return presentationID, err
	}
	pc.PresentationID = presentationID
	return presentationID, nil
}

func (pc *PresentationCreator) CreateSlideMainPoint(slideTitle string) error {
	return CreateSlideMainPoint(
		pc.SlidesClient.GoogleSlidesService.SlidesService,
		pc.SlidesClient.GoogleSlidesService.PresentationsService,
		pc.PresentationID,
		slideTitle)
}

// CreateSlideImageSidebarRight creates a slide for the current
// presentation. `imageID` is optional and will be auto-generated
// if not provided.
func (pc *PresentationCreator) CreateSlideImageSidebarRight(slideTitle, imageID, imageURL, sidebarText string) error {
	slideID, err := pc.SlidesClient.CreateSlideTitleAndBody(
		pc.PresentationID, slideTitle)
	if err != nil {
		return errorsutil.Wrap(err, "PresentationCreator.CreateSlideImageSidebarRight")
	}

	requests, err := CreateSlideImageRequestsSidebarRight(slideID, imageID, imageURL, sidebarText)
	if err != nil {
		return errorsutil.Wrap(err, "PresentationCreator.CreateSlideImageSidebarRight")
	}
	breq := &slides.BatchUpdatePresentationRequest{
		Requests: requests,
	}

	_, err = pc.SlidesClient.BatchUpdate(pc.PresentationID, breq).Do()
	if err != nil {
		return errorsutil.Wrap(err, "PresentationCreator.CreateSlideImageSidebarRight")
	}
	return nil
}
