import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

/**
 * Cari kode pos otomatis berdasarkan wilayah administratif.
 *
 * Sumber: kodepos.vercel.app (pencarian kode pos Indonesia).
 * Karena API sumber hanya menerima satu kata kunci pencarian, kita:
 *   1. query dengan "{kelurahan} {kota}",
 *   2. filter hasil yang nama kelurahannya benar-benar cocok,
 *   3. pilih yang kecamatan/kotanya paling sesuai.
 *
 * Selalu mengembalikan 200 dengan `code` (string kosong bila tidak ketemu)
 * supaya frontend bisa menangani "tidak ditemukan" tanpa error.
 */

const KODEPOS_BASE = 'https://kodepos.vercel.app/search/';

type KodePosItem = {
	code: number | string;
	village: string;
	district: string;
	regency: string;
	province: string;
};

type KodePosResponse = {
	data?: KodePosItem[];
};

const ADMIN_PREFIX =
	/^(kota|kabupaten|kab\.?|kecamatan|kec\.?|kelurahan|desa|administrasi|adm\.?)\s+/i;

/** Normalisasi nama daerah untuk pencocokan yang toleran. */
function normalize(value: string): string {
	let result = value.trim();

	// Buang prefix administratif berulang ("Kota Administrasi ...").
	while (ADMIN_PREFIX.test(result)) {
		result = result.replace(ADMIN_PREFIX, '');
	}

	return result
		.toLowerCase()
		.replace(/[^a-z0-9]/g, '')
		.trim();
}

/** Buang prefix administratif ("Kota", "Kabupaten", ...) untuk kata kunci pencarian. */
function stripPrefix(value: string): string {
	let result = value.trim();

	while (ADMIN_PREFIX.test(result)) {
		result = result.replace(ADMIN_PREFIX, '');
	}

	return result;
}

export const GET: RequestHandler = async ({ url, fetch }) => {
	const kelurahan = (url.searchParams.get('kelurahan') ?? '').trim();
	const kecamatan = (url.searchParams.get('kecamatan') ?? '').trim();
	const kota = (url.searchParams.get('kota') ?? '').trim();

	if (!kelurahan) {
		return json({ code: '' });
	}

	// API sumber tidak mengenal prefix "Kota/Kabupaten" — buang dulu.
	const query = [kelurahan, stripPrefix(kota)].filter(Boolean).join(' ');

	let response: Response;

	try {
		response = await fetch(`${KODEPOS_BASE}?q=${encodeURIComponent(query)}`);
	} catch {
		return json({ code: '' });
	}

	if (!response.ok) {
		return json({ code: '' });
	}

	const payload = (await response.json()) as KodePosResponse;
	const items = payload.data ?? [];

	const targetVillage = normalize(kelurahan);

	// Hanya terima kelurahan yang namanya cocok.
	const candidates = items.filter((item) => normalize(item.village) === targetVillage);

	if (candidates.length === 0) {
		return json({ code: '' });
	}

	const targetDistrict = normalize(kecamatan);
	const targetRegency = normalize(kota);

	// Prioritas: kecamatan cocok, lalu kota cocok, terakhir ambil pertama.
	const best =
		candidates.find(
			(item) =>
				targetDistrict && normalize(item.district) === targetDistrict
		) ??
		candidates.find(
			(item) => targetRegency && normalize(item.regency).includes(targetRegency)
		) ??
		candidates[0];

	return json(
		{ code: String(best.code) },
		{
			headers: {
				'Cache-Control': 'public, max-age=86400, s-maxage=604800'
			}
		}
	);
};
