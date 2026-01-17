import { apiRequest } from './api';

export interface Post {
	id: number;
	user_id: number;
	title: string;
	content: string;
	created_at: string;
	updated_at: string;
	author_name?: string;
	author_email?: string;
	author_image?: string;
}

export async function getAllPosts(): Promise<Post[]> {
	const data = await apiRequest('/api/posts/');
	return data.posts || [];
}

export async function getPost(id: number): Promise<Post> {
	const data = await apiRequest(`/api/posts/${id}`);
	return data.post;
}

export async function getMyPosts(): Promise<Post[]> {
	const data = await apiRequest('/api/posts/my/posts');
	return data.posts || [];
}

export async function createPost(title: string, content: string): Promise<Post> {
	const data = await apiRequest('/api/posts/', {
		method: 'POST',
		body: JSON.stringify({ title, content })
	});
	return data.post;
}

export async function updatePost(id: number, title: string, content: string): Promise<Post> {
	const data = await apiRequest(`/api/posts/${id}`, {
		method: 'PUT',
		body: JSON.stringify({ title, content })
	});
	return data.post;
}

export async function deletePost(id: number): Promise<void> {
	await apiRequest(`/api/posts/${id}`, {
		method: 'DELETE'
	});
}
