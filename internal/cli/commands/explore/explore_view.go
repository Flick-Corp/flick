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
	var out string
	switch m.mode {
	case modeGroups:
		out = m.viewGroups()
	case modePicker:
		out = m.viewPicker()
	default:
		out = m.viewTree()
	}
	if m.width > 0 {
		out = lipgloss.NewStyle().MaxWidth(m.width).Render(out)
	}
	if m.height > 0 {
		out = lipgloss.NewStyle().Height(m.height).MaxHeight(m.height).Render(out)
	}
	return out
}

// chrome: Compute how many terminal lines the non-item parts of a screen occupy
// (fixed title/blank lines plus the help and status blocks) so the item list is
// trimmed enough to always leave room for the footer.
//
// Params:
// - fixed (int): Fixed chrome lines (title, blanks, blank-before-help).
// - help (string): The help bar text.
// - status (string): The status line, empty when absent.
//
// Returns:
// - result1 (int): The number of lines to reserve.
func (m exploreModel) chrome(fixed int, help, status string) int {
	r := fixed + lipgloss.Height(help)
	if status != "" {
		r += lipgloss.Height(status)
	}
	return r + 1 // safety margin to dodge the alt-screen bottom-row clip
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
	if m.height == 0 {
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
	help := "↑/↓ move · → open group · q quit"
	reserve := m.chrome(3, help, m.status) // title + blank + blank-before-help
	start, end := m.visibleRange(len(m.groups), m.groupCursor, reserve)
	for i := start; i < end; i++ {
		g := m.groups[i]
		style := lipgloss.NewStyle()
		if i == m.groupCursor {
			style = style.Bold(true)
		}
		out += "  " + style.Render(fmt.Sprintf("%s (%s)", g.Name, g.Role)) + "\n"
	}
	out += "\n" + exploreHelpStyle.Render(help)
	return appendStatus(out, m.status)
}

// viewTree: Render the folder tree screen.
//
// Returns:
// - result1 (string): The rendered view.
func (m exploreModel) viewTree() string {
	out := exploreTitleStyle.Render("flick - "+m.groupName) + "\n\n"

	if m.creating {
		help := "type a name · enter create · esc cancel"
		start, end := m.visibleRange(len(m.rows), m.cursor, m.chrome(4, help, ""))
		for i := start; i < end; i++ {
			out += m.renderRow(i)
		}
		out += "\n" + "New folder: " + m.nameInput + "▌\n"
		out += exploreHelpStyle.Render(help)
		return appendStatus(out, "")
	}

	help := "↑/↓ move · → open · ← close · d download · esc groups · q quit"
	if m.canManage() {
		help = "↑/↓ move · → open · ← close · d download · u upload · n new folder · x delete · esc groups · q quit"
	}
	reserve := m.chrome(3, help, m.status) // title + blank + blank-before-help
	start, end := m.visibleRange(len(m.rows), m.cursor, reserve)
	for i := start; i < end; i++ {
		out += m.renderRow(i)
	}
	if len(m.rows) == 0 && m.status == "" {
		out += exploreCrumbStyle.Render("Empty group.") + "\n"
	}

	out += "\n" + exploreHelpStyle.Render(help)
	return appendStatus(out, m.status)
}

// renderRow: Render a single flattened tree row (its prefix plus name), styling
// folders, the cursor line and the non-selectable "(empty)" placeholder.
//
// Params:
// - i (int): Index of the row in m.rows.
//
// Returns:
// - result1 (string): The rendered line, newline terminated.
func (m exploreModel) renderRow(i int) string {
	row := m.rows[i]
	if row.placeholder {
		mark := exploreCrumbStyle.Render(row.node.name)
		if i == m.cursor {
			mark = lipgloss.NewStyle().Bold(true).Render(row.node.name)
		}
		return exploreTreeStyle.Render(row.prefix) + mark + "\n"
	}

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
	return exploreTreeStyle.Render(row.prefix) + style.Render(name) + "\n"
}

// viewPicker: Render the local file picker screen.
//
// Returns:
// - result1 (string): The rendered view.
func (m exploreModel) viewPicker() string {
	out := exploreTitleStyle.Render("flick - pick files to upload") + "\n"
	out += exploreCrumbStyle.Render(m.pickerDir) + "\n\n"

	help := "↑/↓ move · → enter dir · ← parent · space select · enter upload · esc cancel"
	reserve := m.chrome(4, help, m.status) // title + crumb + blank + blank-before-help
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
	out += "\n" + exploreHelpStyle.Render(help)
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
