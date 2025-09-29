import { StaticDataTable } from "@/components/data-table/static-data-table"
import { itemColumns } from "@/components/data-table/item-columns"
import { ReadStatusCell } from "@/components/read-status-cell"
import type { Item } from "@/types"
import type { ColumnDef } from "@tanstack/react-table"

interface ItemDataTableProps {
  data: Item[]
  userId: number
}

import { useMemo } from 'react'

export function ItemDataTable({ data, userId }: ItemDataTableProps) {
  // Memoize columns to prevent re-creation on every render
  const customColumns: ColumnDef<Item>[] = useMemo(() => 
    itemColumns.map(column => {
      if (column.accessorKey === 'is_read') {
        return {
          ...column,
          cell: ({ row }) => (
            <ReadStatusCell item={row.original} userId={userId} />
          )
        }
      }
      return column
    }), 
    [userId] // Only re-create when userId changes
  )

  return (
    <StaticDataTable
      columns={customColumns}
      data={data}
    />
  )
}