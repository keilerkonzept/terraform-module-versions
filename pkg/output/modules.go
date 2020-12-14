package output

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"sort"

	junit "github.com/jstemmer/go-junit-report/formatter"
	"github.com/olekukonko/tablewriter"
)

type Modules []Module

func (m Modules) Len() int           { return len(m) }
func (m Modules) Less(i, j int) bool { return m[i].SortKey() < m[j].SortKey() }
func (m Modules) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

type Module struct {
	Path              string `json:"path,omitempty"`
	Name              string `json:"name,omitempty"`
	Type              string `json:"type,omitempty"`
	Source            string `json:"source,omitempty"`
	VersionConstraint string `json:"constraint,omitempty"`
	Version           string `json:"version,omitempty"`
}

func (m *Module) SortKey() string {
	return fmt.Sprint(m.Path, m.Name)
}

func (m Modules) Write(w io.Writer, as Format) error {
	switch as {
	case FormatJSON:
		return m.WriteJSON(w)
	case FormatJSONL:
		return m.WriteJSONL(w)
	case FormatMarkdown:
		return m.WriteMarkdown(w)
	case FormatMarkdownWide:
		return m.WriteMarkdownWide(w)
	case FormatJUnit:
		return m.WriteJUnit(w)
	}
	return nil
}

func (m Modules) WriteJSONL(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	for _, item := range m {
		if err := enc.Encode(item); err != nil {
			return fmt.Errorf("encode json: %w", err)
		}
	}
	return nil
}

func (m Modules) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(m)
}

func (m Modules) WriteMarkdownWide(w io.Writer) error {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Type", "Name", "Constraint", "Version", "Source", "Path"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	rows := make([][]string, 0, len(m))
	for _, item := range m {
		row := []string{item.Type, item.Name, item.VersionConstraint, item.Version, item.Source, item.Path}
		rows = append(rows, row)
	}
	table.AppendBulk(rows)
	table.Render()
	return nil
}

func (m Modules) WriteMarkdown(w io.Writer) error {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Type", "Name", "Constraint", "Version", "Source"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	rows := make([][]string, 0, len(m))
	for _, item := range m {
		row := []string{item.Type, item.Name, item.VersionConstraint, item.Version, item.Source}
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		return fmt.Sprint(rows[i]) > fmt.Sprint(rows[j])
	})
	table.AppendBulk(rows)
	table.Render()
	return nil
}

func (m Modules) WriteJUnit(w io.Writer) error {
	testCases := make([]junit.JUnitTestCase, len(m))

	failures := 0
	for i, module := range m {
		testCase := junit.JUnitTestCase{
			Name:      module.Name,
			Classname: module.Path,
			Time:      "0",
		}
		success := module.Version != "" || module.VersionConstraint != ""
		if module.Type == "git" && module.Version == "" { // special case for git modules: the version annotation is ineffective
			success = false
		}
		if module.Type == "local" { // local modules can't specify versions or constraints
			success = true
		}
		if !success {
			failures++
			testCase.Failure = &junit.JUnitFailure{
				Message:  "Module reference does not explicitly specify a version or version constraint",
				Contents: "",
			}
		}
		testCases[i] = testCase
	}

	suites := junit.JUnitTestSuites{
		Suites: []junit.JUnitTestSuite{
			{
				Time:      "0",
				Tests:     len(m),
				Failures:  failures,
				TestCases: testCases,
			},
		},
	}

	if _, err := fmt.Fprint(w, xml.Header); err != nil {
		return fmt.Errorf("encode junit xml: %w", err)
	}
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(suites); err != nil {
		return fmt.Errorf("encode junit xml: %w", err)
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return fmt.Errorf("encode junit xml: %w", err)
	}
	return nil
}
