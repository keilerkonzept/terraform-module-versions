package output

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/olekukonko/tablewriter"
)

type Updates []Update

type Update struct {
	Path              string `json:"path,omitempty"`
	Name              string `json:"name,omitempty"`
	VersionConstraint string `json:"constraint,omitempty"`
	Version           string `json:"version,omitempty"`
	LatestMatching    string `json:"latestMatching,omitempty"`
	LatestOverall     string `json:"latestOverall,omitempty"`
	MatchingUpdate    bool   `json:"matchingUpdate,omitempty"`
	NonMatchingUpdate bool   `json:"nonMatchingUpdate,omitempty"`
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

func (u Updates) WriteTable(w io.Writer) {
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
	sort.Slice(rows, func(i, j int) bool {
		return fmt.Sprint(rows[i]) > fmt.Sprint(rows[j])
	})
	table.AppendBulk(rows)
	table.Render()
}
