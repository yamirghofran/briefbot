import { createFileRoute } from "@tanstack/react-router";
import { StaticDataTable } from "@/components/data-table/static-data-table";
import { StaticToolbar } from "@/components/data-table/static-toolbar";
import { staticColumns } from "@/components/data-table/static-columns";
import { sampleItems } from "@/components/data-table/static-data";

export const Route = createFileRoute("/")({
  component: App,
});

function App() {
  return (
    <div className="container mx-auto py-10">
      <div className="mb-8">
        <h1 className="text-3xl font-bold tracking-tight">Briefbot</h1>
        <p className="text-muted-foreground mt-2">
          Everything you want to read later in one place.
        </p>
      </div>

      <div className="space-y-4">
        <StaticToolbar />

        <StaticDataTable
          columns={staticColumns}
          data={sampleItems}
        />
      </div>
    </div>
  );
}
