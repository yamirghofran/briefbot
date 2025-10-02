import type { ColumnDef } from "@tanstack/react-table"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { StaticColumnHeader } from "@/components/data-table/static-column-header"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
import { Trash2, Loader2, CheckCircle2, XCircle, Clock } from "lucide-react"
import { format, isToday, isYesterday } from "date-fns"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { podcastApi } from "@/services/api"
import type { Podcast } from "@/types"
import { useState } from "react"

// Helper function to format dates
function formatRelativeDate(date: Date): string {
  if (isToday(date)) {
    return "today"
  } else if (isYesterday(date)) {
    return "yesterday"
  } else {
    const day = date.getDate()
    const month = format(date, "MMMM")

    let suffix = "th"
    if (day < 11 || day > 13) {
      switch (day % 10) {
        case 1:
          suffix = "st"
          break
        case 2:
          suffix = "nd"
          break
        case 3:
          suffix = "rd"
          break
      }
    }

    return `${month} ${day}${suffix}`
  }
}

// Status Badge Component
function StatusBadge({ status }: { status: string }) {
  const statusConfig = {
    pending: { icon: Clock, variant: "secondary" as const, label: "Pending" },
    writing: { icon: Loader2, variant: "secondary" as const, label: "Writing Script" },
    generating: { icon: Loader2, variant: "secondary" as const, label: "Generating Audio" },
    completed: { icon: CheckCircle2, variant: "default" as const, label: "Completed" },
    failed: { icon: XCircle, variant: "destructive" as const, label: "Failed" },
  }

  const config = statusConfig[status as keyof typeof statusConfig] || statusConfig.pending
  const Icon = config.icon

  return (
    <Badge variant={config.variant} className="flex items-center gap-1">
      <Icon className={`h-3 w-3 ${(status === 'writing' || status === 'generating') ? 'animate-spin' : ''}`} />
      {config.label}
    </Badge>
  )
}

// Audio Player Component
function AudioPlayer({ audioUrl, title }: { audioUrl: string | null, title: string }) {
  if (!audioUrl) {
    return <span className="text-gray-400 text-sm">No audio yet</span>
  }

  return (
    <audio controls className="h-10 max-w-md">
      <source src={audioUrl} type="audio/mpeg" />
      Your browser does not support the audio element.
    </audio>
  )
}

// Delete Action Cell Component
function DeleteActionCell({ podcast }: { podcast: Podcast }) {
  const [showDialog, setShowDialog] = useState(false)
  const queryClient = useQueryClient()

  const deleteMutation = useMutation({
    mutationFn: () => podcastApi.deletePodcast(podcast.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["podcasts"] })
      setShowDialog(false)
    },
  })

  return (
    <AlertDialog open={showDialog} onOpenChange={setShowDialog}>
      <AlertDialogTrigger asChild>
        <Button
          variant="ghost"
          size="icon"
          className="h-8 w-8 text-destructive hover:text-destructive hover:bg-destructive/10"
          onClick={(e) => {
            e.stopPropagation()
          }}
        >
          <Trash2 className="h-4 w-4" />
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Are you sure?</AlertDialogTitle>
          <AlertDialogDescription>
            This action cannot be undone. This will permanently delete the podcast
            "{podcast.title}" from your collection.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction asChild>
            <Button variant="destructive" onClick={() => deleteMutation.mutate()}>
              Delete
            </Button>
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}

export const podcastColumns: ColumnDef<Podcast>[] = [
  {
    accessorKey: "title",
    header: ({ column }) => (
      <StaticColumnHeader column={column} title="Title" />
    ),
    cell: ({ row }) => {
      const title = row.getValue("title") as string
      return <div className="max-w-[300px] truncate font-medium">{title}</div>
    },
    enableGlobalFilter: true,
  },
  {
    accessorKey: "status",
    header: ({ column }) => (
      <StaticColumnHeader column={column} title="Status" />
    ),
    cell: ({ row }) => {
      const status = row.getValue("status") as string
      return <StatusBadge status={status} />
    },
  },
  {
    accessorKey: "audio_url",
    header: () => <div>Audio</div>,
    cell: ({ row }) => {
      const audioUrl = row.getValue("audio_url") as string | null
      const title = row.original.title
      return <AudioPlayer audioUrl={audioUrl} title={title} />
    },
  },
  {
    accessorKey: "created_at",
    header: ({ column }) => (
      <StaticColumnHeader column={column} title="Created" />
    ),
    cell: ({ row }) => {
      const dateStr = row.getValue("created_at") as string
      const date = dateStr ? new Date(dateStr) : null

      return (
        <div className="w-[120px] text-sm">
          {date ? formatRelativeDate(date) : <span className="text-gray-400">—</span>}
        </div>
      )
    },
  },
  {
    id: "actions",
    header: () => <div className="text-center">Actions</div>,
    cell: ({ row }) => {
      const podcast = row.original

      return (
        <div className="flex justify-center">
          <DeleteActionCell podcast={podcast} />
        </div>
      )
    },
  },
]
