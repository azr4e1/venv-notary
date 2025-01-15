package graphics

import (
	vn "github.com/azr4e1/venv-notary"
	tea "github.com/charmbracelet/bubbletea"
)

type ListModel struct {
	notary      vn.Notary
	showGlobal  bool
	showLocal   bool
	showVersion bool
}

func (lm ListModel) Init() tea.Cmd {
	return nil
}

func (lm ListModel) View() string {
	return ""
}

func (lm ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return lm, nil
}
