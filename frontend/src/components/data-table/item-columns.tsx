import type { ColumnDef } from "@tanstack/react-table"
import { Badge } from "@/components/ui/badge"
import { StaticColumnHeader } from "@/components/data-table/static-column-header"
import { ReadStatusCell } from "@/components/read-status-cell"
import { Link } from "lucide-react"
import { format, isToday, isYesterday, differenceInDays } from "date-fns"
import { useNavigate } from "@tanstack/react-router"
import type { Item } from "@/types"

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
]