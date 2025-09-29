import { Checkbox } from "@/components/ui/checkbox"
import { useMarkAsRead } from "@/hooks/use-mark-as-read"
import type { Item } from "@/types"

interface ReadStatusCellProps {
  item: Item
  userId: number
}

export function ReadStatusCell({ item, userId }: ReadStatusCellProps) {
  const { markAsRead } = useMarkAsRead(userId)

  const handleCheckedChange = (checked: boolean) => {
    if (item.id) {
      markAsRead(item.id)
    }
  }

  return (
    <div className="flex justify-center items-center h-full">
      <Checkbox
        checked={item.is_read || false}
        onCheckedChange={handleCheckedChange}
        className="data-[state=checked]:bg-black data-[state=checked]:border-black border-gray-400 cursor-pointer"
      />
    </div>
  )
}