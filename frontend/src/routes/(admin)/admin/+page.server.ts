import { error, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { API_URL } from '$app/env/private';

type Venue = {
	id: string;
	slug: string;
};

export const load = (async ({ cookies, fetch, parent }) => {
	await parent();

	const session = cookies.get('session');

	if (!session) {
		redirect(307, '/login');
	}

	let venuesResponse: Response;

	try {
		venuesResponse = await fetch(`${API_URL}/api/venues`, {
			headers: {
				Cookie: `session=${session}`
			}
		});
	} catch {
		error(503, 'Server sedang tidak tersedia');
	}

	if (venuesResponse.status === 401) {
		cookies.delete('session', { path: '/' });
		redirect(307, '/login');
	}

	if (!venuesResponse.ok) {
		error(venuesResponse.status, 'Gagal mengambil venue');
	}

	const venues = (await venuesResponse.json()) as Venue[];
	const venue = venues[0];

	if (!venue) {
		redirect(307, '/admin/venues/new');
	}

	redirect(307, `/admin/venues/${venue.slug}`);
}) satisfies PageServerLoad;

export const actions = {
	logout: async ({ cookies, fetch }) => {
		const session = cookies.get('session');

		if (session) {
			try {
				await fetch(`${API_URL}/api/auth/logout`, {
					method: 'POST',
					headers: {
						Cookie: `session=${session}`
					}
				});
			} catch {
				// Cookie lokal tetap dihapus jika backend sedang tidak tersedia.
			}
		}

		cookies.delete('session', {
			path: '/'
		});

		redirect(303, '/login');
	}
} satisfies Actions;
