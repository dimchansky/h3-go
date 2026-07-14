package main

import (
	"reflect"
	"testing"
)

func TestAPIMapFormsClassifiesPrimaryAndAdditive(t *testing.T) {
	tests := []struct {
		name       string
		cell       string
		primary    []string
		additional []string
	}{
		{
			name:    "plain function and method stay primary",
			cell:    "LatLngToCell; LatLng.Cell",
			primary: []string{"LatLngToCell", "LatLng.Cell"},
		},
		{
			name:       "Append method is additive",
			cell:       "Cell.GridDisk; Cell.AppendGridDisk",
			primary:    []string{"Cell.GridDisk"},
			additional: []string{"Cell.AppendGridDisk"},
		},
		{
			name:       "Append package function is additive",
			cell:       "CompactCells; AppendCompactCells",
			primary:    []string{"CompactCells"},
			additional: []string{"AppendCompactCells"},
		},
		{
			name:       "Seq and Grouped suffixes are additive",
			cell:       "Cell.GridDiskDistances; Cell.AppendGridDiskDistances; Cell.GridDiskDistancesGrouped",
			primary:    []string{"Cell.GridDiskDistances"},
			additional: []string{"Cell.AppendGridDiskDistances", "Cell.GridDiskDistancesGrouped"},
		},
		{
			name:       "Seq package function is additive",
			cell:       "PolygonToCellsExperimental; AppendPolygonToCellsExperimental; PolygonToCellsExperimentalSeq",
			primary:    []string{"PolygonToCellsExperimental"},
			additional: []string{"AppendPolygonToCellsExperimental", "PolygonToCellsExperimentalSeq"},
		},
		{
			name:       "classification uses the method name, not the receiver",
			cell:       "Cell.Children; Cell.AppendChildren; Cell.ChildrenSeq; Cell.ImmediateChildren; Cell.AppendImmediateChildren",
			primary:    []string{"Cell.Children", "Cell.ImmediateChildren"},
			additional: []string{"Cell.AppendChildren", "Cell.ChildrenSeq", "Cell.AppendImmediateChildren"},
		},
		{
			name:    "absorbed annotation stays primary verbatim",
			cell:    "— (garbage collected)",
			primary: []string{"— (garbage collected)"},
		},
		{
			name:    "empty entries are dropped",
			cell:    "Cell.Vertex; ",
			primary: []string{"Cell.Vertex"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			primary, additional := apiMapForms(tt.cell)
			if !reflect.DeepEqual(primary, tt.primary) {
				t.Errorf("primary = %v; want %v", primary, tt.primary)
			}
			if !reflect.DeepEqual(additional, tt.additional) {
				t.Errorf("additional = %v; want %v", additional, tt.additional)
			}
		})
	}
}

func TestCodeifyList(t *testing.T) {
	tests := []struct {
		name    string
		entries []string
		want    string
	}{
		{"empty list renders as a long dash", nil, "—"},
		{"symbols are backticked and comma-joined", []string{"Cell.Children", "AppendChildren"}, "`Cell.Children`, `AppendChildren`"},
		{"absorbed annotation stays readable", []string{"— (sentinel error messages)"}, "— (sentinel error messages)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := codeifyList(tt.entries); got != tt.want {
				t.Errorf("codeifyList(%v) = %q; want %q", tt.entries, got, tt.want)
			}
		})
	}
}
