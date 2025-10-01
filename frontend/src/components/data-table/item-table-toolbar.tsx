import { Table } from "@tanstack/react-table"
import { DataTableFacetedFilter } from "./data-table-faceted-filter"
import { Button } from "@/components/ui/button"
import { X } from "lucide-react"
import type { Item } from "@/types"
import type { Option } from "@/types/data-table"
import { useMemo } from "react"

interface ItemTableToolbarProps {
  table: Table<Item>
  data: Item[]
}

export function ItemTableToolbar({ table, data }: ItemTableToolbarProps) {
  // Extract unique values from the data for filter options
  const filterOptions = useMemo(() => {
    const tags = new Set<string>()
    const types = new Set<string>()
    const platforms = new Set<string>()
    const authors = new Set<string>()

    data.forEach((item) => {
      if (item.tags) {
        item.tags.forEach((tag) => tags.add(tag))
      }
      if (item.type) {
        types.add(item.type)
      }
      if (item.platform) {
        platforms.add(item.platform)
      }
      if (item.authors) {
        item.authors.forEach((author) => authors.add(author))
      }
    })

    return {
      tags: Array.from(tags).sort().map((tag): Option => ({
        label: tag,
        value: tag,
      })),
      types: Array.from(types).sort().map((type): Option => ({
        label: type,
        value: type,
      })),
      platforms: Array.from(platforms).sort().map((platform): Option => ({
        label: platform,
        value: platform,
      })),
      authors: Array.from(authors).sort().map((author): Option => ({
        label: author,
        value: author,
      })),
    }
  }, [data])

  const isFiltered = table.getState().columnFilters.length > 0

  return (
    <div className="flex items-center justify-between gap-2">
      <div className="flex flex-1 flex-wrap items-center gap-2">
        {table.getColumn("tags") && filterOptions.tags.length > 0 && (
          <DataTableFacetedFilter
            column={table.getColumn("tags")}
            title="Tags"
            options={filterOptions.tags}
            multiple
          />
        )}
        {table.getColumn("platform") && filterOptions.platforms.length > 0 && (
          <DataTableFacetedFilter
            column={table.getColumn("platform")}
            title="Platform"
            options={filterOptions.platforms}
            multiple
          />
        )}
        {table.getColumn("type") && filterOptions.types.length > 0 && (
          <DataTableFacetedFilter
            column={table.getColumn("type")}
            title="Type"
            options={filterOptions.types}
            multiple
          />
        )}
        {table.getColumn("authors") && filterOptions.authors.length > 0 && (
          <DataTableFacetedFilter
            column={table.getColumn("authors")}
            title="Authors"
            options={filterOptions.authors}
            multiple
          />
        )}
        {isFiltered && (
          <Button
            variant="ghost"
            onClick={() => table.resetColumnFilters()}
            className="h-8 px-2 lg:px-3"
          >
            Reset
            <X className="ml-2 h-4 w-4" />
          </Button>
        )}
      </div>
    </div>
  )
}
