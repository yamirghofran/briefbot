import { useEffect, useRef, useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { API_BASE_URL } from '@/services/api'

interface ItemUpdateEvent {
  item_id: number
  processing_status: string
  update_type: 'created' | 'processing' | 'completed' | 'failed'
}

/**
 * Hook to subscribe to real-time item updates via Server-Sent Events
 */
export function useItemUpdates(userId: number | undefined) {
  const queryClient = useQueryClient()
  const eventSourceRef = useRef<EventSource | null>(null)
  const [isConnected, setIsConnected] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!userId) return

    console.log(`SSE: Attempting to connect for user ${userId}`)

    let eventSource: EventSource

    try {
      // Create EventSource connection
      eventSource = new EventSource(
        `${API_BASE_URL}/items/user/${userId}/stream`
      )

      eventSourceRef.current = eventSource

      // Handle connection open
      eventSource.onopen = () => {
        console.log('SSE: Connected to item updates stream')
        setIsConnected(true)
        setError(null)
      }

      // Handle item update events
      eventSource.addEventListener('item-update', (event) => {
        try {
          console.log('SSE: Raw event data:', event.data)
          const data: ItemUpdateEvent = JSON.parse(event.data)
          console.log('SSE: Item update received:', data)

          // Invalidate items query to trigger refetch
          queryClient.invalidateQueries({ queryKey: ['items', userId] })

          // Optionally, update specific item in cache
          queryClient.invalidateQueries({
            queryKey: ['item', data.item_id]
          })
        } catch (parseError) {
          console.error('SSE: Error parsing item update:', parseError, 'Raw data:', event.data)
        }
      })

      // Handle generic messages (for keepalive, etc)
      eventSource.onmessage = (event) => {
        console.log('SSE: Generic message received:', event.data)
      }

      // Handle errors
      eventSource.onerror = (event) => {
        console.error('SSE: Connection error occurred', event)
        setIsConnected(false)

        // Check readyState to determine error type
        if (eventSource.readyState === EventSource.CONNECTING) {
          console.log('SSE: Reconnecting...')
          setError('Reconnecting...')
        } else if (eventSource.readyState === EventSource.CLOSED) {
          console.log('SSE: Connection closed')
          setError('Connection closed')
        } else {
          setError('Connection error')
        }
      }
    } catch (err) {
      console.error('SSE: Failed to create EventSource:', err)
      setError(err instanceof Error ? err.message : 'Failed to connect')
    }

    // Cleanup on unmount
    return () => {
      console.log('SSE: Cleaning up connection')
      if (eventSourceRef.current) {
        eventSourceRef.current.close()
        eventSourceRef.current = null
      }
      setIsConnected(false)
    }
  }, [userId, queryClient])

  return {
    isConnected,
    error,
  }
}
