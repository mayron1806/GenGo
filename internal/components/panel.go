package components

import tea "github.com/charmbracelet/bubbletea"

type Panel struct {
	Views []tea.Model
}

func (p Panel) Init() tea.Cmd {
	return nil
}

func (p Panel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return p, tea.Quit
		}
	}
	return p, nil
}

func (p Panel) View() string {
	return ""
}

type PanelOption struct {
	Views int
}

func NewPanel(opt PanelOption) Panel {
	views := make([]tea.Model, opt.Views)
	return Panel{
		Views: views,
	}
}
