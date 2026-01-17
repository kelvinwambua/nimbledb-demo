<script lang="ts">
	import { goto } from '$app/navigation';
	import { register } from '$lib/auth';
	import { toaster } from '$lib/toaster';

	let name = '';
	let email = '';
	let password = '';
	let loading = false;

	interface ErrorResponse {
		error?: string;
		message?: string;
	}

	async function handleSubmit(event: Event) {
		event.preventDefault();
		loading = true;
		try {
			await register(email, password, name);
			toaster.success({
				title: 'Success!',
				description: 'Account created successfully'
			});
			// eslint-disable-next-line svelte/no-navigation-without-resolve
			goto('/');
		} catch (error) {
			const err = error as ErrorResponse;
			const errorMessage = err?.error || err?.message || 'Registration failed';
			toaster.error({
				title: 'Error',
				description: errorMessage
			});
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-surface-50-950 p-4">
	<div class="w-full max-w-md">
		<div class="space-y-6 card bg-surface-100-900 p-8 shadow-xl">
			<div class="space-y-2 text-center">
				<h1 class="h1">Create Account</h1>
				<p class="text-surface-600-400">Join us</p>
			</div>
			<form class="space-y-4" on:submit={handleSubmit}>
				<label class="label">
					<span class="label-text">Name</span>
					<input
						class="input"
						type="text"
						placeholder="Enter your name"
						required
						bind:value={name}
						disabled={loading}
					/>
				</label>
				<label class="label">
					<span class="label-text">Email</span>
					<input
						class="input"
						type="email"
						placeholder="you@example.com"
						required
						bind:value={email}
						disabled={loading}
					/>
				</label>
				<label class="label">
					<span class="label-text">Password</span>
					<input
						class="input"
						type="password"
						placeholder="••••••••"
						required
						bind:value={password}
						disabled={loading}
					/>
				</label>
				<button type="submit" class="btn w-full preset-tonal-primary" disabled={loading}>
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
							<span>Registering...</span>
						</div>
					{:else}
						Register
					{/if}
				</button>
			</form>
			<div class="flex items-center gap-4">
				<hr class="flex-1 opacity-50" />
				<span class="text-sm text-surface-600-400">OR</span>
				<hr class="flex-1 opacity-50" />
			</div>
			<p class="text-center text-sm text-surface-600-400">
				Already have an account?
				<!-- eslint-disable-next-line svelte/no-navigation-without-resolve -->
				<a href="/login" class="anchor text-primary-500 hover:text-primary-600"> Sign in </a>
			</p>
		</div>
	</div>
</div>
