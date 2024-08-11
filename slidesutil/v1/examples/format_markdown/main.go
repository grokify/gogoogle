// Formatting text with the Google Slides API
// Video: https://www.youtube.com/watch?v=_O2aUCJyCoQ
package main

import (
	"context"
	"log"

	"github.com/grokify/googleutil/auth"
	"github.com/grokify/googleutil/slidesutil/v1"
)

const Markdown = "Foo\n* [**Foo**](https://example.com/foo)\n* [**Bar**](http://example.com/bar)\nBar\n* **Foo**\n* **Bar**\n    * Baz"

func main() {
	googHttpClient, err := auth.Setup(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	slidesClient, err := slidesutil.NewSlidesClient(googHttpClient)
	if err != nil {
		log.Fatal(err)
	}

	presentationID, err := slidesClient.CreatePresentation(
		"Slides markdown formatting DEMO",
		"Formatting Markdown",
		"via the Google Slides API")
	if err != nil {
		log.Fatal(err)
	}

	err = slidesClient.CreateSlideMarkdown(
		presentationID, "Markdown Test Slide", Markdown, false)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DONE")
}
