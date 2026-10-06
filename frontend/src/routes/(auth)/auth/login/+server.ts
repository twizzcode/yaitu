import {
	PUBLIC_ROOT_DOMAIN,
	PUBLIC_ROOT_URL
} from '$app/env/public';
import { error, redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

function normalizeReturnPath(value: string | null) {
	const returnPath = value?.trim() || '/';

	if (
		!returnPath.startsWith('/') ||
		returnPath.startsWith('//') ||
		returnPath.includes('\r') ||
		returnPath.includes('\n')
	) {
		return null;
	}

	return returnPath;
}

export const GET = (({ cookies, url }) => {
	const targetHost = url.hostname.toLowerCase();

	if (
		targetHost === PUBLIC_ROOT_DOMAIN ||
		targetHost === `www.${PUBLIC_ROOT_DOMAIN}`
	) {
		redirect(303, '/login');
	}

	const returnPath = normalizeReturnPath(
		url.searchParams.get('next')
	);

	if (!returnPath) {
		error(400, 'Return path tidak valid');
	}

	const state = crypto.randomUUID().replaceAll('-', '');

	cookies.set('sso_state', state, {
		path: '/auth/callback',
		httpOnly: true,
		sameSite: 'lax',
		secure: url.protocol === 'https:',
		maxAge: 5 * 60
	});

	const authorizeURL = new URL(
		'/sso/authorize',
		PUBLIC_ROOT_URL
	);

	authorizeURL.searchParams.set(
		'target_host',
		targetHost
	);
	authorizeURL.searchParams.set(
		'return_path',
		returnPath
	);
	authorizeURL.searchParams.set('state', state);

	redirect(303, authorizeURL, {
		external: [new URL(PUBLIC_ROOT_URL).origin]
	});
}) satisfies RequestHandler;