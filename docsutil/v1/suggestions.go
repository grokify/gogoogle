package docsutil

import (
	"fmt"
	"slices"
	"strings"

	"google.golang.org/api/docs/v1"
)

// SegmentKind classifies a run of text relative to pending suggestions.
type SegmentKind string

const (
	// SegmentOriginal is text no suggestion touches.
	SegmentOriginal SegmentKind = "original"
	// SegmentInsertion is text a suggestion would add.
	SegmentInsertion SegmentKind = "insertion"
	// SegmentDeletion is text a suggestion would remove.
	SegmentDeletion SegmentKind = "deletion"
)

// Segment is a run of text of one kind with the suggestions that affect it.
type Segment struct {
	Kind          SegmentKind `json:"kind"`
	Text          string      `json:"text"`
	SuggestionIDs []string    `json:"suggestion_ids,omitempty"`
}

// Block is one paragraph or table cell with its segments in document order.
type Block struct {
	// Location is "body/p/<n>" for a paragraph or
	// "body/table/<n>/r/<row>/c/<col>" for a table cell, where n is the index
	// of the structural element in the body.
	Location string    `json:"location"`
	Segments []Segment `json:"segments"`
	// Original is the original and deletion segments concatenated.
	Original string `json:"original"`
	// Accepted is the original and insertion segments concatenated.
	Accepted string `json:"accepted"`
}

// SuggestionContent is a document's text split by pending suggestions.
type SuggestionContent struct {
	DocumentID    string   `json:"document_id"`
	Blocks        []Block  `json:"blocks"`
	SuggestionIDs []string `json:"suggestion_ids"` // distinct, document order
}

// ExtractSuggestionContent splits the body of a document fetched with
// suggestions view mode SUGGESTIONS_INLINE into blocks of original, inserted,
// and deleted text. The Docs API does not report who made a suggestion, so
// blocks carry no author.
//
// Document order is preserved: concatenating the original and deletion
// segments of the blocks yields the original text, and concatenating the
// original and insertion segments yields the text as if every suggestion were
// accepted. A run that is both suggested for insertion and for deletion is
// classified as an insertion.
//
// Every paragraph and table cell yields a block, including those without
// suggestions. The input is not modified.
func ExtractSuggestionContent(doc *docs.Document) *SuggestionContent {
	out := &SuggestionContent{
		Blocks:        []Block{},
		SuggestionIDs: []string{},
	}
	if doc == nil {
		return out
	}
	out.DocumentID = doc.DocumentId
	if doc.Body == nil {
		return out
	}
	for i, elem := range doc.Body.Content {
		switch {
		case elem.Paragraph != nil:
			out.add(newBlock(fmt.Sprintf("body/p/%d", i), elem.Paragraph))
		case elem.Table != nil:
			for r, row := range elem.Table.TableRows {
				for c, cell := range row.TableCells {
					var paras []*docs.Paragraph
					for _, ce := range cell.Content {
						if ce.Paragraph != nil {
							paras = append(paras, ce.Paragraph)
						}
					}
					out.add(newBlock(fmt.Sprintf("body/table/%d/r/%d/c/%d", i, r, c), paras...))
				}
			}
		}
	}
	return out
}

func (sc *SuggestionContent) add(b Block) {
	sc.Blocks = append(sc.Blocks, b)
	for _, s := range b.Segments {
		for _, id := range s.SuggestionIDs {
			if !slices.Contains(sc.SuggestionIDs, id) {
				sc.SuggestionIDs = append(sc.SuggestionIDs, id)
			}
		}
	}
}

func newBlock(location string, paras ...*docs.Paragraph) Block {
	b := Block{Location: location, Segments: []Segment{}}
	for _, p := range paras {
		for _, el := range p.Elements {
			if el.TextRun == nil || el.TextRun.Content == "" {
				continue
			}
			b.appendRun(el.TextRun)
		}
	}
	var orig, acc strings.Builder
	for _, s := range b.Segments {
		if s.Kind != SegmentInsertion {
			orig.WriteString(s.Text)
		}
		if s.Kind != SegmentDeletion {
			acc.WriteString(s.Text)
		}
	}
	b.Original, b.Accepted = orig.String(), acc.String()
	return b
}

// appendRun adds a run, merging it into the previous segment when the kind and
// suggestion IDs match.
func (b *Block) appendRun(tr *docs.TextRun) {
	kind, ids := SegmentOriginal, []string(nil)
	switch {
	case len(tr.SuggestedInsertionIds) > 0:
		kind, ids = SegmentInsertion, tr.SuggestedInsertionIds
	case len(tr.SuggestedDeletionIds) > 0:
		kind, ids = SegmentDeletion, tr.SuggestedDeletionIds
	}
	if n := len(b.Segments); n > 0 {
		last := &b.Segments[n-1]
		if last.Kind == kind && slices.Equal(last.SuggestionIDs, ids) {
			last.Text += tr.Content
			return
		}
	}
	b.Segments = append(b.Segments, Segment{Kind: kind, Text: tr.Content, SuggestionIDs: slices.Clone(ids)})
}
