<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { getPost, updatePost, type Post } from '$lib/posts';
	import { toaster } from '$lib/toaster';

	let post: Post | null = null;
	let title = '';
	let content = '';
	let loading = true;
	let saving = false;

	$: postId = parseInt($page.params.id);

	onMount(async () => {
		try {
			post = await getPost(postId);
			title = post.title;
			content = post.content;
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

	async function handleSubmit(event: Event) {
		event.preventDefault();
		saving = true;
		try {
			await updatePost(postId, title, content);
			toaster.success({
				title: 'Success',
				description: 'Post updated successfully'
			});
			goto(`/posts/${postId}`);
		} catch (error) {
			const err = error as { error?: string; message?: string };
			toaster.error({
				title: 'Error',
				description: err?.error || err?.message || 'Failed to update post'
			});
		} finally {
			saving = false;
		}
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
	{:else}
		<div class="mb-6">
			<a href="/posts/{postId}" class="flex items-center gap-2 anchor">
				<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M15 19l-7-7 7-7"
					/>
				</svg>
				Back to post
			</a>
		</div>

		<div class="card bg-surface-100-900 p-8">
			<h1 class="mb-6 h1">Edit Post</h1>
			<form on:submit={handleSubmit} class="space-y-6">
				<label class="label">
					<span class="label-text text-lg">Title</span>
					<input
						class="input"
						type="text"
						placeholder="Enter post title"
						required
						bind:value={title}
						disabled={saving}
					/>
				</label>

				<label class="label">
					<span class="label-text text-lg">Content</span>
					<textarea
						class="textarea"
						rows="12"
						placeholder="Write your post content..."
						required
						bind:value={content}
						disabled={saving}
					></textarea>
				</label>

				<div class="flex gap-4">
					<button type="submit" class="btn flex-1 preset-filled-primary-500" disabled={saving}>
						{#if saving}
							<div class="flex items-center justify-center gap-2">
								<svg
									class="h-5 w-5 animate-spin"
									xmlns="http://www.w3.org/2000/svg"
									fill="none"
									viewBox="0 0 24 24"
								>
									<circle
										class="opacity-25"
										cx="12"
										cy="12"
										r="10"
										stroke="currentColor"
										stroke-width="4"
									></circle>
									<path
										class="opacity-75"
										fill="currentColor"
										d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
									></path>
								</svg>
								<span>Saving...</span>
							</div>
						{:else}
							Save Changes
						{/if}
					</button>
					<a href="/posts/{postId}" class="btn flex-1 preset-tonal-surface"> Cancel </a>
				</div>
			</form>
		</div>
	{/if}
</div>
