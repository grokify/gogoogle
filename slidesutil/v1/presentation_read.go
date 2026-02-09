package slidesutil

import (
	"fmt"
	"strings"

	slides "google.golang.org/api/slides/v1"
)

// FindSlide locates a slide by either index or object ID.
// Returns the slide and its zero-based index, or an error if not found.
func FindSlide(pres *slides.Presentation, index *int, objectID string) (*slides.Page, int, error) {
	if index != nil {
		if *index < 0 || *index >= len(pres.Slides) {
			return nil, 0, fmt.Errorf("slide index %d out of range (0-%d)", *index, len(pres.Slides)-1)
		}
		return pres.Slides[*index], *index, nil
	}

	if objectID != "" {
		for i, slide := range pres.Slides {
			if slide.ObjectId == objectID {
				return slide, i, nil
			}
		}
		return nil, 0, fmt.Errorf("slide with object ID %q not found", objectID)
	}

	return nil, 0, fmt.Errorf("either slide index or object ID must be provided")
}

// ExtractSlideTitle extracts the title text from a slide.
// It looks for shapes with TITLE or CENTERED_TITLE placeholder types.
func ExtractSlideTitle(slide *slides.Page) string {
	for _, elem := range slide.PageElements {
		if elem.Shape == nil || elem.Shape.Placeholder == nil {
			continue
		}
		placeholderType := elem.Shape.Placeholder.Type
		if placeholderType == "TITLE" || placeholderType == "CENTERED_TITLE" {
			return ExtractShapeText(elem.Shape)
		}
	}
	return ""
}

// ExtractTextContent extracts all text content from a page's elements.
// Returns a slice of non-empty text strings.
func ExtractTextContent(page *slides.Page) []string {
	var texts []string
	for _, elem := range page.PageElements {
		text := ExtractElementText(elem)
		if text != "" {
			texts = append(texts, text)
		}
	}
	return texts
}

// ExtractElementText extracts text from a single page element.
func ExtractElementText(elem *slides.PageElement) string {
	if elem.Shape != nil {
		return ExtractShapeText(elem.Shape)
	}
	if elem.Table != nil {
		return ExtractTableText(elem.Table)
	}
	if elem.ElementGroup != nil {
		return ExtractGroupText(elem.ElementGroup)
	}
	return ""
}

// ExtractShapeText extracts text from a shape's text content.
func ExtractShapeText(shape *slides.Shape) string {
	if shape.Text == nil {
		return ""
	}
	return ExtractTextFromTextContent(shape.Text)
}

// ExtractTableText extracts text from all cells in a table.
func ExtractTableText(table *slides.Table) string {
	var parts []string
	for _, row := range table.TableRows {
		for _, cell := range row.TableCells {
			if cell.Text != nil {
				text := ExtractTextFromTextContent(cell.Text)
				if text != "" {
					parts = append(parts, text)
				}
			}
		}
	}
	return strings.Join(parts, " | ")
}

// ExtractGroupText extracts text from all elements in a group.
func ExtractGroupText(group *slides.Group) string {
	var parts []string
	for _, elem := range group.Children {
		text := ExtractElementText(elem)
		if text != "" {
			parts = append(parts, text)
		}
	}
	return strings.Join(parts, " ")
}

// ExtractTextFromTextContent extracts plain text from TextContent.
func ExtractTextFromTextContent(tc *slides.TextContent) string {
	var sb strings.Builder
	for _, elem := range tc.TextElements {
		if elem.TextRun != nil {
			sb.WriteString(elem.TextRun.Content)
		}
	}
	return strings.TrimSpace(sb.String())
}

// ExtractNotesText extracts speaker notes text from a slide.
// Speaker notes are stored in the NotesPage within SlideProperties.
func ExtractNotesText(slide *slides.Page) string {
	if slide.SlideProperties == nil || slide.SlideProperties.NotesPage == nil {
		return ""
	}

	notesPage := slide.SlideProperties.NotesPage
	for _, elem := range notesPage.PageElements {
		if elem.Shape == nil || elem.Shape.Placeholder == nil {
			continue
		}
		if elem.Shape.Placeholder.Type == "BODY" {
			return ExtractShapeText(elem.Shape)
		}
	}
	return ""
}

