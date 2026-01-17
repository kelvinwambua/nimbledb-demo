<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { getPost, deletePost, type Post } from '$lib/posts';
	import { currentUser } from '$lib/auth';
	import { toaster } from '$lib/toaster';

	let post: Post | null = null;
	let loading = true;
	let deleting = false;

	$: postId = parseInt($page.params.id);
	$: isOwner = $currentUser && $currentUser.user && post && $currentUser.user.id === post.user_id;

	onMount(async () => {
		try {
			post = await getPost(postId);
		} catch (error) {
			const err = error as { error?: string; message?: string };
			toaster.error({
				title: 'Error',
				description: err?.error || err?.message || 'Failed to load post'
			});
			goto('/');
		} finally {
			loading = false;
		}
	});

	async function handleDelete() {
		if (!confirm('Are you sure you want to delete this post?')) return;

		deleting = true;
		try {
			await deletePost(postId);
			toaster.success({
				title: 'Success',
				description: 'Post deleted successfully'
			});
			goto('/');
		} catch (error) {
			const err = error as { error?: string; message?: string };
			toaster.error({
				title: 'Error',
				description: err?.error || err?.message || 'Failed to delete post'
			});
		} finally {
			deleting = false;
		}
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}
</script>

<div class="container mx-auto max-w-4xl px-4 py-8">
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
	{:else if post}
		<div class="mb-6">
			<a href="/" class="flex items-center gap-2 anchor">
				<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M15 19l-7-7 7-7"
					/>
				</svg>
				Back to all posts
			</a>
		</div>

		<article class="card bg-surface-100-900 p-8">
			<div class="mb-6 flex items-start justify-between">
				<div class="flex-1">
					<h1 class="mb-4 h1">{post.title}</h1>
					<div class="flex items-center gap-3">
						{#if post.author_image}
							<img src={post.author_image} alt={post.author_name} class="h-12 w-12 rounded-full" />
						{:else}
							<div
								class="flex h-12 w-12 items-center justify-center rounded-full bg-primary-500 font-bold text-white"
							>
								{post.author_name?.charAt(0) || 'U'}
							</div>
						{/if}
						<div>
							<p class="font-medium">{post.author_name || 'Unknown'}</p>
							<p class="text-sm text-surface-600-400">{formatDate(post.created_at)}</p>
						</div>
					</div>
				</div>
				{#if isOwner}
					<div class="flex gap-2">
						<a href="/posts/{post.id}/edit" class="btn preset-filled-primary-500 btn-sm"> Edit </a>
						<button
							on:click={handleDelete}
							disabled={deleting}
							class="btn preset-filled-error-500 btn-sm"
						>
							{deleting ? 'Deleting...' : 'Delete'}
						</button>
					</div>
				{/if}
			</div>

			<div class="prose max-w-none">
				<p class="whitespace-pre-wrap">{post.content}</p>
			</div>

			{#if post.updated_at !== post.created_at}
				<div class="mt-6 border-t border-surface-300-700 pt-6">
					<p class="text-sm text-surface-600-400">
						Last updated: {formatDate(post.updated_at)}
					</p>
				</div>
			{/if}
		</article>
	{/if}
</div>
