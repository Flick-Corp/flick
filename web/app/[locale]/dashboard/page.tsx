import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

export default function DashboardPage() {
  const stats = [
    { label: "Total transfers", value: "—" },
    { label: "Active links", value: "—" },
    { label: "Storage used", value: "—" },
    { label: "Users", value: "—" },
  ]

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-semibold tracking-tight">Overview</h2>
        <p className="text-muted-foreground">Quick glance at your Flick instance.</p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {stats.map((s) => (
          <Card key={s.label}>
            <CardHeader className="pb-2">
              <CardDescription>{s.label}</CardDescription>
              <CardTitle className="text-3xl">{s.value}</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-xs text-muted-foreground">Data will be wired to the Go API.</p>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}
