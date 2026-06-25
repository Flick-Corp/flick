"use client"

import { useCallback, useEffect, useRef, useState } from "react"

// ── Color helper ──────────────────────────────────────────────────
const cc = (c?: string) =>
  c === "dim" ? "text-[#8b949e]"
  : c === "green" ? "text-emerald-400"
  : c === "yellow" ? "text-yellow-400"
  : "text-[#e6edf3]"

// ── Step definitions ──────────────────────────────────────────────
type Color = "dim" | "green" | "yellow"

type Step =
  | { t: "l"; text: string; color?: Color }          // line
  | { t: "p"; text: string; color?: Color }          // typewriter prompt
  | { t: "b" }                                       // blank
  | { t: "q" }                                       // quota bar (fade in)
  | { t: "r" }                                       // upload progress bar

const S: Step[] = [
  { t: "p", text: "$ flick myfile.pdf" },
  { t: "b" },
  { t: "l", text: "This upload contains:" },
  { t: "l", text: "  • myfile.pdf (2.4 MB)", color: "dim" },
  { t: "b" },
  { t: "q" },
  { t: "p", text: "Upload these files? [y/N]: y" },
  { t: "b" },
  { t: "l", text: "Uploading myfile.pdf... (2.4 MB archived)" },
  { t: "r" },
  { t: "b" },
  { t: "l", text: "Code: ocean-tiger-42  [1h left]", color: "yellow" },
  { t: "l", text: "Code copied to clipboard.", color: "green" },
]

// ── Sub-components ────────────────────────────────────────────────

function Cursor({ hide }: { hide?: boolean }) {
  return (
    <span
      className={`ml-[1px] inline-block w-[0.55em] align-text-bottom transition-opacity ${hide ? "opacity-0" : "opacity-100"}`}
      style={{ height: "1.15em", backgroundColor: "#e6edf3" }}
    />
  )
}

function Typewriter({ text, color, onDone }: { text: string; color?: Color; onDone: () => void }) {
  const [i, setI] = useState(0)
  const [cursor, setCursor] = useState(true)

  useEffect(() => {
    if (i >= text.length) {
      const id = setTimeout(onDone, 350)
      return () => clearTimeout(id)
    }
    const id = setTimeout(() => setI((x) => x + 1), 40 + Math.random() * 30)
    return () => clearTimeout(id)
  }, [i, text, onDone])

  useEffect(() => {
    const id = setInterval(() => setCursor((c) => !c), 530)
    return () => clearInterval(id)
  }, [])

  return (
    <div className={cc(color)}>
      {text.slice(0, i)}
      <Cursor hide={i >= text.length ? false : !cursor} />
    </div>
  )
}

function Blank({ onDone }: { onDone: () => void }) {
  useEffect(() => { const id = setTimeout(onDone, 200); return () => clearTimeout(id) }, [onDone])
  return <div className="h-[7px]" />
}

function QuotaBar({ onDone }: { onDone: () => void }) {
  const [show, setShow] = useState(false)
  useEffect(() => { const id = setTimeout(() => setShow(true), 50); return () => clearTimeout(id) }, [])
  useEffect(() => { if (show) { const id = setTimeout(onDone, 700); return () => clearTimeout(id) } }, [show, onDone])

  const bar = "█".repeat(9) + "░".repeat(11)
  return (
    <div className={`text-[#8b949e] transition-opacity duration-500 ${show ? "opacity-100" : "opacity-0"}`}>
      Quota: [{bar}] 125 MB / 256 MB used (49%)
    </div>
  )
}

function UploadProgress({ onDone }: { onDone: () => void }) {
  const [pct, setPct] = useState(0)
  const [show, setShow] = useState(false)

  useEffect(() => {
    const id = setTimeout(() => setShow(true), 50)
    return () => clearTimeout(id)
  }, [])

  useEffect(() => {
    if (!show) return
    const total = 80 // steps
    let p = 0
    const id = setInterval(() => {
      p += 1
      setPct(Math.round((p / total) * 100))
      if (p >= total) {
        clearInterval(id)
        setPct(100)
        const done = setTimeout(onDone, 600)
        return () => clearTimeout(done)
      }
    }, 40)
    return () => clearInterval(id)
  }, [show, onDone])

  const filled = Math.round((pct / 100) * 14)
  const bar = "█".repeat(filled) + "░".repeat(14 - filled)

  return (
    <div className={show ? "opacity-100" : "opacity-0"} style={{ transition: "opacity 0.3s" }}>
      <div className="inline-flex items-center gap-1.5 text-[#8b949e]">
        <span>Uploading</span>
        <span className="text-orange-500 dark:text-orange-400">[{bar}]</span>
        <span>{pct}% ({(pct * 2.4 / 100).toFixed(1)} MB / 2.4 MB)</span>
      </div>
    </div>
  )
}

