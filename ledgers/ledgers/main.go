package main

import(
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		exepath := "C:\\Windows\\system32\\notepad.exe"
    	file := "..\\your_journals\\" + m.choice
    	cmd := exec.Command(exepath, file)
		err := cmd.Start() // non-blocking program run
		if err != nil {
			return fmt.Sprintf("Error: %s", err)
		}
	}
	if m.quitting {
		return quitTextStyle.Render("Come back when you're ready to open your heart")
	}
	return "\n" + m.list.View()
}


func main () {
	/* take in command line arguments
	User inputs:
		ledgers -> make a new ledger named today's date
		ledger "new Doc"  -> make a new ledger named new Doc
	*/
	entries, _ := os.ReadDir("../your_journals")
	var items []list.Item
	for _, e := range entries {
		items = append( items, item(e.Name()) )
	}

	if len(os.Args) == 1 {
		const defaultWidth = 20

		// initializing the TUI
		l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
		l.Title = "Here are your journals so far:"
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(false)
		l.Styles.Title = titleStyle
		l.Styles.PaginationStyle = paginationStyle
		l.Styles.HelpStyle = helpStyle

		m := model{list: l}

		if _, err := tea.NewProgram(m).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}


	} else if len(os.Args) == 2 {
		if os.Args[1] == "-m"{
			// TODO: show me list of all journals in order of last modified
			for _, e := range entries {
				fmt.Println(e.Name())
			}

		} else if os.Args[1] == "-c" {
			// TODO: show me list of all journals in order of last created
			for _, e := range entries {
				fmt.Println(e.Name())
			}
		}
	}
}
