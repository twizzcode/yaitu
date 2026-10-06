/**
 * Utilitas upload ke object storage (Cloudflare R2) dan konversi gambar.
 *
 * Alur upload KTP:
 *   1. Konversi gambar ke WebP di browser (`convertToWebp`).
 *   2. Saat submit, minta presigned URL ke `/api/uploads/presign`.
 *   3. PUT file langsung dari browser ke storage.
 */

type PresignResult = {
	upload_url: string;
	object_key: string;
	public_url?: string;
	expires_in: number;
};

export type UploadedObject = {
	objectKey: string;
	publicUrl: string;
};

/**
 * Konversi gambar apa pun yang didukung browser ke WebP memakai canvas.
 *
 * Bila file sudah WebP, dikembalikan apa adanya. Bila browser gagal
 * men-decode (mis. HEIC di browser non-Safari), error dilempar agar bisa
 * ditangani pemanggil.
 */
export async function convertToWebp(file: File, quality = 0.85): Promise<File> {
	if (file.type === 'image/webp') {
		return file;
	}

	let bitmap: ImageBitmap;

	try {
		bitmap = await createImageBitmap(file);
	} catch {
		throw new Error(
			'Gambar tidak bisa diproses di browser ini. Coba format JPG atau PNG.'
		);
	}

	const canvas = document.createElement('canvas');
	canvas.width = bitmap.width;
	canvas.height = bitmap.height;

	const context = canvas.getContext('2d');

	if (!context) {
		bitmap.close();
		throw new Error('Canvas tidak didukung browser ini');
	}

	context.drawImage(bitmap, 0, 0);
	bitmap.close();

	const blob = await new Promise<Blob | null>((resolve) => {
		canvas.toBlob(resolve, 'image/webp', quality);
	});

	if (!blob) {
		throw new Error('Gagal mengonversi gambar ke WebP');
	}

	const baseName = file.name.replace(/\.[^.]+$/, '') || 'gambar';

	return new File([blob], `${baseName}.webp`, {
		type: 'image/webp'
	});
}

export async function uploadFile(
	file: File,
	kind: 'ktp' | 'logo' = 'ktp'
): Promise<UploadedObject> {
	const presignResponse = await fetch('/api/uploads/presign', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			content_type: file.type,
			kind
		})
	});

	if (!presignResponse.ok) {
		const message = await presignResponse.text();
		throw new Error(message.trim() || 'Gagal membuat URL upload');
	}

	const presign = (await presignResponse.json()) as PresignResult;

	let uploadResponse: Response;

	try {
		uploadResponse = await fetch(presign.upload_url, {
			method: 'PUT',
			headers: {
				'Content-Type': file.type
			},
			body: file
		});
	} catch {
		// Biasanya karena CORS bucket belum mengizinkan origin ini.
		throw new Error(
			'Gagal menghubungi storage. Periksa konfigurasi CORS bucket.'
		);
	}

	if (!uploadResponse.ok) {
		throw new Error('Gagal mengunggah file ke storage');
	}

	return {
		objectKey: presign.object_key,
		publicUrl: presign.public_url ?? ''
	};
}
