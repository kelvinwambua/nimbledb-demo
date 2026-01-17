<script lang="ts">
	import { onMount } from 'svelte';
	import { getMyPosts, type Post } from '$lib/posts';
	import { toaster } from '$lib/toaster';

	let posts: Post[] = [];
	let loading = true;

	onMount(async () => {
		try {
			posts = await getMyPosts();
		} catch (error) {
			const err = error as { error?: string; message?: string };
			toaster.error({
				title: 'Error',
				description: err?.error || err?.message || 'Failed to load your posts'
			});
		} finally {
			loading = false;
		}
	});

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric'
		});
	}
</script>

<div class="container mx-auto px-4 py-8">
	<div class="mb-8 flex items-center justify-between">
		<div>
			<h1 class="mb-2 h1">My Posts</h1>
			<p class="text-surface-600-400">Manage your blog posts</p>
		</div>
		<a href="/posts/create" class="btn preset-filled-primary-500"> Create New Post </a>
	</div>

	{#if loading}
		<div class="flex justify-center py-12">
			<svg
				class="h-12 w-12 animate-spin"
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
			>
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"
				></circle>
				<path
					class="opacity-75"
					fill="currentColor"
					d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
				></path>
			</svg>
		</div>
	{:else if posts.length === 0}
		<div class="card bg-surface-100-900 p-12 text-center">
			<p class="mb-4 text-lg text-surface-600-400">You haven't created any posts yet.</p>
			<a href="/posts/create" class="btn preset-filled-primary-500"> Create Your First Post </a>
		</div>
	{:else}
		<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
			{#each posts as post (post.id)}
				<div class="card bg-surface-100-900 p-6">
					<div class="mb-4">
						<h2 class="mb-2 h3">{post.title}</h2>
						<p class="line-clamp-3 text-surface-600-400">{post.content}</p>
					</div>
					<div
						class="mt-auto flex items-center justify-between border-t border-surface-300-700 pt-4"
					>
						<p class="text-xs text-surface-600-400">{formatDate(post.created_at)}</p>
						<div class="flex gap-2">
							<a href="/posts/{post.id}" class="btn preset-tonal-surface btn-sm"> View </a>
							<a href="/posts/{post.id}/edit" class="btn preset-filled-primary-500 btn-sm">
								Edit
							</a>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
