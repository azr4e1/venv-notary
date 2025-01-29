package graphics

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type errMsg error
type doneMsg int
type actionFunc func(*cobra.Command, []string) func() error

type StatusModel struct {
	spinner        spinner.Model
	waitingMessage string
	exitMessage    string
	errorMessage   string
	errorStyle     lg.Style
	quitting       bool
	action         func() error
	maxWidth       int
	width          int
	ready          bool
}

func (sm StatusModel) Init() tea.Cmd {
	cobraCmd := func() tea.Msg {
		err := sm.action()
		if err != nil {
			return errMsg(err)
		}
		return doneMsg(0)
	}
	cmds := []tea.Cmd{sm.spinner.Tick, cobraCmd}
	return tea.Batch(cmds...)
}

func (sm StatusModel) View() string {
	if sm.errorMessage != "" {
		return errorStyle.Render(truncateLine(sm.errorMessage, sm.width) + "\n")
	}
	if !sm.ready {
		return "\n" + truncateLine("Initializing...", sm.width)
	}
	str := truncateLine(fmt.Sprintf("%s %s\n", sm.spinner.View(), sm.waitingMessage), sm.width)
	if sm.quitting {
		return truncateLine(sm.exitMessage, sm.width) + "\n"
	}
	return str
}

func (sm StatusModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		if msg != nil {
			sm.errorMessage = msg.Error()
		}
		sm.quitting = true
		return sm, tea.Quit
	case doneMsg:
		sm.quitting = true
		return sm, tea.Quit
	case tea.WindowSizeMsg:
		sm.width = min(sm.maxWidth, msg.Width)
		sm.ready = true
		return sm, nil
	default:
		var cmd tea.Cmd
		sm.spinner, cmd = sm.spinner.Update(msg)
		return sm, cmd
	}
}

func newStatus(waitingMessage, exitMessage string, action func() error) StatusModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	sm := StatusModel{
		waitingMessage: waitingMessage,
		exitMessage:    exitMessage,
		action:         action,
		spinner:        s,
		maxWidth:       MaxWidth,
	}
	return sm
}

func StatusMain(waitingMessage, exitMessage string, action actionFunc, setup cobraFunc) cobraFunc {
	return func(cmd *cobra.Command, args []string) error {
		if setup != nil {
			err := setup(cmd, args)
			if err != nil {
				return err
			}
		}
		m := newStatus(waitingMessage, exitMessage, action(cmd, args))
		p := tea.NewProgram(m)
		_, err := p.Run()

		return err
	}
}
