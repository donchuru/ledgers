package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sort"
	"time"
	"bufio"
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/donchuru/ledgers/ledgers_mac/header"
	"github.com/sahilm/fuzzy"
)

var location string

// boilerplate for table view
var baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

// the search box takes every typed character, so the table only moves with keys that can't be typed
var tableKeys = table.KeyMap{
	LineUp:   key.NewBinding(key.WithKeys("up")),
	LineDown: key.NewBinding(key.WithKeys("down")),
	PageUp:   key.NewBinding(key.WithKeys("pgup")),
	PageDown: key.NewBinding(key.WithKeys("pgdown")),
}

type model struct {
	table    table.Model
	search   textinput.Model
	journals []FileDetail // every journal in display order; the table shows the ones matching the search
}

func newModel(journals []FileDetail) model {
	columns := []table.Column{
		{Title: "Date", Width: 15},
		{Title: "Journals", Width: 30},
		{Title: "Tags", Width: 20},
	}

	// initialize table, style it and populate it
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(toRows(journals)),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithKeyMap(tableKeys),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	search := textinput.New()
	search.Prompt = "> "
	search.Placeholder = "Search journals"
	search.Focus()

	return model{table: t, search: search, journals: journals}
}

func (m model) Init() tea.Cmd { return textinput.Blink }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.search.Value() == "" {
				return m, tea.Quit
			}
			m.search.Reset()
			m.applyFilter()
			return m, nil
		case "enter":
			row := m.table.SelectedRow()
			if row == nil { // nothing matches the search
				return m, nil
			}
			exepath := "/usr/bin/open"
			file := fmt.Sprintf(location + "/%s", row[1])
			cmd := exec.Command(exepath, "-a", "TextEdit", file)
			err := cmd.Start() // Non-blocking program run
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return m, nil
			}
			return m, tea.Batch(
				tea.Printf("Opening up '%s' in TextEdit...", row[1]),
			)
		}
	}

	// both get every message: the table only reacts to navigation keys, the search box to typing
	query := m.search.Value()
	var tableCmd, searchCmd tea.Cmd
	m.table, tableCmd = m.table.Update(msg)
	m.search, searchCmd = m.search.Update(msg)
	if m.search.Value() != query {
		m.applyFilter()
	}
	return m, tea.Batch(tableCmd, searchCmd)
}

// applyFilter shows the journals matching the search and selects the best match
func (m *model) applyFilter() {
	m.table.SetRows(toRows(filterJournals(m.journals, m.search.Value())))
	m.table.SetCursor(0)
}

func (m model) View() string {
	return baseStyle.Render(m.search.View()+"\n"+m.table.View()) + "\n" +
		helpStyle.Render("↑/↓ move · enter open · esc clear/quit · ctrl+c quit") + "\n"
}

type FileDetail struct {
	Name         string
	LastModified string
	Tags         string
}

// error checker
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	/* take in command line arguments
	User inputs:
		ledgers -> gives alphabetically sorted table view list of all journals
		ledgers -m -> gives table sorted in order descending order of date
	*/

	// Resolve the home directory
	homeDir, err := os.UserHomeDir()
	check(err)

	// Define the configuration file path
	configFilePath := filepath.Join(homeDir, ".ledgers_config", "init.txt")

	f, _ := os.Open(configFilePath)
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	scanner.Scan()
	location = scanner.Text()

	entries, _ := os.ReadDir(location)
	// var items []list.Item
	// for _, e := range entries {
	// 	items = append(items, item(e.Name()))
	// }

	// fmt.Println(entries)

	var fileDetails []FileDetail
	for _, e := range entries {
		if !e.IsDir() {
			info, err := e.Info()
			if err != nil {
				fmt.Printf("Error getting info for file %s: %s\n", e.Name(), err)
				continue
			}
			// fmt.Println(e.Name())

			tags := extractTags(e.Name(), location + "//")
			fileDetails = append(fileDetails, FileDetail{
				Name:         e.Name(),
				LastModified: info.ModTime().Format("2006-01-02"),
				Tags:         tags,
			})
		}
	}

	// fmt.Println(fileDetails)

	if len(os.Args) == 1 { // sorted by alphabetical order
		
	} else if len(os.Args) == 2 && os.Args[1] == "-m" {
		// Sort by date
		sort.Slice(fileDetails[:], func(i, j int) bool {
			ti, _ := time.Parse("2006-01-02", fileDetails[i].LastModified)
			tj, _ := time.Parse("2006-01-02", fileDetails[j].LastModified)
			return ti.After(tj)
		})
	}

	m := newModel(fileDetails)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}


// helper functions

func toRows(journals []FileDetail) []table.Row {
	rows := make([]table.Row, 0, len(journals))
	for _, j := range journals {
		tags := j.Tags
		if tags == "" {
			tags = "No Tags"
		}
		rows = append(rows, table.Row{j.LastModified, j.Name, tags})
	}
	return rows
}

// filterJournals returns the journals whose title or tags fuzzy-match query, best match first.
// An empty query returns journals in their original order.
func filterJournals(journals []FileDetail, query string) []FileDetail {
	if query == "" {
		return journals
	}
	// titles fill [0, n) and tags fill [n, 2n); matches come best first, so a journal's first match is its best
	n := len(journals)
	targets := make([]string, 2*n)
	for i, j := range journals {
		targets[i], targets[n+i] = j.Name, j.Tags
	}
	seen := make(map[int]bool)
	var matched []FileDetail
	for _, match := range fuzzy.Find(query, targets) {
		if i := match.Index % n; !seen[i] {
			seen[i] = true
			matched = append(matched, journals[i])
		}
	}
	return matched
}

// extractTags returns the journal's tags comma-separated, or "" if it has none
func extractTags(filename string, location_appendage string) string {
	// logic for extracting tags from the file itself
	fileIO, err := os.OpenFile(location_appendage + filename, os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer fileIO.Close()
	rawBytes, err := io.ReadAll(fileIO)
	if err != nil {
		panic(err)
	}
	h, _, ok := header.Parse(string(rawBytes))
	if !ok || len(h.Tags) == 0 {
		return ""
	}
	return strings.Join(h.Tags, ", ")
}