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
	conversionOpts    []string         // List of supported conversions
	conversionChoice  string           // The conversion selected, inititally an empty string
	cursor            int              // which to-do list item our cursor is pointing at
	selected          map[int]struct{} // which to-do items are selected
	inputMeasurment   int              // measurement to convert
	outputMeasurement string           // measurement after conversion
	isCalculated      bool             // Are we done yet?
}

func initialModel() model {
	return model{
		conversionOpts:    []string{"Celsius to Felsius", "Fahrenheit to Felsius", "Felsius to Celsius", "Felsius to Fahrenheit"},
		conversionChoice:  "",
		inputMeasurment:   -42,
		outputMeasurement: "",
		// A map which indicates which conversionOpts are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `conversionOpts` slice, above.
		selected:     make(map[int]struct{}),
		isCalculated: false,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func isConversionChosen(m model) bool {
	return m.conversionChoice != ""
}

func celsiusToFahrenheit(degreesC float32) float32 {
	return (degreesC * 1.8) + 32.0
}

func fahrenheitToCelsius(degreesF float32) float32 {
	return (degreesF - 32) / 1.8
}

func convertTemp(m model) string {
	result := ""
	if m.conversionChoice == m.conversionOpts[0] {
		celsius := float32(m.inputMeasurment)
		fahrenheit := celsiusToFahrenheit(celsius)
		felsius := (celsius + fahrenheit) / 2

		result = fmt.Sprintf("%s%s", strconv.Itoa(int(felsius)), "°ϵ")

	} else if m.conversionChoice == m.conversionOpts[1] {
		fahrenheit := float32(m.inputMeasurment)
		celsius := fahrenheitToCelsius(fahrenheit)
		felsius := (celsius + fahrenheit) / 2

		result = fmt.Sprintf("%s%s", strconv.Itoa(int(felsius)), "°ϵ")

	} else if m.conversionChoice == m.conversionOpts[2] {
		// Felsius to Celsius
		felsius := float32(m.inputMeasurment)
		celsius := ((felsius - 16) * 5) / 7
		result = fmt.Sprintf("%s%s", strconv.Itoa(int(celsius)), "°c")

	} else if m.conversionChoice == m.conversionOpts[3] {
		// Felsius to Fahrenheit
		felsius := float32(m.inputMeasurment)
		fahrenheit := ((felsius * 9) + 80) / 7
		result = fmt.Sprintf("%s%s", strconv.Itoa(int(fahrenheit)), "°c")

	} else {
		//ToDo: Throw an error: "Bad conversion option passed to convertTemp"
	}

	return result
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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
			if !isConversionChosen(m) {
				if m.cursor > 0 {
					m.cursor--
				}
			} else {
				m.inputMeasurment++
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if !isConversionChosen(m) {
				if m.cursor < len(m.conversionOpts)-1 {
					m.cursor++
				}
			} else {
				m.inputMeasurment--
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			// fmt.Print(m.conversionOpts[m.cursor])
			if !isConversionChosen(m) {
				m.conversionChoice = m.conversionOpts[m.cursor]
			} else {
				m.outputMeasurement = convertTemp(m)
				m.isCalculated = true
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

	} else if m.isCalculated == false {
		// The header
		s += fmt.Sprintf("Enter measurement for conversion (from %s)\n\n", m.conversionChoice)

		// Render current value of input
		s += fmt.Sprintf("%s \n", strconv.Itoa(m.inputMeasurment))
	} else {
		// The header
		s += fmt.Sprintf("Result of %s\n\n", m.conversionChoice)

		// Render current value of input
		s += fmt.Sprintf("%s \n", m.outputMeasurement)
	}

	// The footer
	s += "\nPress q to quit.\n"
	// Send the UI for rendering
	return s
}
