import { error, redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

type Court = {
	id: string;
	venue_id: string;
	name: string;
	sport: string;
	price_per_slot: number;
	slot_duration_minutes: number;
	is_active: boolean;
};

export const load = (async ({ cookies, fetch, parent }) => {
	const { venue } = await parent();
	const session = cookies.get('session');

	if (!session) {
		redirect(307, '/login');
	}

	let response: Response;

	try {
		response = await fetch(
			`http://localhost:8080/api/venues/${venue.id}/courts`,
			{
				headers: {
					Cookie: `session=${session}`
				}
			}
		);
	} catch {
		error(503, 'Server sedang tidak tersedia');
	}

	if (response.status === 401) {
		cookies.delete('session', { path: '/' });
		redirect(307, '/login');
	}

	if (response.status === 403) {
		error(403, 'Tidak punya akses ke venue ini');
	}

	if (!response.ok) {
		error(response.status, 'Gagal mengambil lapangan');
	}

	const courts = (await response.json()) as Court[];

	return {
		venue,
		courts
	};
}) satisfies PageServerLoad;
