package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	word       string
	guesses    [6]textinput.Model
	turn       int
	letter_map map[rune]bool
	won        bool
}

func initialModel() model {
	var tis [6]textinput.Model
	for i := 0; i < 6; i++ {
		ti := textinput.New()
		ti.Placeholder = "Guess"
		ti.Focus()
		ti.CharLimit = 6
		ti.Width = 6
		ti.BackgroundStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4"))
		tis[i] = ti
	}
	word := "weedle"
	letter_map := make(map[rune]bool)
	for _, r := range word {
		letter_map[r] = true
	}
	return model{
		word:       word,
		guesses:    tis,
		letter_map: letter_map,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.Type {

		// These keys should exit the program.
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.guesses[m.turn].Value() == m.word {
				m.won = true
				return m, nil
			}
			m.turn++
			if m.turn == 6 {
				return m, tea.Quit
			}
			return m, nil
		}
	}

	m.guesses[m.turn], cmd = m.guesses[m.turn].Update(msg)

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, cmd
}

func (m model) View() string {
	// styles
	base_style := lipgloss.NewStyle().
		Bold(true).
		Background(lipgloss.Color("#7D56F4"))

	var style = lipgloss.NewStyle().
		Width(28).
		Foreground(lipgloss.Color("#FAFAFA")).
		Padding(1, 1, 1, 1).
		Align(lipgloss.Center).
		Inherit(base_style)

	var green = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#77DD77")).
		Inherit(base_style)
	var red = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6961")).
		Inherit(base_style)
	var yellow = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FDFD96")).
		Inherit(base_style)

	// The header
	s := fmt.Sprintf("Guess the 6-letter word.\nGuesses left: %d", 6-m.turn) + "\n\n"

	// Iterate over guesses
	for i := 0; i < 6; i++ {
		// Render the row accordingly
		if m.turn == i {
			s += fmt.Sprintf("%s\n", m.guesses[i].View())
		} else {
			guess := m.guesses[i].Value()
			cs := ""
			for i := 0; i < len(guess); i++ {
				guess_runes := []rune(guess)
				_, ok := m.letter_map[guess_runes[i]]
				if m.word[i] == guess[i] {
					cs += green.Render(fmt.Sprintf("%c", guess[i]))
				} else if ok {
					cs += yellow.Render(fmt.Sprintf("%c", guess[i]))
				} else {
					cs += red.Render(fmt.Sprintf("%c", guess[i]))
				}
			}
			s += cs + "\n"
		}

	}

	// The footer
	if m.won {
		s += "You wonnered!\nPress ctrl-c to quit."
	} else {
		s += "\nPress ctrl-c to give up."
	}

	// Send the UI for rendering
	return style.Render(s)
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
