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
	localHeader      string
	globalHeader     string
	localOnlyHeader  string
	globalOnlyHeader string
	localContent     string
	globalContent    string
}

func (lm ListModel) Init() tea.Cmd {
	return nil
}

func (lm ListModel) View() string {
	header, content := lm.HeaderContent()
	output := lg.JoinVertical(lg.Left, header, content)
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
		case msg.String() == "r":
			lm.Refresh()
		}
	}
	return lm, nil
}

func (lm *ListModel) Refresh() {
	err := lm.notary.GetVenvs()
	if err != nil {
		return
	}
	localContent := createBody(lm.notary, true, true, localHeader, lm.itemStyle, lm.currentItemStyle)
	globalContent := createBody(lm.notary, true, true, globalHeader, lm.itemStyle, lm.currentItemStyle)
	localWidth := lg.Width(localContent)
	globalWidth := lg.Width(globalContent)
	localHead := createHeader(true, true, localHeader, localWidth, lm.activeTabStyle, lm.tabStyle)
	globalHead := createHeader(true, true, globalHeader, globalWidth, lm.activeTabStyle, lm.tabStyle)
	localOnlyHead := createHeader(false, true, localHeader, localWidth, lm.activeTabStyle, lm.tabStyle)
	globalOnlyHead := createHeader(true, false, globalHeader, globalWidth, lm.activeTabStyle, lm.tabStyle)

	lm.localHeader = localHead
	lm.localContent = localContent
	lm.globalHeader = globalHead
	lm.globalContent = globalContent
	lm.localOnlyHeader = localOnlyHead
	lm.globalOnlyHeader = globalOnlyHead
}

func (lm ListModel) HeaderContent() (string, string) {
	var header, content string
	if lm.showGlobal && !lm.showLocal {
		return lm.globalOnlyHeader, lm.globalContent
	}
	if !lm.showGlobal && lm.showLocal {
		return lm.localOnlyHeader, lm.localContent
	}
	if lm.environmentType == globalHeader {
		header, content = lm.globalHeader, lm.globalContent
	} else {
		header, content = lm.localHeader, lm.localContent
	}

	return header, content
}

func newListModel(localVenv, globalVenv bool) (tea.Model, error) {
	notary, err := vn.NewNotary()
	if err != nil {
		return ListModel{}, err
	}
	environmentType := globalHeader
	if localVenv && !globalVenv {
		environmentType = localHeader
	}
	lm := ListModel{
		notary:           notary,
		showGlobal:       globalVenv,
		showLocal:        localVenv,
		environmentType:  environmentType,
		activeTabStyle:   activeTab,
		tabStyle:         tab,
		itemStyle:        itemStyle,
		currentItemStyle: currentItemStyle,
	}
	lm.Refresh()
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
