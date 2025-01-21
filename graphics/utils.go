package graphics

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	vn "github.com/azr4e1/venv-notary"
	lg "github.com/charmbracelet/lipgloss"
)

type headerType int

const (
	globalHeader headerType = iota
	localHeader
	ignoreHeader
)

const (
	LocalName             = "Local Environments"
	GlobalName            = "Global Environments"
	ReplaceVersion string = "py"
	truncateRatio         = 0.5
	truncateChar          = "…"
)

func fillLine(header string, width int, lineStyle lg.Style) string {
	gap := lineStyle.Render(strings.Repeat(" ", max(0, width-lg.Width(header)-2)))
	header = lg.JoinHorizontal(lg.Bottom, header, gap)
	return header
}

func createLocalHeader(width int, activeStyle, inactiveStyle lg.Style) string {
	header := activeStyle.Render(LocalName)
	return fillLine(header, width, inactiveStyle)
}

func createGlobalHeader(width int, activeStyle, inactiveStyle lg.Style) string {
	header := activeStyle.Render(GlobalName)
	return fillLine(header, width, inactiveStyle)
}

func createActiveHeader(activeHeader headerType, contentWidth, rowWidth int, activeStyle, inactiveStyle lg.Style) string {
	globalStyle := inactiveStyle
	localStyle := inactiveStyle
	if activeHeader == globalHeader {
		globalStyle = activeStyle
	}
	if activeHeader == localHeader {
		localStyle = activeStyle
	}
	header := lg.JoinHorizontal(
		lg.Top,
		globalStyle.Render(GlobalName),
		localStyle.Render(LocalName),
	)
	for i := 0; i < min(len(LocalName), len(GlobalName)) && lg.Width(header) > rowWidth; i++ {
		header = lg.JoinHorizontal(
			lg.Top,
			globalStyle.Render(GlobalName[:len(GlobalName)-1-i]+truncateChar),
			localStyle.Render(LocalName[:len(LocalName)-1-i]+truncateChar),
		)
	}
	return fillLine(header, contentWidth, inactiveStyle)
}

func printGlobal(notary vn.Notary, width int, itemStyle, currentItemStyle lg.Style) string {
	items := make(map[string][]string)
	names := make(map[string]string)
	for _, name := range notary.ListGlobal() {
		clnName, version := vn.ExtractVersion(filepath.Base(name))
		version = strings.ReplaceAll(version, ReplaceVersion, "")
		oldVersions, ok := items[clnName]
		if !ok {
			oldVersions = []string{}
		}
		oldVersions = append(oldVersions, version)
		items[clnName] = oldVersions
		names[name] = clnName
	}

	return prettyPrintList(notary, width, names, items, itemStyle, currentItemStyle)
}

func printLocal(notary vn.Notary, width int, itemStyle, currentItemStyle lg.Style) string {
	items := make(map[string][]string)
	names := make(map[string]string)
	for _, name := range notary.ListLocal() {
		clnName, version := vn.ExtractVersion(filepath.Base(name))
		version = strings.ReplaceAll(version, ReplaceVersion, "")
		clnNameWithHash := clnName
		clnName = vn.RemoveHash(clnName)
		hashVal := clnNameWithHash[len(clnName)+1:]
		clnName = fmt.Sprintf("%s-%s", clnName, hashVal[:4])
		oldVersions, ok := items[clnName]
		if !ok {
			oldVersions = []string{}
		}
		oldVersions = append(oldVersions, version)
		items[clnName] = oldVersions
		names[name] = clnName
	}

	return prettyPrintList(notary, width, names, items, itemStyle, currentItemStyle)
}

func prettyPrintList(notary vn.Notary, width int, nameMap map[string]string, items map[string][]string, itemStyle, currentItemStyle lg.Style) string {
	activeVenv, _ := notary.GetActiveEnv()
	activeName := nameMap[activeVenv.Path]
	_, activeVersion := vn.ExtractVersion(activeVenv.Path)
	activeVersion = strings.ReplaceAll(activeVersion, ReplaceVersion, "")

	names := []string{}
	for n := range items {
		names = append(names, n)
	}

	nameWidth := int(truncateRatio * float64(width))
	versionWidth := width - nameWidth
	nameBlock := prettyPrintEnv(names, nameWidth, activeName, itemStyle, currentItemStyle)
	versionBlock := prettyPrintVersion(names, versionWidth, items, activeName, activeVersion, itemStyle, currentItemStyle)
	return lg.JoinHorizontal(lg.Center, nameBlock, versionBlock)
}

func prettyPrintEnv(names []string, width int, activeName string, itemStyle, currentItemStyle lg.Style) string {
	slices.SortFunc(names, vn.AlphanumericSort)

	coloredNames := []string{}
	for _, n := range names {
		// check if needs to be truncated
		if len(n) > width && width > 0 {
			n = n[:width-1] + truncateChar
		}
		el := itemStyle.Render(n)
		if n == activeName {
			el = currentItemStyle.Render(n)
		}
		coloredNames = append(coloredNames, el)
	}
	nameBlock := venvBlockStyle.Render(strings.Join(coloredNames, "\n"))

	return nameBlock
}

func prettyPrintVersion(names []string, width int, items map[string][]string, activeName, activeVersion string, itemStyle, currentItemStyle lg.Style) string {
	versionBlockElements := []string{}
	for _, name := range names {
		versions := items[name]
		coloredVersions := []string{}
		slices.SortFunc(versions, vn.SemanticVersioningSort)
		for _, v := range versions {
			el := itemStyle.Render(v)
			if name == activeName && v == activeVersion {
				el = currentItemStyle.Render(v)
			}
			coloredVersions = append(coloredVersions, el)
		}
		// check if needs to be truncated
		versionElement := "(" + strings.Join(coloredVersions, " ") + ")"
		for i := 0; lg.Width(versionElement) > width && width > 0 && i < len(coloredVersions); i++ {
			versionElement = "(" + strings.Join(coloredVersions[:len(coloredVersions)-i-1], " ") + truncateChar + ")"
		}

		versionBlockElements = append(versionBlockElements, versionElement)
	}
	versionBlock := versionBlockStyle.Align(lg.Right).Render(strings.Join(versionBlockElements, "\n"))

	return versionBlock
}
