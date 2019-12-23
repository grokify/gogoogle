package slidesutil

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grokify/gotilla/fmt/fmtutil"
	"google.golang.org/api/slides/v1"
)

type CommonMarkData struct {
	text  string
	lines []TextLine
}

func NewCommonMarkData(cm string) CommonMarkData {
	cmd := CommonMarkData{
		text:  cm,
		lines: Split(cm)}
	cmd.Inflate()
	return cmd
}

func (cmd *CommonMarkData) Inflate() {
	prefix := int64(0)
	for i, line := range cmd.lines {
		line = InflateCommonMarkToGoogleSlides(line)
		// Set Line Start and End
		line.GoogleIndexStart = prefix
		line.GoogleIndexEnd = prefix + int64(len(line.GoogleSlideText))
		prefix += int64(len(line.GoogleSlideText)) + 1
		line.Inflated = true
		cmd.lines[i] = line
	}
}

func (cmd *CommonMarkData) GoogleSlideTextString() string {
	googTexts := []string{}
	for _, line := range cmd.lines {
		googTexts = append(googTexts, line.GoogleSlideText)
	}
	return strings.Join(googTexts, "\n")
}

func (cmd *CommonMarkData) Lines() []TextLine {
	return cmd.lines
}

func (cmd *CommonMarkData) LineCount() int {
	return len(cmd.lines)
}

func Split(cm string) []TextLine {
	lines := []TextLine{}
	raw := strings.Split(cm, "\n")
	fmt.Println(cm)
	fmtutil.PrintJSON(raw)
	for _, line := range raw {
		lines = append(lines, TextLine{CommonMarkText: line})
	}
	return lines
}

type TextLine struct {
	CommonMarkText   string
	GoogleSlideText  string
	IsBullet         bool
	IsBold           bool
	Links            []LinkInfo
	Inflated         bool
	GoogleIndexStart int64
	GoogleIndexEnd   int64
}

var (
	rxLink        *regexp.Regexp = regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	rxBulletFirst *regexp.Regexp = regexp.MustCompile(`^[*-]\s+`)
	rxBullet      *regexp.Regexp = regexp.MustCompile(`^(\s\s\s\s)?[*-]\s+`)
	rxBold        *regexp.Regexp = regexp.MustCompile(`^(\s*)\*\*(.*)\*\*\s*$`)
)

func InflateCommonMarkToGoogleSlides(cml TextLine) TextLine {
	googLine := cml.CommonMarkText
	links := rxLink.FindAllStringSubmatch(googLine, -1)
	if len(links) > 0 {
		for _, mx := range links {
			fullMatch := mx[0]
			textDisp := mx[1]
			fullMatchRx := regexp.MustCompile(regexp.QuoteMeta(fullMatch))
			googLine = fullMatchRx.ReplaceAllString(googLine, textDisp)
		}
	}
	if rxBulletFirst.MatchString(googLine) {
		cml.IsBullet = true
		googLine = rxBulletFirst.ReplaceAllString(googLine, "")
	} else if rxBullet.MatchString(googLine) {
		cml.IsBullet = true
		googLine = rxBullet.ReplaceAllString(googLine, "\t")
	}
	if rxBold.MatchString(googLine) {
		cml.IsBold = true
		googLine = rxBold.ReplaceAllString(googLine, "$1$2")
	}
	if len(links) > 0 {
		linksInfos := []LinkInfo{}
		for _, link := range links {
			if len(link) < 3 {
				continue
			}
			linkText := link[1]
			linkText = rxBold.ReplaceAllString(linkText, "$1$2")
			linkURL := link[2]
			idxStart := strings.Index(googLine, linkText)
			if idxStart >= 0 {
				linkRange := &slides.Range{
					Type:       RangeTypeFixedRange,
					StartIndex: cml.GoogleIndexStart + int64(idxStart),
					EndIndex:   cml.GoogleIndexStart + int64(idxStart+len(linkText))}
				linksInfos = append(linksInfos, LinkInfo{
					URL:   linkURL,
					Range: linkRange})
			}
		}
		cml.Links = linksInfos
	}
	cml.GoogleSlideText = googLine
	return cml
}

func LineStartEndIndexes(cml TextLine, priorLength int64) (int64, int64) {
	start := int64(0) + priorLength
	finish := int64(len(cml.GoogleSlideText)) + priorLength
	return start, finish
}

func CommonMarkDataToRequests(textboxID string, cmd CommonMarkData, underlineLinks bool) []*slides.Request {
	reqs := []*slides.Request{
		{
			InsertText: &slides.InsertTextRequest{
				ObjectId: textboxID,
				Text:     cmd.GoogleSlideTextString(),
			},
		},
	}
	lines := cmd.Lines()
	for _, line := range lines {
		if line.IsBold {
			reqs = append(reqs, UpdateTextStyleRequestBold(textboxID, line.GoogleIndexStart, line.GoogleIndexEnd))
		}
		if line.IsBullet {
			reqs = append(reqs, UpdateTextStyleRequestBullet(textboxID, line.GoogleIndexStart, line.GoogleIndexEnd))
		}
		for _, linkInfo := range line.Links {
			reqs = append(reqs, UpdateTextStyleRequestLinkURL(textboxID, linkInfo.URL, linkInfo.Range, underlineLinks))
		}
	}
	if 1 == 0 {
		reqs = append(reqs, UpdateTextStyleRequestFontSize(textboxID, slides.Dimension{
			Magnitude: float64(14),
			Unit:      "PT"}))
	}
	return reqs
}

type LinkInfo struct {
	URL   string
	Range *slides.Range
}
