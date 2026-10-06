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
		return { user: null };
	}

	try {
		const response = await fetch(`${API_URL}/api/auth/me`, {
			headers: {
				Cookie: `session=${session}`
			}
		});

		if (response.status === 401) {
			cookies.delete('session', { path: '/' });
			return { user: null };
		}

		if (!response.ok) {
			return { user: null };
		}

		return { user: (await response.json()) as User };
	} catch {
		return { user: null };
	}
}) satisfies LayoutServerLoad;
