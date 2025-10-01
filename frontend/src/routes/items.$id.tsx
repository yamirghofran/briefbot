import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { itemApi } from "@/services/api";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { ArrowLeft, Calendar, ExternalLink, Tag, User, Edit2, X, Check, Trash2 } from "lucide-react";
import { formatRelativeDate } from "@/lib/date-utils";
import { useState } from "react";

export const Route = createFileRoute("/items/$id")({
	component: ItemDetailPage,
});

function ItemDetailPage() {
	const navigate = useNavigate();
	const { id } = Route.useParams();
	const queryClient = useQueryClient();
	const [isEditing, setIsEditing] = useState(false);
	const [showDeleteDialog, setShowDeleteDialog] = useState(false);
	const [editForm, setEditForm] = useState({
		title: "",
		summary: "",
		tags: [] as string[],
		authors: [] as string[],
	});

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

	const patchMutation = useMutation({
		mutationFn: (data: Partial<{ title: string; summary: string; tags: string[]; authors: string[] }>) =>
			itemApi.patchItem(parseInt(id), data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["item", id] });
			setIsEditing(false);
		},
	});

	const deleteMutation = useMutation({
		mutationFn: () => itemApi.deleteItem(parseInt(id)),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["items"] });
			navigate({ to: "/" });
		},
	});

	const handleEdit = () => {
		setEditForm({
			title: item?.title || "",
			summary: item?.summary || "",
			tags: item?.tags || [],
			authors: item?.authors || [],
		});
		setIsEditing(true);
	};

	const handleCancel = () => {
		setIsEditing(false);
	};

	const handleSave = () => {
		const updates: Partial<{ title: string; summary: string; tags: string[]; authors: string[] }> = {};

		if (editForm.title !== item?.title) updates.title = editForm.title;
		if (editForm.summary !== item?.summary) updates.summary = editForm.summary;
		if (JSON.stringify(editForm.tags) !== JSON.stringify(item?.tags)) updates.tags = editForm.tags;
		if (JSON.stringify(editForm.authors) !== JSON.stringify(item?.authors)) updates.authors = editForm.authors;

		patchMutation.mutate(updates);
	};

	const handleDelete = () => {
		deleteMutation.mutate();
		setShowDeleteDialog(false);
	};

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
				<div className="flex items-center justify-between">
					<Button
						onClick={() => navigate({ to: "/" })}
						variant="ghost"
						className="mb-4"
					>
						<ArrowLeft className="mr-2 h-4 w-4" />
						Back to Items
					</Button>
					{!isEditing ? (
						<div className="flex gap-2">
							<Button onClick={handleEdit} variant="outline" size="sm">
								<Edit2 className="mr-2 h-4 w-4" />
								Edit
							</Button>
							<Button onClick={() => setShowDeleteDialog(true)} variant="destructive" size="sm">
								<Trash2 className="mr-2 h-4 w-4" />
								Delete
							</Button>
						</div>
					) : (
						<div className="flex gap-2">
							<Button onClick={handleSave} variant="default" size="sm" disabled={patchMutation.isPending}>
								<Check className="mr-2 h-4 w-4" />
								Save
							</Button>
							<Button onClick={handleCancel} variant="outline" size="sm">
								<X className="mr-2 h-4 w-4" />
								Cancel
							</Button>
						</div>
					)}
				</div>
			</div>

			{/* Delete Confirmation Dialog */}
			<AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
				<AlertDialogContent>
					<AlertDialogHeader>
						<AlertDialogTitle>Are you sure?</AlertDialogTitle>
						<AlertDialogDescription>
							This action cannot be undone. This will permanently delete the item
							"{item?.title || 'Untitled'}" from your collection.
						</AlertDialogDescription>
					</AlertDialogHeader>
					<AlertDialogFooter>
						<AlertDialogCancel>Cancel</AlertDialogCancel>
						<AlertDialogAction asChild>
							<Button variant="destructive" onClick={handleDelete}>
								Delete
							</Button>
						</AlertDialogAction>
					</AlertDialogFooter>
				</AlertDialogContent>
			</AlertDialog>

			{/* Item Details - Clean Layout */}
			<div className="space-y-8">
				{/* Title and Main Info */}
				<div>
					<div className="flex items-start justify-between mb-4">
						<div className="space-y-2 flex-1">
							{isEditing ? (
								<Input
									value={editForm.title}
									onChange={(e) => setEditForm({ ...editForm, title: e.target.value })}
									className="text-3xl font-bold"
									placeholder="Item title"
								/>
							) : (
								<h1 className="text-3xl font-bold text-foreground">
									{item.title || "Untitled Item"}
								</h1>
							)}
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
				{(item.authors && item.authors.length > 0) || isEditing ? (
					<div>
						<h2 className="text-xl font-semibold flex items-center gap-2 mb-3">
							<User className="h-5 w-5" />
							Authors
						</h2>
						{isEditing ? (
							<div className="space-y-2">
								<div className="flex flex-wrap gap-2">
									{editForm.authors.map((author, index) => (
										<Badge key={index} variant="secondary" className="flex items-center gap-1 pr-1">
											{author}
											<button
												onClick={() =>
													setEditForm({
														...editForm,
														authors: editForm.authors.filter((_, i) => i !== index),
													})
												}
												className="ml-1 hover:bg-background/50 rounded-sm p-0.5"
											>
												<X className="h-3 w-3" />
											</button>
										</Badge>
									))}
								</div>
								<div className="flex gap-2">
									<Input
										placeholder="Add author and press Enter"
										onKeyDown={(e) => {
											if (e.key === "Enter") {
												e.preventDefault();
												const input = e.currentTarget;
												const value = input.value.trim();
												if (value && !editForm.authors.includes(value)) {
													setEditForm({
														...editForm,
														authors: [...editForm.authors, value],
													});
													input.value = "";
												}
											}
										}}
									/>
								</div>
							</div>
						) : (
							<div className="flex flex-wrap gap-2">
								{item.authors.map((author, index) => (
									<Badge key={index} variant="secondary">
										{author}
									</Badge>
								))}
							</div>
						)}
					</div>
				) : null}

				{/* Tags */}
				{(item.tags && item.tags.length > 0) || isEditing ? (
					<div>
						<h2 className="text-xl font-semibold flex items-center gap-2 mb-3">
							<Tag className="h-5 w-5" />
							Tags
						</h2>
						{isEditing ? (
							<div className="space-y-2">
								<div className="flex flex-wrap gap-2">
									{editForm.tags.map((tag, index) => (
										<Badge key={index} variant="secondary" className="flex items-center gap-1 pr-1">
											{tag}
											<button
												onClick={() =>
													setEditForm({
														...editForm,
														tags: editForm.tags.filter((_, i) => i !== index),
													})
												}
												className="ml-1 hover:bg-background/50 rounded-sm p-0.5"
											>
												<X className="h-3 w-3" />
											</button>
										</Badge>
									))}
								</div>
								<div className="flex gap-2">
									<Input
										placeholder="Add tag and press Enter"
										onKeyDown={(e) => {
											if (e.key === "Enter") {
												e.preventDefault();
												const input = e.currentTarget;
												const value = input.value.trim();
												if (value && !editForm.tags.includes(value)) {
													setEditForm({
														...editForm,
														tags: [...editForm.tags, value],
													});
													input.value = "";
												}
											}
										}}
									/>
								</div>
							</div>
						) : (
							<div className="flex flex-wrap gap-2">
								{item.tags.map((tag, index) => (
									<Badge key={index} variant="secondary">
										{tag}
									</Badge>
								))}
							</div>
						)}
					</div>
				) : null}

				{/* Summary */}
				{item.summary || isEditing ? (
					<div>
						<h2 className="text-xl font-semibold mb-3">Summary</h2>
						{isEditing ? (
							<Textarea
								value={editForm.summary}
								onChange={(e) => setEditForm({ ...editForm, summary: e.target.value })}
								placeholder="Enter summary"
								rows={6}
								className="w-full"
							/>
						) : (
							<p className="text-muted-foreground leading-relaxed">
								{item.summary}
							</p>
						)}
					</div>
				) : null}
			</div>
		</div>
	);
}
