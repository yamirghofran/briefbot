import { useAppForm } from '@/hooks/demo.form'
import { itemApi } from '@/services/api'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { z } from 'zod'

interface UrlSubmissionFormProps {
  userId: number
  onSuccess?: () => void
}

export default function UrlSubmissionForm({ userId, onSuccess }: UrlSubmissionFormProps) {
  const queryClient = useQueryClient()

  // Mutation for submitting URL
  const submitUrlMutation = useMutation({
    mutationFn: (url: string) => itemApi.submitUrl({ user_id: userId, url }),
    onSuccess: () => {
      toast.success('URL submitted successfully!')
      queryClient.invalidateQueries({ queryKey: ['items', userId] })
      onSuccess?.()
    },
    onError: (error) => {
      toast.error('Failed to submit URL. Please try again.')
      console.error('URL submission error:', error)
    },
  })

  const form = useAppForm({
    defaultValues: {
      url: '',
    },
    validators: {
      onChange: z.object({
        url: z.string().url('Please enter a valid URL'),
      }),
    },
    onSubmit: async ({ value }) => {
      try {
        await submitUrlMutation.mutateAsync(value.url)
        form.reset()
      } catch (error) {
        console.error('URL submission error:', error)
      }
    },
  })

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault()
        e.stopPropagation()
        form.handleSubmit()
      }}
      className="space-y-4"
    >
      <form.AppField name="url">
        {(field) => (
           <field.TextField 
             placeholder="https://example.com/article"
             disabled={submitUrlMutation.isPending}
           />
        )}
      </form.AppField>

      <div className="flex justify-end">
        <form.AppForm>
          <form.SubscribeButton 
            label={submitUrlMutation.isPending ? 'Processing...' : 'Submit URL'} 
          />
        </form.AppForm>
      </div>
    </form>
  )
}