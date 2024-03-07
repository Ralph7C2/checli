package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle           = lipgloss.NewStyle().Padding(1, 2)
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type model struct {
	list list.Model
}

func newModel(items []list.Item) *model {
	itemList := list.New(items, newItemDelegate(), 0, 0)
	return &model{
		list: itemList,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	return appStyle.Render(m.list.View())
}

func main() {

	var items []list.Item
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		items = append(items, item{item: scanner.Text()})
	}

	if _, err := tea.NewProgram(newModel(items), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

type item struct {
	item string
}

func (i item) Title() string {
	return i.item
}

func (i item) Description() string {
	return ""
}

func (i item) FilterValue() string {
	return i.item
}

func newItemDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var selectedItem string

		if i, ok := m.SelectedItem().(item); ok {
			selectedItem = i.item
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
				index := m.Index()
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					return tea.Quit
				}
				return m.NewStatusMessage(statusMessageStyle("Removed " + selectedItem))
			}
		}
		return nil
	}

	return d
}
