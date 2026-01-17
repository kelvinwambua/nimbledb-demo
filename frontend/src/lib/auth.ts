import { apiRequest } from './api';
import { browser } from '$app/environment';
import { writable } from 'svelte/store';

export const currentUser = writable<any>(null);

export async function register(email: string, password: string, name: string) {
	const response = await apiRequest('/api/auth/register', {
		method: 'POST',
		body: JSON.stringify({ email, password, name }),
		skipAuthRedirect: true
	});
	currentUser.set(response);
	return response;
}

export async function login(email: string, password: string) {
	const response = await apiRequest('/api/auth/login', {
		method: 'POST',
		body: JSON.stringify({ email, password }),
		skipAuthRedirect: true
	});
	currentUser.set(response);
	return response;
}

export async function logout() {
	if (browser) {
		document.cookie = 'nimbledb-test_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
	}
	currentUser.set(null);
}

export async function getUser() {
	const response = await apiRequest('/api/auth/me', {
		method: 'GET'
	});
	currentUser.set(response);
	return response;
}
