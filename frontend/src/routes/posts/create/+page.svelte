<script lang="ts">
	import { goto } from '$app/navigation';
	import { createPost } from '$lib/posts';
	import { toaster } from '$lib/toaster';

	let title = '';
	let content = '';
	let loading = false;

	async function handleSubmit(event: Event) {
		event.preventDefault();
		loading = true;
		try {
			const post = await createPost(title, content);
			toaster.success({
				title: 'Success',
				description: 'Post created successfully'
			});
			goto(`/posts/${post.id}`);
		} catch (error) {
			const err = error as { error?: string; message?: string };
			toaster.error({
				title: 'Error',
				description: err?.error || err?.message || 'Failed to create post'
			});
		} finally {
			loading = false;
		}
	}
</script>

<div class="container mx-auto max-w-4xl px-4 py-8">
	<div class="mb-6">
		<a href="/" class="flex items-center gap-2 anchor">
			<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
			</svg>
			Back to all posts
		</a>
	</div>

	<div class="card bg-surface-100-900 p-8">
		<h1 class="mb-6 h1">Create New Post</h1>
		<form on:submit={handleSubmit} class="space-y-6">
			<label class="label">
				<span class="label-text text-lg">Title</span>
				<input
					class="input"
					type="text"
					placeholder="Enter post title"
					required
					bind:value={title}
					disabled={loading}
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
					disabled={loading}
				></textarea>
			</label>

			<div class="flex gap-4">
				<button type="submit" class="btn flex-1 preset-filled-primary-500" disabled={loading}>
					{#if loading}
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
							<span>Creating...</span>
						</div>
					{:else}
						Create Post
					{/if}
				</button>
				<a href="/" class="btn flex-1 preset-tonal-surface"> Cancel </a>
			</div>
		</form>
	</div>
</div>
