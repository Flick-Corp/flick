import type { AuthSession } from "@/lib/api"

// The signed-in session (user + token) is kept in localStorage. The API stores
// sessions server-side; this only remembers which token belongs to this browser.
const STORAGE_KEY = "flick.session"

export function saveSession(session: AuthSession): void {
  if (typeof window === "undefined") return
  window.localStorage.setItem(STORAGE_KEY, JSON.stringify(session))
}

export function loadSession(): AuthSession | null {
  if (typeof window === "undefined") return null

  const raw = window.localStorage.getItem(STORAGE_KEY)
  if (!raw) return null

  try {
    const parsed = JSON.parse(raw) as AuthSession
    if (!parsed || typeof parsed.token !== "string" || !parsed.user) return null
    return parsed
  } catch {
    return null
  }
}

export function clearSession(): void {
  if (typeof window === "undefined") return
  window.localStorage.removeItem(STORAGE_KEY)
}
