package main

import (
	"fmt"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	conversionOpts   []string         // List of supported conversions
	conversionChoice string           // The conversion selected, inititally an empty string
	cursor           int              // which to-do list item our cursor is pointing at
	selected         map[int]struct{} // which to-do items are selected
	inputMeasurment  int              // measurement to convert
}

func initialModel() model {
	return model{
		conversionOpts:   []string{"Celsius to Felsius", "Fahrenheit to Felsius", "Felsius to Celsius", "Felsius to Fahrenheit"},
		conversionChoice: "",
		inputMeasurment:  10,
		// A map which indicates which conversionOpts are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `conversionOpts` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if m.conversionChoice == "" {
		switch msg := msg.(type) {

		// Is it a key press?
		case tea.KeyMsg:

			// Cool, what was the actual key pressed?
			switch msg.String() {

			// These keys should exit the program.
			case "ctrl+c", "q":
				return m, tea.Quit

			// The "up" and "k" keys move the cursor up
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}

			// The "down" and "j" keys move the cursor down
			case "down", "j":
				if m.cursor < len(m.conversionOpts)-1 {
					m.cursor++
				}

			// The "enter" key and the spacebar (a literal space) toggle
			// the selected state for the item that the cursor is pointing at.
			case "enter", " ":
				// fmt.Print(m.conversionOpts[m.cursor])
				m.conversionChoice = m.conversionOpts[m.cursor]
			}
		}
	} else {
		switch msg := msg.(type) {

		// Is it a key press?
		case tea.KeyMsg:

			// Cool, what was the actual key pressed?
			switch msg.String() {

			// These keys should exit the program.
			case "ctrl+c", "q":
				return m, tea.Quit

			// The "up" and "k" keys move the cursor up
			case "up", "k":
				m.inputMeasurment++

			// The "down" and "j" keys move the cursor down
			case "down", "j":
				m.inputMeasurment--

			// The "enter" key and the spacebar (a literal space) toggle
			// the selected state for the item that the cursor is pointing at.
			case "enter", " ":
				// ToDo: Run conversion & display output
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := ""

	if m.conversionChoice == "" {
		// The header
		s += "What conversion do you want to do?\n\n"

		// Iterate over our conversionOpts
		for i, cOpt := range m.conversionOpts {

			// Is the cursor pointing at this conversion option (cOpt)?
			cursor := " " // no cursor
			if m.cursor == i {
				cursor = ">" // cursor!
			}

			// Calculate option number
			optionNo := strconv.Itoa(i + 1)
			// Render the row
			s += fmt.Sprintf("%s %s. %s\n", cursor, optionNo, cOpt)
		}

	} else {
		// The header
		s += fmt.Sprintf("Enter measurement for conversion (from %s)\n\n", m.conversionChoice)

		// Render current value of input
		s += fmt.Sprintf("%s \n", strconv.Itoa(m.inputMeasurment))
	}

	// The footer
	s += "\nPress q to quit.\n"
	// Send the UI for rendering
	return s
}
