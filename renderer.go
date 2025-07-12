package main

import (
	"regexp"
	"strings"

	"github.com/cucumber/messages/go/v22"
	"github.com/willf/pad/utf8"
	"golang.org/x/exp/slices"
)

const INDENT = "  "

var docStringLineRegexp = regexp.MustCompile("(^|\n)([^\n])")
var spaceRegexp = regexp.MustCompile(`\s+`)

type renderer struct {
	builder  *strings.Builder
	depth    int
	comments []*messages.Comment
}

func NewRenderer() *renderer {
	return &renderer{&strings.Builder{}, 0, nil}
}

func (r *renderer) Render(d *messages.GherkinDocument) string {
	r.comments = append([]*messages.Comment{}, d.Comments...)
	slices.Reverse(r.comments)

	r.renderFeature(d.Feature)

	return r.builder.String()
}

func (r *renderer) renderFeature(f *messages.Feature) {
	r.renderTags(f.Tags)

	r.writeHeadline("Feature", f.Name, f.Location)

	r.depth++
	defer func() { r.depth-- }()

	r.writeDescription(f.Description)

	for i, c := range f.Children {
		if c.Background != nil {
			r.renderBackground(c.Background)
		}

		if c.Scenario != nil {
			r.renderScenario(c.Scenario)
		}

		if c.Rule != nil {
			r.renderRule(c.Rule)
		}

		if i != len(f.Children)-1 {
			r.writeLine("")
		}
	}
}

func (r *renderer) renderBackground(b *messages.Background) {
	r.writeHeadline("Background", b.Name, b.Location)

	r.depth++
	defer func() { r.depth-- }()

	r.writeDescription(b.Description)
	r.renderSteps(b.Steps)
}

func (r *renderer) renderTags(ts []*messages.Tag) {
	if len(ts) > 0 {
		ss := []string{}

		for _, t := range ts {
			ss = append(ss, t.Name)
		}

		slices.Sort(ss)

		r.writeLine(strings.Join(ss, " "))
	}
}

func (r *renderer) renderScenario(s *messages.Scenario) {
	r.renderTags(s.Tags)

	t := "Scenario"

	if len(s.Examples) != 0 {
		t += " Outline"
	}

	r.writeHeadline(t, s.Name, s.Location)

	r.depth++
	defer func() { r.depth-- }()

	r.writeDescription(s.Description)
	r.renderSteps(s.Steps)

	if len(s.Examples) != 0 {
		r.writeLine("")
		r.renderExamples(s.Examples)
	}
}

func (r *renderer) renderRule(l *messages.Rule) {
	r.renderTags(l.Tags)
	r.writeHeadline("Rule", l.Name, l.Location)

	r.depth++
	defer func() { r.depth-- }()

	r.writeDescription(l.Description)

	for i, c := range l.Children {
		if c.Background != nil {
			r.renderBackground(c.Background)
		}

		if c.Scenario != nil {
			r.renderScenario(c.Scenario)
		}

		if i != len(l.Children)-1 {
			r.writeLine("")
		}
	}
}

func (r *renderer) renderSteps(ss []*messages.Step) {
	for _, s := range ss {
		r.renderStep(s)
	}
}

func (r *renderer) renderDocString(d *messages.DocString) {
	r.depth++
	defer func() { r.depth-- }()

	r.writeLine(`"""` + d.MediaType)

	if d.Content != "" {
		r.builder.WriteString(
			docStringLineRegexp.
				ReplaceAllString(d.Content, "$1"+r.padding()+"$2") + "\n",
		)
	}

	r.writeLine(`"""`)
}

func (r *renderer) renderStep(s *messages.Step) {
	r.renderComments(s.Location)
	r.writeLine(strings.TrimSpace(s.Keyword) + " " + s.Text)

	if s.DocString != nil {
		r.renderDocString(s.DocString)
	}

	r.depth++
	defer func() {
		r.depth--
		r.renderAfterComments(stepLastLocation(s))
	}()

	if s.DataTable != nil {
		r.renderDataTable(s.DataTable)
	}
}

func (r *renderer) renderExamples(es []*messages.Examples) {
	for i, e := range es {
		r.renderTags(e.Tags)
		r.writeHeadline("Examples", e.Name, e.Location)

		r.depth++
		r.writeDescription(e.Description)
		r.renderExampleTable(e.TableHeader, e.TableBody)
		r.depth--

		if i != len(es)-1 {
			r.writeLine("")
		}
	}
}

func (r renderer) renderExampleTable(h *messages.TableRow, rs []*messages.TableRow) {
	ws := r.getCellWidths(append([]*messages.TableRow{h}, rs...))

	r.renderCells(h.Cells, ws)

	for _, t := range rs {
		r.renderCells(t.Cells, ws)
	}
}

func (r renderer) renderDataTable(t *messages.DataTable) {
	ws := r.getCellWidths(t.Rows)

	for _, t := range t.Rows {
		r.renderCells(t.Cells, ws)
	}
}

func (r renderer) renderCells(cs []*messages.TableCell, ws []int) {
	s := "|"

	for i, c := range cs {
		s += " " + utf8.Right(r.escapeCellValue(c), ws[i], " ") + " |"
	}

	r.writeLine(s)
}
func (renderer) escapeCellValue(c *messages.TableCell) string {
	s := c.Value

	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "|", "\\|")

	return s
}

func (r *renderer) renderComments(l *messages.Location) {
	for len(r.comments) > 0 && r.comments[len(r.comments)-1].Location.Line <= l.Line {
		i := len(r.comments) - 1
		r.renderComment(r.comments[i])
		r.comments = r.comments[:i]
	}
}

func (r *renderer) renderAfterComments(l *messages.Location) {
	for len(r.comments) > 0 && r.comments[len(r.comments)-1].Location.Line <= l.Line+1 {
		i := len(r.comments) - 1
		c := r.comments[i]
		r.renderComment(c)
		r.comments = r.comments[:i]
		l = c.Location
	}
}

func (r *renderer) renderComment(c *messages.Comment) {
	r.writeLine(normalizeText(c.Text))
}

func (r renderer) getCellWidths(rs []*messages.TableRow) []int {
	ws := make([]int, len(rs[0].Cells))

	for _, row := range rs {
		for i, c := range row.Cells {
			if w := len(r.escapeCellValue(c)); w > ws[i] {
				ws[i] = w
			}
		}
	}

	return ws
}

func (r renderer) writeDescription(s string) {
	if s != "" {
		r.writeLine(strings.TrimSpace(s))
		r.writeLine("")
	}
}

func (r *renderer) writeHeadline(s, t string, l *messages.Location) {
	r.renderComments(l)

	s += ":"

	if t != "" {
		s += " " + normalizeText(t)
	}

	r.writeLine(s)
}

func (r *renderer) writeLine(s string) {
	if s != "" {
		s = r.padding() + s
	}

	_, err := r.builder.WriteString(s + "\n")

	if err != nil {
		panic(err)
	}
}

func (r renderer) padding() string {
	return strings.Repeat(INDENT, r.depth)
}

func stepLastLocation(s *messages.Step) *messages.Location {
	l := s.Location

	if s.DocString != nil {
		ll := *s.DocString.Location
		ll.Line += int64(strings.Count(s.DocString.Content, "\n")) + 2
		l = &ll
	}

	if s.DataTable != nil {
		ll := *s.DataTable.Location
		ll.Line += int64(len(s.DataTable.Rows) - 1)
		l = &ll

	}

	return l
}

func normalizeText(s string) string {
	return spaceRegexp.ReplaceAllString(strings.TrimSpace(s), " ")
}
