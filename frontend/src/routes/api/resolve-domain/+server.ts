import { json } from '@sveltejs/kit';
import { API_URL } from '$app/env/private';
import type { RequestHandler } from './$types';

/**
 * Proxy resolusi domain ke backend.
 *
 * `hooks.ts` (reroute) berjalan di browser dan server, jadi tidak boleh
 * mengakses private env langsung. Route ini berjalan server-only dan
 * meneruskan permintaan ke backend memakai API_URL.
 */
export const GET: RequestHandler = async ({ url, fetch }) => {
	const hostname = (url.searchParams.get('hostname') ?? '').trim().toLowerCase();

	if (!hostname) {
		return json({ venue_slug: null });
	}

	const base = API_URL ?? 'http://localhost:8080';

	let response: Response;

	try {
		response = await fetch(
			`${base}/api/public/domains/resolve?hostname=${encodeURIComponent(hostname)}`
		);
	} catch {
		return json({ venue_slug: null });
	}

	if (!response.ok) {
		return json({ venue_slug: null });
	}

	const data = (await response.json()) as { venue_slug?: string };

	return json({ venue_slug: data.venue_slug ?? null });
};
