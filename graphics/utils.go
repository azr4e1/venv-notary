package graphics

import (
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
	LocalName  = "Local Environments"
	GlobalName = "Global Environments"
)

const (
	ReplaceVersion string = "py"
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

func createActiveHeader(activeHeader headerType, width int, activeStyle, inactiveStyle lg.Style) string {
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
	return fillLine(header, width, inactiveStyle)
}

func printGlobal(notary vn.Notary, itemStyle, currentItemStyle lg.Style) string {
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

	return prettyPrintList(notary, names, items, itemStyle, currentItemStyle)
}

func printLocal(notary vn.Notary, itemStyle, currentItemStyle lg.Style) string {
	items := make(map[string][]string)
	names := make(map[string]string)
	for _, name := range notary.ListLocal() {
		clnName, version := vn.ExtractVersion(filepath.Base(name))
		version = strings.ReplaceAll(version, ReplaceVersion, "")
		clnName = vn.RemoveHash(clnName)
		oldVersions, ok := items[clnName]
		if !ok {
			oldVersions = []string{}
		}
		oldVersions = append(oldVersions, version)
		items[clnName] = oldVersions
		names[name] = clnName
	}

	return prettyPrintList(notary, names, items, itemStyle, currentItemStyle)
}

func prettyPrintList(notary vn.Notary, nameMap map[string]string, items map[string][]string, itemStyle, currentItemStyle lg.Style) string {
	activeVenv, _ := notary.GetActiveEnv()
	activeName := nameMap[activeVenv.Path]
	_, activeVersion := vn.ExtractVersion(activeVenv.Path)
	activeVersion = strings.ReplaceAll(activeVersion, ReplaceVersion, "")

	names := []string{}
	for n := range items {
		names = append(names, n)
	}

	nameBlock := prettyPrintEnv(names, activeName, itemStyle, currentItemStyle)
	versionBlock := prettyPrintVersion(names, items, activeName, activeVersion, itemStyle, currentItemStyle)
	return lg.JoinHorizontal(lg.Center, nameBlock, versionBlock)
}

func prettyPrintEnv(names []string, activeName string, itemStyle, currentItemStyle lg.Style) string {
	slices.SortFunc(names, vn.AlphanumericSort)

	coloredNames := []string{}
	for _, n := range names {
		el := itemStyle.Render(n)
		if n == activeName {
			el = currentItemStyle.Render(n)
		}
		coloredNames = append(coloredNames, el)
	}
	nameBlock := venvBlockStyle.Render(strings.Join(coloredNames, "\n"))

	return nameBlock
}

func prettyPrintVersion(names []string, items map[string][]string, activeName, activeVersion string, itemStyle, currentItemStyle lg.Style) string {
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
		versionBlockElements = append(versionBlockElements, "("+strings.Join(coloredVersions, " ")+")")
	}
	versionBlock := versionBlockStyle.Render(strings.Join(versionBlockElements, "\n"))

	return versionBlock
}
