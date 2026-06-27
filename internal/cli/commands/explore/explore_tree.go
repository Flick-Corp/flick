/*
** FLICK PROJECT, 2026
** flick/internal/cli/commands/explore/explore_tree
** File description:
** Tree and picker helpers of the group explorer: build nodes from a level,
** locate a node, flatten the tree into rows and list a local directory.
 */

package explore

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// childrenFrom: Build tree nodes from a level's folders and files.
//
// Params:
// - folders ([]exploreFolder): The sub-folders at this level.
// - files ([]exploreFile): The files at this level.
//
// Returns:
// - result1 ([]*exploreNode): The folder and file nodes, folders first.
func childrenFrom(folders []exploreFolder, files []exploreFile) []*exploreNode {
	out := make([]*exploreNode, 0, len(folders)+len(files))
	for _, f := range folders {
		out = append(out, &exploreNode{id: f.ID, name: f.Name, isFolder: true})
	}
	for _, f := range files {
		out = append(out, &exploreNode{name: f.name, code: f.code, uploadID: f.id})
	}
	return out
}

// preserveExpansion: Carry the expanded/loaded state and already-loaded children
// of the previous nodes onto the freshly reloaded ones, matched by id, so that
// reloading a level (after an upload or a new folder) does not collapse the
// folders the user had open.
//
// Params:
// - old ([]*exploreNode): The nodes before the reload.
// - fresh ([]*exploreNode): The newly fetched nodes for the same level.
//
// Returns:
// - result1 ([]*exploreNode): fresh, with expansion state carried over.
func preserveExpansion(old, fresh []*exploreNode) []*exploreNode {
	for _, n := range fresh {
		if !n.isFolder {
			continue
		}
		for _, o := range old {
			if o.isFolder && o.id == n.id {
				n.expanded = o.expanded
				n.loaded = o.loaded
				n.children = o.children
				break
			}
		}
	}
	return fresh
}

// findNode: Find the folder node carrying the given id in the tree.
//
// Params:
// - nodes ([]*exploreNode): The nodes to search, recursively.
// - id (string): The folder id to look for.
//
// Returns:
// - result1 (*exploreNode): The matching node, or nil when not found.
func findNode(nodes []*exploreNode, id string) *exploreNode {
	for _, n := range nodes {
		if n.isFolder && n.id == id {
			return n
		}
		if found := findNode(n.children, id); found != nil {
			return found
		}
	}
	return nil
}

// rebuild: Flatten the expanded tree into the visible rows and clamp the cursor.
func (m *exploreModel) rebuild() {
	m.rows = nil
	var walk func(nodes []*exploreNode, ancestorsLast []bool, parentID string)
	walk = func(nodes []*exploreNode, ancestorsLast []bool, parentID string) {
		for i, n := range nodes {
			last := i == len(nodes)-1

			var prefix strings.Builder
			for _, parentLast := range ancestorsLast {
				if parentLast {
					prefix.WriteString("    ")
				} else {
					prefix.WriteString("│   ")
				}
			}
			if last {
				prefix.WriteString("└── ")
			} else {
				prefix.WriteString("├── ")
			}

			m.rows = append(m.rows, exploreRow{node: n, prefix: prefix.String(), parentID: parentID})
			if n.isFolder && n.expanded {
				if len(n.children) == 0 {
					var p strings.Builder
					for _, parentLast := range append(ancestorsLast, last) {
						if parentLast {
							p.WriteString("    ")
						} else {
							p.WriteString("│   ")
						}
					}
					p.WriteString("└── ")
					m.rows = append(m.rows, exploreRow{node: emptyMarker, prefix: p.String(), parentID: n.id, placeholder: true})
				}
				walk(n.children, append(ancestorsLast, last), n.id)
			}
		}
	}
	walk(m.roots, nil, "")

	if m.cursor >= len(m.rows) {
		m.cursor = len(m.rows) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

// targetFolder: The group folder an action (upload, new folder) applies to.
//
// Returns:
// - result1 (string): The target folder id, or "" for the group root.
func (m exploreModel) targetFolder() string {
	if len(m.rows) == 0 {
		return m.currentID
	}
	return m.rows[m.cursor].parentID
}

// loadPicker: List the visible entries of a local directory for the picker.
//
// Params:
// - dir (string): The directory to read.
//
// Returns:
// - result1 ([]pickerItem): The directory entries, folders first then by name.
// - result2 (error): An error if the directory cannot be read.
func loadPicker(dir string) ([]pickerItem, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	items := make([]pickerItem, 0, len(entries))
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		items = append(items, pickerItem{name: e.Name(), path: filepath.Join(dir, e.Name()), isDir: e.IsDir()})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].isDir != items[j].isDir {
			return items[i].isDir
		}
		return items[i].name < items[j].name
	})
	return items, nil
}
