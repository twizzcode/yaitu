import { error, redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

type BookingDetail = {
	id: string;
	court_id: string;
	court_name: string;
	venue_name: string;
	venue_slug: string;
	timezone: string;
	customer_id: string;
	starts_at: string;
	ends_at: string;
	total_amount: number;
	status:
		| 'pending_payment'
		| 'confirmed'
		| 'cancelled'
		| 'expired'
		| 'completed';
	payment_expires_at: string | null;
	created_at: string;
};

export const load = (async ({ cookies, fetch, url }) => {
	const session = cookies.get('session');

	if (!session) {
		const next = url.pathname + url.search;

		redirect(
			307,
			`/auth/login?next=${encodeURIComponent(next)}`
		);
	}

	let response: Response;

	try {
		response = await fetch(
			'http://localhost:8080/api/bookings',
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
		cookies.delete('session', {
			path: '/'
		});

		const next = url.pathname + url.search;

		redirect(
			307,
			`/auth/login?next=${encodeURIComponent(next)}`
		);
	}

	if (!response.ok) {
		error(response.status, 'Gagal mengambil daftar booking');
	}

	const bookings = (await response.json()) as BookingDetail[];

	return {
		bookings
	};
}) satisfies PageServerLoad;