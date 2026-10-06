import { error, redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { API_URL } from '$app/env/private';

type StartResponse = {
	authorization_url: string;
};

export const GET = (async ({ fetch, url }) => {
	const next = url.searchParams.get('next') ?? '/admin';
	const endpoint = new URL(
		`${API_URL}/api/auth/google/start`
	);

	endpoint.searchParams.set('next', next);

	let response: Response;

	try {
		response = await fetch(endpoint);
	} catch {
		error(503, 'Server sedang tidak tersedia');
	}

	if (!response.ok) {
		const message = await response.text();
		error(response.status, message.trim() || 'Gagal memulai login Google');
	}

	const { authorization_url: authorizationURL } =
		(await response.json()) as StartResponse;

	if (!authorizationURL) {
		error(500, 'Authorization URL Google tidak ditemukan');
	}

	redirect(303, authorizationURL, { external: true });
}) satisfies RequestHandler;
