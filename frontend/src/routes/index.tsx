import { createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { ItemDataTable } from "@/components/item-data-table";
import { itemApi } from "@/services/api";
import { UrlSubmissionDialog } from "@/components/url-submission-dialog";
import { Button } from "@/components/ui/button";
import { useItemUpdates } from "@/hooks/use-item-updates";
import type { Item } from "@/types";

export const Route = createFileRoute("/")({
  component: App,
});

function App() {
  // For now, we'll use a hardcoded user ID. In a real app, this would come from auth context
  const userId = 1; // TODO: Get this from auth context

  // Subscribe to real-time item updates via SSE
  const { isConnected } = useItemUpdates(userId);

  // Fetch real items data (no polling - updates come via SSE)
  const { data: items = [], isLoading, error, refetch } = useQuery({
    queryKey: ['items', userId],
    queryFn: () => itemApi.getItemsByUser(userId),
    enabled: !!userId,
    staleTime: Infinity, // Data is always fresh since we get SSE updates
    refetchOnWindowFocus: true, // Refetch on window focus as fallback
  });

  return (
    <div className="container mx-auto py-10">
      <div className="mb-8">
        <h1 className="text-3xl font-bold tracking-tight">Briefbot</h1>
        <p className="text-muted-foreground mt-2">
          Everything you want to read later in one place.
        </p>
      </div>

      {/* Data Table Section */}
      <div className="space-y-4">
        {isLoading ? (
          <div className="flex items-center justify-center h-64">
            <div className="text-center">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto mb-4"></div>
              <p className="text-muted-foreground">Loading your items...</p>
            </div>
          </div>
        ) : error ? (
          <div className="text-center p-8 text-red-500">
            <p>Error loading items: {error.message}</p>
            <Button 
              onClick={() => window.location.reload()} 
              variant="outline" 
              className="mt-4"
            >
              Retry
            </Button>
          </div>
        ) : items.length === 0 ? (
          <div className="flex flex-col items-center gap-4 rounded-lg border border-dashed p-8 text-center text-muted-foreground">
            <p>No items yet. Add one to kick things off.</p>
            <UrlSubmissionDialog userId={userId} />
          </div>
        ) : (
          <ItemDataTable
            data={items}
            userId={userId}
          />
        )}
      </div>
    </div>
  );
}
