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
	MaxWidth  = 120
)

type cobraFunc func(*cobra.Command, []string) error

type ListModel struct {
	notary          vn.Notary
	showGlobal      bool
	showLocal       bool
	pythonVersion   string
	environmentType headerType
	windowWidth     int
	windowHeight    int
	ready           bool
	error           error

	MaxHeight int
	MaxWidth  int

	itemStyle        lg.Style
	currentItemStyle lg.Style
	activeTabStyle   lg.Style
	tabStyle         lg.Style
	tabGap           lg.Style
	errorStyle       lg.Style

	viewport viewport.Model

	localHeader      string
	globalHeader     string
	localOnlyHeader  string
	globalOnlyHeader string
	localContent     string
	globalContent    string
}

func (lm ListModel) Init() tea.Cmd {
	if lm.error != nil {
		return func() tea.Msg { return errMsg(lm.error) }
	}
	return nil
}

func (lm ListModel) View() string {
	if lm.error != nil {
		return lm.errorStyle.Render(lm.error.Error() + "\n")
	}
	var header, content string
	header = lm.headerView()
	if lm.showGlobal == lm.showLocal {
		if !lm.ready {
			return "\nInitializing..."
		}
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
	case errMsg:
		return lm, tea.Quit
	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyCtrlC || msg.String() == "q" || msg.String() == "esc" || msg.String() == "enter":
			return lm, tea.Quit
		case msg.String() == "l":
			lm.Local()
		case msg.String() == "g":
			lm.Global()
		case msg.Type == tea.KeyTab || msg.Type == tea.KeyShiftTab || msg.Type == tea.KeyRight || msg.Type == tea.KeyLeft:
			lm.SwitchContent()
		case msg.String() == "r":
			lm.Refresh()
		}
	case tea.WindowSizeMsg:
		headerHeight := lg.Height(lm.headerView())
		windowHeight := msg.Height - headerHeight - 3
		lm.windowHeight = windowHeight
		lm.windowWidth = msg.Width

		if !lm.ready {
			lm.Refresh()
			content := lm.contentView()
			contentHeight := lg.Height(content)
			lm.viewport = viewport.New(min(msg.Width, lm.MaxWidth), min(windowHeight, lm.MaxHeight, contentHeight))
			lm.viewport.YPosition = headerHeight + 3
			lm.viewport.SetContent(content)
			lm.ready = true
		} else {
			lm.Refresh()
			lm.resetViewport()
		}
	}
	lm.viewport, cmd = lm.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return lm, tea.Batch(cmds...)
}

func (lm *ListModel) resetViewport() {
	viewport := lm.viewport

	content := lm.contentView()
	contentHeight := lg.Height(content)

	viewport.Width = min(lm.windowWidth, lm.MaxWidth)
	viewport.Height = min(lm.windowHeight, lm.MaxHeight, contentHeight)
	viewport.SetContent(content)

	lm.viewport = viewport
}

func (lm *ListModel) SwitchContent() {
	if lm.environmentType == localHeader {
		lm.Global()
	} else {
		lm.Local()
	}
}

func (lm *ListModel) Global() {
	lm.environmentType = globalHeader
	lm.resetViewport()
}

func (lm *ListModel) Local() {
	lm.environmentType = localHeader
	lm.resetViewport()
}

func (lm *ListModel) Refresh() {
	err := lm.notary.GetVenvs()
	if err != nil {
		return
	}
	width := min(lm.windowWidth, lm.MaxWidth) - 4 // account for padding
	localContent := printLocal(lm.notary, width, lm.pythonVersion, lm.itemStyle, lm.currentItemStyle)
	globalContent := printGlobal(lm.notary, width, lm.pythonVersion, lm.itemStyle, lm.currentItemStyle)
	localWidth := lg.Width(localContent)
	globalWidth := lg.Width(globalContent)
	localActiveHeader := createActiveHeader(localHeader, localWidth, width, lm.activeTabStyle, lm.tabStyle)
	globalActiveHeader := createActiveHeader(globalHeader, globalWidth, width, lm.activeTabStyle, lm.tabStyle)
	localOnlyHeader := createLocalHeader(localWidth, lm.activeTabStyle, lm.tabStyle)
	globalOnlyHeader := createGlobalHeader(globalWidth, lm.activeTabStyle, lm.tabStyle)

	lm.localHeader = localActiveHeader
	lm.localContent = localContent
	lm.globalHeader = globalActiveHeader
	lm.globalContent = globalContent
	lm.localOnlyHeader = localOnlyHeader
	lm.globalOnlyHeader = globalOnlyHeader

	if lm.environmentType == globalHeader {
		lm.Global()
	} else {
		lm.Local()
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

func newListModel(localVenv, globalVenv bool, pythonExec string) (tea.Model, error) {
	notary, err := vn.NewNotary()
	if err != nil {
		return ListModel{}, err
	}
	environmentType := globalHeader
	if localVenv && !globalVenv {
		environmentType = localHeader
	}

	var pythonVersion string
	if pythonExec != "" {
		pythonVersion, err = vn.PythonVersion(pythonExec)
	}
	lm := ListModel{
		notary:           notary,
		showGlobal:       globalVenv,
		showLocal:        localVenv,
		pythonVersion:    pythonVersion,
		environmentType:  environmentType,
		activeTabStyle:   activeTab,
		tabStyle:         tab,
		itemStyle:        itemStyle,
		currentItemStyle: currentItemStyle,
		MaxHeight:        MaxHeight,
		MaxWidth:         MaxWidth,
		errorStyle:       errorStyle,
		error:            err,
	}
	lm.Refresh()
	return lm, nil
}

func ListMain(localVenv, globalVenv *bool, pythonExec *string, stdout io.Writer) cobraFunc {
	return func(cmd *cobra.Command, args []string) error {
		m, err := newListModel(*localVenv, *globalVenv, *pythonExec)
		if err != nil {
			return err
		}
		if *localVenv != *globalVenv {
			fmt.Fprint(stdout, m.View())
			return nil
		}
		p := tea.NewProgram(m)
		_, err = p.Run()

		return err
	}
}
