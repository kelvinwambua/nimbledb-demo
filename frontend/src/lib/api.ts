import { goto } from '$app/navigation';
import { browser } from '$app/environment';

const API_URL = import.meta.env.VITE_API_URL;

interface FetchOptions extends RequestInit {
	skipAuthRedirect?: boolean;
}

export async function apiRequest(endpoint: string, options: FetchOptions = {}) {
	const { skipAuthRedirect = false, ...fetchOptions } = options;

	const response = await fetch(`${API_URL}${endpoint}`, {
		credentials: 'include',
		headers: {
			'Content-Type': 'application/json',
			...fetchOptions.headers
		},
		...fetchOptions
	});

	if (response.status === 401 && !skipAuthRedirect && browser) {
		const currentPath = window.location.pathname;
		const isAuthPage = currentPath.startsWith('/login') || currentPath.startsWith('/register');

		if (!isAuthPage) {
			// eslint-disable-next-line svelte/no-navigation-without-resolve
			goto('/login');
		}
	}

	if (!response.ok) {
		const errorData = await response.json().catch(() => ({
			message: 'An error occurred'
		}));
		throw errorData;
	}

	return response.json();
}
