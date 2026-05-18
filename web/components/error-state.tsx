"use client"

import { AlertTriangle, RefreshCw } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

type ErrorStateProps = {
  title: string
  description: string
  details?: string
  retryLabel?: string
  onRetry?: () => void
}

export function ErrorState({ title, description, details, retryLabel, onRetry }: ErrorStateProps) {
  return (
    <div className="w-full">
      <Card className="border-destructive/40">
        <CardHeader>
          <div className="flex items-start gap-3">
            <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-md bg-destructive/10 text-destructive">
              <AlertTriangle className="h-5 w-5" />
            </span>
            <div className="space-y-1">
              <CardTitle>{title}</CardTitle>
              <CardDescription>{description}</CardDescription>
            </div>
          </div>
        </CardHeader>
        {(details || onRetry) && (
          <CardContent className="space-y-4">
            {details && (
              <pre className="overflow-auto rounded-md border bg-muted p-3 text-xs whitespace-pre-wrap break-words">
                {details}
              </pre>
            )}
            {onRetry && retryLabel && (
              <div className="flex justify-end">
                <Button variant="outline" onClick={onRetry}>
                  <RefreshCw className="h-4 w-4" />
                  {retryLabel}
                </Button>
              </div>
            )}
          </CardContent>
        )}
      </Card>
    </div>
  )
}
