import { createFileRoute } from "@tanstack/react-router";
import { StaticDataTable } from "@/components/data-table/static-data-table";
import { StaticToolbar } from "@/components/data-table/static-toolbar";
import { staticColumns } from "@/components/data-table/static-columns";
import { sampleItems } from "@/components/data-table/static-data";

export const Route = createFileRoute("/demo/table-styled")({
  component: DemoTableStyled,
});

function DemoTableStyled() {
  return (
    <div className="container mx-auto py-10">
      <div className="mb-8">
        <h1 className="text-3xl font-bold tracking-tight">Items Data Table</h1>
        <p className="text-muted-foreground mt-2">
          A beautifully styled data table for Items with URL, content type, and read status.
        </p>
      </div>

      <div className="space-y-4">
        <StaticToolbar userId={undefined} />

        <StaticDataTable
          columns={staticColumns}
          data={sampleItems}
        />
      </div>

      <div className="mt-8 rounded-lg border p-4">
        <h3 className="font-medium mb-2">Table Schema:</h3>
        <ul className="text-sm text-muted-foreground space-y-1">
          <li>✅ Title, URL, Content Type</li>
          <li>✅ Platform, Authors, Tags</li>
          <li>✅ Read/Unread Status</li>
          <li>✅ Created/Modified Dates</li>
          <li>✅ Text Content & Summary</li>
        </ul>
      </div>
    </div>
  );
}
