import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { itemApi } from "@/services/api";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ArrowLeft, Calendar, ExternalLink, Tag, User } from "lucide-react";
import { formatRelativeDate } from "@/lib/date-utils";

export const Route = createFileRoute("/items/$id")({
	component: ItemDetailPage,
});

function ItemDetailPage() {
	const navigate = useNavigate();
	const { id } = Route.useParams();

	const {
		data: item,
		isLoading,
		error,
	} = useQuery({
		queryKey: ["item", id],
		queryFn: async () => {
			const result = await itemApi.getItem(parseInt(id));
			return result;
		},
		enabled: !!id,
	});

	if (isLoading) {
		return (
			<div className="container mx-auto py-8 flex items-center justify-center min-h-screen">
				<div className="text-center">
					<div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto mb-4"></div>
					<p className="text-muted-foreground">Loading item details...</p>
				</div>
			</div>
		);
	}

	if (error) {
		return (
			<div className="container mx-auto py-8">
				<div className="text-center">
					<h1 className="text-2xl font-bold text-red-600 mb-4">
						Error Loading Item
					</h1>
					<p className="text-muted-foreground mb-2">
						{(error as Error).message}
					</p>
					<p className="text-muted-foreground mb-6">
						The item you're looking for doesn't exist or you don't have access
						to it.
					</p>
					<Button onClick={() => navigate({ to: "/" })} variant="outline">
						<ArrowLeft className="mr-2 h-4 w-4" />
						Back to Items
					</Button>
				</div>
			</div>
		);
	}

	if (!item) {
		return (
			<div className="container mx-auto py-8">
				<div className="text-center">
					<h1 className="text-2xl font-bold text-red-600 mb-4">
						Item Not Found
					</h1>
					<p className="text-muted-foreground mb-6">
						The item you're looking for doesn't exist.
					</p>
					<Button onClick={() => navigate({ to: "/" })} variant="outline">
						<ArrowLeft className="mr-2 h-4 w-4" />
						Back to Items
					</Button>
				</div>
			</div>
		);
	}

	return (
		<div className="container mx-auto py-8 max-w-4xl">
			{/* Header with back button */}
			<div className="mb-8">
				<Button
					onClick={() => navigate({ to: "/" })}
					variant="ghost"
					className="mb-4"
				>
					<ArrowLeft className="mr-2 h-4 w-4" />
					Back to Items
				</Button>
			</div>

			{/* Item Details - Clean Layout */}
			<div className="space-y-8">
				{/* Title and Main Info */}
				<div>
					<div className="flex items-start justify-between mb-4">
						<div className="space-y-2 flex-1">
							<h1 className="text-3xl font-bold text-foreground">
								{item.title || "Untitled Item"}
							</h1>
							{item.url && (
								<p className="text-sm text-muted-foreground">
									<a
										href={item.url}
										target="_blank"
										rel="noopener noreferrer"
										className="text-primary hover:text-primary/80 inline-flex items-center gap-1"
									>
										<ExternalLink className="h-3 w-3" />
										{item.url}
									</a>
								</p>
							)}
						</div>
						{item.is_read && (
							<Badge variant="default" className="ml-4">
								Read
							</Badge>
						)}
					</div>

					{/* Metadata */}
					<div className="flex flex-wrap items-center gap-4 text-sm text-muted-foreground">
						{item.platform && (
							<div className="flex items-center gap-1">
								<span className="font-medium">Platform:</span>
								<Badge variant="outline">{item.platform}</Badge>
							</div>
						)}

						{item.type && (
							<div className="flex items-center gap-1">
								<span className="font-medium">Type:</span>
								<Badge variant="outline">{item.type}</Badge>
							</div>
						)}

						{item.created_at && (
							<div className="flex items-center gap-1">
								<Calendar className="h-3 w-3" />
								<span>{formatRelativeDate(new Date(item.created_at))}</span>
							</div>
						)}
					</div>
				</div>

				{/* Authors */}
				{item.authors && item.authors.length > 0 && (
					<div>
						<h2 className="text-xl font-semibold flex items-center gap-2 mb-3">
							<User className="h-5 w-5" />
							Authors
						</h2>
						<div className="flex flex-wrap gap-2">
							{item.authors.map((author, index) => (
								<Badge key={index} variant="secondary">
									{author}
								</Badge>
							))}
						</div>
					</div>
				)}

				{/* Tags */}
				{item.tags && item.tags.length > 0 && (
					<div>
						<h2 className="text-xl font-semibold flex items-center gap-2 mb-3">
							<Tag className="h-5 w-5" />
							Tags
						</h2>
						<div className="flex flex-wrap gap-2">
							{item.tags.map((tag, index) => (
								<Badge key={index} variant="secondary">
									{tag}
								</Badge>
							))}
						</div>
					</div>
				)}

				{/* Summary */}
				{item.summary && (
					<div>
						<h2 className="text-xl font-semibold mb-3">Summary</h2>
						<p className="text-muted-foreground leading-relaxed">
							{item.summary}
						</p>
					</div>
				)}
			</div>
		</div>
	);
}
