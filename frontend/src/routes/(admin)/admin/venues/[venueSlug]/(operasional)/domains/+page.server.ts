import { error, fail, redirect } from '@sveltejs/kit';
import { PUBLIC_ROOT_DOMAIN } from '$app/env/public';
import { SERVER_PUBLIC_IP, API_URL } from '$app/env/private';
import type { Actions, PageServerLoad } from './$types';

type Domain = {
	id: string;
	hostname: string;
	type: 'platform' | 'custom';
	status: 'pending' | 'active' | 'failed';
	verification_token: string;
	verified_at: string | null;
	created_at: string;
	is_apex: boolean;
};

export const load = (async ({ cookies, fetch, params, parent }) => {
	await parent();

	const session = cookies.get('session');

	if (!session) {
		redirect(307, '/login');
	}

	let response: Response;

	try {
		response = await fetch(
			`${API_URL}/api/venues/${encodeURIComponent(params.venueSlug)}/domains`,
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

	if (response.status === 404) {
		error(404, 'Venue tidak ditemukan');
	}

	if (!response.ok) {
		error(response.status, 'Gagal mengambil daftar domain');
	}

	const domains = (await response.json()) as Domain[];

	return {
		domains,
		serverIp: SERVER_PUBLIC_IP ?? '',
		rootDomain: PUBLIC_ROOT_DOMAIN
	};
}) satisfies PageServerLoad;

export const actions = {
	add: async ({ request, cookies, fetch, params }) => {
		const session = cookies.get('session');

		if (!session) {
			redirect(303, '/login');
		}

		const data = await request.formData();
		const hostname = String(data.get('hostname') ?? '')
			.trim()
			.toLowerCase();

		if (!hostname) {
			return fail(400, { hostname, message: 'Domain wajib diisi' });
		}

		let response: Response;

		try {
			response = await fetch(
				`${API_URL}/api/venues/${encodeURIComponent(params.venueSlug)}/domains`,
				{
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
						Cookie: `session=${session}`
					},
					body: JSON.stringify({ hostname })
				}
			);
		} catch {
			return fail(503, { hostname, message: 'Server sedang tidak tersedia' });
		}

		if (response.status === 401) {
			cookies.delete('session', { path: '/' });
			redirect(303, '/login');
		}

		if (!response.ok) {
			const message = await response.text();
			return fail(response.status, {
				hostname,
				message: message.trim() || 'Gagal menambah domain'
			});
		}

		return { success: true, message: 'Domain berhasil ditambahkan' };
	},

	verify: async ({ request, cookies, fetch, params }) => {
		const session = cookies.get('session');

		if (!session) {
			redirect(303, '/login');
		}

		const data = await request.formData();
		const domainId = String(data.get('domain_id') ?? '').trim();

		if (!domainId) {
			return fail(400, { message: 'Domain tidak valid' });
		}

		let response: Response;

		try {
			response = await fetch(
				`${API_URL}/api/venues/${encodeURIComponent(params.venueSlug)}/domains/${encodeURIComponent(domainId)}/verify`,
				{
					method: 'POST',
					headers: {
						Cookie: `session=${session}`
					}
				}
			);
		} catch {
			return fail(503, { message: 'Server sedang tidak tersedia' });
		}

		if (response.status === 401) {
			cookies.delete('session', { path: '/' });
			redirect(303, '/login');
		}

		if (!response.ok) {
			return fail(422, {
				message:
					'Verifikasi gagal. Pastikan record DNS TXT sudah benar dan tunggu propagasi DNS.'
			});
		}

		return { success: true, message: 'Domain berhasil diverifikasi' };
	},

	remove: async ({ request, cookies, fetch, params }) => {
		const session = cookies.get('session');

		if (!session) {
			redirect(303, '/login');
		}

		const data = await request.formData();
		const domainId = String(data.get('domain_id') ?? '').trim();

		if (!domainId) {
			return fail(400, { message: 'Domain tidak valid' });
		}

		let response: Response;

		try {
			response = await fetch(
				`${API_URL}/api/venues/${encodeURIComponent(params.venueSlug)}/domains/${encodeURIComponent(domainId)}`,
				{
					method: 'DELETE',
					headers: {
						Cookie: `session=${session}`
					}
				}
			);
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
				message: message.trim() || 'Gagal menghapus domain'
			});
		}

		return { success: true, message: 'Domain berhasil dihapus' };
	}
} satisfies Actions;
