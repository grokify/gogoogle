package docsutil

import (
	"strings"

	"google.golang.org/api/docs/v1"
)

// DocumentContent represents extracted content from a Google Doc.
type DocumentContent struct {
	Title      string            `json:"title"`
	DocumentID string            `json:"document_id"`
	Sections   []SectionContent  `json:"sections"`
	Images     []ImageContent    `json:"images,omitempty"`
	Tables     []TableContent    `json:"tables,omitempty"`
	Lists      []ListContent     `json:"lists,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Footers    map[string]string `json:"footers,omitempty"`
}

// SectionContent represents a section of text content.
type SectionContent struct {
	Type    string `json:"type"` // "heading", "paragraph", "list_item"
	Level   int    `json:"level,omitempty"`
	Text    string `json:"text"`
	StyleID string `json:"style_id,omitempty"`
}

// ImageContent represents an embedded image.
type ImageContent struct {
	ObjectID    string `json:"object_id"`
	ContentURI  string `json:"content_uri"`
	SourceURI   string `json:"source_uri,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// TableContent represents a table in the document.
type TableContent struct {
	Rows    int        `json:"rows"`
	Columns int        `json:"columns"`
	Cells   [][]string `json:"cells"`
}

// ListContent represents a list in the document.
type ListContent struct {
	ListID string   `json:"list_id"`
	Items  []string `json:"items"`
}

// ExtractDocumentContent extracts all text and structural content from a document.
func ExtractDocumentContent(doc *docs.Document) *DocumentContent {
	content := &DocumentContent{
		Title:      doc.Title,
		DocumentID: doc.DocumentId,
		Headers:    make(map[string]string),
		Footers:    make(map[string]string),
	}

	if doc.Body != nil {
		content.Sections, content.Images, content.Tables = extractBodyContent(doc.Body)
	}

	// Extract headers
	for headerID, header := range doc.Headers {
		content.Headers[headerID] = extractStructuralElementsText(header.Content)
	}

	// Extract footers
	for footerID, footer := range doc.Footers {
		content.Footers[footerID] = extractStructuralElementsText(footer.Content)
	}

	// Extract lists
	content.Lists = extractLists(doc)

	return content
}

// ExtractPlainText extracts all text from a document as a single string.
func ExtractPlainText(doc *docs.Document) string {
	if doc.Body == nil {
		return ""
	}

	var sb strings.Builder
	for _, elem := range doc.Body.Content {
		sb.WriteString(extractStructuralElementText(elem))
	}
	return sb.String()
}

// ExtractTextByParagraph extracts text organized by paragraphs.
func ExtractTextByParagraph(doc *docs.Document) []string {
	if doc.Body == nil {
		return nil
	}

	var paragraphs []string
	for _, elem := range doc.Body.Content {
		if elem.Paragraph != nil {
			text := extractParagraphText(elem.Paragraph)
			if text != "" {
				paragraphs = append(paragraphs, text)
			}
		}
	}
	return paragraphs
}

// extractBodyContent extracts content from the document body.
func extractBodyContent(body *docs.Body) ([]SectionContent, []ImageContent, []TableContent) {
	var sections []SectionContent
	var images []ImageContent
	var tables []TableContent

	for _, elem := range body.Content {
		if elem.Paragraph != nil {
			section := extractParagraphSection(elem.Paragraph)
			if section.Text != "" {
				sections = append(sections, section)
			}

			// Extract images from paragraph
			for _, pe := range elem.Paragraph.Elements {
				if pe.InlineObjectElement != nil {
					img := extractInlineImage(pe.InlineObjectElement)
					if img != nil {
						images = append(images, *img)
					}
				}
			}
		}

		if elem.Table != nil {
			table := extractTable(elem.Table)
			tables = append(tables, table)
		}
	}

	return sections, images, tables
}

// extractParagraphSection extracts a section from a paragraph.
func extractParagraphSection(para *docs.Paragraph) SectionContent {
	section := SectionContent{
		Type: "paragraph",
		Text: extractParagraphText(para),
	}

	if para.ParagraphStyle != nil {
		style := para.ParagraphStyle.NamedStyleType
		switch {
		case strings.HasPrefix(style, "HEADING_"):
			section.Type = "heading"
			// Extract heading level from style name (HEADING_1 -> 1)
			if len(style) > 8 {
				level := int(style[8] - '0')
				if level >= 1 && level <= 6 {
					section.Level = level
				}
			}
		case style == "TITLE":
			section.Type = "heading"
			section.Level = 1
		case style == "SUBTITLE":
			section.Type = "heading"
			section.Level = 2
		}
		section.StyleID = style
	}

	return section
}

// extractParagraphText extracts plain text from a paragraph.
func extractParagraphText(para *docs.Paragraph) string {
	var sb strings.Builder
	for _, elem := range para.Elements {
		if elem.TextRun != nil {
			sb.WriteString(elem.TextRun.Content)
		}
	}
	return strings.TrimSpace(sb.String())
}

// extractStructuralElementText extracts text from a structural element.
func extractStructuralElementText(elem *docs.StructuralElement) string {
	var sb strings.Builder

	if elem.Paragraph != nil {
		sb.WriteString(extractParagraphText(elem.Paragraph))
		sb.WriteString("\n")
	}

	if elem.Table != nil {
		for _, row := range elem.Table.TableRows {
			for _, cell := range row.TableCells {
				sb.WriteString(extractStructuralElementsText(cell.Content))
				sb.WriteString("\t")
			}
			sb.WriteString("\n")
		}
	}

	if elem.SectionBreak != nil {
		sb.WriteString("\n---\n")
	}

	return sb.String()
}

// extractStructuralElementsText extracts text from multiple structural elements.
func extractStructuralElementsText(elements []*docs.StructuralElement) string {
	var sb strings.Builder
	for _, elem := range elements {
		sb.WriteString(extractStructuralElementText(elem))
	}
	return strings.TrimSpace(sb.String())
}

// extractInlineImage extracts image information from an inline object element.
func extractInlineImage(elem *docs.InlineObjectElement) *ImageContent {
	if elem.InlineObjectId == "" {
		return nil
	}
	return &ImageContent{
		ObjectID: elem.InlineObjectId,
	}
}

// extractTable extracts table content.
func extractTable(table *docs.Table) TableContent {
	tc := TableContent{
		Rows:    int(table.Rows),
		Columns: int(table.Columns),
	}

	for _, row := range table.TableRows {
		var rowCells []string
		for _, cell := range row.TableCells {
			cellText := extractStructuralElementsText(cell.Content)
			rowCells = append(rowCells, cellText)
		}
		tc.Cells = append(tc.Cells, rowCells)
	}

	return tc
}

// extractLists extracts list information from the document.
func extractLists(doc *docs.Document) []ListContent {
	if doc.Lists == nil {
		return nil
	}

	var lists []ListContent
	for listID := range doc.Lists {
		lc := ListContent{
			ListID: listID,
		}

		// Find list items in document body
		if doc.Body != nil {
			for _, elem := range doc.Body.Content {
				if elem.Paragraph != nil && elem.Paragraph.Bullet != nil {
					if elem.Paragraph.Bullet.ListId == listID {
						text := extractParagraphText(elem.Paragraph)
						if text != "" {
							lc.Items = append(lc.Items, text)
						}
					}
				}
			}
		}

		if len(lc.Items) > 0 {
			lists = append(lists, lc)
		}
	}

	return lists
}

// EnrichImagesWithURIs populates image URIs from the document's inline objects.
func EnrichImagesWithURIs(content *DocumentContent, doc *docs.Document) {
	if doc.InlineObjects == nil {
		return
	}

	for i := range content.Images {
		objID := content.Images[i].ObjectID
		if obj, ok := doc.InlineObjects[objID]; ok {
			if obj.InlineObjectProperties != nil && obj.InlineObjectProperties.EmbeddedObject != nil {
				eo := obj.InlineObjectProperties.EmbeddedObject
				content.Images[i].ContentURI = eo.ImageProperties.ContentUri
				content.Images[i].SourceURI = eo.ImageProperties.SourceUri
				content.Images[i].Title = eo.Title
				content.Images[i].Description = eo.Description
			}
		}
	}
}
