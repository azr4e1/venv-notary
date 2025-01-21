package graphics

import (
	// "strings"

	"fmt"
	"io"

	vn "github.com/azr4e1/venv-notary"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type cobraFunc func(*cobra.Command, []string) error

type ListModel struct {
	notary           vn.Notary
	showGlobal       bool
	showLocal        bool
	showVersion      string
	environmentType  headerType
	itemStyle        lg.Style
	currentItemStyle lg.Style
	activeTabStyle   lg.Style
	tabStyle         lg.Style
	tabGap           lg.Style
	viewport         viewport.Model
}

func (lm ListModel) Init() tea.Cmd {
	return nil
}

func (lm ListModel) View() string {
	body := createBody(lm.notary, lm.showGlobal, lm.showLocal, lm.environmentType, lm.itemStyle, lm.currentItemStyle)
	width := lg.Width(body)
	header := createHeader(lm.showGlobal, lm.showLocal, lm.environmentType, width, lm.activeTabStyle, lm.tabStyle)
	output := lg.JoinVertical(lg.Left, header, body)
	return lg.NewStyle().Padding(1, 0).Render(output)
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
		notary:           notary,
		showGlobal:       globalVenv,
		showLocal:        localVenv,
		environmentType:  globalHeader,
		activeTabStyle:   activeTab,
		tabStyle:         tab,
		itemStyle:        itemStyle,
		currentItemStyle: currentItemStyle,
	}
	return lm, nil
}

func ListMain(localVenv, globalVenv *bool, stdout io.Writer) cobraFunc {
	return func(cmd *cobra.Command, args []string) error {
		m, err := newListModel(*localVenv, *globalVenv)
		if err != nil {
			return err
		}
		if *localVenv != *globalVenv {
			fmt.Fprintln(stdout, m.View())
			return nil
		}
		p := tea.NewProgram(m)
		_, err = p.Run()

		return err
	}
}
