import { createFileRoute } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { itemApi } from '@/services/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { useState } from 'react'
import type { CreateItemRequest } from '@/types'

export const Route = createFileRoute('/items/')({
  component: ItemsPage,
})

function ItemsPage() {
  const queryClient = useQueryClient()
  const [showForm, setShowForm] = useState(false)
  const [selectedUserId, setSelectedUserId] = useState<number>(1) // Default user ID
  const [filter, setFilter] = useState<'all' | 'unread'>('all')
  const [formData, setFormData] = useState<CreateItemRequest>({
    user_id: 1,
    url: '',
    text_content: '',
    summary: '',
    type: '',
    platform: '',
    tags: [],
    authors: [],
  })

  // Fetch items based on filter
  const { data: items = [], isLoading, error, refetch } = useQuery({
    queryKey: ['items', selectedUserId, filter],
    queryFn: async () => {
      if (filter === 'unread') {
        return await itemApi.getUnreadItemsByUser(selectedUserId)
      } else {
        return await itemApi.getItemsByUser(selectedUserId)
      }
    },
  })

  const createItemMutation = useMutation({
    mutationFn: itemApi.createItem,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['items', selectedUserId] })
      setShowForm(false)
      setFormData({ user_id: selectedUserId, url: '', text_content: '', summary: '', type: '', platform: '', tags: [], authors: [] })
    },
  })

  const markAsReadMutation = useMutation({
    mutationFn: itemApi.markItemAsRead,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['items', selectedUserId] })
    },
  })

  const deleteItemMutation = useMutation({
    mutationFn: itemApi.deleteItem,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['items', selectedUserId] })
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (formData.url && formData.text_content) {
      const data = {
        ...formData,
        tags: formData.tags?.length ? formData.tags : undefined,
        authors: formData.authors?.length ? formData.authors : undefined,
      }
      createItemMutation.mutate(data)
    }
  }

  const handleMarkAsRead = (itemId: number) => {
    markAsReadMutation.mutate(itemId)
  }

  const handleDeleteItem = (itemId: number) => {
    if (confirm('Are you sure you want to delete this item?')) {
      deleteItemMutation.mutate(itemId)
    }
  }

  if (isLoading) return <div>Loading...</div>
  if (error) return <div>Error loading items</div>

  return (
    <div className="container mx-auto p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Items</h1>
        <Button onClick={() => setShowForm(!showForm)}>
          {showForm ? 'Cancel' : 'Create Item'}
        </Button>
      </div>

      <div className="mb-6 flex gap-4 items-center">
        <div>
          <Label htmlFor="userId">User ID</Label>
          <Input
            id="userId"
            type="number"
            value={selectedUserId}
            onChange={(e) => setSelectedUserId(parseInt(e.target.value) || 1)}
            className="w-24"
          />
        </div>
        <div className="flex gap-2">
          <Button
            variant={filter === 'all' ? 'default' : 'outline'}
            onClick={() => setFilter('all')}
          >
            All Items
          </Button>
          <Button
            variant={filter === 'unread' ? 'default' : 'outline'}
            onClick={() => setFilter('unread')}
          >
            Unread Only
          </Button>
        </div>
        <Button onClick={() => refetch()} variant="outline">
          Refresh
        </Button>
      </div>

      {showForm && (
        <div className="bg-white p-6 rounded-lg shadow-md mb-6">
          <h2 className="text-lg font-semibold mb-4">Create New Item</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <Label htmlFor="user_id">User ID</Label>
              <Input
                id="user_id"
                type="number"
                value={formData.user_id}
                onChange={(e) => setFormData({ ...formData, user_id: parseInt(e.target.value) || 1 })}
                required
              />
            </div>
            <div>
              <Label htmlFor="url">URL</Label>
              <Input
                id="url"
                type="url"
                value={formData.url}
                onChange={(e) => setFormData({ ...formData, url: e.target.value })}
                required
              />
            </div>
             <div>
               <Label htmlFor="text_content">Text Content</Label>
               <Textarea
                 id="text_content"
                 value={formData.text_content}
                 onChange={(e) => setFormData({ ...formData, text_content: e.target.value })}
                 rows={4}
                 required
               />
             </div>
             <div>
               <Label htmlFor="summary">Summary</Label>
               <Textarea
                 id="summary"
                 value={formData.summary || ''}
                 onChange={(e) => setFormData({ ...formData, summary: e.target.value })}
                 rows={2}
               />
             </div>
             <div>
               <Label htmlFor="type">Type</Label>
               <Input
                 id="type"
                 type="text"
                 value={formData.type || ''}
                 onChange={(e) => setFormData({ ...formData, type: e.target.value })}
               />
             </div>
             <div>
               <Label htmlFor="platform">Platform</Label>
               <Input
                 id="platform"
                 type="text"
                 value={formData.platform || ''}
                 onChange={(e) => setFormData({ ...formData, platform: e.target.value })}
               />
             </div>
             <div>
               <Label htmlFor="tags">Tags (comma-separated)</Label>
               <Input
                 id="tags"
                 type="text"
                 value={formData.tags?.join(', ') || ''}
                 onChange={(e) => setFormData({ ...formData, tags: e.target.value.split(',').map(t => t.trim()).filter(t => t) })}
               />
             </div>
             <div>
               <Label htmlFor="authors">Authors (comma-separated)</Label>
               <Input
                 id="authors"
                 type="text"
                 value={formData.authors?.join(', ') || ''}
                 onChange={(e) => setFormData({ ...formData, authors: e.target.value.split(',').map(a => a.trim()).filter(a => a) })}
               />
             </div>
            <Button type="submit" disabled={createItemMutation.isPending}>
              {createItemMutation.isPending ? 'Creating...' : 'Create Item'}
            </Button>
          </form>
        </div>
      )}

      <div className="bg-white rounded-lg shadow-md">
        <div className="p-6">
          <h2 className="text-lg font-semibold mb-4">
            {filter === 'unread' ? 'Unread Items' : 'All Items'} ({items.length})
          </h2>
          {items.length === 0 ? (
            <p className="text-gray-500">
              {filter === 'unread' ? 'No unread items found.' : 'No items found. Create an item to get started.'}
            </p>
          ) : (
            <div className="space-y-4">
              {items.map((item) => (
                <div key={item.id} className="border rounded p-4 hover:bg-gray-50">
                  <div className="flex justify-between items-start mb-2">
                    <div className="flex-1">
                      <h3 className="font-medium text-blue-600 hover:underline">
                        <a href={item.url || '#'} target="_blank" rel="noopener noreferrer">
                          {item.url}
                        </a>
                      </h3>
                      <p className="text-sm text-gray-600 mt-1">
                        User ID: {item.user_id} | Status: {item.is_read ? 'Read' : 'Unread'}
                      </p>
                    </div>
                    <div className="space-x-2">
                      {!item.is_read && (
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleMarkAsRead(item.id)}
                          disabled={markAsReadMutation.isPending}
                        >
                          Mark as Read
                        </Button>
                      )}
                      <Button variant="outline" size="sm">
                        View
                      </Button>
                      <Button variant="outline" size="sm" className="text-red-600">
                        Edit
                      </Button>
                      <Button
                        variant="outline"
                        size="sm"
                        className="text-red-600"
                        onClick={() => handleDeleteItem(item.id)}
                        disabled={deleteItemMutation.isPending}
                      >
                        Delete
                      </Button>
                    </div>
                  </div>
                  {item.text_content && (
                    <div className="mt-3">
                      <p className="text-sm text-gray-700 line-clamp-3">{item.text_content}</p>
                    </div>
                  )}
                   {item.summary && (
                     <div className="mt-3 p-3 bg-gray-50 rounded">
                       <h4 className="text-sm font-medium text-gray-900 mb-1">Summary:</h4>
                       <p className="text-sm text-gray-700">{item.summary}</p>
                     </div>
                   )}
                   {item.type && (
                     <div className="mt-2">
                       <span className="text-xs text-gray-500">Type: {item.type}</span>
                     </div>
                   )}
                   {item.platform && (
                     <div className="mt-2">
                       <span className="text-xs text-gray-500">Platform: {item.platform}</span>
                     </div>
                   )}
                   {item.tags && item.tags.length > 0 && (
                     <div className="mt-2">
                       <span className="text-xs text-gray-500">Tags: {item.tags.join(', ')}</span>
                     </div>
                   )}
                   {item.authors && item.authors.length > 0 && (
                     <div className="mt-2">
                       <span className="text-xs text-gray-500">Authors: {item.authors.join(', ')}</span>
                     </div>
                   )}
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}