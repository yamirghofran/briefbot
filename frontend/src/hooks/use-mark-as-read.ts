import { useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { itemApi } from '@/services/api'
import type { Item } from '@/types'
import { useCallback, useRef } from 'react'

export function useMarkAsRead(userId: number) {
  const queryClient = useQueryClient()
  const debounceRef = useRef<NodeJS.Timeout | null>(null)

  const markAsReadMutation = useMutation({
    mutationFn: (itemId: number) => itemApi.toggleItemReadStatus(itemId),
    // Optimistic update - immediately update UI before server response
    onMutate: async (itemId: number) => {
      // Cancel any outgoing refetches
      await queryClient.cancelQueries({ queryKey: ['items', userId] })
      
      // Snapshot the previous value
      const previousItems = queryClient.getQueryData(['items', userId])
      
      // Optimistically update to the new value
      queryClient.setQueryData(['items', userId], (oldData: Item[] | undefined) => {
        if (!oldData) return []
        return oldData.map(item => 
          item.id === itemId ? { ...item, is_read: !item.is_read } : item
        )
      })
      
      // Return a context object with the snapshotted value
      return { previousItems }
    },
    onError: (error, variables, context) => {
      // If the mutation fails, roll back to the previous value
      if (context?.previousItems) {
        queryClient.setQueryData(['items', userId], context.previousItems)
      }
      toast.error('Failed to update read status')
      console.error('Mark as read error:', error)
    },
    onSettled: (data, error, variables) => {
      // Show success message after successful update
      if (data && !error) {
        toast.success(data.is_read ? 'Marked as read' : 'Marked as unread')
      }
    },
  })

  // Debounced mark as read function to prevent rapid clicks
  const debouncedMarkAsRead = useCallback((itemId: number) => {
    if (debounceRef.current) {
      clearTimeout(debounceRef.current)
    }
    
    debounceRef.current = setTimeout(() => {
      markAsReadMutation.mutate(itemId)
    }, 100) // 100ms debounce
  }, [markAsReadMutation])

  return {
    markAsRead: debouncedMarkAsRead,
    isLoading: markAsReadMutation.isPending,
  }
}