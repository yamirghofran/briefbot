import type { Column } from "@tanstack/react-table";

export function getCommonPinningStyles<TData>({
  column,
  withBorder = false,
}: {
  column: Column<TData>;
  withBorder?: boolean;
}): React.CSSProperties {
  // Simplified version without actual pinning functionality
  // This is just for styling compatibility
  return {};
}