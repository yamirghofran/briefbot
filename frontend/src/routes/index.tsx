import { createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { StaticToolbar } from "@/components/data-table/static-toolbar";
import { ItemDataTable } from "@/components/item-data-table";
import { itemApi } from "@/services/api";
import { Button } from "@/components/ui/button";
import type { Item } from "@/types";

export const Route = createFileRoute("/")({
  component: App,
});

function App() {
  // For now, we'll use a hardcoded user ID. In a real app, this would come from auth context
  const userId = 1; // TODO: Get this from auth context

  // Fetch real items data
  const { data: items = [], isLoading, error, refetch } = useQuery({
    queryKey: ['items', userId],
    queryFn: () => itemApi.getItemsByUser(userId),
    enabled: !!userId,
    refetchInterval: 3000, // Poll every 3 seconds to check for status updates
    staleTime: 5000, // Consider data fresh for 5 seconds
    cacheTime: 60000, // Keep cache for 1 minute
    refetchOnWindowFocus: false, // Don't refetch on window focus
    refetchOnReconnect: false, // Don't refetch on reconnect
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
        <StaticToolbar userId={userId} />

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
          <div className="text-center p-8 text-muted-foreground">
            <p>No items found. Submit a URL above to get started!</p>
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
