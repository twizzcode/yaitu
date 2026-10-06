import { error, redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { API_URL } from '$app/env/private';

type CallbackResponse = {
	session_token: string;
	expires_at: string;
	next_path: string;
};

export const GET = (async ({ cookies, fetch, url }) => {
	const providerError = url.searchParams.get('error');

	if (providerError) {
		error(400, 'Login Google dibatalkan atau ditolak');
	}

	const code = url.searchParams.get('code')?.trim() ?? '';
	const state = url.searchParams.get('state')?.trim() ?? '';

	if (!code || !state) {
		error(400, 'Authorization code dan state tidak ditemukan');
	}

	let response: Response;

	try {
		response = await fetch(
			`${API_URL}/api/auth/google/callback`,
			{
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ code, state })
			}
		);
	} catch {
		error(503, 'Server sedang tidak tersedia');
	}

	if (!response.ok) {
		const message = await response.text();
		error(response.status, message.trim() || 'Login Google gagal');
	}

	const result = (await response.json()) as CallbackResponse;
	const expiresAt = new Date(result.expires_at);

	if (!result.session_token || Number.isNaN(expiresAt.getTime())) {
		error(500, 'Session Google tidak valid');
	}

	cookies.set('session', result.session_token, {
		path: '/',
		httpOnly: true,
		sameSite: 'lax',
		secure: url.protocol === 'https:',
		expires: expiresAt
	});

	redirect(303, result.next_path);
}) satisfies RequestHandler;
