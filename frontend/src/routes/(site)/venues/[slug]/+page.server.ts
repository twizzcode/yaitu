import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

type PublicCourt = {
	id: string;
	name: string;
	sport: string;
	price_per_slot: number;
	slot_duration_minutes: number;
};

type PublicVenue = {
	id: string;
	name: string;
	slug: string;
	address: string;
	timezone: string;
	whatsapp: string;
	courts: PublicCourt[];
};

export const load = (async ({ fetch, params }) => {
	let response: Response;

	try {
		response = await fetch(
			`http://localhost:8080/api/public/venues/${encodeURIComponent(params.slug)}`
		);
	} catch {
		error(503, 'Server sedang tidak tersedia');
	}

	if (response.status === 404) {
		error(404, 'Venue tidak ditemukan');
	}

	if (!response.ok) {
		error(response.status, 'Gagal mengambil data venue');
	}

	const venue = (await response.json()) as PublicVenue;

	return {
		venue
	};
}) satisfies PageServerLoad;