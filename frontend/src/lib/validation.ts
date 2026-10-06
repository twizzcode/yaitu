/**
 * Validasi bersama untuk form pendaftaran venue.
 * Semua fungsi mengembalikan pesan error (string) bila tidak valid,
 * atau `null` bila nilainya valid.
 */

const SLUG_PATTERN = /^[a-z0-9]+(?:-[a-z0-9]+)*$/;

/** Slug venue/username: huruf kecil, angka, dan tanda hubung. */
export function validateUsername(value: string): string | null {
	const slug = value.trim().toLowerCase();

	if (!slug) {
		return 'Alamat publik wajib diisi';
	}

	if (slug.length < 3 || slug.length > 63) {
		return 'Alamat publik harus 3 sampai 63 karakter';
	}

	if (!SLUG_PATTERN.test(slug)) {
		return 'Hanya boleh huruf kecil, angka, dan tanda hubung';
	}

	return null;
}

/** NIK: tepat 16 digit angka. */
export function validateNik(value: string): string | null {
	const nik = value.trim();

	if (!nik) {
		return 'NIK wajib diisi';
	}

	if (!/^\d{16}$/.test(nik)) {
		return 'NIK harus 16 digit angka';
	}

	return null;
}

/** Nomor HP/WhatsApp Indonesia. */
export function validateWhatsapp(value: string): string | null {
	const raw = value.trim();

	if (!raw) {
		return 'Nomor HP / WhatsApp wajib diisi';
	}

	const digits = raw.replace(/[^\d]/g, '');

	if (digits.startsWith('62')) {
		if (digits.length < 11 || digits.length > 15) {
			return 'Nomor WhatsApp tidak valid';
		}
		return null;
	}

	if (digits.startsWith('0')) {
		if (digits.length < 10 || digits.length > 14) {
			return 'Nomor WhatsApp tidak valid';
		}
		return null;
	}

	return 'Nomor harus diawali 08 atau 62';
}

/** Kode pos: 5 digit angka. */
export function validateKodePos(value: string): string | null {
	const kodePos = value.trim();

	if (!kodePos) {
		return 'Kode pos wajib diisi';
	}

	if (!/^\d{5}$/.test(kodePos)) {
		return 'Kode pos harus 5 digit angka';
	}

	return null;
}

/** Ubah teks menjadi slug yang aman untuk alamat publik. */
export function slugify(value: string): string {
	return value
		.toLowerCase()
		.trim()
		.replace(/[^a-z0-9]+/g, '-')
		.replace(/^-+|-+$/g, '')
		.slice(0, 63);
}
