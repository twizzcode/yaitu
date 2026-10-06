import { error } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

type User = {
	id: string;
	name: string;
	email: string;
};

export const load = (async ({ cookies, fetch }) => {
	const session = cookies.get('session');

	if (!session) {
		return {
			user: null
		};
	}

	let response: Response;

	try {
		response = await fetch(
			'http://localhost:8080/api/auth/me',
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
		cookies.delete('session', {
			path: '/'
		});

		return {
			user: null
		};
	}

	if (!response.ok) {
		error(
			response.status,
			'Gagal mengambil data pengguna'
		);
	}

	const user = (await response.json()) as User;

	return {
		user
	};
}) satisfies LayoutServerLoad;