function ActiveStep({ step, onDone }: { step: Step; onDone: () => void }) {
  switch (step.t) {
    case "l": return <StaticLine step={step} onDone={onDone} />
    case "p": return <Typewriter text={step.text} color={step.color} onDone={onDone} />
    case "b": return <Blank onDone={onDone} />
    case "q": return <QuotaBar onDone={onDone} />
    case "r": return <UploadProgress onDone={onDone} />
  }
}

// ── Main component ────────────────────────────────────────────────
export function CliDemo() {
  const [idx, setIdx] = useState(0)
  const [cursor, setCursor] = useState(true)
  const [done, setDone] = useState(false)
  const scrollRef = useRef<HTMLDivElement>(null)

  // Blinking cursor for final state
  useEffect(() => {
    if (!done) return
    const id = setInterval(() => setCursor((c) => !c), 530)
    return () => clearInterval(id)
  }, [done])

  // Auto-scroll
  useEffect(() => {
    if (scrollRef.current) scrollRef.current.scrollTop = scrollRef.current.scrollHeight
  })

  // Loop
  useEffect(() => {
    if (!done) return
    const id = setTimeout(() => { setIdx(0); setDone(false) }, 3500)
    return () => clearTimeout(id)
  }, [done])

  const advance = useCallback(() => {
    setIdx((i) => {
      const next = i + 1
      if (next >= S.length) setDone(true)
      return next
    })
  }, [])

  return (
    <div className="group relative">
      {/* Glow */}
      <div className="pointer-events-none absolute -inset-4 z-0 rounded-2xl bg-gradient-to-b from-primary/20 via-primary/5 to-transparent opacity-0 blur-2xl transition-opacity duration-1000 group-hover:opacity-100" />

      <div className="relative z-10 overflow-hidden rounded-xl border bg-[#0d1117] shadow-2xl shadow-black/30 dark:shadow-black/60">
        {/* Title bar */}
        <div className="flex items-center gap-2 border-b border-[#21262d] px-4 py-[11px]">
          <div className="flex gap-1.5">
            <span className="h-3 w-3 rounded-full bg-[#ff5f57]" />
            <span className="h-3 w-3 rounded-full bg-[#febc2e]" />
            <span className="h-3 w-3 rounded-full bg-[#28c840]" />
          </div>
          <span className="ml-2 text-xs font-medium text-[#8b949e]">flick — upload</span>
        </div>

        {/* Terminal body */}
        <div
          ref={scrollRef}
          className="overflow-auto p-4 pb-6 font-mono text-sm leading-relaxed [&::-webkit-scrollbar]:hidden"
          style={{ maxHeight: 420, minHeight: 340 }}
        >
          {/* Rendered steps go here, each one handles its own lifecycle */}
          <div className="flex flex-col gap-0">
            {Array.from({ length: idx }, (_, i) => {
              const step = S[i]
              switch (step.t) {
                case "b": return <div key={i} className="h-[7px]" />
                case "q": return (
                  <div key={i} className="text-[#8b949e]">
                    Quota: [{"█".repeat(9)}{"░".repeat(11)}] 125 MB / 256 MB used (49%)
                  </div>
                )
                case "r": return (
                  <div key={i} className="inline-flex items-center gap-1.5 text-[#8b949e]">
                    <span>Uploading</span>
                    <span className="text-orange-500 dark:text-orange-400">[{"█".repeat(14)}]</span>
                    <span>100% (2.4 MB / 2.4 MB)</span>
                  </div>
                )
                default: return (
                  <div key={i} className={cc((step as any).color)}>
                    {(step as any).text || ""}
                    {step.t === "p" && <Cursor />}
                  </div>
                )
              }
            })}

            {/* Active step */}
            {!done && idx < S.length && <ActiveStep key={idx} step={S[idx]} onDone={advance} />}

            {/* Final cursor */}
            {done && <Cursor hide={cursor} />}
          </div>
        </div>
      </div>
    </div>
  )
}

function StaticLine({ step, onDone }: { step: Step; onDone?: () => void }) {
  useEffect(() => {
    if (onDone) { const id = setTimeout(onDone, 300); return () => clearTimeout(id) }
  }, [onDone, step])
  return (
    <div className={cc((step as any).color)}>
      {(step as any).text || ""}
    </div>
  )
}
