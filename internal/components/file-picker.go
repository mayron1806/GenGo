package components

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var filePickerHelpKeys keyMap = keyMap{
	"up":       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "move up")),
	"down":     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "move down")),
	"select":   key.NewBinding(key.WithKeys("+"), key.WithHelp("+", "select dir/file")),
	"unselect": key.NewBinding(key.WithKeys("-"), key.WithHelp("-", "unselect dir/file")),
	"opendir":  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "open dir")),
	"goback":   key.NewBinding(key.WithKeys("backspace"), key.WithHelp("backspace", "go back")),
	"quit":     key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q", "quit")),
	"help":     key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
}

type itemType uint8

const (
	itemFile itemType = iota
	itemFolder
)

type filePickerItem struct {
	title    string
	itemType itemType
	path     string
	selected bool
}

func (i filePickerItem) Title() string       { return i.title }
func (i filePickerItem) FilterValue() string { return i.title }
func (i filePickerItem) ItemType() itemType  { return i.itemType }

type filePickerItemDelegate struct{}

func (d filePickerItemDelegate) Height() int                               { return 1 }
func (d filePickerItemDelegate) Spacing() int                              { return 0 }
func (d filePickerItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d filePickerItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(filePickerItem)
	if !ok {
		return
	}
	var icon string
	if i.itemType == itemFolder {
		icon = "📁"
	} else {
		icon = "📄"
	}
	str := fmt.Sprintf("%s %s", icon, i.title)
	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			if i.selected {
				return selectedItemStyle.Render("> " + strings.Join(s, " "))
			} else {
				return hoveredItemStyle.Render("> " + strings.Join(s, " "))
			}
		}
	} else if i.selected {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("  " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(str))
}

type selectedFilePickerItemDelegate struct{}

func (d selectedFilePickerItemDelegate) Height() int                               { return 1 }
func (d selectedFilePickerItemDelegate) Spacing() int                              { return 0 }
func (d selectedFilePickerItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d selectedFilePickerItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(filePickerItem)
	if !ok {
		return
	}
	var icon string
	if i.itemType == itemFolder {
		icon = "📁"
	} else {
		icon = "📄"
	}
	str := fmt.Sprintf("%s %s", icon, i.path)
	fn := itemStyle.Render

	fmt.Fprint(w, fn(str))
}

type FilePicker struct {
	RootDir      string
	currentDir   string
	list         list.Model
	selectedList list.Model
	width        int
	height       int
	help         Help
}

func (m *FilePicker) setWidth(width int) {
	m.width = width
}
func (m *FilePicker) getWidth() int {
	return m.width
}
func (m *FilePicker) getBlockWidth() int {
	return m.width/2 - 2
}
func (m *FilePicker) setHeight(height int) {
	m.height = height
}
func (m *FilePicker) getHeight() int {
	return m.height
}
func (m *FilePicker) getBlockHeight() int {
	offset := 4
	return m.height - offset
}

func (m FilePicker) Init() tea.Cmd {
	return nil
}
func (m FilePicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setWidth(msg.Width)
		m.setHeight(msg.Height)
		m.list.SetSize(m.getBlockWidth(), m.getBlockHeight()-4)
		m.selectedList.SetSize(m.getBlockWidth(), m.getBlockHeight()-4)
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.help.keys["quit"]):
			return m, tea.Quit
		case key.Matches(msg, m.help.keys["help"]):
			m.help.Toggle()
		case key.Matches(msg, m.help.keys["opendir"]):
			if m.currentItem().itemType == itemFolder {
				m.enterDir(m.currentItem().title)
			}
		case key.Matches(msg, m.help.keys["goback"]):
			m.goBack()
		case key.Matches(msg, m.help.keys["select"]):
			m.addCurrentToList()
		case key.Matches(msg, m.help.keys["unselect"]):
			m.removeCurrentFromList()
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	m.selectedList, _ = m.selectedList.Update(msg) // Atualiza o selectedList
	return m, cmd
}
func (m FilePicker) View() string {
	colunmStyle := lipgloss.
		NewStyle().
		Width(m.getBlockWidth()).
		Height(m.getBlockHeight()).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("170"))

	var filePicker string
	filePicker += m.list.View()
	filePicker += "\n"
	filePicker += m.help.View()
	filePicker = colunmStyle.Render(filePicker)

	var selectedFiles string
	selectedFiles += m.selectedList.View() // Renderiza a lista selecionada
	selectedFiles = colunmStyle.Render(selectedFiles)

	main := lipgloss.
		JoinHorizontal(lipgloss.Left, filePicker, selectedFiles)

	return main
}

