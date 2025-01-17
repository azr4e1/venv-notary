package graphics

import (
	vn "github.com/azr4e1/venv-notary"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type ListModel struct {
	notary      vn.Notary
	showGlobal  bool
	showLocal   bool
	showVersion bool
	itemStyle   lg.Style
	boxStyle    lg.Style
}

func (lm ListModel) Init() tea.Cmd {
	return nil
}

func (lm ListModel) View() string {
	global := ""
	local := ""
	var output string
	if lm.showGlobal {
		globalVenvs := lm.notary.ListGlobal()
	}
	return output
}

func (lm ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return lm, nil
}
