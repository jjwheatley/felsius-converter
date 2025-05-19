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
	conversionOpts    []string // List of supported conversions
	conversionChoice  string   // The conversion selected, inititally an empty string
	cursor            int      // which item the cursor is pointed at
	inputMeasurment   int      // measurement to convert
	outputMeasurement string   // measurement after conversion
	isCalculated      bool     // Are we done yet?
}

func initialModel() model {
	return model{
		conversionOpts:    []string{"Celsius to Felsius", "Fahrenheit to Felsius", "Felsius to Celsius", "Felsius to Fahrenheit"},
		conversionChoice:  "",
		inputMeasurment:   0,
		outputMeasurement: "",
		isCalculated:      false,
	}
}

func (m model) Init() tea.Cmd {
	return nil //`nil` = "no current I/O"
}

func isConversionChosen(m model) bool {
	return m.conversionChoice != ""
}

func celsiusToFahrenheit(degreesC float32) float32 {
	return (degreesC * 1.8) + 32.0
}

func celsiusToFelsius(degreesC float32) float32 {
	return (degreesC + celsiusToFahrenheit(degreesC)) / 2
}

func fahrenheitToCelsius(degreesF float32) float32 {
	return (degreesF - 32) / 1.8
}

func fahrenheitToFelsius(degreesF float32) float32 {
	return (fahrenheitToCelsius(degreesF) + degreesF) / 2
}

func felsiusToCelsius(degreesFelcius float32) float32 {
	return ((degreesFelcius - 16) * 5) / 7
}

func felsiusToFahrenheit(degreesFelcius float32) float32 {
	return ((degreesFelcius * 9) + 80) / 7
}

func convertTemp(m model) string {
	result := ""
	if m.conversionChoice == m.conversionOpts[0] {
		celcius := m.inputMeasurment
		felsius := celsiusToFelsius(float32(celcius))
		result = fmt.Sprintf("%s°c = %s °ϵ", strconv.Itoa(int(celcius)), strconv.Itoa(int(felsius)))

	} else if m.conversionChoice == m.conversionOpts[1] {
		fahrenheit := m.inputMeasurment
		felsius := fahrenheitToFelsius(float32(fahrenheit))
		result = fmt.Sprintf("%s°f = %s °ϵ", strconv.Itoa(int(fahrenheit)), strconv.Itoa(int(felsius)))

	} else if m.conversionChoice == m.conversionOpts[2] {
		felsius := m.inputMeasurment
		celsius := felsiusToCelsius(float32(m.inputMeasurment))
		result = fmt.Sprintf("%s°ϵ = %s °c", strconv.Itoa(int(felsius)), strconv.Itoa(int(celsius)))

	} else if m.conversionChoice == m.conversionOpts[3] {
		felsius := m.inputMeasurment
		fahrenheit := felsiusToFahrenheit(float32(m.inputMeasurment))
		result = fmt.Sprintf("%s°ϵ = %s °f", strconv.Itoa(int(felsius)), strconv.Itoa(int(fahrenheit)))
	}

	return result
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg: //Type check: was the input type a key?

		switch msg.String() { //Which key was pressed
		case "ctrl+c", "q":
			return m, tea.Quit //exit the program

		case "up", "k":
			if !isConversionChosen(m) {
				if m.cursor > 0 {
					m.cursor--
				}
			} else if m.outputMeasurement == "" { //Freeze input once output is calculated
				m.inputMeasurment++
			}

		case "down", "j":
			if !isConversionChosen(m) {
				if m.cursor < len(m.conversionOpts)-1 {
					m.cursor++
				}
			} else if m.outputMeasurement == "" { //Freeze input once output is calculated
				m.inputMeasurment--
			}

		case "enter", " ": //" " = spacebar
			if !isConversionChosen(m) {
				m.conversionChoice = m.conversionOpts[m.cursor]
			} else {
				m.outputMeasurement = convertTemp(m)
				m.isCalculated = true
				return m, tea.Quit //exit the program
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
			cursor := " " // no cursor
			if m.cursor == i {
				cursor = ">" // cursor!
			}

			optionNo := strconv.Itoa(i + 1)                         // Calculate option number
			s += fmt.Sprintf("%s %s. %s\n", cursor, optionNo, cOpt) // Render row
		}
	} else if !m.isCalculated {
		s += fmt.Sprintf("Enter measurement for conversion (from %s) using the arrow keys.\n\n", m.conversionChoice) // Header
		s += fmt.Sprintf("%s \n", strconv.Itoa(m.inputMeasurment))                                                   // Render current value of input
	} else {
		s += fmt.Sprintf("Result of %s\n\n", m.conversionChoice) //Header
		s += fmt.Sprintf("%s \n", m.outputMeasurement)           // Render current value of input
	}

	s += "\nPress q to quit.\n" //Footer
	return s                    // Send text to UI for rendering
}
