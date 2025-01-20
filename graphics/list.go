package graphics

import (
	"strings"

	vn "github.com/azr4e1/venv-notary"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type cobraFunc func(*cobra.Command, []string) error

type ListModel struct {
	notary           vn.Notary
	showGlobal       bool
	showLocal        bool
	showVersion      bool
	itemStyle        lg.Style
	currentItemStyle lg.Style
	boxStyle         lg.Style
}

func (lm ListModel) Init() tea.Cmd {
	return nil
}

func (lm ListModel) View() string {
	// global := ""
	// local := ""
	// var output string
	// if lm.showGlobal {
	// 	globalVenvs := lm.notary.ListGlobal()
	// }
	if lm.showGlobal {

	}
	localOutput := strings.Join(lm.notary.ListLocal(), "\n")
	globalOutput := strings.Join(lm.notary.ListGlobal(), "\n")
	output := lg.JoinHorizontal(lg.Center, globalOutput, localOutput)
	return output
}

func (lm ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyCtrlC || msg.String() == "q":
			return lm, tea.Quit
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
		notary:     notary,
		showGlobal: globalVenv,
		showLocal:  localVenv,
	}
	return lm, nil
}

func ListMain(localVenv, globalVenv bool) cobraFunc {
	return func(cmd *cobra.Command, args []string) error {
		m, err := newListModel(localVenv, globalVenv)
		if err != nil {
			return err
		}
		p := tea.NewProgram(m)
		_, err = p.Run()

		return err
	}
}
