/*
** FLICK PROJECT, 2026
** flick/internal/cli/commands/explore/explore_view
** File description:
** Rendering of the group explorer screens: groups list, folder tree and local
** file picker, plus the shared status line helper.
 */

package explore

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// View: Render the active screen.
//
// Returns:
// - result1 (string): The rendered view.
func (m exploreModel) View() string {
	switch m.mode {
	case modeGroups:
		return m.viewGroups()
	case modePicker:
		return m.viewPicker()
	default:
		return m.viewTree()
	}
}

// visibleRange: Compute the visible slice of items centred around the cursor
// so the help bar and status line always fit on screen.
//
// Params:
//   - count (int): Total number of items.
//   - cursor (int): Current cursor position in the item list.
//   - reserve (int): Number of terminal lines reserved for UI chrome (title,
//     help bar, status line, etc.).
//
// Returns:
// - start (int): First index of the visible slice (inclusive).
// - end (int): Last index of the visible slice (exclusive).
func (m exploreModel) visibleRange(count, cursor, reserve int) (start, end int) {
	if m.height == 0 || count <= reserve {
		return 0, count
	}
	max := m.height - reserve
	if max < 1 {
		max = 1
	}
	if count <= max {
		return 0, count
	}
	half := max / 2
	start = cursor - half
	if start < 0 {
		start = 0
	}
	end = start + max
	if end > count {
		end = count
		start = count - max
	}
	return
}

// viewGroups: Render the groups selection screen.
//
// Returns:
// - result1 (string): The rendered view.
func (m exploreModel) viewGroups() string {
	out := exploreTitleStyle.Render("flick - your groups") + "\n\n"
	if len(m.groups) == 0 {
		out += exploreCrumbStyle.Render("You don't belong to any group.") + "\n"
	}
	reserve := 4 // title(1) + blank(1) + blank-before-help(1) + help(1)
	if m.status != "" {
		reserve++ // status line
	}
	start, end := m.visibleRange(len(m.groups), m.groupCursor, reserve)
	for i := start; i < end; i++ {
		g := m.groups[i]
		style := lipgloss.NewStyle()
		if i == m.groupCursor {
			style = style.Bold(true)
		}
		out += "  " + style.Render(fmt.Sprintf("%s (%s)", g.Name, g.Role)) + "\n"
	}
	out += "\n" + exploreHelpStyle.Render("↑/↓ move · → open group · q quit")
	return appendStatus(out, m.status)
}

// viewTree: Render the folder tree screen.
//
// Returns:
// - result1 (string): The rendered view.
func (m exploreModel) viewTree() string {
	out := exploreTitleStyle.Render("flick - "+m.groupName) + "\n\n"

	if m.creating {
		// Layout: title(1) + blank(1) + items + blank(1) + prompt(1) + help(1) = 5
		start, end := m.visibleRange(len(m.rows), m.cursor, 5)
		for i := start; i < end; i++ {
			row := m.rows[i]
			name := row.node.name
			if row.node.isFolder {
				name += "/"
			}
			style := lipgloss.NewStyle()
			if row.node.isFolder {
				style = style.Foreground(exploreFolderClr)
			}
			if i == m.cursor {
				style = style.Bold(true)
			}
			out += exploreTreeStyle.Render(row.prefix) + style.Render(name) + "\n"
		}
		out += "\n" + "New folder: " + m.nameInput + "▌\n"
		out += exploreHelpStyle.Render("type a name · enter create · esc cancel")
		return appendStatus(out, "")
	}

	reserve := 4 // title(1) + blank(1) + blank-before-help(1) + help(1)
	if m.status != "" {
		reserve++ // status line
	}
	start, end := m.visibleRange(len(m.rows), m.cursor, reserve)
	for i := start; i < end; i++ {
		row := m.rows[i]
		name := row.node.name
		if row.node.isFolder {
			name += "/"
		}

		style := lipgloss.NewStyle()
		if row.node.isFolder {
			style = style.Foreground(exploreFolderClr)
		}
		if i == m.cursor {
			style = style.Bold(true)
		}

		out += exploreTreeStyle.Render(row.prefix) + style.Render(name) + "\n"
	}
	if len(m.rows) == 0 && m.status == "" {
		out += exploreCrumbStyle.Render("Empty group.") + "\n"
	}

	help := "↑/↓ move · → open · ← close · d download · esc groups · q quit"
	if m.canManage() {
		help = "↑/↓ move · → open · ← close · d download · u upload · n new folder · x delete · esc groups · q quit"
	}
	out += "\n" + exploreHelpStyle.Render(help)
	return appendStatus(out, m.status)
}

// viewPicker: Render the local file picker screen.
//
// Returns:
// - result1 (string): The rendered view.
func (m exploreModel) viewPicker() string {
	out := exploreTitleStyle.Render("flick - pick files to upload") + "\n"
	out += exploreCrumbStyle.Render(m.pickerDir) + "\n\n"

	reserve := 5 // title(1) + crumb(1) + blank(1) + blank-before-help(1) + help(1)
	if m.status != "" {
		reserve++ // status line
	}
	start, end := m.visibleRange(len(m.pickerItems), m.pickerCursor, reserve)
	for i := start; i < end; i++ {
		item := m.pickerItems[i]
		name := item.name
		if item.isDir {
			name += "/"
		}

		style := lipgloss.NewStyle()
		if item.isDir {
			style = style.Foreground(exploreFolderClr)
		}
		if m.pickerSelected[item.path] {
			style = style.Foreground(exploreSelectClr)
		}
		if i == m.pickerCursor {
			style = style.Bold(true)
		}
		out += "  " + style.Render(name) + "\n"
	}
	out += "\n" + exploreHelpStyle.Render("↑/↓ move · → enter dir · ← parent · space select · enter upload · esc cancel")
	return appendStatus(out, m.status)
}

// appendStatus: Append a status or error line to the view when present.
//
// Params:
// - out (string): The view rendered so far.
// - status (string): The status line, empty to append nothing.
//
// Returns:
// - result1 (string): The view with the status line appended.
func appendStatus(out, status string) string {
	if status != "" {
		out += "\n" + status
	}
	return out
}
