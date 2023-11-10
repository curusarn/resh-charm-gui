package main

// A simple program that opens the alternate screen buffer then counts down
// from 5 and then exits.

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	ticker int

	init       bool
	mouseEvent tea.MouseEvent

	textInput textinput.Model
	err       error

	table table.Model

	height int
	width  int

	data *DataHolder
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	data := NewDataHolder()
	return model{
		ticker: 7,

		textInput: ti,
		err:       nil,

		table: data.GetInitialTable(),

		data: data,
	}
}

type tickMsg time.Time

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

func (m model) Init() tea.Cmd {
	// return tea.Batch(tick(), tea.EnterAltScreen)
	return textinput.Blink
	// return nil
}

func (m model) handleKeyUpdate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.Type {
	case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
		return m, tea.Quit
	case tea.KeyUp, tea.KeyDown:
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	default:
		m.textInput, cmd = m.textInput.Update(msg)
		val := m.textInput.Value()
		m.table.SetRows(m.data.GetRows(val))
		return m, cmd
	}
}

func (m model) handleWindowSizeUpdate(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.height = msg.Height
	m.width = msg.Width

	cmdWidth := msg.Width - 4 - 10 - 6 - 2
	tableHeight := msg.Height - 16
	columns := []table.Column{
		{Title: "Time", Width: 4},
		{Title: "Directory", Width: 10},
		{Title: "Command", Width: cmdWidth},
	}
	m.table.SetColumns(columns)
	m.table.SetHeight(tableHeight)
	return m, tea.ClearScreen
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSizeUpdate(msg)

	case tea.MouseMsg:
		m.init = true
		m.mouseEvent = tea.MouseEvent(msg)
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyUpdate(msg)

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil

		// case tickMsg:
		// 	m.ticker--
		// 	if m.ticker <= 0 {
		// 		return m, tea.Quit
		// 	}
		// 	return m, tick()

	}
	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf(
		"\n\n"+
			"    Hi. This program will exit in %d seconds...\n\n"+
			"    Do mouse stuff. When you're done press q to quit.\n\n",
		m.ticker)

	if m.init {
		e := m.mouseEvent
		s += fmt.Sprintf("    (X: %d, Y: %d) %s\n\n", e.X, e.Y, e)
	}
	s += fmt.Sprintf(
		"    What’s your favorite Pokémon?\n\n    %s\n\n    %s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n\n" +
		baseStyle.Render(m.table.View()) + "\n\n" +
		fmt.Sprintf("    Let's go to %s! (selected)", m.table.SelectedRow()[1]) + "\n\n"

	return s
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (d *DataHolder) GetInitialTable() table.Model {
	t := table.New(
		table.WithColumns(d.GetColumns(80)),
		table.WithRows(d.GetRows("")),
		table.WithFocused(true),
		table.WithHeight(20),
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

	return t
}

type DataRow struct {
	Time      time.Time
	Directory string
	Command   string
}

type DataHolder struct {
	Rows []DataRow
}

func (d *DataHolder) GetColumns(windowWidth int) []table.Column {
	cmdWidth := windowWidth - 4 - 10 - 6 - 2
	columns := []table.Column{
		{Title: "Time", Width: 4},
		{Title: "Directory", Width: 10},
		{Title: "Command", Width: cmdWidth},
	}
	return columns
}

func ToTableRow(r DataRow) table.Row {
	return table.Row{"now", r.Directory, r.Command}
}

func (d *DataHolder) GetRows(query string) []table.Row {
	rows := []table.Row{}
	for _, row := range d.Rows {
		if len(query) == 0 || strings.Contains(row.Command, query) {
			rows = append(rows, ToTableRow(row))
		}
	}
	return rows
}

func NewDataHolder() *DataHolder {
	rows := []DataRow{
		{Time: time.Now(), Directory: "~/git/betterstack", Command: "git push --force"},
		{Time: time.Now(), Directory: "~/git/betterstack", Command: "git commit"},
		{Time: time.Now(), Directory: "~/git/betterstack", Command: "git rebase"},
		{Time: time.Now(), Directory: "~/git/betterstack", Command: "git stash"},
		{Time: time.Now(), Directory: "~/git/betterstack", Command: "git pull"},
		{Time: time.Now(), Directory: "~/git/logtail", Command: "git push"},
		{Time: time.Now(), Directory: "~/git/logtail", Command: "git commit"},
		{Time: time.Now(), Directory: "~/git/logtail", Command: "git rebase"},
		{Time: time.Now(), Directory: "~/git/logtail", Command: "git merge"},
		{Time: time.Now(), Directory: "~/git/logtail", Command: "git cherry-pick"},
		{Time: time.Now(), Directory: "~/git/uptime", Command: "git commit"},
		{Time: time.Now(), Directory: "~/git/uptime", Command: "git push --force"},
		{Time: time.Now(), Directory: "~/git/uptime", Command: "git merge"},
		{Time: time.Now(), Directory: "~/git/uptime", Command: "git push"},
	}

	return &DataHolder{
		Rows: rows,
	}
}
