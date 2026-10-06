import { fail, redirect } from '@sveltejs/kit';
import { PUBLIC_ROOT_DOMAIN } from '$app/env/public';
import type { Actions, PageServerLoad } from './$types';
import { API_URL } from '$app/env/private';

type CreatedVenue = {
	id: string;
	slug: string;
};

type Venue = {
	id: string;
	name: string;
	slug: string;
};

export const load = (async ({ cookies, fetch, parent }) => {
	const { user } = await parent();

	const session = cookies.get('session');

	if (!session) {
		redirect(307, '/login');
	}

	let venuesResponse: Response;

	try {
		venuesResponse = await fetch(`${API_URL}/api/venues`, {
			headers: {
				Cookie: `session=${session}`
			}
		});
	} catch {
		venuesResponse = new Response(null, { status: 503 });
	}

	if (venuesResponse.status === 401) {
		cookies.delete('session', { path: '/' });
		redirect(307, '/login');
	}

	const venues = venuesResponse.ok
		? ((await venuesResponse.json()) as Venue[])
		: [];

	return {
		user,
		hasExistingVenue: venues.length > 0,
		rootDomain: PUBLIC_ROOT_DOMAIN
	};
}) satisfies PageServerLoad;

/**
 * Action pendaftaran venue.
 *
 * Data inti venue (nama, slug, alamat, WhatsApp) dibuat lewat POST /api/venues.
 * Data pribadi (NIK, KTP, wilayah) belum dipersist karena endpoint backend-nya
 * belum tersedia. Saat endpoint tersedia, lengkapi langkah bertanda TODO tanpa
 * mengubah kontrak form di sisi klien.
 *
 * Rencana alur backend lanjutan:
 *   1. Upload KTP           -> POST /api/uploads (multipart, type=ktp) -> { key }
 *   2. Simpan data pribadi  -> PATCH /api/venues/{id}/profile (nik, wa, ktp_key, alamat)
 */
export const actions = {
	default: async ({ request, cookies, fetch }) => {
		const session = cookies.get('session');

		if (!session) {
			redirect(303, '/login');
		}

		const data = await request.formData();

		const namaVenue = String(data.get('namaVenue') ?? '').trim();
		const slug = String(data.get('slug') ?? '').trim().toLowerCase();
		const namaLengkap = String(data.get('namaLengkap') ?? '').trim();
		const nik = String(data.get('nik') ?? '').trim();
		const nomorWhatsApp = String(data.get('nomorWhatsApp') ?? '').trim();
		const alamat = String(data.get('alamat') ?? '').trim();
		const kota = String(data.get('kota') ?? '').trim();
		const kecamatan = String(data.get('kecamatan') ?? '').trim();
		const kelurahan = String(data.get('kelurahan') ?? '').trim();
		const provinsi = String(data.get('provinsi') ?? '').trim();
		const kodePos = String(data.get('kodePos') ?? '').trim();
		const ktpKey = String(data.get('ktp_key') ?? '').trim();

		const values = {
			namaVenue,
			slug,
			namaLengkap,
			nik,
			nomorWhatsApp,
			alamat,
			kota,
			kecamatan,
			kelurahan,
			provinsi,
			kodePos
		};

		if (
			!namaVenue ||
			!slug ||
			!namaLengkap ||
			!nik ||
			!nomorWhatsApp ||
			!alamat ||
			!kota ||
			!provinsi ||
			!kodePos
		) {
			return fail(400, { ...values, message: 'Lengkapi semua data yang wajib diisi' });
		}

		if (!ktpKey) {
			return fail(400, {
				...values,
				message: 'Foto KTP wajib diunggah'
			});
		}

		let response: Response;

		try {
			response = await fetch(`${API_URL}/api/venues`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Cookie: `session=${session}`
				},
				body: JSON.stringify({
					name: namaVenue,
					slug,
					address: alamat,
					timezone: 'Asia/Jakarta',
					whatsapp: nomorWhatsApp,
					ktp_key: ktpKey,
					owner_name: namaLengkap,
					owner_nik: nik,
					province: provinsi,
					city: kota,
					district: kecamatan,
					village: kelurahan,
					postal_code: kodePos
				})
			});
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
				message: message.trim() || 'Gagal membuat venue'
			});
		}

		const venue = (await response.json()) as CreatedVenue;

		redirect(303, `/admin/venues/${venue.slug}`);
	}
} satisfies Actions;
