import { error, json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { API_URL } from '$app/env/private';

type PresignResponse = {
	upload_url: string;
	object_key: string;
	expires_in: number;
};

/**
 * Proxy pembuatan presigned URL upload ke backend.
 *
 * Cookie session bersifat HttpOnly sehingga browser tidak bisa memanggil
 * backend langsung. Endpoint ini meneruskan permintaan beserta cookie, lalu
 * mengembalikan presigned URL yang dipakai browser untuk upload langsung ke
 * object storage (Cloudflare R2).
 */
export const POST: RequestHandler = async ({ request, cookies, fetch }) => {
	const session = cookies.get('session');

	if (!session) {
		error(401, 'Belum login');
	}

	let body: { content_type?: string; kind?: string };

	try {
		body = await request.json();
	} catch {
		error(400, 'Body tidak valid');
	}

	let response: Response;

	try {
		response = await fetch(`${API_URL}/api/uploads/presign`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Cookie: `session=${session}`
			},
			body: JSON.stringify({
				content_type: body.content_type ?? '',
				kind: body.kind ?? 'ktp'
			})
		});
	} catch {
		error(503, 'Server sedang tidak tersedia');
	}

	if (!response.ok) {
		const message = await response.text();
		error(response.status, message.trim() || 'Gagal membuat URL upload');
	}

	const result = (await response.json()) as PresignResponse;

	return json(result);
};
