import type { ColumnDef } from "@tanstack/react-table"
import { Badge } from "@/components/ui/badge"
import { Checkbox } from "@/components/ui/checkbox"
import { StaticColumnHeader } from "@/components/data-table/static-column-header"
import type { Item } from "@/components/data-table/static-data"

export const staticColumns: ColumnDef<Item>[] = [
  {
    accessorKey: "title",
    header: ({ column }) => (
      <StaticColumnHeader column={column} title="Title" />
    ),
    cell: ({ row }) => (
      <div className="max-w-[300px] truncate font-medium">
        {row.getValue("title")}
      </div>
    ),
  },
  {
    accessorKey: "type",
    header: ({ column }) => (
      <StaticColumnHeader column={column} title="Type" />
    ),
    cell: ({ row }) => {
      const type = row.getValue("type") as string
      
      return (
        <Badge variant="outline">
          {type}
        </Badge>
      )
    },
  },
  {
    accessorKey: "platform",
    header: ({ column }) => (
      <StaticColumnHeader column={column} title="Platform" />
    ),
    cell: ({ row }) => (
      <div className="w-[100px]">
        {row.getValue("platform")}
      </div>
    ),
  },
  {
    accessorKey: "authors",
    header: ({ column }) => (
      <StaticColumnHeader column={column} title="Authors" />
    ),
    cell: ({ row }) => {
      const authors = row.getValue("authors") as string[]
      
      return (
        <div className="w-[150px]">
          {authors.join(", ")}
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
      const isRead = row.getValue("is_read") as boolean
      
      return (
        <div className="flex justify-center items-center h-full">
          <Checkbox
            checked={isRead}
            className="data-[state=checked]:bg-black data-[state=checked]:border-black border-gray-400"
            onCheckedChange={(checked) => {
              // Visual only - no state change in static table
              console.log(`Item ${row.original.id} read status would be:`, checked)
            }}
          />
        </div>
      )
    },
  },
  {
    accessorKey: "created_at",
    header: ({ column }) => (
      <StaticColumnHeader column={column} title="Created" />
    ),
    cell: ({ row }) => {
      const date = row.getValue("created_at") as Date
      
      return (
        <div className="w-[100px]">
          {date.toLocaleDateString()}
        </div>
      )
    },
  },
]