import { error, fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { API_URL } from '$app/env/private';

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

type AvailabilitySlot = {
	starts_at: string;
	ends_at: string;
	available: boolean;
};

type Availability = {
	court_id: string;
	venue_id: string;
	date: string;
	timezone: string;
	slot_duration_minutes: number;
	price_per_slot: number;
	slots: AvailabilitySlot[];
};

type Booking = {
	id: string;
	court_id: string;
	customer_id: string;
	starts_at: string;
	ends_at: string;
	total_amount: number;
	status: string;
	payment_expires_at: string | null;
	created_at: string;
};

function currentDateInTimezone(timezone: string) {
	const parts = new Intl.DateTimeFormat('en-CA', {
		timeZone: timezone,
		year: 'numeric',
		month: '2-digit',
		day: '2-digit'
	}).formatToParts();

	const values = Object.fromEntries(
		parts.map((part) => [part.type, part.value])
	);

	return `${values.year}-${values.month}-${values.day}`;
}

export const load = (async ({ fetch, params, url }) => {
	let venueResponse: Response;

	try {
		venueResponse = await fetch(
			`${API_URL}/api/public/venues/${encodeURIComponent(params.slug)}`
		);
	} catch {
		error(503, 'Server sedang tidak tersedia');
	}

	if (venueResponse.status === 404) {
		error(404, 'Venue tidak ditemukan');
	}

	if (!venueResponse.ok) {
		error(venueResponse.status, 'Gagal mengambil data venue');
	}

	const venue = (await venueResponse.json()) as PublicVenue;

	const court = venue.courts.find(
		(item) => item.id === params.courtId
	);

	if (!court) {
		error(404, 'Lapangan tidak ditemukan');
	}

	const requestedDate = url.searchParams.get('date');
	const date = requestedDate || currentDateInTimezone(venue.timezone);

	if (!/^\d{4}-\d{2}-\d{2}$/.test(date)) {
		error(400, 'Format tanggal tidak valid');
	}

	let availabilityResponse: Response;

	try {
		availabilityResponse = await fetch(
			`${API_URL}/api/courts/${encodeURIComponent(court.id)}/availability?date=${encodeURIComponent(date)}`
		);
	} catch {
		error(503, 'Server sedang tidak tersedia');
	}

	if (availabilityResponse.status === 400) {
		const message = await availabilityResponse.text();
		error(400, message.trim() || 'Tanggal tidak valid');
	}

	if (availabilityResponse.status === 404) {
		error(404, 'Lapangan tidak ditemukan');
	}

	if (!availabilityResponse.ok) {
		error(
			availabilityResponse.status,
			'Gagal mengambil jadwal lapangan'
		);
	}

	const availability =
		(await availabilityResponse.json()) as Availability;

	return {
		venue,
		court,
		date,
		availability
	};
}) satisfies PageServerLoad;

export const actions = {
	default: async ({ request, cookies, fetch, params, url }) => {
		const data = await request.formData();

		const startsAt = String(data.get('starts_at') ?? '').trim();

		if (!startsAt) {
			return fail(400, {
				message: 'Pilih slot terlebih dahulu'
			});
		}

		const parsedStartsAt = new Date(startsAt);

		if (Number.isNaN(parsedStartsAt.getTime())) {
			return fail(400, {
				message: 'Waktu booking tidak valid'
			});
		}

		const session = cookies.get('session');

		if (!session) {
			const next =
				url.pathname +
				url.search;

			redirect(
				303,
				`/auth/login?next=${encodeURIComponent(next)}`
			);
		}

		let response: Response;

		try {
			response = await fetch(
				`${API_URL}/api/bookings`,
				{
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
						Cookie: `session=${session}`
					},
					body: JSON.stringify({
						court_id: params.courtId,
						starts_at: startsAt
					})
				}
			);
		} catch {
			return fail(503, {
				message: 'Server sedang tidak tersedia'
			});
		}

		if (response.status === 401) {
			cookies.delete('session', {
				path: '/'
			});

			const next =
				url.pathname +
				url.search;

			redirect(
				303,
				`/auth/login?next=${encodeURIComponent(next)}`
			);
		}

		if (!response.ok) {
			const message = await response.text();

			return fail(response.status, {
				message: message.trim() || 'Gagal membuat booking'
			});
		}

		const booking = (await response.json()) as Booking;

		redirect(
			303,
			`/bookings/${booking.id}`
		);
	}
} satisfies Actions;