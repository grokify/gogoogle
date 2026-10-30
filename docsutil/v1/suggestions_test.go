package docsutil

import (
	"encoding/json"
	"strings"
	"testing"

	"google.golang.org/api/docs/v1"
)

// suggestionDoc is a documents.get response (SUGGESTIONS_INLINE) with a
// paragraph carrying an insertion and a deletion, an untouched paragraph, and
// a table whose second cell carries an insertion.
const suggestionDoc = `{
  "documentId": "doc1",
  "body": {"content": [
    {"sectionBreak": {}},
    {"paragraph": {"elements": [
      {"textRun": {"content": "We will "}},
      {"textRun": {"content": "quickly ", "suggestedDeletionIds": ["s.del1"]}},
      {"textRun": {"content": "train "}},
      {"textRun": {"content": "all ", "suggestedInsertionIds": ["s.ins1"]}},
      {"textRun": {"content": "staff.\n"}}
    ]}},
    {"paragraph": {"elements": [{"textRun": {"content": "Untouched.\n"}}]}},
    {"table": {"tableRows": [
      {"tableCells": [
        {"content": [{"paragraph": {"elements": [{"textRun": {"content": "Metric\n"}}]}}]},
        {"content": [{"paragraph": {"elements": [
          {"textRun": {"content": "Monthly", "suggestedInsertionIds": ["s.ins2"]}},
          {"textRun": {"content": " count\n", "suggestedInsertionIds": ["s.ins2"]}}
        ]}}]}
      ]}
    ]}}
  ]}
}`

func loadDoc(t *testing.T, raw string) *docs.Document {
	t.Helper()
	var d docs.Document
	if err := json.Unmarshal([]byte(raw), &d); err != nil {
		t.Fatal(err)
	}
	return &d
}

func TestExtractSuggestionContent(t *testing.T) {
	got := ExtractSuggestionContent(loadDoc(t, suggestionDoc))
	if got.DocumentID != "doc1" {
		t.Errorf("DocumentID = %q", got.DocumentID)
	}
	wantLoc := []string{"body/p/1", "body/p/2", "body/table/3/r/0/c/0", "body/table/3/r/0/c/1"}
	if len(got.Blocks) != len(wantLoc) {
		t.Fatalf("got %d blocks, want %d", len(got.Blocks), len(wantLoc))
	}
	for i, b := range got.Blocks {
		if b.Location != wantLoc[i] {
			t.Errorf("block %d location = %q, want %q", i, b.Location, wantLoc[i])
		}
	}
	p := got.Blocks[0]
	if p.Original != "We will quickly train staff.\n" {
		t.Errorf("original = %q", p.Original)
	}
	if p.Accepted != "We will train all staff.\n" {
		t.Errorf("accepted = %q", p.Accepted)
	}
	kinds := make([]SegmentKind, len(p.Segments))
	for i, s := range p.Segments {
		kinds[i] = s.Kind
	}
	want := []SegmentKind{SegmentOriginal, SegmentDeletion, SegmentOriginal, SegmentInsertion, SegmentOriginal}
	if strings.Join(kindStrings(kinds), ",") != strings.Join(kindStrings(want), ",") {
		t.Errorf("kinds = %v, want %v", kinds, want)
	}
	// Adjacent runs with the same kind and IDs merge into one segment.
	cell := got.Blocks[3]
	if len(cell.Segments) != 1 || cell.Segments[0].Text != "Monthly count\n" {
		t.Errorf("cell segments = %+v", cell.Segments)
	}
	if cell.Original != "" || cell.Accepted != "Monthly count\n" {
		t.Errorf("cell original %q accepted %q", cell.Original, cell.Accepted)
	}
	if ids := strings.Join(got.SuggestionIDs, ","); ids != "s.del1,s.ins1,s.ins2" {
		t.Errorf("SuggestionIDs = %s", ids)
	}
}

func kindStrings(ks []SegmentKind) []string {
	out := make([]string, len(ks))
	for i, k := range ks {
		out[i] = string(k)
	}
	return out
}

// INV-4: original plus deletion segments rebuild the original text, original
// plus insertion segments rebuild the accepted text.
func TestExtractSuggestionContentInvariant(t *testing.T) {
	doc := loadDoc(t, suggestionDoc)
	got := ExtractSuggestionContent(doc)
	var orig, acc string
	for _, b := range got.Blocks {
		for _, s := range b.Segments {
			if s.Kind != SegmentInsertion {
				orig += s.Text
			}
			if s.Kind != SegmentDeletion {
				acc += s.Text
			}
		}
	}
	if want := "We will quickly train staff.\nUntouched.\nMetric\n"; orig != want {
		t.Errorf("original = %q, want %q", orig, want)
	}
	if want := "We will train all staff.\nUntouched.\nMetric\nMonthly count\n"; acc != want {
		t.Errorf("accepted = %q, want %q", acc, want)
	}
}

func TestExtractSuggestionContentDoesNotMutateInput(t *testing.T) {
	doc := loadDoc(t, suggestionDoc)
	before, err := doc.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	got := ExtractSuggestionContent(doc)
	got.Blocks[0].Segments[1].SuggestionIDs[0] = "changed"
	after, err := doc.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Error("input document was modified")
	}
}

func TestExtractSuggestionContentEmpty(t *testing.T) {
	for name, doc := range map[string]*docs.Document{"nil": nil, "no body": {DocumentId: "x"}} {
		got := ExtractSuggestionContent(doc)
		if got.Blocks == nil || got.SuggestionIDs == nil || len(got.Blocks) != 0 {
			t.Errorf("%s: want empty non-nil slices, got %+v", name, got)
		}
	}
}
