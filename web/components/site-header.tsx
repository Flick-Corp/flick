"use client"

import { ArrowDownLeft, ArrowUpRight, LayoutDashboard, LogIn, UserRound } from "lucide-react"
import { useTranslations } from "next-intl"
import { useEffect, useState } from "react"

import { ThemeToggle } from "@/components/theme-toggle"
import { Button } from "@/components/ui/button"
import { type AuthSession } from "@/lib/api"
import { loadSession } from "@/lib/auth"
import { Link, usePathname } from "@/i18n/navigation"

export default function SiteHeader() {
  const t = useTranslations("Header")
  const pathname = usePathname()

  const [session, setSession] = useState<AuthSession | null>(null)

  // Re-read the session on every navigation so the header reflects login/logout
  // without a full page reload.
  useEffect(() => {
    setSession(loadSession())
  }, [pathname])

  if (pathname.startsWith("/dashboard")) {
    return null
  }

  return (
    <header className="border-b">
      <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-6">
        <Link href="/" className="flex items-center gap-2">
          <span className="flex h-8 w-8 items-center justify-center rounded-md bg-primary text-primary-foreground">
            <ArrowUpRight className="h-5 w-5" />
          </span>
          <span className="text-lg font-semibold">flick</span>
        </Link>

        <nav className="flex items-center gap-4">
          <Button asChild variant="ghost">
            <Link href="/dashboard">
              <LayoutDashboard className="h-4 w-4" />
              {t("dashboard")}
            </Link>
          </Button>
          <Button asChild>
            <Link href="/send">
              <ArrowUpRight className="h-4 w-4" />
              {t("send")}
            </Link>
          </Button>
          <Button asChild variant="outline">
            <Link href="/receive">
              <ArrowDownLeft className="h-4 w-4" />
              {t("receive")}
            </Link>
          </Button>
          {session ? (
            <Button asChild variant="ghost">
              <Link href="/profile">
                <UserRound className="h-4 w-4" />
                {session.user.username || t("profile")}
              </Link>
            </Button>
          ) : (
            <Button asChild variant="ghost">
              <Link href="/login">
                <LogIn className="h-4 w-4" />
                {t("login")}
              </Link>
            </Button>
          )}
          <ThemeToggle />
        </nav>
      </div>
    </header>
  )
}
