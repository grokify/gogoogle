package slidesutil

import (
	"net/http"

	"github.com/pkg/errors"
	"google.golang.org/api/slides/v1"
)

type PresentationCreator struct {
	SlidesClient   *SlidesClient
	Filename       string
	Title          string
	Subtitle       string
	PresentationID string
}

func NewPresentationCreator(googHttpClient *http.Client) (PresentationCreator, error) {
	pc := PresentationCreator{}
	slidesClient, err := NewSlidesClient(googHttpClient)
	if err != nil {
		return pc, err
	}
	pc.SlidesClient = slidesClient
	return pc, nil
}

func (pc *PresentationCreator) Create(filename, title, subtitle string) error {
	presentationID, err := pc.SlidesClient.CreatePresentation(
		filename, title, subtitle)
	if err != nil {
		return err
	}
	pc.Filename = filename
	pc.Title = title
	pc.Subtitle = subtitle
	pc.PresentationID = presentationID
	return nil
}

func (pc *PresentationCreator) CreateSlideImageSidebarRight(slideTitle, imageID, imageURL, sidebarText string) error {
	slideID, err := pc.SlidesClient.CreateSlideTitleAndBody(
		pc.PresentationID, slideTitle)
	if err != nil {
		return errors.Wrap(err, "PresentationCreator.CreateSlideImageSidebarRight")
	}

	requests, err := CreateSlideImageRequestsSidebarRight(slideID, imageID, imageURL, sidebarText)
	if err != nil {
		return errors.Wrap(err, "PresentationCreator.CreateSlideImageSidebarRight")
	}
	breq := &slides.BatchUpdatePresentationRequest{
		Requests: requests,
	}

	_, err = pc.SlidesClient.BatchUpdate(pc.PresentationID, breq).Do()
	if err != nil {
		return errors.Wrap(err, "PresentationCreator.CreateSlideImageSidebarRight")
	}
	return nil
}
