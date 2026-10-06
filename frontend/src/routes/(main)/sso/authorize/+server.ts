import {
	PUBLIC_ROOT_DOMAIN,
	PUBLIC_ROOT_URL
} from '$app/env/public';
import { error, redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { API_URL } from '$app/env/private';

type AuthorizeResponse = {
	code: string;
	target_host: string;
	return_path: string;
	expires_in: number;
};

const statePattern = /^[A-Za-z0-9_-]{32,128}$/;

export const GET = (async ({ cookies, fetch, url }) => {
	const targetHost = (
		url.searchParams.get('target_host') ?? ''
	)
		.trim()
		.toLowerCase();

	const returnPath =
		url.searchParams.get('return_path') ?? '/';

	const state = url.searchParams.get('state') ?? '';

	if (url.hostname !== PUBLIC_ROOT_DOMAIN) {
		const centralURL = new URL(
			'/sso/authorize',
			PUBLIC_ROOT_URL
		);

		centralURL.searchParams.set(
			'target_host',
			targetHost
		);
		centralURL.searchParams.set(
			'return_path',
			returnPath
		);
		centralURL.searchParams.set('state', state);

		redirect(307, centralURL, {
			external: [new URL(PUBLIC_ROOT_URL).origin]
		});
	}

	if (!targetHost) {
		error(400, 'Target hostname wajib diisi');
	}

	if (!statePattern.test(state)) {
		error(400, 'State SSO tidak valid');
	}

	const session = cookies.get('session');

	if (!session) {
		const next = url.pathname + url.search;

		redirect(
			307,
			`/login?next=${encodeURIComponent(next)}`
		);
	}

	let response: Response;

	try {
		response = await fetch(
			`${API_URL}/api/sso/authorize`,
			{
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Cookie: `session=${session}`
				},
				body: JSON.stringify({
					target_host: targetHost,
					return_path: returnPath
				})
			}
		);
	} catch {
		error(503, 'Server sedang tidak tersedia');
	}

	if (response.status === 401) {
		cookies.delete('session', {
			path: '/'
		});

		const next = url.pathname + url.search;

		redirect(
			307,
			`/login?next=${encodeURIComponent(next)}`
		);
	}

	if (!response.ok) {
		const message = await response.text();

		error(
			response.status,
			message.trim() || 'Gagal membuat authorization code'
		);
	}

	const authorization =
		(await response.json()) as AuthorizeResponse;

	const rootURL = new URL(PUBLIC_ROOT_URL);
	const callbackURL = new URL('/auth/callback', rootURL);

	callbackURL.hostname = authorization.target_host;
	callbackURL.searchParams.set(
		'code',
		authorization.code
	);
	callbackURL.searchParams.set('state', state);

	redirect(303, callbackURL, {
		external: true
	});
}) satisfies RequestHandler;