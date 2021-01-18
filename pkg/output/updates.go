package output

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"

	junit "github.com/jstemmer/go-junit-report/formatter"
	"github.com/olekukonko/tablewriter"
)

type Updates []Update

func (u Updates) Len() int           { return len(u) }
func (u Updates) Less(i, j int) bool { return u[i].SortKey() < u[j].SortKey() }
func (u Updates) Swap(i, j int)      { u[i], u[j] = u[j], u[i] }

type Update struct {
	Path              string `json:"path,omitempty"`
	Name              string `json:"name,omitempty"`
	Source            string `json:"source,omitempty"`
	VersionConstraint string `json:"constraint,omitempty"`
	Version           string `json:"version,omitempty"`
	LatestMatching    string `json:"latestMatching,omitempty"`
	LatestOverall     string `json:"latestOverall,omitempty"`
	MatchingUpdate    bool   `json:"matchingUpdate,omitempty"`
	NonMatchingUpdate bool   `json:"nonMatchingUpdate,omitempty"`
}

func (u *Update) SortKey() string {
	return fmt.Sprint(u.Path, u.Name)
}

func (u Updates) Format(w io.Writer, as Format) error {
	switch as {
	case FormatJSON:
		return u.WriteJSON(w)
	case FormatJSONL:
		return u.WriteJSONL(w)
	case FormatMarkdown:
		return u.WriteMarkdown(w)
	case FormatMarkdownWide:
		return u.WriteMarkdownWide(w)
	case FormatJUnit:
		return u.WriteJUnit(w)
	}
	return nil
}

func (u Updates) WriteJSONL(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	for _, item := range u {
		if err := enc.Encode(item); err != nil {
			return fmt.Errorf("encode json: %w", err)
		}
	}
	return nil
}

func (u Updates) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(u)
}

func (u Updates) WriteMarkdownWide(w io.Writer) error {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Update?", "Name", "Path", "Source", "Constraint", "Version", "Latest matching", "Latest"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	rows := make([][]string, 0, len(u))
	for _, item := range u {
		update := ""
		switch {
		case item.MatchingUpdate:
			update = "Y"
		case item.NonMatchingUpdate:
			update = "(Y)"
		case item.Version == "":
			update = "?"
		}
		row := []string{update, item.Name, item.Path, item.Source, item.VersionConstraint, item.Version, item.LatestMatching, item.LatestOverall}
		rows = append(rows, row)
	}
	table.AppendBulk(rows)
	table.Render()
	return nil
}

func (u Updates) WriteMarkdown(w io.Writer) error {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Update?", "Name", "Constraint", "Version", "Latest matching", "Latest"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	rows := make([][]string, 0, len(u))
	for _, item := range u {
		update := ""
		switch {
		case item.MatchingUpdate:
			update = "Y"
		case item.NonMatchingUpdate:
			update = "(Y)"
		case item.Version == "":
			update = "?"
		}
		row := []string{update, item.Name, item.VersionConstraint, item.Version, item.LatestMatching, item.LatestOverall}
		rows = append(rows, row)
	}
	table.AppendBulk(rows)
	table.Render()
	return nil
}

func (u Updates) WriteJUnit(w io.Writer) error {
	testCases := make([]junit.JUnitTestCase, len(u))

	failures := 0
	for i, update := range u {
		testCase := junit.JUnitTestCase{
			Name:      update.Name,
			Classname: update.Path,
			Time:      "0",
		}
		success := !update.MatchingUpdate
		if !success {
			failures++
			testCase.Failure = &junit.JUnitFailure{
				Message:  fmt.Sprintf("Module version can be updated to %v (from %v)", update.LatestMatching, update.Version),
				Contents: "",
			}
		}
		testCases[i] = testCase
	}

	suites := junit.JUnitTestSuites{
		Suites: []junit.JUnitTestSuite{
			{
				Time:      "0",
				Tests:     len(u),
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
