package slidesutil

import (
	"context"
	"fmt"

	"google.golang.org/api/slides/v1"
)

// GetPresentation retrieves a presentation by ID.
func (gsc *GoogleSlidesService) GetPresentation(ctx context.Context, presentationID string) (*slides.Presentation, error) {
	pres, err := gsc.PresentationsService.Get(presentationID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("slidesutil: get presentation %s: %w", presentationID, err)
	}
	return pres, nil
}
