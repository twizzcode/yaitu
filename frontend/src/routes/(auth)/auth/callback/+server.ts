import { error, redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

type ExchangeResponse = {
	session_token: string;
	return_path: string;
	expires_at: string;
};

const statePattern = /^[A-Za-z0-9_-]{32,128}$/;

export const GET = (async ({ cookies, fetch, url }) => {
	const code = (
		url.searchParams.get('code') ?? ''
	).trim();

	const returnedState =
		url.searchParams.get('state') ?? '';

	const storedState = cookies.get('sso_state') ?? '';

	if (!code) {
		error(400, 'Authorization code tidak ditemukan');
	}

	if (
		!statePattern.test(returnedState) ||
		!statePattern.test(storedState) ||
		returnedState !== storedState
	) {
		cookies.delete('sso_state', {
			path: '/auth/callback'
		});

		error(400, 'State SSO tidak valid');
	}

	const targetHost = url.hostname.toLowerCase();

	let response: Response;

	try {
		response = await fetch(
			'http://localhost:8080/api/sso/exchange',
			{
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					code,
					target_host: targetHost
				})
			}
		);
	} catch {
		error(503, 'Server sedang tidak tersedia');
	}

	if (!response.ok) {
		cookies.delete('sso_state', {
			path: '/auth/callback'
		});

		const message = await response.text();

		error(
			response.status,
			message.trim() || 'Gagal menukar authorization code'
		);
	}

	const exchange =
		(await response.json()) as ExchangeResponse;

	const expiresAt = new Date(exchange.expires_at);

	if (Number.isNaN(expiresAt.getTime())) {
		error(500, 'Waktu kedaluwarsa session tidak valid');
	}

	cookies.set(
		'session',
		exchange.session_token,
		{
			path: '/',
			httpOnly: true,
			sameSite: 'lax',
			secure: url.protocol === 'https:',
			expires: expiresAt
		}
	);

	cookies.delete('sso_state', {
		path: '/auth/callback'
	});

	redirect(303, exchange.return_path);
}) satisfies RequestHandler;