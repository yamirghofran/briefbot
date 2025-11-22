import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { podcastApi, API_BASE_URL } from '@/services/api'
import { StaticDataTable } from '@/components/data-table/static-data-table'
import { podcastColumns } from '@/components/data-table/podcast-columns'
import { Podcast } from '@/types'
import { useEffect, useState } from 'react'

export const Route = createFileRoute('/items/podcasts/')({
  component: PodcastsPage,
})

function PodcastsPage() {
  const userId = 1 // TODO: Get from auth context
  const [podcasts, setPodcasts] = useState<Podcast[]>([])

  // Fetch podcasts
  const { data, isLoading, error } = useQuery({
    queryKey: ['podcasts', userId],
    queryFn: () => podcastApi.getPodcastsByUser(userId),
    refetchInterval: 5000, // Refetch every 5 seconds to update status
  })

  useEffect(() => {
    if (data?.podcasts) {
      setPodcasts(data.podcasts)
    }
  }, [data])

  // Setup SSE for real-time updates
  useEffect(() => {
    const eventSource = new EventSource(`${API_BASE_URL}/podcasts/user/${userId}/stream`)

    eventSource.addEventListener('podcast-update', (event) => {
      const updateData = JSON.parse(event.data)
      console.log('Podcast update received:', updateData)

      // Update the specific podcast in the list
      setPodcasts(prev => prev.map(podcast =>
        podcast.id === updateData.podcast_id
          ? { ...podcast, status: updateData.status }
          : podcast
      ))
    })

    eventSource.onerror = (error) => {
      console.error('SSE Error:', error)
      eventSource.close()
    }

    return () => {
      eventSource.close()
    }
  }, [userId])

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-lg">Loading podcasts...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-lg text-red-600">Error loading podcasts</div>
      </div>
    )
  }

  return (
    <div className="container mx-auto py-10">
      <div className="mb-8">
        <h1 className="text-3xl font-bold">My Podcasts</h1>
        <p className="text-muted-foreground mt-2">
          View and manage your generated podcasts
        </p>
      </div>

      <StaticDataTable
        columns={podcastColumns}
        data={podcasts}
      />
    </div>
  )
}
