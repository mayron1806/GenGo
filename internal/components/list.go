package components

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	defaultListHeight = 14
	defaultListWidth  = 20
)

type listItem string

func (i listItem) FilterValue() string { return "" }

type listItemDelegate struct{}

func (d listItemDelegate) Height() int                             { return 1 }
func (d listItemDelegate) Spacing() int                            { return 0 }
func (d listItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d listItemDelegate) Render(w io.Writer, m list.Model, index int, li list.Item) {
	i, ok := li.(listItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type List struct {
	List     list.Model
	Choice   string
	Quitting bool
}

func (m List) Init() tea.Cmd {
	return nil
}
func (m List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.Quitting = true
			return m, tea.Quit

		case tea.KeyEnter:
			item := m.List.SelectedItem()
			m.Choice = item.FilterValue()
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}
func (m List) View() string {
	if m.Choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("Selected: %s", m.Choice))
	}
	if m.Quitting {
		return quitTextStyle.Render("Bye!")
	}
	return "\n" + m.List.View()
}

type ListOptions struct {
	Title  string
	Items  []string
	Width  *int
	Height *int
}

func NewList(options ListOptions) *List {
	if options.Height == nil {
		options.Height = &defaultListHeight
	}
	if options.Width == nil {
		options.Width = &defaultListWidth
	}
	items := make([]list.Item, len(options.Items))
	for index, it := range options.Items {
		items[index] = listItem(it)
	}
	l := list.New(items, listItemDelegate{}, *options.Width, *options.Height)
	l.Title = options.Title
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle

	return &List{
		List: l,
	}
}
