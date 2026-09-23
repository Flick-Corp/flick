"use client"

import { useTranslations } from "next-intl"

import { GithubIcon } from "@/components/icons"
import { legal } from "@/lib/legal"
import { Link, usePathname } from "@/i18n/navigation"

export default function SiteFooter() {
  const t = useTranslations("Footer")
  const pathname = usePathname()

  if (pathname.startsWith("/dashboard")) {
    return null
  }

  return (
    <footer className="border-t border-border">
      <div className="mx-auto flex max-w-6xl flex-col gap-5 px-6 py-7 text-sm text-muted-foreground sm:flex-row sm:items-center sm:justify-between">
        <nav className="flex flex-wrap items-center gap-x-6 gap-y-2" aria-label={t("legalLabel")}>
          <Link href="/legal" className="transition-colors hover:text-foreground">
            {t("legal")}
          </Link>
          <Link href="/terms" className="transition-colors hover:text-foreground">
            {t("terms")}
          </Link>
          <Link href="/privacy" className="transition-colors hover:text-foreground">
            {t("privacy")}
          </Link>
          <a
            href={legal.repoUrl}
            target="_blank"
            rel="noreferrer"
            className="inline-flex items-center gap-1.5 transition-colors hover:text-foreground"
          >
            <GithubIcon className="size-4" />
            GitHub
          </a>
        </nav>

        <span>{t("rights")}</span>
      </div>
    </footer>
  )
}
