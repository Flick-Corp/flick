import type { Metadata } from "next"
import { getTranslations, setRequestLocale } from "next-intl/server"

import { LegalDocument } from "@/components/legal-document"
import { fillLegal, legal } from "@/lib/legal"

export async function generateMetadata({ params }: { params: Promise<{ locale: string }> }): Promise<Metadata> {
  const { locale } = await params
  const t = await getTranslations({ locale, namespace: "Legal" })

  return {
    title: `${t("title")} · ${legal.siteName}`,
    description: fillLegal(t("summary")) ?? t("summary"),
  }
}

export default async function LegalPage({ params }: { params: Promise<{ locale: string }> }) {
  const { locale } = await params
  setRequestLocale(locale)

  return <LegalDocument namespace="Legal" />
}
