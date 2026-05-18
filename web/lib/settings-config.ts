export type SettingFieldType = "text" | "textarea" | "number" | "switch" | "select"

export type SettingOption = { value: string; label: string }

export type SettingDependency = {
  key: string
  equals: string | number | boolean
}

export type SettingField = {
  key: string
  type: SettingFieldType
  defaultValue: string | number | boolean
  placeholder?: string
  options?: SettingOption[]
  hasDescription?: boolean
  dependsOn?: SettingDependency
  notAvailable?: boolean
}

export type SettingSection = {
  id: string
  hasDescription?: boolean
  fields: SettingField[]
}

export const settingsSections: SettingSection[] = [
  {
    id: "storage",
    hasDescription: true,
    fields: [
      {
        key: "persistence",
        type: "switch",
        defaultValue: true,
        hasDescription: true,
        notAvailable: true,
      },
    ],
  },
  {
    id: "uploads",
    hasDescription: true,
    fields: [
      {
        key: "maxFileSizeMb",
        type: "number",
        defaultValue: 1000,
      },
      {
        key: "defaultExpiration",
        type: "text",
        defaultValue: "15m",
        hasDescription: true,
      },
      {
        key: "maxExpiration",
        type: "text",
        defaultValue: "4h",
      },
      {
        key: "allowMultipleDownloads",
        type: "switch",
        defaultValue: false,
        hasDescription: true,
      },
      {
        key: "defaultDownloadCount",
        type: "number",
        defaultValue: 1,
        dependsOn: { key: "allowMultipleDownloads", equals: true },
        hasDescription: true,
      },
      {
        key: "maxDownloadCount",
        type: "number",
        defaultValue: 5,
        dependsOn: { key: "allowMultipleDownloads", equals: true },
      },
    ],
  },
  {
    id: "security",
    hasDescription: true,
    fields: [
      {
        key: "requirePassword",
        type: "switch",
        defaultValue: false,
        hasDescription: true,
        notAvailable: true,
      },
    ],
  },
  {
    id: "antiAbuse",
    hasDescription: true,
    fields: [
      {
        key: "activateRateLimit",
        type: "switch",
        defaultValue: true,
        hasDescription: true,
      },
      {
        key: "maxGenerationKeyPerHour",
        type: "number",
        defaultValue: 60,
        dependsOn: { key: "activateRateLimit", equals: true },
        hasDescription: true,
      },
      {
        key: "maxUploadPerHourPerKey",
        type: "number",
        defaultValue: 10,
        dependsOn: { key: "activateRateLimit", equals: true },
        hasDescription: true,
      },
      {
        key: "maxUploadPerHourPerIP",
        type: "number",
        defaultValue: 30,
        dependsOn: { key: "activateRateLimit", equals: true },
        hasDescription: true,
      },
      {
        key: "maxUploadPerHour",
        type: "number",
        defaultValue: 1000,
        dependsOn: { key: "activateRateLimit", equals: true },
        hasDescription: true,
      },
    ],
  },
]
