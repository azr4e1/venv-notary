package graphics

import (
	lg "github.com/charmbracelet/lipgloss"
)

var (
	headerActive   = lg.NewStyle().Italic(true).Foreground(lg.Color("105"))
	headerInactive = lg.NewStyle().Foreground(lg.Color("254"))
	boxStyle       = lg.NewStyle().BorderStyle(lg.RoundedBorder())
)
