package ui

import (
	table "github.com/calyptia/go-bubble-table"
	"github.com/charmbracelet/lipgloss"
)

var (
	style = lipgloss.NewStyle().Padding(1)
)

type Table struct {
	table table.Model
}

func NewTable(header []string, data [][]interface{}, w, h int) (Table, error) {
	t := Table{}

	top, right, bottom, left := style.GetPadding()
	w = w - left - right
	h = h - top - bottom

	tbl := table.New(header, w, h)

	var rows []table.Row
	for _, d := range data {
		var r table.SimpleRow
		t := append(r, append(table.SimpleRow{}, d...))
		rows = append(rows, t)
	}
	tbl.SetRows(rows)

	t.table = tbl

	return t, nil
}

func (t Table) Init() {}

func (t Table) Update(msg string) {}

func (t Table) View() string {
	return style.Render(t.table.View())
}
