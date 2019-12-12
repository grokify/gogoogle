// Formatting text with the Google Slides API
// Video: https://www.youtube.com/watch?v=_O2aUCJyCoQ
package main

import (
	"github.com/grokify/googleutil/slidesutil/v1"
	log "github.com/sirupsen/logrus"

	slidesutilexamples "github.com/grokify/googleutil/slidesutil/v1/examples"
)

const Markdown = "Foo\n* [**Foo**](https://example.com/foo)\n* [**Bar**](http://example.com/bar)\nBar\n* **Foo**\n* **Bar**\n    * Baz"

func main() {
	gss, err := slidesutilexamples.Setup()
	if err != nil {
		log.Fatal(err)
	}
	srv := gss.SlidesSerivce
	psv := gss.PresentationsService

	presentationID, err := slidesutil.CreatePresentation(srv, psv,
		"Slides markdown formatting DEMO",
		"Formatting Markdown",
		"via the Google Slides API")
	if err != nil {
		log.Fatal(err)
	}

	err = slidesutil.CreateSlideMarkdown(srv, psv,
		presentationID, "Markdown Test Slide", Markdown)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("DONE")
}
