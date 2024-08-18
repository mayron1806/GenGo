package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	defaultMaxLength int = 156
	defaultMinLength int = 0
	defaultWidth     int = 20
)

type InputText struct {
	TextInput textinput.Model
	Err       string
	label     string
	Quitting  bool
}

func (m InputText) Init() tea.Cmd {
	return textinput.Blink
}
func (m InputText) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.Quitting = true
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.Err = msg.Error()
		return m, nil
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}
func (m InputText) View() string {
	if m.Quitting {
		return quitTextStyle.Render("Good Bye!")
	}
	var str string
	str += fmt.Sprintf("%s\n", titleStyle.Render(m.label))
	str += itemStyle.Render(m.TextInput.View())
	str += fmt.Sprintf("%s\n", infoTextStyle.Render("Press enter to submit or ctrl+c to quit"))
	if m.Err != "" {
		str += errorTextStyle.Render(m.Err)
	}
	return str
}

type InputTextOptions struct {
	Placeholder string
	Width       *int
	Label       string
	MinLength   int
	MaxLength   int
}

func NewInputText(options InputTextOptions) *InputText {

	if options.Width == nil {
		options.Width = &defaultWidth
	}
	ti := textinput.New()
	ti.Focus()
	if options.Placeholder != "" {
		ti.Placeholder = options.Placeholder
	}
	ti.Width = *options.Width
	ti.CharLimit = options.MaxLength
	ti.Validate = func(s string) error {
		if len(s) < options.MinLength || len(s) > options.MaxLength {
			return fmt.Errorf("Please enter between %d and %d characters", options.MinLength, options.MaxLength)
		}
		return nil
	}
	return &InputText{
		TextInput: ti,
		Err:       "",
		label:     options.Label,
	}
}
