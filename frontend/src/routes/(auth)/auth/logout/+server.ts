import { error, redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const POST = (async ({ cookies, fetch }) => {
	const session = cookies.get('session');

	if (!session) {
		redirect(303, '/');
	}

	let response: Response;

	try {
		response = await fetch(
			'http://localhost:8080/api/auth/logout',
			{
				method: 'POST',
				headers: {
					Cookie: `session=${session}`
				}
			}
		);
	} catch {
		error(503, 'Server sedang tidak tersedia');
	}

	if (!response.ok) {
		const message = await response.text();

		error(
			response.status,
			message.trim() || 'Gagal keluar'
		);
	}

	cookies.delete('session', {
		path: '/'
	});

	redirect(303, '/');
}) satisfies RequestHandler;