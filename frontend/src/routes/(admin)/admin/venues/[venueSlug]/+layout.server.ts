import { error, redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';
import { API_URL } from '$app/env/private';

type Venue = {
	id: string;
	name: string;
	slug: string;
	address: string;
	timezone: string;
	whatsapp: string;
	is_active: boolean;
	role: 'owner' | 'admin';
	ktp_key: string;
	ktp_url: string;
	owner_name: string;
	owner_nik: string;
	province: string;
	city: string;
	district: string;
	village: string;
	postal_code: string;
	description: string;
	logo_key: string;
	logo_url: string;
};

export const load = (async ({ cookies, fetch, params, parent }) => {
	await parent();

	const session = cookies.get('session');

	if (!session) {
		redirect(307, '/login');
	}

	const headers = {
		Cookie: `session=${session}`
	};

	let venuesResponse: Response;

	try {
		venuesResponse = await fetch(`${API_URL}/api/venues`, {
			headers
		});
	} catch {
		error(503, 'Server sedang tidak tersedia');
	}

	if (venuesResponse.status === 401) {
		cookies.delete('session', { path: '/' });
		redirect(307, '/login');
	}

	if (!venuesResponse.ok) {
		error(venuesResponse.status, 'Gagal mengambil daftar venue');
	}

	const venues = (await venuesResponse.json()) as Venue[];

	const venue = venues.find((item) => item.slug === params.venueSlug);

	if (!venue) {
		error(404, 'Venue tidak ditemukan');
	}

	return {
		venues,
		venue
	};
}) satisfies LayoutServerLoad;
