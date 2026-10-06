import { error, fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

export const load = (async ({ parent }) => {
	const { venue } = await parent();
	return { venue };
}) satisfies PageServerLoad;

export const actions = {
	default: async ({ request, cookies, fetch, params }) => {
		const session = cookies.get('session');

		if (!session) {
			redirect(303, '/login');
		}

		const venueResponse = await fetch(
			`http://localhost:8080/api/venues/${encodeURIComponent(params.venueSlug)}`,
			{
				headers: {
					Cookie: `session=${session}`
				}
			}
		);

		if (venueResponse.status === 401) {
			cookies.delete('session', { path: '/' });
			redirect(303, '/login');
		}

		if (!venueResponse.ok) {
			error(venueResponse.status, 'Venue tidak ditemukan');
		}

		const venue = (await venueResponse.json()) as { id: string };

		const data = await request.formData();

		const name = String(data.get('name') ?? '').trim();
		const sport = String(data.get('sport') ?? '')
			.trim()
			.toLowerCase();
		const pricePerSlot = Number(data.get('price_per_slot'));
		const slotDurationMinutes = Number(data.get('slot_duration_minutes'));

		const formValues = {
			name,
			sport,
			pricePerSlot,
			slotDurationMinutes
		};

		if (!name || !sport) {
			return fail(400, {
				...formValues,
				message: 'Nama dan jenis olahraga wajib diisi'
			});
		}

		if (!Number.isInteger(pricePerSlot) || pricePerSlot < 0) {
			return fail(400, {
				...formValues,
				message: 'Harga harus berupa angka bulat dan tidak boleh negatif'
			});
		}

		if (
			!Number.isInteger(slotDurationMinutes) ||
			slotDurationMinutes < 15 ||
			slotDurationMinutes > 1440
		) {
			return fail(400, {
				...formValues,
				message: 'Durasi slot harus 15 sampai 1440 menit'
			});
		}

		let response: Response;

		try {
			response = await fetch(
				`http://localhost:8080/api/venues/${venue.id}/courts`,
				{
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
						Cookie: `session=${session}`
					},
					body: JSON.stringify({
						name,
						sport,
						price_per_slot: pricePerSlot,
						slot_duration_minutes: slotDurationMinutes
					})
				}
			);
		} catch {
			return fail(503, {
				...formValues,
				message: 'Server sedang tidak tersedia'
			});
		}

		if (response.status === 401) {
			cookies.delete('session', { path: '/' });
			redirect(303, '/login');
		}

		if (!response.ok) {
			const message = await response.text();

			return fail(response.status, {
				...formValues,
				message: message.trim() || 'Gagal membuat lapangan'
			});
		}

		redirect(303, `/admin/venues/${params.venueSlug}/courts`);
	}
} satisfies Actions;
