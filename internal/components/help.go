package components

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type keyMap map[string]key.Binding

var (
	defaultHelpSecureMarginSize = 10
	defaultHelpKey              = key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help"))
)

func (k keyMap) ShortHelp() []key.Binding {
	keys := make([]key.Binding, 0, len(k))
	for _, binding := range k {
		keys = append(keys, binding)
	}
	return keys
}

func (k keyMap) FullHelp() [][]key.Binding {
	keys := make([][]key.Binding, 0)
	row := make([]key.Binding, 0, 4)
	for _, binding := range k {
		row = append(row, binding)
		if len(row) == 4 {
			keys = append(keys, row)
			row = make([]key.Binding, 0, 4)
		}
	}
	if len(row) > 0 {
		keys = append(keys, row)
	}
	return keys
}

type Help struct {
	keys     keyMap
	help     help.Model
	quitting bool
	width    int
}

func (m *Help) Toggle() {
	m.help.ShowAll = !m.help.ShowAll
}
func (m Help) Init() tea.Cmd {
	return nil
}

func (m Help) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
		m.width = msg.Width
	case tea.KeyMsg:
		if key.Matches(msg, m.keys["quit"]) {
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Help) View() string {
	if m.quitting {
		return ""
	}
	var keys keyMap = make(map[string]key.Binding)
	if !m.help.ShowAll { // show only the help key
		keys["help"] = defaultHelpKey
	} else { // show all keys
		keys = m.keys
	}
	helpView := m.help.View(keys)
	return helpView
}

func NewHelp(keys keyMap) Help {
	return Help{
		keys: keys,
		help: help.New(),
	}
}
