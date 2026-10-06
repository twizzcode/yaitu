import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { API_URL } from '$app/env/private';
import { validateUsername } from '#lib/validation.js';

/**
 * Cek ketersediaan alamat publik (slug) venue.
 *
 * Untuk sekarang memakai endpoint publik yang sudah ada:
 *   GET /api/public/venues/{slug}
 *   - 200 -> slug sudah dipakai
 *   - 404 -> slug tersedia
 *
 * Saat backend menyediakan endpoint khusus (mis. GET /api/venues/slug-available),
 * cukup ganti pemanggilan di bawah tanpa mengubah kontrak respons frontend.
 */

export const GET: RequestHandler = async ({ url, fetch }) => {
	const slug = (url.searchParams.get('slug') ?? '')
		.trim()
		.toLowerCase();

	const validationError = validateUsername(slug);

	if (validationError) {
		return json({
			status: 'invalid',
			available: false,
			message: validationError
		});
	}

	let response: Response;

	try {
		response = await fetch(`${API_URL}/api/public/venues/${slug}`);
	} catch {
		return json(
			{
				status: 'error',
				available: false,
				message: 'Gagal memeriksa alamat publik'
			},
			{ status: 503 }
		);
	}

	if (response.status === 404) {
		return json({
			status: 'available',
			available: true,
			message: 'Alamat publik tersedia'
		});
	}

	if (response.ok) {
		return json({
			status: 'taken',
			available: false,
			message: 'Alamat publik sudah digunakan'
		});
	}

	return json(
		{
			status: 'error',
			available: false,
			message: 'Gagal memeriksa alamat publik'
		},
		{ status: 502 }
	);
};
