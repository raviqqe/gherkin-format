package main

import (
	"regexp"
	"strings"

	"github.com/cucumber/messages/go/v22"
	"github.com/willf/pad/utf8"
)

const INDENT = "  "

type renderer struct {
	*strings.Builder
	depth int
}

func newRenderer() *renderer {
	return &renderer{&strings.Builder{}, 0}
}

func (r *renderer) Render(d *messages.GherkinDocument) string {
	r.renderFeature(d.Feature)

	return r.Builder.String()
}

func (r *renderer) renderFeature(f *messages.Feature) {
	r.writeHeadline("Feature", f.Name)

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
	r.writeHeadline("Background", b.Name)

	r.depth++
	defer func() { r.depth-- }()

	r.writeDescription(b.Description)
	r.renderSteps(b.Steps)
}

func (r *renderer) renderScenario(s *messages.Scenario) {
	t := "Scenario"

	if len(s.Examples) != 0 {
		t += " Outline"
	}

	r.writeHeadline(t, s.Name)

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
	r.writeHeadline("Rule", l.Name)

	r.depth++
	defer func() { r.depth-- }()

	r.writeDescription(l.Description)

	for _, c := range l.Children {
		r.writeLine("")

		if c.Background != nil {
			r.renderBackground(c.Background)
		}

		if c.Scenario != nil {
			r.renderScenario(c.Scenario)
		}
	}
}

func (r *renderer) renderSteps(ss []*messages.Step) {
	for i, s := range ss {
		r.renderStep(s, i == len(ss)-1)
	}
}

func (r *renderer) renderDocString(d *messages.DocString) {
	r.writeLine(`"""` + d.MediaType)

	if d.Content != "" {
		r.WriteString(
			regexp.MustCompile("(^|\n)([^\n])").
				ReplaceAllString(d.Content, "$1"+r.padding()+"$2") + "\n",
		)
	}

	r.writeLine(`"""`)
}

func (r *renderer) renderStep(s *messages.Step, last bool) {
	r.writeLine(strings.TrimSpace(s.Keyword) + " " + s.Text)

	if s.DocString != nil {
		r.renderDocString(s.DocString)
	}

	if s.DataTable != nil {
		r.renderDataTable(s.DataTable)
	}
}

func (r *renderer) renderExamples(es []*messages.Examples) {
	r.writeHeadline("Examples", "")

	r.depth++
	defer func() { r.depth-- }()

	for _, e := range es {
		if e.Name != "" {
			r.writeLine("")
			r.writeHeadline(e.Name, "")
		}

		r.writeDescription(e.Description)
		r.renderExampleTable(e.TableHeader, e.TableBody)
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
		s += " " + utf8.Right(c.Value, ws[i], " ") + " |"
	}

	r.writeLine(s)
}

func (renderer) getCellWidths(rs []*messages.TableRow) []int {
	ws := make([]int, len(rs[0].Cells))

	for _, r := range rs {
		for i, c := range r.Cells {
			if w := len(c.Value); w > ws[i] {
				ws[i] = w
			}
		}
	}

	return ws
}

func (r renderer) writeDescription(s string) {
	if s != "" {
		r.writeLine("")
		r.writeLine(strings.TrimSpace(s))
	}
}

func (r renderer) writeHeadline(s, t string) {
	s += ":"

	if t != "" {
		s += " " + t
	}

	r.writeLine(s)
}

func (r renderer) writeLine(s string) {
	if s != "" {
		s = r.padding() + s
	}

	_, err := r.WriteString(s + "\n")

	if err != nil {
		panic(err)
	}
}

func (r renderer) padding() string {
	return strings.Repeat(INDENT, r.depth)
}
