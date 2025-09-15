import * as React from "react"
import type { Column } from "@tanstack/react-table"
import { ArrowUpDown, ChevronDown, ChevronUp } from "lucide-react"

import { Button } from "@/components/ui/button"

interface StaticColumnHeaderProps<TData, TValue>
  extends React.HTMLAttributes<HTMLDivElement> {
  column: Column<TData, TValue>
  title: string
}

export function StaticColumnHeader<TData, TValue>({
  column,
  title,
  className,
}: StaticColumnHeaderProps<TData, TValue>) {
  const [sortOrder, setSortOrder] = React.useState<'asc' | 'desc' | false>(false)

  const handleSort = () => {
    if (!column.getCanSort()) return
    
    if (sortOrder === false) {
      setSortOrder('asc')
    } else if (sortOrder === 'asc') {
      setSortOrder('desc')
    } else {
      setSortOrder(false)
    }
  }

  if (!column.getCanSort()) {
    return <div className={className}>{title}</div>
  }

  return (
    <div className={className}>
      <Button
        variant="ghost"
        onClick={handleSort}
        className="-ml-3 h-8 data-[state=open]:bg-accent"
      >
        <span className="text-sm font-medium">{title}</span>
        {sortOrder === 'asc' && <ChevronUp className="ml-2 h-4 w-4" />}
        {sortOrder === 'desc' && <ChevronDown className="ml-2 h-4 w-4" />}
        {sortOrder === false && <ArrowUpDown className="ml-2 h-4 w-4" />}
      </Button>
    </div>
  )
}