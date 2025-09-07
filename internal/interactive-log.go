package internal

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2)

	statusStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1)
)

type InteractiveLogModel struct {
	table        table.Model
	entries      []Entry
	showHelp     bool
	searchMode   bool
	searchQuery  string
	filterMode   string
	width        int
	height       int
	pageSize     int
	currentPage  int
	totalEntries int
	totalPages   int
}

func NewInteractiveLogModel(pageSize int) InteractiveLogModel {
	columns := []table.Column{
		{Title: "date", Width: 16},
		{Title: "mood", Width: 12},
		{Title: "intensity", Width: 9},
		{Title: "message", Width: 30},
		{Title: "tags", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(15),
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

	return InteractiveLogModel{
		table:        t,
		entries:      []Entry{},
		showHelp:     false,
		searchMode:   false,
		filterMode:   "all",
		width:        80,
		height:       24,
		pageSize:     pageSize,
		currentPage:  0,
		totalEntries: 0,
		totalPages:   0,
	}
}

type tickMsg time.Time
type entriesLoadedMsg struct {
	entries      []Entry
	totalEntries int
	currentPage  int
}

func (m InteractiveLogModel) Init() tea.Cmd {
	return loadEntries(m.pageSize, m.currentPage, m.filterMode, m.searchQuery)
}

func (m InteractiveLogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		return m, loadEntries(m.pageSize, m.currentPage, m.filterMode, m.searchQuery)

	case entriesLoadedMsg:
		m.entries = msg.entries
		m.totalEntries = msg.totalEntries
		m.currentPage = msg.currentPage
		if m.totalEntries > 0 {
			// i don't like the repeated type conversion in math.Ceil
			m.totalPages = (m.totalEntries + m.pageSize - 1) / m.pageSize
		} else {
			m.totalPages = 0
		}
		m.updateTableRows()
		return m, nil

	case tea.KeyMsg:
		if m.searchMode {
			return m.handleSearchInput(msg)
		}
		return m.handleNormalInput(msg)
	}

	return m, cmd
}

func (m InteractiveLogModel) View() string {
	if m.showHelp {
		return m.helpView()
	}

	var s strings.Builder

	title := titleStyle.Render("🎭 moodgit interactive")
	stats := fmt.Sprintf("total: %d entries | filter: %s", m.totalEntries, m.filterMode)

	if m.totalPages > 0 {
		stats += fmt.Sprintf(" | page: %d/%d", m.currentPage+1, m.totalPages)
	}

	header := title + "  " + stats
	if m.searchMode {
		header += fmt.Sprintf(" | search: %s_", m.searchQuery)
	} else if m.searchQuery != "" {
		header += fmt.Sprintf(" | search: \"%s\"", m.searchQuery)
	}

	s.WriteString(header)
	s.WriteString("\n")

	s.WriteString(baseStyle.Render(m.table.View()))
	s.WriteString("\n")

	status := "↑/↓,j/k: navigate | q: quit | /: search | f: filter | r: refresh | ←/→: page | ?: help"
	s.WriteString(statusStyle.Width(m.width).Render(status))

	return s.String()
}

func loadEntries(pageSize int, page int, filter string, search string) tea.Cmd {
	return func() tea.Msg {
		offset := page * pageSize
		entries, totalCount, err := getFilteredHistory(pageSize, offset, filter, search)
		if err != nil {
			return entriesLoadedMsg{entries: []Entry{}, totalEntries: 0, currentPage: page}
		}
		return entriesLoadedMsg{entries: entries, totalEntries: totalCount, currentPage: page}
	}
}

func (m *InteractiveLogModel) updateTableRows() {
	rows := []table.Row{}
	for _, entry := range m.entries {
		intensityStr := fmt.Sprintf("%02d/10", entry.Intensity)

		tagsStr := ""
		if len(entry.Tags) > 0 {
			tagsStr = strings.Join(entry.Tags, ", ")
		}

		message := entry.Message
		if len(message) > 28 {
			message = message[:25] + "..."
		}

		rows = append(rows, table.Row{
			entry.CreatedAt.Format("2006/01/02 15:04"),
			entry.Mood,
			intensityStr,
			message,
			tagsStr,
		})
	}
	m.table.SetRows(rows)
}

func (m InteractiveLogModel) handleNormalInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "?":
		m.showHelp = !m.showHelp

	case "/":
		m.searchMode = true
		m.searchQuery = ""

	case "r":
		return m, loadEntries(m.pageSize, m.currentPage, m.filterMode, m.searchQuery)

	case "left", "h":
		if m.currentPage > 0 {
			m.currentPage--
			return m, loadEntries(m.pageSize, m.currentPage, m.filterMode, m.searchQuery)
		}

	case "right", "l":
		if m.currentPage < m.totalPages-1 {
			m.currentPage++
			return m, loadEntries(m.pageSize, m.currentPage, m.filterMode, m.searchQuery)
		}

	case "f":
		filters := []string{"all", "happy", "sad", "angry", "anxious", "excited", "calm", "stressed", "tired", "neutral"}
		for i, filter := range filters {
			if filter == m.filterMode {
				m.filterMode = filters[(i+1)%len(filters)]
				break
			}
		}
		m.currentPage = 0
		return m, loadEntries(m.pageSize, m.currentPage, m.filterMode, m.searchQuery)
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m InteractiveLogModel) handleSearchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.searchMode = false
		m.currentPage = 0
		return m, loadEntries(m.pageSize, m.currentPage, m.filterMode, m.searchQuery)

	case "esc":
		m.searchMode = false
		m.searchQuery = ""
		m.currentPage = 0
		return m, loadEntries(m.pageSize, m.currentPage, m.filterMode, m.searchQuery)

	case "backspace":
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
		}

	default:
		if len(msg.String()) == 1 {
			m.searchQuery += msg.String()
		}
	}

	return m, nil
}

func (m InteractiveLogModel) helpView() string {
	help := `🎭 moodgit interactive log - help

navigation:
  ↑/↓, j/k       move cursor up/down
  page up/down   page through entries (10 rows)
  home/end       go to first/last entry
  ←/→, h/l       previous/next page

actions:
  /              enter search mode
  f              cycle through mood filters
  r              refresh entries
  ?              toggle this help
  q, ctrl+c      quit

search mode:
  type to search in messages and tags
  enter          apply search
  esc            cancel search

pagination:
  page size: ` + fmt.Sprintf("%d", m.pageSize) + ` entries per page
  current: page ` + fmt.Sprintf("%d/%d", m.currentPage+1, m.totalPages) + `
  total entries: ` + fmt.Sprintf("%d", m.totalEntries) + `

current filter: ` + m.filterMode + `
current search: "` + m.searchQuery + `"

press ? again to return to the table view.`

	return helpStyle.Render(help)
}

func StartInteractiveLog(pageSize int) error {
	p := tea.NewProgram(
		NewInteractiveLogModel(pageSize),
		tea.WithAltScreen(),
	)

	_, err := p.Run()
	return err
}