// GetElementType returns a human-readable type name for a page element.
func GetElementType(elem *slides.PageElement) string {
	switch {
	case elem.Shape != nil:
		if elem.Shape.Placeholder != nil && elem.Shape.Placeholder.Type != "" && elem.Shape.Placeholder.Type != "NONE" {
			return "placeholder:" + strings.ToLower(elem.Shape.Placeholder.Type)
		}
		if elem.Shape.ShapeType != "" && elem.Shape.ShapeType != "TYPE_UNSPECIFIED" {
			return "shape:" + strings.ToLower(elem.Shape.ShapeType)
		}
		return "shape"
	case elem.Image != nil:
		return "image"
	case elem.Video != nil:
		return "video"
	case elem.Table != nil:
		return "table"
	case elem.Line != nil:
		return "line"
	case elem.SheetsChart != nil:
		return "chart"
	case elem.WordArt != nil:
		return "wordart"
	case elem.SpeakerSpotlight != nil:
		return "speaker_spotlight"
	case elem.ElementGroup != nil:
		return "group"
	default:
		return "unknown"
	}
}

// GetElementDescription returns a description for a page element.
func GetElementDescription(elem *slides.PageElement) string {
	if elem.Description != "" {
		return elem.Description
	}
	if elem.Title != "" {
		return elem.Title
	}

	// For shapes with text, return a preview
	if elem.Shape != nil {
		text := ExtractShapeText(elem.Shape)
		if text != "" {
			if len(text) > 50 {
				return text[:50] + "..."
			}
			return text
		}
	}

	return ""
}

// GetImageURL returns the content URL for an image element, if any.
func GetImageURL(elem *slides.PageElement) string {
	if elem.Image != nil && elem.Image.ContentUrl != "" {
		return elem.Image.ContentUrl
	}
	return ""
}

// ImageInfo represents an image extracted from a slide.
type ImageInfo struct {
	ObjectID   string `json:"object_id"`
	ContentURL string `json:"content_url"`
	SourceURL  string `json:"source_url,omitempty"`
	AltText    string `json:"alt_text,omitempty"`
}

// ExtractImages extracts all images from a page.
func ExtractImages(page *slides.Page) []ImageInfo {
	var images []ImageInfo
	for _, elem := range page.PageElements {
		if elem.Image != nil && elem.Image.ContentUrl != "" {
			info := ImageInfo{
				ObjectID:   elem.ObjectId,
				ContentURL: elem.Image.ContentUrl,
				SourceURL:  elem.Image.SourceUrl,
			}
			// Use element description/title as alt text
			if elem.Description != "" {
				info.AltText = elem.Description
			} else if elem.Title != "" {
				info.AltText = elem.Title
			}
			images = append(images, info)
		}
	}
	return images
}

// SlideContent represents the extracted content of a single slide.
type SlideContent struct {
	Index       int         `json:"index"`
	ObjectID    string      `json:"object_id"`
	Title       string      `json:"title,omitempty"`
	TextContent []string    `json:"text_content"`
	Images      []ImageInfo `json:"images,omitempty"`
	Notes       string      `json:"notes,omitempty"`
}

// PresentationContent represents the extracted content of an entire presentation.
type PresentationContent struct {
	PresentationID string         `json:"presentation_id"`
	Title          string         `json:"title"`
	Locale         string         `json:"locale,omitempty"`
	SlideCount     int            `json:"slide_count"`
	Slides         []SlideContent `json:"slides"`
}

// ExtractPresentationContent extracts all content from a presentation.
// If includeNotes is true, speaker notes are included for each slide.
func ExtractPresentationContent(pres *slides.Presentation, includeNotes bool) PresentationContent {
	content := PresentationContent{
		PresentationID: pres.PresentationId,
		Title:          pres.Title,
		Locale:         pres.Locale,
		SlideCount:     len(pres.Slides),
		Slides:         make([]SlideContent, len(pres.Slides)),
	}

	for i, slide := range pres.Slides {
		sc := SlideContent{
			Index:       i,
			ObjectID:    slide.ObjectId,
			Title:       ExtractSlideTitle(slide),
			TextContent: ExtractTextContent(slide),
			Images:      ExtractImages(slide),
		}
		if includeNotes {
			sc.Notes = ExtractNotesText(slide)
		}
		content.Slides[i] = sc
	}

	return content
}
