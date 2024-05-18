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

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// boilerplate for table view
var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			exepath := "C:\\Windows\\system32\\notepad.exe"
			file := fmt.Sprintf("..\\your_journals\\%s", m.table.SelectedRow()[1])
			cmd := exec.Command(exepath, file)
			err := cmd.Start() // Non-blocking program run
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return m, nil
			}
			return m, tea.Batch(
				tea.Printf("Opening up %s in Notepad...", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

type FileDetail struct {
	Name         string
	LastModified string
	Tags         string
}

func main() {
	/* take in command line arguments
	User inputs:
		ledgers -> gives alphabetically sorted table view list of all journals
		ledgers -m -> gives table sorted in order descending order of date
	*/
	entries, _ := os.ReadDir("../your_journals")
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

			// fetch location where all journals are stored
			f, _ := os.Open("../config/init.txt")
			scanner := bufio.NewScanner(f)
			scanner.Scan()
			scanner.Scan()
			location := scanner.Text()

			tags := extractTags(e.Name(), location + "\\")
			fileDetails = append(fileDetails, FileDetail{
				Name:         e.Name(),
				LastModified: info.ModTime().Format("2006-01-02"),
				Tags:         tags,
			})
		}
	}

	// fmt.Println(fileDetails)

	columns := []table.Column{
		{Title: "Date", Width: 15},
		{Title: "Journals", Width: 30},
		{Title: "Tags", Width: 20},
	}

	if len(os.Args) == 1 { // sorted by alphabetical order
		
	} else if len(os.Args) == 2 && os.Args[1] == "-m" {
		// Sort by date
		sort.Slice(fileDetails[:], func(i, j int) bool {
			ti, _ := time.Parse("2006-01-02", fileDetails[i].LastModified)
			tj, _ := time.Parse("2006-01-02", fileDetails[j].LastModified)
			return ti.After(tj)
		})
	}

	rows := []table.Row{}
	for _, file := range fileDetails {
		rows = append(rows, table.Row{file.LastModified, file.Name, file.Tags})
	}

	// initialize table, style it and populate it
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
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

	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}


// helper functions
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
	lines := strings.Split(string(rawBytes), "\n")
	firstline := lines[0]
	if strings.Contains(firstline, "tags:"){
		firstlineSliced := strings.Split(firstline, ":")
		return strings.ReplaceAll(firstlineSliced[1], " ", "")
	}

	return "No Tags"
}