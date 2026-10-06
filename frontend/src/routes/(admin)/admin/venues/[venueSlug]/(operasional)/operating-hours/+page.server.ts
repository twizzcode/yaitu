import { error, fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

type OperatingHour = {
	day_of_week: number;
	opens_at: string | null;
	closes_at: string | null;
	is_closed: boolean;
};

const defaultHours = (): OperatingHour[] =>
	Array.from({ length: 7 }, (_, index) => ({
		day_of_week: index + 1,
		opens_at: '08:00',
		closes_at: '22:00',
		is_closed: false
	}));

export const load = (async ({ cookies, fetch, parent }) => {
	const { venue } = await parent();

	const session = cookies.get('session');

	if (!session) {
		redirect(307, '/login');
	}

	let response: Response;

	try {
		response = await fetch(
			`http://localhost:8080/api/venues/${venue.id}/operating-hours`,
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
		error(response.status, 'Gagal mengambil jam operasional');
	}

	const result = (await response.json()) as {
		hours: OperatingHour[];
	};

	return {
		venue,
		hours: result.hours.length === 0 ? defaultHours() : result.hours
	};
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

		const venue = (await venueResponse.json()) as { id: string; slug: string };

		const data = await request.formData();

		const hours: OperatingHour[] = [];

		for (let day = 1; day <= 7; day++) {
			const isClosed = data.has(`closed_${day}`);
			const opensAt = String(data.get(`opens_at_${day}`) ?? '').trim();
			const closesAt = String(data.get(`closes_at_${day}`) ?? '').trim();

			if (!isClosed && (!opensAt || !closesAt)) {
				return fail(400, {
					hours,
					message: `Jam buka dan tutup hari ke-${day} wajib diisi`
				});
			}

			hours.push({
				day_of_week: day,
				opens_at: isClosed ? null : opensAt,
				closes_at: isClosed ? null : closesAt,
				is_closed: isClosed
			});
		}

		let response: Response;

		try {
			response = await fetch(
				`http://localhost:8080/api/venues/${venue.id}/operating-hours`,
				{
					method: 'PUT',
					headers: {
						'Content-Type': 'application/json',
						Cookie: `session=${session}`
					},
					body: JSON.stringify({ hours })
				}
			);
		} catch {
			return fail(503, {
				hours,
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
				hours,
				message: message.trim() || 'Gagal menyimpan jam operasional'
			});
		}

		redirect(303, `/admin/venues/${venue.slug}/operating-hours`);
	}
} satisfies Actions;
