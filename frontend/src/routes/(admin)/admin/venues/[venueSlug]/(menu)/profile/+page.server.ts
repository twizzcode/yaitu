import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
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

		const data = await request.formData();

		const name = String(data.get('name') ?? '').trim();
		const address = String(data.get('address') ?? '').trim();
		const whatsapp = String(data.get('whatsapp') ?? '').trim();
		const ownerName = String(data.get('owner_name') ?? '').trim();
		const ownerNik = String(data.get('owner_nik') ?? '').trim();
		const province = String(data.get('provinsi') ?? '').trim();
		const city = String(data.get('kota') ?? '').trim();
		const district = String(data.get('kecamatan') ?? '').trim();
		const village = String(data.get('kelurahan') ?? '').trim();
		const postalCode = String(data.get('kodePos') ?? '').trim();
		const description = String(data.get('description') ?? '').trim();
		const ktpKey = String(data.get('ktp_key') ?? '').trim();
		const logoKey = String(data.get('logo_key') ?? '').trim();

		const values = {
			name,
			address,
			whatsapp,
			ownerName,
			ownerNik,
			province,
			city,
			district,
			village,
			postalCode,
			description
		};

		if (!name || !address) {
			return fail(400, {
				...values,
				message: 'Nama dan alamat wajib diisi'
			});
		}

		if (ownerNik && !/^\d{16}$/.test(ownerNik)) {
			return fail(400, {
				...values,
				message: 'NIK harus 16 digit angka'
			});
		}

		if (postalCode && !/^\d{5}$/.test(postalCode)) {
			return fail(400, {
				...values,
				message: 'Kode pos harus 5 digit angka'
			});
		}

		if (description.length > 500) {
			return fail(400, {
				...values,
				message: 'Deskripsi maksimal 500 karakter'
			});
		}

		let response: Response;

		try {
			response = await fetch(
				`${API_URL}/api/venues/${encodeURIComponent(params.venueSlug)}`,
				{
					method: 'PUT',
					headers: {
						'Content-Type': 'application/json',
						Cookie: `session=${session}`
					},
					body: JSON.stringify({
						name,
						address,
						whatsapp,
						ktp_key: ktpKey,
						owner_name: ownerName,
						owner_nik: ownerNik,
						province,
						city,
						district,
						village,
						postal_code: postalCode,
						description,
						logo_key: logoKey
					})
				}
			);
		} catch {
			return fail(503, {
				...values,
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
				...values,
				message: message.trim() || 'Gagal menyimpan profil'
			});
		}

		const venue = (await response.json()) as Venue;

		return {
			success: true,
			message: 'Profil berhasil disimpan',
			venue
		};
	}
} satisfies Actions;
