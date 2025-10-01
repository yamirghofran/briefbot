import type { ColumnDef } from "@tanstack/react-table"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { StaticColumnHeader } from "@/components/data-table/static-column-header"
import { ReadStatusCell } from "@/components/read-status-cell"
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
import { Link, Trash2 } from "lucide-react"
import { format, isToday, isYesterday, differenceInDays } from "date-fns"
import { useNavigate } from "@tanstack/react-router"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { itemApi } from "@/services/api"
import type { Item } from "@/types"
import { useState } from "react"

// Helper function to format dates as "today", "yesterday", or "September 5th"
function formatRelativeDate(date: Date): string {
  if (isToday(date)) {
    return "today"
  } else if (isYesterday(date)) {
    return "yesterday"
  } else {
    // Format as "September 5th" with ordinal suffix
    const day = date.getDate()
    const month = format(date, "MMMM")
    
    // Add ordinal suffix (st, nd, rd, th)
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

// Title Cell Component with navigation
function TitleCell({ title, url, id }: { title: string | null, url: string | null, id: number }) {
  const navigate = useNavigate()
  
  const handleTitleClick = () => {
    navigate({ to: `/items/${id}` })
  }
  
  return (
    <div className="max-w-[300px] truncate font-medium flex items-center gap-2">
      {url && (
        <a 
          href={url} 
          target="_blank" 
          rel="noopener noreferrer"
          className="text-blue-600 hover:text-blue-800 flex-shrink-0"
          onClick={(e) => e.stopPropagation()}
        >
          <Link className="h-4 w-4" />
        </a>
      )}
      <button
        onClick={handleTitleClick}
        className="text-left hover:text-blue-600 hover:underline cursor-pointer truncate flex-1"
      >
        {title || <span className="text-gray-400">Untitled</span>}
      </button>
    </div>
  )
}

// Delete Action Cell Component
function DeleteActionCell({ item }: { item: Item }) {
  const [showDialog, setShowDialog] = useState(false)
  const queryClient = useQueryClient()

  const deleteMutation = useMutation({
    mutationFn: () => itemApi.deleteItem(item.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["items"] })
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
            This action cannot be undone. This will permanently delete the item
            "{item.title || 'Untitled'}" from your collection.
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

export const itemColumns: ColumnDef<Item>[] = [
  {
    accessorKey: "title",
    header: ({ column }) => (
      <StaticColumnHeader column={column} title="Title" />
    ),
    cell: ({ row }) => {
      const title = row.getValue("title") as string
      const url = row.original.url
      const id = row.original.id

      return <TitleCell title={title} url={url} id={id} />
    },
    enableGlobalFilter: true,
  },
  {
    accessorKey: "tags",
    header: ({ column }) => (
      <StaticColumnHeader column={column} title="Tags" />
    ),
    cell: ({ row }) => {
      const tags = row.getValue("tags") as string[]

      if (!tags || tags.length === 0) {
        return <span className="text-gray-400">—</span>
      }

      // Show up to 2 tags, then show +x for the rest
      const displayTags = tags.slice(0, 2)
      const remainingCount = Math.max(0, tags.length - 2)

      return (
        <div className="flex items-center gap-1 flex-wrap">
          {displayTags.map((tag, index) => (
            <Badge key={index} variant="secondary" className="text-xs">
              {tag}
            </Badge>
          ))}
          {remainingCount > 0 && (
            <Badge variant="outline" className="text-xs">
              +{remainingCount}
            </Badge>
          )}
        </div>
      )
    },
    filterFn: (row, id, filterValue) => {
      const tags = row.getValue(id) as string[]
      if (!tags || tags.length === 0) return false
      if (!filterValue || filterValue.length === 0) return true
      return filterValue.some((filter: string) => tags.includes(filter))
    },
    enableGlobalFilter: true,
  },
  {
    id: "status",
    header: ({ column }) => (
      <StaticColumnHeader column={column} title="Status" />
    ),
    cell: ({ row }) => {
      const item = row.original
      const hasContent = item.text_content && item.summary
      
      return (
        <Badge variant={hasContent ? "default" : "secondary"}>
          {hasContent ? "Completed" : "Processing"}
        </Badge>
      )
    },
  },
  {
    accessorKey: "type",
    header: ({ column }) => (
      <StaticColumnHeader column={column} title="Type" />
    ),
    cell: ({ row }) => {
      const type = row.getValue("type") as string

      return type ? (
        <Badge variant="outline">
          {type}
        </Badge>
      ) : (
        <span className="text-gray-400">—</span>
      )
    },
    filterFn: (row, id, filterValue) => {
      const type = row.getValue(id) as string
      if (!type) return false
      if (!filterValue || filterValue.length === 0) return true
      return filterValue.includes(type)
    },
    enableGlobalFilter: true,
  },
  {
    accessorKey: "platform",
    header: ({ column }) => (
      <StaticColumnHeader column={column} title="Platform" />
    ),
    cell: ({ row }) => {
      const platform = row.getValue("platform") as string
      return (
        <div className="w-[100px]">
          {platform || <span className="text-gray-400">—</span>}
        </div>
      )
    },
    filterFn: (row, id, filterValue) => {
      const platform = row.getValue(id) as string
      if (!platform) return false
      if (!filterValue || filterValue.length === 0) return true
      return filterValue.includes(platform)
    },
    enableGlobalFilter: true,
  },
  {
    accessorKey: "authors",
    header: ({ column }) => (
      <StaticColumnHeader column={column} title="Authors" />
    ),
    cell: ({ row }) => {
      const authors = row.getValue("authors") as string[]

      if (!authors || authors.length === 0) {
        return <span className="text-gray-400">—</span>
      }

      return (
        <div className="flex items-center gap-1 flex-wrap">
          <Badge variant="secondary" className="text-xs">
            {authors[0]}
          </Badge>
          {authors.length > 1 && (
            <Badge variant="outline" className="text-xs">
              +{authors.length - 1}
            </Badge>
          )}
        </div>
      )
    },
    filterFn: (row, id, filterValue) => {
      const authors = row.getValue(id) as string[]
      if (!authors || authors.length === 0) return false
      if (!filterValue || filterValue.length === 0) return true
      return filterValue.some((filter: string) => authors.includes(filter))
    },
    enableGlobalFilter: true,
  },
  {
    accessorKey: "is_read",
    header: ({ column }) => (
      <div className="flex justify-center">
        <StaticColumnHeader column={column} title="Read" />
      </div>
    ),
    cell: ({ row }) => {
      // This will be handled by the parent component that provides userId
      return null // Will be replaced by a custom cell renderer
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
      const item = row.original

      return (
        <div className="flex justify-center">
          <DeleteActionCell item={item} />
        </div>
      )
    },
  },
]