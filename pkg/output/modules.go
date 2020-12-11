package output

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/olekukonko/tablewriter"
)

type Modules []Module

type Module struct {
	Path              string `json:"path,omitempty"`
	Name              string `json:"name,omitempty"`
	Type              string `json:"type,omitempty"`
	Source            string `json:"source,omitempty"`
	VersionConstraint string `json:"constraint,omitempty"`
	Version           string `json:"version,omitempty"`
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

func (m Modules) WriteTable(w io.Writer) {
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
}
