package slidesutil

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
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
	service, err := slides.New(gsc.httpClient)
	if err != nil {
		return errors.Wrap(err, "Unable to create slides.Service")
	}
	gsc.SlidesService = service
	gsc.PresentationsService = slides.NewPresentationsService(service)
	return nil
}
