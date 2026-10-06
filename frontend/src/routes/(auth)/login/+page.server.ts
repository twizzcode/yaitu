import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { API_URL } from '$app/env/private';

export const actions = {
	default: async ({ request, cookies, fetch, url }) => {
		const data = await request.formData();

		const email = String(data.get('email') ?? '')
			.trim()
			.toLowerCase();

		const password = String(data.get('password') ?? '');

		if (!email || !password) {
			return fail(400, {
				email,
				message: 'Email dan password wajib diisi'
			});
		}

		let response: Response;

		try {
			response = await fetch(`${API_URL}/api/auth/login`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					email,
					password
				})
			});
		} catch {
			return fail(503, {
				email,
				message: 'Server sedang tidak tersedia'
			});
		}

		if (!response.ok) {
			const message = await response.text();

			return fail(response.status, {
				email,
				message: message.trim() || 'Login gagal'
			});
		}

		const setCookie = response.headers.get('set-cookie');
		const sessionPrefix = 'session=';

		if (!setCookie?.startsWith(sessionPrefix)) {
			return fail(500, {
				email,
				message: 'Session tidak ditemukan'
			});
		}

		const token = setCookie
			.slice(sessionPrefix.length)
			.split(';', 1)[0];

		cookies.set('session', token, {
			path: '/',
			httpOnly: true,
			sameSite: 'lax',
			secure: false,
			maxAge: 60 * 60 * 24 * 7
		});

		const next = url.searchParams.get('next');
		
		const destination =
			next?.startsWith('/') && !next.startsWith('//')
				? next
				: '/admin';
		
		redirect(303, destination);
	}
} satisfies Actions;