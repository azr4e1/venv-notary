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
	ReplaceVersion string = "py"
)

func createHeader(showGlobal, showLocal bool, activeHeader headerType, activeStyle, inactiveStyle lg.Style) string {
	localName := "Local Environments"
	globalName := "Global Environments"
	if showLocal && !showGlobal {
		return activeStyle.Render(localName)
	}
	if showGlobal && !showLocal {
		return activeStyle.Render(globalName)
	}
	globalStyle := inactiveStyle
	localStyle := inactiveStyle
	if activeHeader == globalHeader {
		globalStyle = activeStyle
	}
	if activeHeader == localHeader {
		localStyle = activeStyle
	}
	header := strings.Join([]string{
		globalStyle.Render(globalName),
		localStyle.Render(localName),
	}, " | ")
	return header
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
			names[name] = clnName
		}
		oldVersions = append(oldVersions, version)
		items[clnName] = oldVersions
	}

	return prettyPrint(notary, names, items, itemStyle, currentItemStyle)
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
			names[name] = clnName
		}
		oldVersions = append(oldVersions, version)
		items[clnName] = oldVersions
	}

	return prettyPrint(notary, names, items, itemStyle, currentItemStyle)
}

func prettyPrint(notary vn.Notary, nameMap map[string]string, items map[string][]string, itemStyle, currentItemStyle lg.Style) string {
	names := []string{}
	coloredNames := []string{}
	activeVenv, _ := notary.GetActiveEnv()
	activeName := nameMap[activeVenv.Path]
	_, activeVersion := vn.ExtractVersion(activeVenv.Path)
	activeVersion = strings.ReplaceAll(activeVersion, ReplaceVersion, "")
	for _, n := range nameMap {
		names = append(names, n)
	}
	slices.SortFunc(names, vn.AlphanumericSort)
	for _, n := range names {
		el := itemStyle.Render(n)
		if n == activeName {
			el = currentItemStyle.Render(n)
		}
		coloredNames = append(coloredNames, el)
	}
	nameBlock := lg.NewStyle().PaddingRight(1).Render(strings.Join(coloredNames, "\n"))

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
	versionBlock := lg.NewStyle().PaddingLeft(1).Render(strings.Join(versionBlockElements, "\n"))
	return lg.JoinHorizontal(lg.Center, nameBlock, versionBlock)
}
