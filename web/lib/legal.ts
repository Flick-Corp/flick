// Facts quoted by /legal, /terms and /privacy. The service limits mirror DefaultConfig
// in internal/api/serverconfig and SessionDuration in internal/api/routes/account.
export const legal = {
  siteName: "Flick",
  siteHost: "flick.d3l.tech",
  siteUrl: "https://flick.d3l.tech",
  repoUrl: "https://github.com/Flick-Corp/flick",

  updatedAt: "2026-09-22",

  publisher: {
    names: "Matteo TUTTI, Valentino ZARRILLO, Matthias YILDIRIM, Antoine MALGAT",
    publicationDirector: "Matteo TUTTI",
    address: "KMØ, 30 rue François Spoerry, 68100 Mulhouse, France",
    phone: "",
    email: "matteo.tutti@epitech.eu",
  },

  host: {
    name: "TekitUp",
    email: "contact@tekitup.net",
    address: "",
    phone: "",
  },

  authority: {
    name: "Commission nationale de l'informatique et des libertés (CNIL)",
    address: "3 place de Fontenoy, TSA 80715, 75334 Paris Cedex 07",
    url: "https://www.cnil.fr",
  },

  service: {
    maxFileSizeMb: 1000,
    defaultExpirationMinutes: 15,
    maxExpirationHours: 4,
    defaultDownloadCount: 1,
    maxDownloadCount: 5,
    anonymousQuotaGb: 1,
    userQuotaGb: 5,
    groupQuotaGb: 10,
    sessionDurationDays: 7,
    abuseAckHours: 72,
    mailRetentionMonths: 12,
  },
} as const

export const legalValues: Record<string, string> = {
  site: legal.siteName,
  siteHost: legal.siteHost,
  siteUrl: legal.siteUrl,
  repo: legal.repoUrl,
  publisher: legal.publisher.names,
  director: legal.publisher.publicationDirector,
  publisherAddress: legal.publisher.address,
  publisherPhone: legal.publisher.phone,
  contactEmail: legal.publisher.email,
  host: legal.host.name,
  hostEmail: legal.host.email,
  hostAddress: legal.host.address,
  hostPhone: legal.host.phone,
  authority: legal.authority.name,
  authorityAddress: legal.authority.address,
  authorityUrl: legal.authority.url,
  maxFileSize: String(legal.service.maxFileSizeMb),
  defaultExpiration: String(legal.service.defaultExpirationMinutes),
  maxExpiration: String(legal.service.maxExpirationHours),
  defaultDownloads: String(legal.service.defaultDownloadCount),
  maxDownloads: String(legal.service.maxDownloadCount),
  anonymousQuota: String(legal.service.anonymousQuotaGb),
  userQuota: String(legal.service.userQuotaGb),
  groupQuota: String(legal.service.groupQuotaGb),
  sessionDays: String(legal.service.sessionDurationDays),
  abuseAck: String(legal.service.abuseAckHours),
  mailMonths: String(legal.service.mailRetentionMonths),
}

export const legalRoutes = [
  { href: "/legal", namespace: "Legal" },
  { href: "/terms", namespace: "Terms" },
  { href: "/privacy", namespace: "Privacy" },
] as const

export function fillLegal(text: string): string | null {
  let missing = false

  const filled = text.replace(/%(\w+)%/g, (match, key: string) => {
    const value = legalValues[key]
    if (value === undefined) return match
    if (value === "") {
      missing = true
      return ""
    }
    return value
  })

  return missing ? null : filled
}
