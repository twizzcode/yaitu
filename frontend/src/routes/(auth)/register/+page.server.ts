import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions = {
	default: async ({ request, fetch }) => {
		const data = await request.formData();

		const name = String(data.get('name') ?? '').trim();
		const email = String(data.get('email') ?? '')
			.trim()
			.toLowerCase();
		const password = String(data.get('password') ?? '');

		if (!name || !email || !password) {
			return fail(400, {
				name,
				email,
				message: 'Semua field wajib diisi'
			});
		}

		if (password.length < 8) {
			return fail(400, {
				name,
				email,
				message: 'Password minimal 8 karakter'
			});
		}

		let response: Response;

		try {
			response = await fetch('http://localhost:8080/api/auth/register', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					name,
					email,
					password
				})
			});
		} catch {
			return fail(503, {
				name,
				email,
				message: 'Server sedang tidak tersedia'
			});
		}

		if (!response.ok) {
			const message = await response.text();

			return fail(response.status, {
				name,
				email,
				message: message.trim() || 'Registrasi gagal'
			});
		}

		redirect(303, '/login');
	}
} satisfies Actions;