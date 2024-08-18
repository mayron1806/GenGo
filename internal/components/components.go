package components

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/mayron1806/gengo/config"

	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg error
)

var (
	titleStyle = lipgloss.NewStyle().Padding(0, 1).Background(lipgloss.Color("170")).Bold(true)

	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("62"))
	hoveredItemStyle  = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))

	quitTextStyle = lipgloss.NewStyle().Margin(1, 0, 2, 4)

	errorTextStyle = lipgloss.NewStyle().Margin(1, 0, 2, 4).Foreground(lipgloss.Color("1"))

	infoTextStyle = lipgloss.NewStyle().Margin(1, 0, 2, 4).Foreground(lipgloss.Color("3"))

	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

var logger = config.GetLogger()
