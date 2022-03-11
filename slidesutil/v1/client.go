package slidesutil

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grokify/mogo/errors/errorsutil"
	"google.golang.org/api/option"
	"google.golang.org/api/slides/v1"
)

type GoogleSlidesService struct {
	httpClient           *http.Client
	SlidesService        *slides.Service
	PresentationsService *slides.PresentationsService
}

func NewGoogleSlidesService(httpClient *http.Client) (*GoogleSlidesService, error) {
	gsc := GoogleSlidesService{}
	err := gsc.SetHTTPClient(httpClient)
	return &gsc, err
}

func (gsc *GoogleSlidesService) SetHTTPClient(httpClient *http.Client) error {
	if httpClient == nil {
		return fmt.Errorf("httpClient parameter canot be nil")
	}
	gsc.httpClient = httpClient
	service, err := slides.NewService(context.Background(), option.WithHTTPClient(httpClient))
	if err != nil {
		return errorsutil.Wrap(err, "unable to create slides.Service")
	}
	gsc.SlidesService = service
	gsc.PresentationsService = slides.NewPresentationsService(service)
	return nil
}

type SlidesClient struct {
	GoogleSlidesService *GoogleSlidesService
}

func NewSlidesClient(googHTTPClient *http.Client) (*SlidesClient, error) {
	sc := &SlidesClient{}
	gss, err := NewGoogleSlidesService(googHTTPClient)
	if err != nil {
		return nil, err
	}
	sc.GoogleSlidesService = gss
	return sc, nil
}

func (sc *SlidesClient) CreatePresentation(
	filename, titleText, subtitleText string) (string, error) {
	return CreatePresentation(
		sc.GoogleSlidesService.SlidesService,
		sc.GoogleSlidesService.PresentationsService,
		filename, titleText, subtitleText)
}

func (sc *SlidesClient) CreateEmptyPresentation(name string) (string, error) {
	return CreateEmptyPresentation(
		sc.GoogleSlidesService.PresentationsService, name)
}

// BatchUpdate is a convenience function to make calling `BatchUpdate`
// less verbose.
func (sc *SlidesClient) BatchUpdate(presentationID string, batchupdatepresentationrequest *slides.BatchUpdatePresentationRequest) *slides.PresentationsBatchUpdateCall {
	return sc.GoogleSlidesService.PresentationsService.BatchUpdate(
		presentationID, batchupdatepresentationrequest)
}

// CreateSlideTitleAndBody is a convenience function.
func (sc *SlidesClient) CreateSlideTitleAndBody(presentationID string, filename string) (string, error) {
	return CreateSlideTitleAndBody(
		sc.GoogleSlidesService.SlidesService,
		sc.GoogleSlidesService.PresentationsService,
		presentationID, filename)
}

func (sc *SlidesClient) CreateSlideMarkdown(presentationID, titleText, bodyMarkdown string, underlineLinks bool) error {
	return CreateSlideMarkdown(
		sc.GoogleSlidesService.SlidesService,
		sc.GoogleSlidesService.PresentationsService,
		presentationID, titleText, bodyMarkdown, underlineLinks)
}
