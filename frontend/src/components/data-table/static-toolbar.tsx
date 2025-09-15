import * as React from "react"
import { Search, Filter } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"

interface StaticToolbarProps extends React.HTMLAttributes<HTMLDivElement> {}

export function StaticToolbar({ className, ...props }: StaticToolbarProps) {
  const [searchValue, setSearchValue] = React.useState("")

  return (
    <div
      className="flex items-center justify-between"
      {...props}
    >
      <div className="flex flex-1 items-center space-x-2">
        <div className="relative">
          <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Search..."
            value={searchValue}
            onChange={(event) => setSearchValue(event.target.value)}
            className="pl-8 w-[150px] lg:w-[250px]"
          />
        </div>
        <Button variant="outline" size="sm">
          <Filter className="mr-2 h-4 w-4" />
          Filter
        </Button>
      </div>
      
      <div className="flex items-center space-x-2">
        <Button variant="outline" size="sm">
          Columns
        </Button>
        <Button variant="outline" size="sm">
          Export
        </Button>
      </div>
    </div>
  )
}