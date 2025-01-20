package graphics

import (
	// "strings"

	vn "github.com/azr4e1/venv-notary"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type cobraFunc func(*cobra.Command, []string) error

type ListModel struct {
	notary              vn.Notary
	showGlobal          bool
	showLocal           bool
	showVersion         bool
	itemStyle           lg.Style
	currentItemStyle    lg.Style
	boxStyle            lg.Style
	activeHeaderStyle   lg.Style
	inactiveHeaderStyle lg.Style
	environmentType     headerType
}

func (lm ListModel) Init() tea.Cmd {
	return nil
}

func (lm ListModel) View() string {
	header := createHeader(lm.showGlobal, lm.showLocal, lm.environmentType, lm.activeHeaderStyle, lm.inactiveHeaderStyle)
	var body string
	if lm.environmentType == globalHeader {
		body = printGlobal(lm.notary, lm.itemStyle, lm.currentItemStyle)
	} else {
		body = printLocal(lm.notary, lm.itemStyle, lm.currentItemStyle)
	}
	output := lg.JoinVertical(lg.Left, header, body)
	return output
}

func (lm ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyCtrlC || msg.String() == "q":
			return lm, tea.Quit
		case msg.String() == "l":
			lm.environmentType = localHeader
		case msg.String() == "g":
			lm.environmentType = globalHeader
		case msg.Type == tea.KeyTab:
			if lm.environmentType == localHeader {
				lm.environmentType = globalHeader
			} else {
				lm.environmentType = localHeader
			}
		case msg.Type == tea.KeyShiftTab:
			if lm.environmentType == localHeader {
				lm.environmentType = globalHeader
			} else {
				lm.environmentType = localHeader
			}
		}
	}
	return lm, nil
}

func newListModel(localVenv, globalVenv bool) (tea.Model, error) {
	notary, err := vn.NewNotary()
	if err != nil {
		return ListModel{}, err
	}
	lm := ListModel{
		notary:              notary,
		showGlobal:          globalVenv,
		showLocal:           localVenv,
		environmentType:     globalHeader,
		activeHeaderStyle:   activeHeaderStyle,
		inactiveHeaderStyle: inactiveHeaderStyle,
		itemStyle:           itemStyle,
		currentItemStyle:    currentItemStyle,
	}
	return lm, nil
}

func ListMain(localVenv, globalVenv *bool) cobraFunc {
	return func(cmd *cobra.Command, args []string) error {
		m, err := newListModel(*localVenv, *globalVenv)
		if err != nil {
			return err
		}
		p := tea.NewProgram(m)
		_, err = p.Run()

		return err
	}
}
