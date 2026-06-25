import { ArrowDownLeft, ArrowUpRight, Clock, Download, KeyRound, UserPlus, Zap } from "lucide-react"
import { useTranslations } from "next-intl"

import { CliDemo } from "@/components/cli-demo"
import { MouseMist } from "@/components/mouse-mist"
import { Button } from "@/components/ui/button"
import { Card } from "@/components/ui/card"
import { Link } from "@/i18n/navigation"

export default function Page() {
  const tHero = useTranslations("Hero")
  const tFeatures = useTranslations("Features")
  const tHow = useTranslations("HowItWorks")
  const tCTA = useTranslations("CallToAction")

  const features = [
    { icon: Zap, title: tFeatures("instantTitle"), body: tFeatures("instantBody") },
    { icon: Download, title: tFeatures("receiveTitle"), body: tFeatures("receiveBody") },
    { icon: Clock, title: tFeatures("expirationTitle"), body: tFeatures("expirationBody") },
    { icon: KeyRound, title: tFeatures("protectionTitle"), body: tFeatures("protectionBody") },
  ]

  const steps = [
    { title: tHow("step1Title"), body: tHow("step1Body") },
    { title: tHow("step2Title"), body: tHow("step2Body") },
    { title: tHow("step3Title"), body: tHow("step3Body") },
  ]

  return (
    <main className="relative mx-auto max-w-6xl px-4 py-12 sm:px-6 sm:py-16">
      <MouseMist />

      {/* Hero */}
      <section className="flex flex-col items-center text-center">
        <h1 className="max-w-3xl text-4xl font-bold tracking-tight sm:text-5xl md:text-6xl">
          {tHero("titleStart")} <span className="text-primary">{tHero("titleHighlight")}</span>
        </h1>
        <p className="mt-6 max-w-xl text-lg text-muted-foreground">{tHero("description")}</p>

        <div className="mt-8 flex flex-col gap-3 sm:flex-row">
          <Button asChild size="lg" className="h-14 px-8 text-lg">
            <Link href="/send">
              <ArrowUpRight className="size-6" />
              {tHero("ctaSend")}
            </Link>
          </Button>
          <Button asChild size="lg" variant="outline" className="h-14 px-8 text-lg">
            <Link href="/receive">
              <ArrowDownLeft className="size-6" />
              {tHero("ctaReceive")}
            </Link>
          </Button>
        </div>
      </section>

      {/* Features */}
      <section className="mt-20 grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
        {features.map((feature) => {
          const Icon = feature.icon
          return (
            <Card key={feature.title} className="p-8 text-center">
              <span className="mx-auto flex h-16 w-16 items-center justify-center rounded-2xl bg-primary/10 text-primary">
                <Icon className="h-7 w-7" />
              </span>
              <h3 className="mt-4 text-sm font-semibold">{feature.title}</h3>
              <p className="mt-1.5 text-sm text-muted-foreground">{feature.body}</p>
            </Card>
          )
        })}
      </section>

      {/* CLI Demo */}
      <section className="mt-20">
        <div className="mx-auto max-w-3xl">
          <CliDemo />
        </div>
      </section>

      {/* How it works */}
      <section className="mt-20">
        <div className="mx-auto max-w-2xl text-center">
          <p className="text-xs font-semibold tracking-widest text-primary uppercase">{tHow("eyebrow")}</p>
          <h2 className="mt-3 text-3xl font-bold tracking-tight sm:text-4xl">{tHow("title")}</h2>
        </div>

        <div className="mt-10 grid grid-cols-1 gap-6 md:grid-cols-3">
          {steps.map((step, index) => (
            <Card key={step.title} className="relative p-8">
              <span className="font-mono text-5xl font-bold text-primary">{String(index + 1).padStart(2, "0")}</span>
              <h3 className="mt-4 text-lg font-semibold">{step.title}</h3>
              <p className="mt-2 text-sm text-muted-foreground">{step.body}</p>
            </Card>
          ))}
        </div>
      </section>

      {/* CTA */}
      <section className="mt-20 mb-16">
        <div className="mx-auto max-w-xl rounded-2xl border bg-card p-10 text-center shadow-sm">
          <h2 className="text-2xl font-bold tracking-tight sm:text-3xl">{tCTA("title")}</h2>
          <p className="mt-3 text-muted-foreground">{tCTA("body")}</p>
          <Button asChild size="lg" className="mt-8 h-12 px-8 text-base">
            <Link href="/register">
              <UserPlus className="size-5" />
              {tCTA("button")}
            </Link>
          </Button>
        </div>
      </section>
    </main>
  )
}
