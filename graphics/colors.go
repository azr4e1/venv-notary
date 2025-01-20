package graphics

import (
	lg "github.com/charmbracelet/lipgloss"
)

var (
	activeHeaderStyle   = lg.NewStyle().Bold(true).Foreground(lg.Color("105"))
	inactiveHeaderStyle = lg.NewStyle().Foreground(lg.Color("242"))
	boxStyle            = lg.NewStyle().BorderStyle(lg.RoundedBorder())
	itemStyle           = lg.NewStyle().Foreground(lg.Color("15"))
	currentItemStyle    = lg.NewStyle().Italic(true).Foreground(lg.Color("214"))
)
