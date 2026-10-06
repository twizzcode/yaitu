import { error, redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';
import { API_URL } from '$app/env/private';

type User = {
	id: string;
	name: string;
	email: string;
};

export const load = (async ({ cookies, fetch }) => {
	const session = cookies.get('session');

	if (!session) {
		redirect(307, '/login');
	}

	let response: Response;

	try {
		response = await fetch(`${API_URL}/api/auth/me`, {
			headers: {
				Cookie: `session=${session}`
			}
		});
	} catch {
		error(503, 'Server sedang tidak tersedia');
	}

	if (response.status === 401) {
		cookies.delete('session', {
			path: '/'
		});

		redirect(307, '/login');
	}

	if (!response.ok) {
		error(response.status, 'Gagal mengambil data pengguna');
	}

	const user = (await response.json()) as User;

	return {
		user
	};
}) satisfies LayoutServerLoad;