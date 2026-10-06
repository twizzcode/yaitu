import { error, fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { API_URL } from '$app/env/private';

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

type Payment = {
	id: string;
	booking_id: string;
	order_id: string;
	status: string;
	gross_amount: number;
	qr_url: string;
	payment_type: string;
	expires_at: string | null;
	transaction_id: string;
};

export const load = (async ({ cookies, fetch, params, url }) => {
	const session = cookies.get('session');

	if (!session) {
		const next = url.pathname + url.search;

		redirect(
			307,
			`/auth/login?next=${encodeURIComponent(next)}`
		);
	}

	const headers = { Cookie: `session=${session}` };

	let response: Response;
	let paymentResponse: Response;

	try {
		[response, paymentResponse] = await Promise.all([
			fetch(
				`${API_URL}/api/bookings/${encodeURIComponent(params.bookingId)}`,
				{ headers }
			),
			fetch(
				`${API_URL}/api/bookings/${encodeURIComponent(params.bookingId)}/payment`,
				{ headers }
			)
		]);
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

	if (response.status === 404) {
		error(404, 'Booking tidak ditemukan');
	}

	if (!response.ok) {
		error(response.status, 'Gagal mengambil detail booking');
	}

	const booking = (await response.json()) as BookingDetail;
	const payment = paymentResponse.ok
		? ((await paymentResponse.json()) as Payment)
		: null;

	return {
		booking,
		payment
	};
}) satisfies PageServerLoad;

export const actions = {
	pay: async ({ cookies, fetch, params }) => {
		const session = cookies.get('session');

		if (!session) {
			redirect(303, '/login');
		}

		let response: Response;

		try {
			response = await fetch(`${API_URL}/api/payments`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Cookie: `session=${session}`
				},
				body: JSON.stringify({ booking_id: params.bookingId })
			});
		} catch {
			return fail(503, { message: 'Server sedang tidak tersedia' });
		}

		if (response.status === 401) {
			cookies.delete('session', { path: '/' });
			redirect(303, '/login');
		}

		if (!response.ok) {
			const message = await response.text();
			return fail(response.status, {
				message: message.trim() || 'Gagal membuat pembayaran'
			});
		}

		const payment = (await response.json()) as Payment;

		return {
			success: true,
			payment
		};
	},

	// Dipakai polling klien: mengembalikan status pembayaran terbaru.
	status: async ({ cookies, fetch, params }) => {
		const session = cookies.get('session');

		if (!session) {
			return fail(401, { message: 'belum login' });
		}

		let response: Response;

		try {
			response = await fetch(
				`${API_URL}/api/bookings/${encodeURIComponent(params.bookingId)}/payment`,
				{
					headers: { Cookie: `session=${session}` }
				}
			);
		} catch {
			return fail(503, { message: 'Server sedang tidak tersedia' });
		}

		if (!response.ok) {
			return fail(response.status, { message: 'Gagal mengambil status' });
		}

		const payment = (await response.json()) as Payment;

		return { success: true, payment };
	}
} satisfies Actions;
