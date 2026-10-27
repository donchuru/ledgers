package main

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

var journals = []FileDetail{
	{Name: "2024-05-17"},
	{Name: "on going back to work", Tags: "career, life"},
	{Name: "random.txt", Tags: "random"},
	{Name: "today.txt"},
	{Name: "dreams", Tags: "dream interpretations"},
}

func names(journals []FileDetail) []string {
	var out []string
	for _, j := range journals {
		out = append(out, j.Name)
	}
	return out
}

func rowNames(rows []table.Row) []string {
	var out []string
	for _, r := range rows {
		out = append(out, r[1])
	}
	return out
}

func TestFilterJournals(t *testing.T) {
	tests := []struct {
		name, query string
		want        []string
	}{
		{"empty query keeps order", "", names(journals)},
		{"title substring", "work", []string{"on going back to work"}},
		{"title abbreviation", "gbw", []string{"on going back to work"}},
		{"tag only", "career", []string{"on going back to work"}},
		{"ignores case", "WORK", []string{"on going back to work"}},
		{"best match first, once per journal", "ra", []string{"random.txt", "dreams"}},
		{"no match", "zzz", nil},
		{"untagged journals have no tags to match", "tags", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := names(filterJournals(journals, tt.query))

			if !slices.Equal(got, tt.want) {
				t.Errorf("filterJournals(%q) = %q, want %q", tt.query, got, tt.want)
			}
		})
	}
}

func TestExtractTags(t *testing.T) {
	dir := t.TempDir() + "/"
	os.WriteFile(filepath.Join(dir, "tagged"), []byte("---\nTitle: tagged\nTags: career, life\n---\nhi\n"), 0600)
	os.WriteFile(filepath.Join(dir, "untagged"), []byte("hi\n"), 0600)

	if got := extractTags("tagged", dir); got != "career, life" {
		t.Errorf("extractTags(tagged) = %q, want %q", got, "career, life")
	}
	if got := extractTags("untagged", dir); got != "" {
		t.Errorf("extractTags(untagged) = %q, want empty", got)
	}
}

func TestToRowsShowsPlaceholderForNoTags(t *testing.T) {
	rows := toRows([]FileDetail{{Name: "today.txt", LastModified: "2024-05-17"}})

	if want := (table.Row{"2024-05-17", "today.txt", "No Tags"}); !slices.Equal(rows[0], want) {
		t.Errorf("toRows = %q, want %q", rows[0], want)
	}
}

func update(m model, msg tea.Msg) (model, tea.Cmd) {
	next, cmd := m.Update(msg)
	return next.(model), cmd
}

func TestTypingFiltersAndEscClears(t *testing.T) {
	m := newModel(journals)

	m, _ = update(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("work")})
	if got := rowNames(m.table.Rows()); !slices.Equal(got, []string{"on going back to work"}) {
		t.Fatalf("rows after typing = %q", got)
	}

	m, _ = update(m, tea.KeyMsg{Type: tea.KeyEsc})
	if got := rowNames(m.table.Rows()); !slices.Equal(got, names(journals)) {
		t.Fatalf("rows after esc = %q, want all journals", got)
	}

	m, _ = update(m, tea.KeyMsg{Type: tea.KeyDown})
	if m.table.Cursor() != 1 {
		t.Errorf("cursor after down = %d, want 1", m.table.Cursor())
	}
}

func TestEscWithEmptySearchQuits(t *testing.T) {
	m := newModel(journals)

	_, cmd := update(m, tea.KeyMsg{Type: tea.KeyEsc})

	if cmd == nil {
		t.Fatal("esc returned no command, want tea.Quit")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Errorf("esc returned %T, want tea.QuitMsg", cmd())
	}
}

func TestEnterWithNoMatchesDoesNothing(t *testing.T) {
	m := newModel(journals)
	m, _ = update(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("zzz")})

	_, cmd := update(m, tea.KeyMsg{Type: tea.KeyEnter})

	if cmd != nil {
		t.Errorf("enter with no matches returned a command")
	}
}
