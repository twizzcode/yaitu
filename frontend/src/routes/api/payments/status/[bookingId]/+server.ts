import { json } from '@sveltejs/kit';
import { API_URL } from '$app/env/private';
import type { RequestHandler } from './$types';

/**
 * Proxy status pembayaran untuk polling klien.
 *
 * Mengembalikan JSON payment terbaru untuk sebuah booking, atau
 * `{ payment: null }` bila belum ada.
 */
export const GET: RequestHandler = async ({ params, cookies, fetch }) => {
	const session = cookies.get('session');

	if (!session) {
		return json({ payment: null }, { status: 401 });
	}

	const base = API_URL ?? 'http://localhost:8080';

	let response: Response;

	try {
		response = await fetch(
			`${base}/api/bookings/${encodeURIComponent(params.bookingId)}/payment`,
			{
				headers: { Cookie: `session=${session}` }
			}
		);
	} catch {
		return json({ payment: null }, { status: 503 });
	}

	if (response.status === 404) {
		return json({ payment: null });
	}

	if (!response.ok) {
		return json({ payment: null }, { status: response.status });
	}

	const payment = await response.json();

	return json({ payment });
};
