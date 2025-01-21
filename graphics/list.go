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

const (
	MaxHeight = 130
	MaxWidth  = 96
)

type cobraFunc func(*cobra.Command, []string) error

type ListModel struct {
	notary          vn.Notary
	showGlobal      bool
	showLocal       bool
	showVersion     string
	environmentType headerType
	windowWidth     int
	windowHeight    int
	ready           bool

	itemStyle        lg.Style
	currentItemStyle lg.Style
	activeTabStyle   lg.Style
	tabStyle         lg.Style
	tabGap           lg.Style

	viewport viewport.Model

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
	var header, content string
	header = lm.headerView()
	if lm.showGlobal == lm.showLocal {
		content = lm.viewport.View()
	} else {
		content = lm.contentView()
	}
	output := lg.JoinVertical(lg.Left, header, content)
	output = lg.NewStyle().Padding(1, 0).Render(output)
	return output
}

func (lm ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyCtrlC || msg.String() == "q":
			return lm, tea.Quit
		case msg.String() == "l":
			lm.local()
		case msg.String() == "g":
			lm.global()
		case msg.Type == tea.KeyTab || msg.Type == tea.KeyShiftTab || msg.Type == tea.KeyRight || msg.Type == tea.KeyLeft:
			lm.switchContent()
		case msg.String() == "r":
			lm.refresh()
		}
	case tea.WindowSizeMsg:
		headerHeight := lg.Height(lm.headerView())
		content := lm.contentView()
		contentHeight := lg.Height(content)
		windowHeight := msg.Height - headerHeight - 2
		lm.windowHeight = windowHeight
		lm.windowWidth = msg.Width

		if !lm.ready {
			lm.viewport = viewport.New(min(msg.Width, MaxWidth), min(windowHeight, MaxHeight, contentHeight))
			lm.viewport.YPosition = headerHeight + 2
			lm.viewport.SetContent(content)
			lm.ready = true
		} else {
			lm.ResetViewport()
		}
	}
	lm.viewport, cmd = lm.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return lm, tea.Batch(cmds...)
}

func (lm *ListModel) ResetViewport() {
	viewport := lm.viewport

	content := lm.contentView()
	contentHeight := lg.Height(content)

	viewport.Width = min(lm.windowWidth, MaxWidth)
	viewport.Height = min(lm.windowHeight, MaxHeight, contentHeight)
	viewport.SetContent(content)

	lm.viewport = viewport
}

func (lm *ListModel) switchContent() {
	if lm.environmentType == localHeader {
		lm.global()
	} else {
		lm.local()
	}
}

func (lm *ListModel) global() {
	lm.environmentType = globalHeader
	lm.ResetViewport()
}

func (lm *ListModel) local() {
	lm.environmentType = localHeader
	lm.ResetViewport()
}

func (lm *ListModel) refresh() {
	err := lm.notary.GetVenvs()
	if err != nil {
		return
	}
	localContent := printLocal(lm.notary, lm.itemStyle, lm.currentItemStyle)
	globalContent := printGlobal(lm.notary, lm.itemStyle, lm.currentItemStyle)
	localWidth := lg.Width(localContent)
	globalWidth := lg.Width(globalContent)
	localActiveHeader := createActiveHeader(localHeader, localWidth, lm.activeTabStyle, lm.tabStyle)
	globalActiveHeader := createActiveHeader(globalHeader, globalWidth, lm.activeTabStyle, lm.tabStyle)
	localOnlyHeader := createLocalHeader(localWidth, lm.activeTabStyle, lm.tabStyle)
	globalOnlyHeader := createGlobalHeader(globalWidth, lm.activeTabStyle, lm.tabStyle)

	lm.localHeader = localActiveHeader
	lm.localContent = localContent
	lm.globalHeader = globalActiveHeader
	lm.globalContent = globalContent
	lm.localOnlyHeader = localOnlyHeader
	lm.globalOnlyHeader = globalOnlyHeader

	if lm.environmentType == globalHeader {
		lm.global()
	} else {
		lm.local()
	}
}

func (lm ListModel) headerView() string {
	var header string
	if lm.showGlobal && !lm.showLocal {
		return lm.globalOnlyHeader
	}
	if !lm.showGlobal && lm.showLocal {
		return lm.localOnlyHeader
	}
	if lm.environmentType == globalHeader {
		header = lm.globalHeader
	} else {
		header = lm.localHeader
	}

	return header
}

func (lm ListModel) contentView() string {
	var content string
	if lm.showGlobal && !lm.showLocal {
		return lm.globalContent
	}
	if !lm.showGlobal && lm.showLocal {
		return lm.localContent
	}
	if lm.environmentType == globalHeader {
		content = lm.globalContent
	} else {
		content = lm.localContent
	}

	return content
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
	lm.refresh()
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
