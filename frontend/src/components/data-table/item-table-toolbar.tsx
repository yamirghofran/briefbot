import { Table } from "@tanstack/react-table"
import { DataTableFacetedFilter } from "./data-table-faceted-filter"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { UrlSubmissionDialog } from "@/components/url-submission-dialog"
import { X, Search, Zap, Check, Loader2 } from "lucide-react"
import { useMutation } from "@tanstack/react-query"
import { toast } from "sonner"
import { digestApi } from "@/services/api"
import type { Item } from "@/types"
import type { Option } from "@/types/data-table"
import { useMemo, useState } from "react"

interface ItemTableToolbarProps {
  table: Table<Item>
  data: Item[]
  userId?: number
}

export function ItemTableToolbar({ table, data, userId }: ItemTableToolbarProps) {
  const [showSuccess, setShowSuccess] = useState(false)

  // Mutation for triggering integrated digest
  const triggerDigestMutation = useMutation({
    mutationFn: () => {
      if (userId) {
        return digestApi.triggerIntegratedDigestForUser(userId)
      } else {
        return digestApi.triggerIntegratedDigest()
      }
    },
    onSuccess: () => {
      toast.success(
        "Digest processing started! You'll receive an email when ready.",
      )
      setShowSuccess(true)
      // Reset success state after 3 seconds
      setTimeout(() => {
        setShowSuccess(false)
      }, 3000)
    },
    onError: (error) => {
      toast.error("Failed to trigger integrated digest")
      console.error("Digest trigger error:", error)
    },
  })

  const handleTriggerDigest = () => {
    triggerDigestMutation.mutate()
  }

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

  const isFiltered = table.getState().columnFilters.length > 0 || table.getState().globalFilter

  return (
    <div className="flex flex-col gap-4">
      {/* Top Bar: Search and Actions */}
      <div className="flex items-center justify-between gap-4">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Search items..."
            value={(table.getState().globalFilter as string) ?? ""}
            onChange={(event) => table.setGlobalFilter(event.target.value)}
            className="pl-8"
          />
        </div>
        <div className="flex items-center gap-2">
          {userId && <UrlSubmissionDialog userId={userId} />}
          <Button
            onClick={handleTriggerDigest}
            disabled={triggerDigestMutation.isPending || showSuccess}
            variant="default"
            size="sm"
            className={`transition-all duration-300 ${
              showSuccess
                ? "bg-green-600 hover:bg-green-600"
                : "bg-black hover:bg-gray-800"
            } text-white`}
          >
            {triggerDigestMutation.isPending ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Creating Digest...
              </>
            ) : showSuccess ? (
              <>
                <Check className="mr-2 h-4 w-4" />
                Digest Started!
              </>
            ) : (
              <>
                <Zap className="mr-2 h-4 w-4" />
                Trigger Digest
              </>
            )}
          </Button>
        </div>
      </div>

      {/* Filters */}
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
              onClick={() => {
                table.resetColumnFilters()
                table.setGlobalFilter("")
              }}
              className="h-8 px-2 lg:px-3"
            >
              Reset
              <X className="ml-2 h-4 w-4" />
            </Button>
          )}
        </div>
      </div>
    </div>
  )
}
