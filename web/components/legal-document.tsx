import { Fragment, type ReactNode } from "react"
import { useLocale, useTranslations } from "next-intl"

import { fillLegal, legal, legalRoutes } from "@/lib/legal"
import { Link } from "@/i18n/navigation"

export type LegalBlock =
  | { type: "p"; text: string }
  | { type: "ul"; items: string[] }
  | { type: "dl"; items: { term: string; desc: string }[] }
  | { type: "note"; text: string }

export type LegalSection = {
  id: string
  heading: string
  blocks: LegalBlock[]
}

const LINK = /(https?:\/\/[^\s,;]+|[\w.+-]+@[\w-]+\.[\w.-]+)/g

function renderText(text: string, keyPrefix: string): ReactNode[] {
  return text.split(LINK).map((part, index) => {
    if (index % 2 === 0) {
      return part ? <Fragment key={`${keyPrefix}-${index}`}>{part}</Fragment> : null
    }

    const trailing = part.match(/[.)]+$/)?.[0] ?? ""
    const target = trailing ? part.slice(0, -trailing.length) : part
    const isEmail = target.includes("@")

    return (
      <Fragment key={`${keyPrefix}-${index}`}>
        <a
          href={isEmail ? `mailto:${target}` : target}
          {...(isEmail ? {} : { target: "_blank", rel: "noreferrer" })}
          className="text-primary underline underline-offset-4 transition-colors hover:text-foreground"
        >
          {target}
        </a>
        {trailing}
      </Fragment>
    )
  })
}

function Block({ block, id }: { block: LegalBlock; id: string }) {
  if (block.type === "p" || block.type === "note") {
    const text = fillLegal(block.text)
    if (!text) return null

    if (block.type === "note") {
      return (
        <p className="mt-4 rounded-xl border border-border bg-accent/25 p-4 text-[15px] leading-relaxed text-muted-foreground">
          {renderText(text, id)}
        </p>
      )
    }

    return <p className="mt-4 text-[15px] leading-relaxed text-muted-foreground">{renderText(text, id)}</p>
  }

  if (block.type === "ul") {
    const items = block.items.map(fillLegal).filter((item): item is string => item !== null)
    if (items.length === 0) return null

    return (
      <ul className="mt-4 space-y-2.5">
        {items.map((item, index) => (
          <li
            key={`${id}-li-${index}`}
            className="relative pl-5 text-[15px] leading-relaxed text-muted-foreground before:absolute before:top-[0.62em] before:left-0 before:size-1.5 before:rounded-full before:bg-primary/70 before:content-['']"
          >
            {renderText(item, `${id}-li-${index}`)}
          </li>
        ))}
      </ul>
    )
  }

  const rows = block.items
    .map((item) => ({ term: fillLegal(item.term), desc: fillLegal(item.desc) }))
    .filter((row): row is { term: string; desc: string } => row.term !== null && row.desc !== null)
  if (rows.length === 0) return null

  return (
    <dl className="mt-5 divide-y divide-border overflow-hidden rounded-xl border border-border">
      {rows.map((row, index) => (
        <div key={`${id}-dt-${index}`} className="grid gap-1 p-4 sm:grid-cols-[200px_minmax(0,1fr)] sm:gap-5">
          <dt className="font-heading text-sm font-semibold">{row.term}</dt>
          <dd className="text-[15px] leading-relaxed text-muted-foreground">
            {renderText(row.desc, `${id}-dd-${index}`)}
          </dd>
        </div>
      ))}
    </dl>
  )
}

export function LegalDocument({ namespace }: { namespace: "Legal" | "Terms" | "Privacy" }) {
  const t = useTranslations(namespace)
  const tRoot = useTranslations()
  const locale = useLocale()

  const sections = t.raw("sections") as LegalSection[]
  const updated = new Intl.DateTimeFormat(locale, { dateStyle: "long" }).format(new Date(legal.updatedAt))
  const others = legalRoutes.filter((route) => route.namespace !== namespace)

  return (
    <main className="relative mx-auto max-w-6xl px-6 pb-20">
      <header className="pt-16 sm:pt-20">
        <p className="font-heading text-xs font-semibold tracking-[0.14em] text-primary uppercase">{t("eyebrow")}</p>
        <h1 className="mt-4 max-w-3xl font-heading text-4xl font-bold tracking-tight sm:text-5xl">{t("title")}</h1>
        <p className="mt-5 max-w-2xl text-lg leading-relaxed text-muted-foreground">
          {fillLegal(t("intro")) ?? t("intro")}
        </p>
        <p className="mt-6 font-mono text-xs tracking-wide text-muted-foreground uppercase">
          {t("updated", { date: updated })}
        </p>
      </header>

      <div className="mt-14 grid gap-x-14 gap-y-10 lg:grid-cols-[minmax(0,1fr)_240px]">
        <article className="min-w-0">
          {sections.map((section) => (
            <section
              key={section.id}
              id={section.id}
              className="scroll-mt-24 border-t border-border pt-8 first:border-t-0 first:pt-0 [&+&]:mt-10"
            >
              <h2 className="font-heading text-xl font-bold tracking-tight">
                {fillLegal(section.heading) ?? section.heading}
              </h2>
              {section.blocks.map((block, index) => (
                <Block key={`${section.id}-${index}`} block={block} id={`${section.id}-${index}`} />
              ))}
            </section>
          ))}

          <nav className="mt-14 grid gap-4 sm:grid-cols-2">
            {others.map((route) => (
              <Link
                key={route.href}
                href={route.href}
                className="rounded-xl border border-border bg-accent/20 p-5 transition-colors hover:border-ring/50"
              >
                <span className="font-heading font-semibold">
                  {fillLegal(tRoot(`${route.namespace}.title`)) ?? tRoot(`${route.namespace}.title`)}
                </span>
                <span className="mt-1 block text-sm leading-relaxed text-muted-foreground">
                  {fillLegal(tRoot(`${route.namespace}.summary`)) ?? tRoot(`${route.namespace}.summary`)}
                </span>
              </Link>
            ))}
          </nav>
        </article>

        <aside className="order-first lg:order-last">
          <nav className="lg:sticky lg:top-24" aria-label={t("tocLabel")}>
            <p className="font-heading text-xs font-semibold tracking-[0.14em] text-muted-foreground uppercase">
              {t("tocLabel")}
            </p>
            <ol className="mt-4 space-y-2 border-l border-border">
              {sections.map((section) => (
                <li key={section.id}>
                  <a
                    href={`#${section.id}`}
                    className="-ml-px block border-l border-transparent pl-4 text-sm leading-snug text-muted-foreground transition-colors hover:border-primary hover:text-foreground"
                  >
                    {fillLegal(section.heading) ?? section.heading}
                  </a>
                </li>
              ))}
            </ol>
          </nav>
        </aside>
      </div>
    </main>
  )
}
