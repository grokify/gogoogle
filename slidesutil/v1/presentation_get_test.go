package slidesutil

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"google.golang.org/api/option"
	"google.golang.org/api/slides/v1"
)

func TestGetPresentation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path != "/v1/presentations/p1" {
			w.WriteHeader(http.StatusNotFound)
		}
		body := `{"presentationId":"p1","title":"Deck"}`
		if r.URL.Path != "/v1/presentations/p1" {
			body = `{"error":{"code":404,"message":"gone"}}`
		}
		if _, err := w.Write([]byte(body)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	t.Cleanup(srv.Close)
	svc, err := slides.NewService(context.Background(), option.WithHTTPClient(srv.Client()), option.WithEndpoint(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	gsc := &GoogleSlidesService{SlidesService: svc, PresentationsService: slides.NewPresentationsService(svc)}
	pres, err := gsc.GetPresentation(context.Background(), "p1")
	if err != nil || pres.Title != "Deck" {
		t.Fatalf("got %+v, %v", pres, err)
	}
	if _, err := gsc.GetPresentation(context.Background(), "missing"); err == nil || !strings.Contains(err.Error(), "missing") {
		t.Errorf("err = %v, want one naming the presentation", err)
	}
}
