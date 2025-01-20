package graphics

import (
	lg "github.com/charmbracelet/lipgloss"
)

// adpative colors
var (
	highlight         = lg.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveHighlight = lg.Color("242")
	itemColor         = lg.Color("15")
	currentItemColor  = lg.Color("214")
)

// item styles
var (
	itemStyle        = lg.NewStyle().Foreground(lg.Color("15"))
	currentItemStyle = lg.NewStyle().Italic(true).Foreground(lg.Color("214"))
	tab              = lg.NewStyle().
				Border(lg.NormalBorder(), false, false, true, false).
				BorderForeground(highlight).Foreground(inactiveHighlight).Padding(0, 1)
	activeTab         = tab.Foreground(highlight).Bold(true)
	venvBlockStyle    = lg.NewStyle().Padding(0, 1)
	versionBlockStyle = lg.NewStyle().PaddingLeft(1)
)

// header style
// var (
// 	activeTabBorder = lg.Border{
// 		Top:         "─",
// 		Bottom:      " ",
// 		Left:        "│",
// 		Right:       "│",
// 		TopLeft:     "╭",
// 		TopRight:    "╮",
// 		BottomLeft:  "┘",
// 		BottomRight: "└",
// 	}

// 	tabBorder = lg.Border{
// 		Top:         "─",
// 		Bottom:      "─",
// 		Left:        "│",
// 		Right:       "│",
// 		TopLeft:     "╭",
// 		TopRight:    "╮",
// 		BottomLeft:  "┴",
// 		BottomRight: "┴",
// 	}

// 	tab = lg.NewStyle().
// 		Border(tabBorder, true).
// 		BorderForeground(highlight).
// 		Padding(0, 1).Foreground(inactiveHighlight)

// 	activeTab = tab.Border(activeTabBorder, true).Bold(true).Foreground(highlight)

// 	tabGap = tab.
// 		BorderTop(false).
// 		BorderLeft(false).
// 		BorderRight(false)
// )
