package main

import (
	"fmt"
	"io"
	"sort"

	"github.com/olekukonko/tablewriter"
)

func listPrettyPrint(w io.Writer, output []outputList) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Type", "Name", "Constraint", "Version", "Source"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	var rows [][]string
	for _, o := range output {
		row := []string{o.Type, o.Name, o.VersionConstraint, o.Version, o.Source}
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		return fmt.Sprint(rows[i]) > fmt.Sprint(rows[j])
	})
	table.AppendBulk(rows)
	table.Render()
}

func updatePrettyPrint(w io.Writer, output []outputUpdates) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Update?", "Name", "Constraint", "Version", "Latest matching", "Latest"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	var rows [][]string
	for _, o := range output {
		var update = ""
		switch {
		case o.MatchingUpdate:
			update = "Y"
		case o.NonMatchingUpdate:
			update = "(Y)"
		case o.ConstraintUpdate:
			update = "?"
		}
		row := []string{update, o.Name, o.VersionConstraint, o.Version, o.Latest, o.LatestWithoutConstraint}
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		return fmt.Sprint(rows[i]) > fmt.Sprint(rows[j])
	})
	table.AppendBulk(rows)
	table.Render()
}
