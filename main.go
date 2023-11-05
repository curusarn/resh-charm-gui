package main

// A simple program that opens the alternate screen buffer then counts down
// from 5 and then exits.

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	ticker int

	init       bool
	mouseEvent tea.MouseEvent

	textInput textinput.Model
	err       error
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		ticker: 7,

		textInput: ti,
		err:       nil,
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

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := message.(type) {
	case tea.MouseMsg:
		m.init = true
		m.mouseEvent = tea.MouseEvent(msg)
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		default:
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

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
	) + "\n"

	return s
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
