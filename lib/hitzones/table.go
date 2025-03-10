package hitzones

import (
	"slices"

	"github.com/jgabriel1/mh-weakness-bot/lib/element"
)

type Column struct {
	Element element.Element
	Data    []int
}

type Table struct {
	Columns []*Column
}

func NewTable() *Table {
	return &Table{Columns: []*Column{}}
}

func (t *Table) AddColumn(el element.Element) int {
	t.Columns = append(t.Columns, &Column{Element: el, Data: []int{}})
	return len(t.Columns) - 1
}

func (t *Table) AddValueToColumn(colIndex int, value int) {
	if col := t.GetColumnAt(colIndex); col != nil {
		col.Data = append(col.Data, value)
	}
}

func (t *Table) GetColumnAt(index int) *Column {
	if len(t.Columns) < index+1 {
		return nil
	}
	return t.Columns[index]
}

func (t *Table) GetWeaknessElement() (element.Element, int) {
	var maximum int
	var el element.Element = element.Unknown
	for _, col := range t.Columns {
		colMaximum := slices.Max(col.Data)
		if colMaximum > maximum {
			maximum = colMaximum
			el = col.Element
		}
	}
	return el, maximum
}