func (m *FilePicker) enterDir(dirName string) {
	newDir := filepath.Join(m.currentDir, dirName)
	files, _ := os.ReadDir(newDir)
	m.currentDir = newDir
	m.loadFiles(files)
}
func (m *FilePicker) goBack() {
	newDir := filepath.Dir(m.currentDir)
	files, _ := os.ReadDir(newDir)
	m.currentDir = newDir
	m.loadFiles(files)
}
func (m FilePicker) currentItem() filePickerItem {
	return m.list.Items()[m.list.Index()].(filePickerItem)
}
func (m *FilePicker) addCurrentToList() {
	currentItemPath := filepath.Join(m.currentDir, m.currentItem().title)
	alreadySelected := false
	for _, item := range m.selectedList.Items() {
		if item.(filePickerItem).path == currentItemPath {
			alreadySelected = true
			break
		}
	}

	if alreadySelected {
		return
	}
	m.selectedList.InsertItem(len(m.selectedList.Items()), filePickerItem{
		title:    m.currentItem().title,
		itemType: m.currentItem().ItemType(),
		path:     currentItemPath,
		selected: true,
	})
	for i, item := range m.list.Items() {
		logger.Debug("selected: ", item.(filePickerItem).path)
		if item.(filePickerItem).path == currentItemPath {
			m.list.SetItem(i, filePickerItem{
				title:    m.currentItem().title,
				itemType: m.currentItem().ItemType(),
				path:     currentItemPath,
				selected: true,
			})
			break
		}
	}
}
func (m *FilePicker) removeCurrentFromList() {
	currentItemPath := filepath.Join(m.currentDir, m.currentItem().title)

	for i, item := range m.selectedList.Items() {
		if item.(filePickerItem).path == currentItemPath {
			m.selectedList.RemoveItem(i)
			break
		}
	}
	for i, item := range m.list.Items() {
		if item.(filePickerItem).path == currentItemPath {
			m.list.SetItem(i, filePickerItem{
				title:    m.currentItem().title,
				itemType: m.currentItem().ItemType(),
				path:     currentItemPath,
				selected: false,
			})
			break
		}
	}

}

func (m *FilePicker) loadFiles(files []fs.DirEntry) {
	items := make([]list.Item, len(files))
	for i, file := range files {
		itemType := itemFile
		if file.IsDir() {
			itemType = itemFolder
		}
		selected := false
		for _, item := range m.selectedList.Items() {
			if item.(filePickerItem).path == filepath.Join(m.currentDir, file.Name()) {
				selected = true
				break
			}
		}
		items[i] = filePickerItem{
			title:    file.Name(),
			itemType: itemType,
			path:     filepath.Join(m.currentDir, file.Name()),
			selected: selected,
		}
	}
	m.list.SetItems(items)
}

func NewFilePicker(path string) FilePicker {
	filePicker := FilePicker{
		RootDir: path,
		help:    NewHelp(filePickerHelpKeys),
	}

	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	items := make([]list.Item, len(files))
	for i, file := range files {
		itemType := itemFile
		if file.IsDir() {
			itemType = itemFolder
		}
		items[i] = filePickerItem{
			title:    file.Name(),
			itemType: itemType,
			path:     filepath.Join(path, file.Name()),
			selected: false,
		}
	}

	filesList := list.New(items, filePickerItemDelegate{}, 20, 14)
	filesList.SetShowHelp(false)
	filesList.SetShowStatusBar(false)
	filesList.SetFilteringEnabled(false)
	filesList.Styles.Title = titleStyle
	filesList.Styles.PaginationStyle = paginationStyle
	filesList.Styles.HelpStyle = helpStyle
	filePicker.list = filesList

	selectedFilesList := list.New(make([]list.Item, 0), selectedFilePickerItemDelegate{}, 20, 14)
	selectedFilesList.Title = "Selected Files"
	selectedFilesList.SetShowHelp(false)
	selectedFilesList.SetShowStatusBar(false)
	selectedFilesList.SetFilteringEnabled(false)
	selectedFilesList.Styles.Title = titleStyle
	selectedFilesList.Styles.PaginationStyle = paginationStyle
	selectedFilesList.Styles.HelpStyle = helpStyle
	filePicker.selectedList = selectedFilesList
	return filePicker
}
