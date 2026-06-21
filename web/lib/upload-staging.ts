import { type UploadEntry } from "@/lib/api"

// UploadItem: one entry in a staging list. A loose file is a single-entry item;
// a folder keeps every file with its relative path so it can be zipped with its
// structure intact, just like the CLI. Shared by the send page and group uploads.
export interface UploadItem {
  id: string
  name: string
  isFolder: boolean
  entries: UploadEntry[]
  size: number
}

// formatBytes: Human-readable byte size (French units, matching the send page).
export function formatBytes(bytes: number): string {
  if (bytes === 0) return "0 o"
  const units = ["o", "Ko", "Mo", "Go"]
  const i = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1)
  return `${(bytes / 1024 ** i).toFixed(i === 0 ? 0 : 1)} ${units[i]}`
}

let itemCounter = 0
function nextId(): string {
  itemCounter += 1
  return `item-${itemCounter}`
}

// fileItem: Wrap a single loose file as a staging item.
export function fileItem(file: File): UploadItem {
  return { id: nextId(), name: file.name, isFolder: false, entries: [{ path: file.name, file }], size: file.size }
}

// folderItem: Wrap a folder (its flat entry list) as a staging item.
export function folderItem(name: string, entries: UploadEntry[]): UploadItem {
  const size = entries.reduce((total, entry) => total + entry.file.size, 0)
  return { id: nextId(), name, isFolder: true, entries, size }
}

// folderItemFromInputFiles: Build a folder staging item from the FileList of a
// <input webkitdirectory>. Each file carries its webkitRelativePath, e.g.
// "myfolder/sub/a.txt"; the first segment is the folder name.
export function folderItemFromInputFiles(files: File[]): UploadItem {
  const top = files[0].webkitRelativePath.split("/")[0] || "folder"
  const entries: UploadEntry[] = files.map((file) => ({
    path: file.webkitRelativePath || file.name,
    file,
  }))
  return folderItem(top, entries)
}

// fileFromEntry / walkDirectory: read a dropped FileSystemEntry tree into a flat
// UploadEntry list, preserving each file's path relative to the dropped folder.
function fileFromEntry(entry: FileSystemFileEntry): Promise<File> {
  return new Promise((resolve, reject) => entry.file(resolve, reject))
}

async function walkDirectory(dir: FileSystemDirectoryEntry, prefix: string, out: UploadEntry[]): Promise<void> {
  const reader = dir.createReader()
  const readBatch = () => new Promise<FileSystemEntry[]>((resolve, reject) => reader.readEntries(resolve, reject))

  // readEntries returns the directory in batches; keep reading until it drains.
  for (let batch = await readBatch(); batch.length > 0; batch = await readBatch()) {
    for (const child of batch) {
      const childPath = `${prefix}/${child.name}`
      if (child.isFile) {
        const file = await fileFromEntry(child as FileSystemFileEntry)
        out.push({ path: childPath, file })
      } else if (child.isDirectory) {
        await walkDirectory(child as FileSystemDirectoryEntry, childPath, out)
      }
    }
  }
}

// itemsFromDataTransfer: Turn a drop's DataTransferItemList into staging items,
// walking dropped folders via the Entry API.
export async function itemsFromDataTransfer(list: DataTransferItemList): Promise<UploadItem[]> {
  // webkitGetAsEntry() must be called synchronously while the event is live, so
  // collect every entry first, then traverse the directories asynchronously.
  const entries: FileSystemEntry[] = []
  for (let i = 0; i < list.length; i++) {
    const entry = list[i].webkitGetAsEntry()
    if (entry) entries.push(entry)
  }

  const items: UploadItem[] = []
  for (const entry of entries) {
    if (entry.isFile) {
      const file = await fileFromEntry(entry as FileSystemFileEntry)
      items.push(fileItem(file))
    } else if (entry.isDirectory) {
      const collected: UploadEntry[] = []
      await walkDirectory(entry as FileSystemDirectoryEntry, entry.name, collected)
      if (collected.length > 0) items.push(folderItem(entry.name, collected))
    }
  }
  return items
}
