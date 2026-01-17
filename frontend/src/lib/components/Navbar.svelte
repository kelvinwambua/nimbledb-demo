<script lang="ts">
	import { Avatar } from '@skeletonlabs/skeleton-svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { logout, getUser, currentUser } from '$lib/auth';
	import { toaster } from '$lib/toaster';
	import { onMount } from 'svelte';

	let showDropdown = false;

	onMount(async () => {
		try {
			await getUser().catch(() => null);
		} catch {
			currentUser.set(null);
		}
	});

	async function handleLogout() {
		logout();
		toaster.success({
			title: 'Success',
			description: 'Logged out successfully'
		});
		goto('/login');
	}

	function toggleDropdown() {
		showDropdown = !showDropdown;
	}

	function closeDropdown() {
		showDropdown = false;
	}
</script>

<nav class="border-b border-surface-300-700 bg-surface-100-900">
	<div class="container mx-auto px-4">
		<div class="flex h-16 items-center justify-between">
			<a href="/" class="text-xl font-bold text-primary-500">Blog</a>
			<div class="flex items-center gap-4">
				<a
					href="/"
					class="text-surface-600-400 hover:text-primary-500"
					class:text-primary-500={$page.url.pathname === '/'}
				>
					Home
				</a>
				{#if $currentUser && $currentUser.user}
					<a
						href="/posts/my"
						class="text-surface-600-400 hover:text-primary-500"
						class:text-primary-500={$page.url.pathname === '/posts/my'}
					>
						My Posts
					</a>
					<a href="/posts/create" class="btn preset-filled-primary-500 btn-sm"> Create Post </a>
					<div class="relative">
						<button on:click={toggleDropdown} class="flex items-center gap-2">
							<Avatar class="size-10">
								{#if $currentUser.user.image}
									<Avatar.Image src={$currentUser.user.image} alt={$currentUser.user.name} />
								{/if}
								<Avatar.Fallback>{$currentUser.user.name?.charAt(0) || 'U'}</Avatar.Fallback>
							</Avatar>
						</button>
						{#if showDropdown}
							<div class="absolute right-0 z-50 mt-2 w-48 card bg-surface-100-900 shadow-xl">
								<div class="border-b border-surface-300-700 p-4">
									<p class="font-medium">{$currentUser.user.name}</p>
									<p class="text-sm text-surface-600-400">{$currentUser.user.email}</p>
								</div>
								<div class="p-2">
									<button
										on:click={() => {
											closeDropdown();
											handleLogout();
										}}
										class="w-full rounded px-4 py-2 text-left text-surface-600-400 hover:bg-surface-200-800"
									>
										Logout
									</button>
								</div>
							</div>
						{/if}
					</div>
				{:else}
					<a href="/login" class="text-surface-600-400 hover:text-primary-500"> Login </a>
					<a href="/register" class="btn preset-filled-primary-500 btn-sm"> Sign Up </a>
				{/if}
			</div>
		</div>
	</div>
</nav>

<svelte:window
	on:click={(e) => {
		const target = e.target as HTMLElement;
		if (!target.closest('.relative')) {
			closeDropdown();
		}
	}}
/>
