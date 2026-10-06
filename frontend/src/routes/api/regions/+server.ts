import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

/**
 * Proxy data wilayah Indonesia (provinsi → kota → kecamatan → kelurahan).
 *
 * Diproksikan lewat server sendiri supaya:
 * - tidak ada masalah CORS di browser,
 * - penyedia data pihak ketiga tidak melihat intent user secara langsung,
 * - bentuk respons bisa dinormalisasi untuk frontend.
 *
 * Sumber: wilayah.id (https://wilayah.id/api)
 */

const WILAYAH_BASE = 'https://wilayah.id/api';

type WilayahItem = {
	code: string;
	name: string;
};

type WilayahResponse = {
	data?: WilayahItem[];
};

/**
 * Terjemahkan path internal menjadi URL wilayah.id.
 *
 * /provinces                     -> /provinces.json
 * /regencies/{provinceCode}      -> /regencies/{provinceCode}.json
 * /districts/{regencyCode}       -> /districts/{regencyCode}.json
 * /villages/{districtCode}       -> /villages/{districtCode}.json
 */
function resolveUpstreamPath(path: string): string | null {
	const clean = path.trim().replace(/^\/+|\/+$/g, '');

	if (clean === 'provinces') {
		return '/provinces.json';
	}

	const match = clean.match(
		/^(regencies|districts|villages)\/([0-9.]+)$/
	);

	if (!match) {
		return null;
	}

	const [, resource, code] = match;

	return `/${resource}/${code}.json`;
}

export const GET: RequestHandler = async ({ url, fetch }) => {
	const path = url.searchParams.get('path') ?? '';

	const upstreamPath = resolveUpstreamPath(path);

	if (!upstreamPath) {
		return json(
			{ data: [], error: 'Path wilayah tidak valid' },
			{ status: 400 }
		);
	}

	let response: Response;

	try {
		response = await fetch(`${WILAYAH_BASE}${upstreamPath}`);
	} catch {
		return json(
			{ data: [], error: 'Gagal memuat data wilayah' },
			{ status: 502 }
		);
	}

	if (!response.ok) {
		return json(
			{ data: [], error: 'Gagal memuat data wilayah' },
			{ status: response.status === 404 ? 404 : 502 }
		);
	}

	const payload = (await response.json()) as WilayahResponse;

	return json(
		{ data: payload.data ?? [] },
		{
			headers: {
				// Data wilayah jarang berubah — cache di browser & CDN.
				'Cache-Control': 'public, max-age=86400, s-maxage=604800'
			}
		}
	);
};
