import { useState } from 'react'
import { Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import UrlSubmissionForm from '@/components/url-submission-form'

interface UrlSubmissionDialogProps {
  userId: number
}

export function UrlSubmissionDialog({ userId }: UrlSubmissionDialogProps) {
  const [open, setOpen] = useState(false)

  const handleSuccess = () => {
    setOpen(false)
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline" size="sm">
          <Plus className="mr-2 h-4 w-4" />
          Add URL
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Add New URL</DialogTitle>
          <DialogDescription>
            Submit a URL to process and add to your reading list
          </DialogDescription>
        </DialogHeader>
        <UrlSubmissionForm 
          userId={userId} 
          onSuccess={handleSuccess}
        />
      </DialogContent>
    </Dialog>
  )
}