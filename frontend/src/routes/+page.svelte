<script lang="ts">
	import { onMount } from 'svelte';
	import { getAllPosts, type Post } from '$lib/posts';
	import { currentUser } from '$lib/auth';
	import { toaster } from '$lib/toaster';
	let posts: Post[] = [];
	let loading = true;
	onMount(async () => {
		try {
			posts = await getAllPosts();
		} catch (error) {
			const err = error as { error?: string; message?: string };
			toaster.error({
				title: 'Error',
				description: err?.error || err?.message || 'Failed to load posts'
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
	<div class="mb-8">
		<h1 class="mb-2 h1">All Posts</h1>
		<p class="text-surface-600-400">Browse all blog posts</p>
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
			<p class="text-lg text-surface-600-400">No posts yet. Be the first to create one!</p>
		</div>
	{:else}
		<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
			{#each posts as post (post.id)}
				<a
					href="/posts/{post.id}"
					class="card bg-surface-100-900 p-6 transition-shadow hover:shadow-lg"
				>
					<div class="mb-4">
						<h2 class="mb-2 h3">{post.title}</h2>
						<p class="line-clamp-3 text-surface-600-400">{post.content}</p>
					</div>
					<div class="mt-auto flex items-center gap-3 border-t border-surface-300-700 pt-4">
						{#if post.author_image}
							<img src={post.author_image} alt={post.author_name} class="h-10 w-10 rounded-full" />
						{:else}
							<div
								class="flex h-10 w-10 items-center justify-center rounded-full bg-primary-500 font-bold text-white"
							>
								{post.author_name?.charAt(0) || 'U'}
							</div>
						{/if}
						<div class="flex-1">
							<p class="text-sm font-medium">{post.author_name || 'Unknown'}</p>
							<p class="text-xs text-surface-600-400">{formatDate(post.created_at)}</p>
						</div>
						{#if $currentUser && $currentUser.user && $currentUser.user.id === post.user_id}
							<span class="badge preset-filled-primary-500">Your Post</span>
						{/if}
					</div>
				</a>
			{/each}
		</div>
	{/if}
</div>
