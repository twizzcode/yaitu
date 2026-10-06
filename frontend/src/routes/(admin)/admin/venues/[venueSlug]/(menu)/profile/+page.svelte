<script lang="ts">
	import { onDestroy, tick } from 'svelte';
	import AlertTriangleIcon from '@lucide/svelte/icons/alert-triangle';
	import Trash2Icon from '@lucide/svelte/icons/trash-2';
	import UploadIcon from '@lucide/svelte/icons/upload';
	import XIcon from '@lucide/svelte/icons/x';
	import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';
	import { Button } from '#lib/components/ui/button/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '#lib/components/ui/card/index.js';
	import * as Field from '#lib/components/ui/field/index.js';
	import { Input } from '#lib/components/ui/input/index.js';
	import AdminHeroSection from '#lib/components/admin/admin-hero-section.svelte';
	import AddressFields from '#lib/components/venue/address-fields.svelte';
	import type { AddressSelection } from '#lib/components/venue/types.js';
	import { toast } from 'svelte-sonner';
	import { convertToWebp, uploadFile } from '#lib/upload.js';
	import { PUBLIC_ROOT_DOMAIN } from '$app/env/public';

	let { data, form } = $props();

	const confirmationText = 'HAPUS USAHA';
	const MAX_BYTES = 5 * 1024 * 1024;
	const ALLOWED = ['image/jpeg', 'image/png', 'image/webp', 'image/heic', 'image/heif'];

	type ProfileForm = {
		namaProfile: string;
		username: string;
		deskripsi: string;
		namaLengkap: string;
		nik: string;
		nomorWhatsApp: string;
	} & AddressSelection;

	function buildInitialForm(): ProfileForm {
		return {
			namaProfile: data.venue.name,
			username: data.venue.slug,
			deskripsi: data.venue.description,
			namaLengkap: data.venue.owner_name,
			nik: data.venue.owner_nik,
			nomorWhatsApp: data.venue.whatsapp,
			alamat: data.venue.address,
			kota: data.venue.city,
			kecamatan: data.venue.district,
			kelurahan: data.venue.village,
			provinsi: data.venue.province,
			kodePos: data.venue.postal_code
		};
	}

	const initialForm = buildInitialForm();

	let form_ = $state<ProfileForm>({ ...initialForm });
	let errors = $state<Record<string, string>>({});
	let logoFile = $state<File | null>(null);
	let ktpFile = $state<File | null>(null);
	let logoPreview = $state('');
	let ktpPreview = $state('');
	let logoObjectKey = $state('');
	let ktpObjectKey = $state('');
	let logoProcessing = $state(false);
	let ktpProcessing = $state(false);
	let submitting = $state(false);
	let status = $state('');
	let deleteDialog: HTMLDialogElement;
	let confirmation = $state('');

	// Objek tersimpan dari database.
	const savedLogoURL = $derived(data.venue.logo_url);
	const savedKtpURL = $derived(data.venue.ktp_url);

	let isDirty = $derived(
		logoFile !== null ||
			ktpFile !== null ||
			(Object.keys(initialForm) as (keyof ProfileForm)[]).some(
				(key) => form_[key] !== initialForm[key]
			)
	);

	// Sinkronkan ulang setelah submit sukses.
	$effect(() => {
		if (form?.success && form.venue) {
			const v = form.venue;
			form_ = {
				namaProfile: v.name,
				username: v.slug,
				deskripsi: v.description,
				namaLengkap: v.owner_name,
				nik: v.owner_nik,
				nomorWhatsApp: v.whatsapp,
				alamat: v.address,
				kota: v.city,
				kecamatan: v.district,
				kelurahan: v.village,
				provinsi: v.province,
				kodePos: v.postal_code
			};
			clearLogoPreview();
			clearKtpPreview();
			status = 'Profil berhasil disimpan.';
			submitting = false;
			toast.success('Profil berhasil disimpan');
		}

		if (form?.message && !form?.success) {
			status = form.message;
			submitting = false;
			toast.error(form.message);
		}
	});

	onDestroy(() => {
		if (logoPreview) URL.revokeObjectURL(logoPreview);
		if (ktpPreview) URL.revokeObjectURL(ktpPreview);
	});

	function updateField(field: keyof ProfileForm, value: string) {
		form_[field] = value;
		delete errors[field];
		status = '';
	}

	async function processImage(file: File, type: 'logo' | 'ktp') {
		if (!ALLOWED.includes(file.type)) {
			const message = 'Format harus JPG, PNG, WEBP, atau HEIC';
			errors = { ...errors, [type]: message };
			toast.error(message);
			return;
		}
		if (file.size > MAX_BYTES) {
			const message = 'Ukuran file maksimal 5 MB';
			errors = { ...errors, [type]: message };
			toast.error(message);
			return;
		}

		const next = { ...errors };
		delete next[type];
		errors = next;

		if (type === 'logo') {
			logoProcessing = true;
		} else {
			ktpProcessing = true;
		}

		try {
			const webp = await convertToWebp(file);

			if (type === 'logo') {
				if (logoPreview) URL.revokeObjectURL(logoPreview);
				logoFile = webp;
				logoPreview = URL.createObjectURL(webp);
				logoObjectKey = '';
			} else {
				if (ktpPreview) URL.revokeObjectURL(ktpPreview);
				ktpFile = webp;
				ktpPreview = URL.createObjectURL(webp);
				ktpObjectKey = '';
			}
			status = '';
		} catch (err) {
			const message = err instanceof Error ? err.message : 'Gagal memproses gambar';
			errors = { ...errors, [type]: message };
			toast.error(message);
		} finally {
			if (type === 'logo') {
				logoProcessing = false;
			} else {
				ktpProcessing = false;
			}
		}
	}

	function handleImageChange(event: Event, type: 'logo' | 'ktp') {
		const input = event.currentTarget as HTMLInputElement;
		const file = input.files?.[0];
		if (!file) return;

		processImage(file, type).then(() => {
			input.value = '';
		});
	}

	function clearLogoPreview() {
		if (logoPreview) URL.revokeObjectURL(logoPreview);
		logoFile = null;
		logoPreview = '';
		logoObjectKey = '';
	}

	function clearKtpPreview() {
		if (ktpPreview) URL.revokeObjectURL(ktpPreview);
		ktpFile = null;
		ktpPreview = '';
		ktpObjectKey = '';
	}

	function validate() {
		const next: Record<string, string> = {};
		if (!form_.namaProfile.trim()) next.namaProfile = 'Nama lapangan wajib diisi';
		if (!form_.alamat.trim()) next.alamat = 'Alamat wajib diisi';
		if (
			form_.nomorWhatsApp &&
			!/^\+?[0-9]{9,15}$/.test(form_.nomorWhatsApp.replace(/[\s-]/g, ''))
		) {
			next.nomorWhatsApp = 'Nomor WhatsApp tidak valid';
		}
		if (form_.nik && !/^\d{16}$/.test(form_.nik)) {
			next.nik = 'NIK harus terdiri dari 16 digit';
		}
		if (form_.kodePos && !/^\d{5}$/.test(form_.kodePos)) {
			next.kodePos = 'Kode pos harus terdiri dari 5 digit';
		}
		errors = next;
		return Object.keys(next).length === 0;
	}

	async function uploadPending(): Promise<boolean> {
		// Upload logo/KTP baru bila ada. Mengembalikan false bila gagal.
		if (logoFile && !logoObjectKey) {
			const toastId = toast.loading('Mengunggah logo...');
			try {
				const uploaded = await uploadFile(logoFile, 'logo');
				logoObjectKey = uploaded.objectKey;
				toast.success('Logo berhasil diunggah', { id: toastId });
			} catch (err) {
				const message = err instanceof Error ? err.message : 'Gagal mengunggah logo';
				toast.error(message, { id: toastId });
				return false;
			}
		}

		if (ktpFile && !ktpObjectKey) {
			const toastId = toast.loading('Mengunggah foto KTP...');
			try {
				const uploaded = await uploadFile(ktpFile, 'ktp');
				ktpObjectKey = uploaded.objectKey;
				toast.success('Foto KTP berhasil diunggah', { id: toastId });
			} catch (err) {
				const message = err instanceof Error ? err.message : 'Gagal mengunggah foto KTP';
				toast.error(message, { id: toastId });
				return false;
			}
		}

		return true;
	}

	async function submit(event: SubmitEvent) {
		if (!validate()) {
			event.preventDefault();
			status = 'Periksa kembali data yang belum valid.';
			toast.error('Periksa kembali data yang belum valid');
			return;
		}

		if (logoFile || ktpFile) {
			event.preventDefault();
			submitting = true;

			const ok = await uploadPending();

			if (!ok) {
				submitting = false;
				return;
			}

			await tick();

			// Kirim form ke server action dengan object key terbaru.
			formElement?.requestSubmit();
		} else {
			submitting = true;
		}
	}

	let formElement = $state<HTMLFormElement | null>(null);

	function closeDeleteDialog() {
		confirmation = '';
		deleteDialog.close();
	}
</script>

<svelte:head><title>Profil | LapanganKu.id</title></svelte:head>

<main class="mx-auto max-w-7xl space-y-5">
	<AdminHeroSection
	    title="Kelola Profil Lapangan Anda"
        description="Perbarui informasi lapangan, data pemilik, logo, alamat, dan dokumen verifikasi dalam satu halaman yang rapi."
    />

	<form bind:this={formElement} id="profil-form" method="POST" class="space-y-5" onsubmit={submit}>
		<input type="hidden" name="logo_key" value={logoObjectKey} />
		<input type="hidden" name="ktp_key" value={ktpObjectKey} />

		<Card class="w-full overflow-hidden">
			<CardContent>
				<div class="flex flex-col items-start justify-between gap-5 sm:flex-row sm:items-center">
					<div>
						<h2 class="font-semibold">Simpan Perubahan</h2>
						<p class="text-sm text-muted-foreground">
							{isDirty ? 'Anda memiliki perubahan yang belum disimpan.' : 'Belum ada perubahan untuk disimpan.'}
						</p>
						{#if status}<p class="mt-1 text-xs text-muted-foreground" role="status">{status}</p>{/if}
					</div>
					<Button type="submit" size="lg" disabled={!isDirty || submitting || logoProcessing || ktpProcessing} class="w-full sm:w-auto">
						{#if submitting || logoProcessing || ktpProcessing}
							<LoaderCircleIcon class="animate-spin" />
							Memproses...
						{:else}
							Simpan Perubahan
						{/if}
					</Button>
				</div>
			</CardContent>
		</Card>

		<div class="grid gap-4 lg:grid-cols-[minmax(0,320px)_1fr]">
			<Card class="w-full overflow-hidden">
				<CardContent class="space-y-6">
					<div>
						<h2 class="font-semibold">Logo Lapangan</h2>
						<p class="text-sm text-muted-foreground">Tampil di halaman publik lapangan Anda.</p>
					</div>
					<Field.Field>
						<input id="logo" type="file" accept="image/jpeg,image/png,image/webp,image/heic,image/heif" class="hidden" onchange={(event) => handleImageChange(event, 'logo')} />
						{#if logoPreview}
							<div class="relative aspect-square w-full overflow-hidden rounded-xl border border-input">
								<img src={logoPreview} alt="Preview logo" class="size-full object-cover" />
								<Button type="button" variant="secondary" size="icon-sm" onclick={clearLogoPreview} class="absolute top-2 right-2" aria-label="Hapus logo">
									<XIcon />
								</Button>
							</div>
							<label for="logo" class="inline-flex cursor-pointer items-center justify-center rounded-lg border bg-background px-3 py-1.5 text-sm font-medium hover:bg-muted">Ganti Logo</label>
						{:else if savedLogoURL}
							<div class="relative aspect-square w-full overflow-hidden rounded-xl border border-input">
								<img src={savedLogoURL} alt="Logo tersimpan" class="size-full object-cover" />
							</div>
							<label for="logo" class="inline-flex cursor-pointer items-center justify-center rounded-lg border bg-background px-3 py-1.5 text-sm font-medium hover:bg-muted">Ganti Logo</label>
						{:else}
							<label for="logo" class="flex aspect-square w-full cursor-pointer flex-col items-center justify-center gap-2 rounded-xl border border-dashed border-input text-center transition-colors hover:bg-muted/50">
								<UploadIcon class="size-6 text-muted-foreground" />
								<span class="px-4 text-sm text-muted-foreground">Klik untuk pilih logo</span>
								<span class="text-xs text-muted-foreground">JPG/PNG, maks 5 MB</span>
							</label>
						{/if}
						{#if errors.logo}<Field.Error>{errors.logo}</Field.Error>{/if}
					</Field.Field>
				</CardContent>
			</Card>

			<Card class="w-full overflow-hidden">
				<CardContent class="space-y-6">
					<div>
						<h2 class="font-semibold">Detail Lapangan</h2>
						<p class="text-sm text-muted-foreground">Nama, tautan, dan deskripsi lapangan.</p>
					</div>
					<Field.Group>
						<Field.Field data-invalid={!!errors.namaProfile}>
							<Field.Label for="namaProfile">Nama Lapangan *</Field.Label>
							<Input id="namaProfile" name="name" value={form_.namaProfile} oninput={(event) => updateField('namaProfile', event.currentTarget.value)} aria-invalid={!!errors.namaProfile} />
							{#if errors.namaProfile}<Field.Error>{errors.namaProfile}</Field.Error>{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="username">Username</Field.Label>
							<Input id="username" value={form_.username} disabled readonly class="bg-muted/50" />
							<Field.Description>Alamat publik: {form_.username}.{PUBLIC_ROOT_DOMAIN}</Field.Description>
						</Field.Field>
						<Field.Field>
							<Field.Label for="deskripsi">Deskripsi</Field.Label>
							<textarea id="deskripsi" name="description" maxlength="500" value={form_.deskripsi} oninput={(event) => updateField('deskripsi', event.currentTarget.value)} class="min-h-28 w-full rounded-md border border-input bg-input/20 px-3 py-2 text-sm outline-none focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/30" placeholder="Ceritakan tentang lapangan Anda"></textarea>
							<Field.Description>Maksimal 500 karakter.</Field.Description>
						</Field.Field>
					</Field.Group>
				</CardContent>
			</Card>
		</div>

		<Card class="w-full overflow-hidden">
			<CardContent class="space-y-6">
				<div>
					<h2 class="font-semibold">Alamat</h2>
					<p class="text-sm text-muted-foreground">Lokasi lapangan agar pelanggan mudah menemukan Anda.</p>
				</div>
				<AddressFields
					value={{
						alamat: form_.alamat,
						kota: form_.kota,
						kecamatan: form_.kecamatan,
						kelurahan: form_.kelurahan,
						provinsi: form_.provinsi,
						kodePos: form_.kodePos
					}}
					onChange={(field, value) => updateField(field, value)}
					errors={errors}
				/>
				<Field.Field data-invalid={!!errors.alamat}>
					<Field.Label for="alamat">Alamat Lengkap *</Field.Label>
					<textarea id="alamat" name="address" value={form_.alamat} oninput={(event) => updateField('alamat', event.currentTarget.value)} aria-invalid={!!errors.alamat} class="min-h-20 w-full rounded-md border border-input bg-input/20 px-3 py-2 text-sm outline-none focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/30"></textarea>
					{#if errors.alamat}<Field.Error>{errors.alamat}</Field.Error>{/if}
				</Field.Field>
			</CardContent>
		</Card>

		<div class="grid gap-4 lg:grid-cols-[minmax(0,320px)_1fr]">
			<Card class="w-full overflow-hidden">
				<CardContent class="space-y-6">
					<div><h2 class="font-semibold">Foto KTP</h2><p class="text-sm text-muted-foreground">Foto KTP pemilik untuk keperluan verifikasi.</p></div>
					<Field.Field>
						<input id="ktp" type="file" accept="image/jpeg,image/png,image/webp,image/heic,image/heif" class="hidden" onchange={(event) => handleImageChange(event, 'ktp')} />
						{#if ktpPreview}
							<img src={ktpPreview} alt="Preview KTP" class="h-44 w-full rounded-xl border border-input object-cover" />
							<div class="flex gap-2">
								<label for="ktp" class="flex flex-1 cursor-pointer items-center justify-center rounded-lg border px-3 py-2 text-sm font-medium hover:bg-muted">Ganti Foto KTP</label>
								<Button type="button" variant="outline" size="icon" onclick={clearKtpPreview} aria-label="Batalkan perubahan foto KTP">
									<XIcon />
								</Button>
							</div>
						{:else if savedKtpURL}
							<img src={savedKtpURL} alt="Foto KTP tersimpan" class="h-44 w-full rounded-xl border border-input object-cover" />
							<label for="ktp" class="flex cursor-pointer items-center justify-center rounded-lg border px-3 py-2 text-sm font-medium hover:bg-muted">Ganti Foto KTP</label>
							<p class="text-xs text-muted-foreground">Foto KTP tersimpan di storage.</p>
						{:else}
							<label for="ktp" class="flex h-44 cursor-pointer flex-col items-center justify-center gap-2 rounded-xl border border-dashed border-input text-center hover:bg-muted/50">
								<UploadIcon class="size-5 text-muted-foreground" /><span class="text-xs text-muted-foreground">Pilih foto KTP, maks 5 MB</span>
							</label>
						{/if}
						{#if errors.ktp}<Field.Error>{errors.ktp}</Field.Error>{/if}
					</Field.Field>
				</CardContent>
			</Card>

			<Card class="w-full overflow-hidden">
				<CardContent class="space-y-6">
					<div><h2 class="font-semibold">Data Diri</h2><p class="text-sm text-muted-foreground">Data pemilik venue untuk keperluan verifikasi.</p></div>
					<Field.Group>
						<div class="grid gap-5 sm:grid-cols-2">
							<Field.Field><Field.Label for="namaLengkap">Nama Lengkap</Field.Label><Input id="namaLengkap" name="owner_name" value={form_.namaLengkap} oninput={(event) => updateField('namaLengkap', event.currentTarget.value)} /></Field.Field>
							<Field.Field data-invalid={!!errors.nik}><Field.Label for="nik">NIK</Field.Label><Input id="nik" name="owner_nik" inputmode="numeric" maxlength={16} value={form_.nik} oninput={(event) => updateField('nik', event.currentTarget.value.replace(/\D/g, '').slice(0, 16))} aria-invalid={!!errors.nik} />{#if errors.nik}<Field.Error>{errors.nik}</Field.Error>{/if}</Field.Field>
						</div>
						<Field.Field data-invalid={!!errors.nomorWhatsApp}><Field.Label for="nomorWhatsApp">Nomor HP / WhatsApp</Field.Label><Input id="nomorWhatsApp" name="whatsapp" type="tel" value={form_.nomorWhatsApp} oninput={(event) => updateField('nomorWhatsApp', event.currentTarget.value)} aria-invalid={!!errors.nomorWhatsApp} />{#if errors.nomorWhatsApp}<Field.Error>{errors.nomorWhatsApp}</Field.Error>{/if}</Field.Field>
					</Field.Group>
				</CardContent>
			</Card>
		</div>
	</form>

	<Card class="border border-destructive/35 bg-destructive/[0.035] ring-destructive/15">
		<CardHeader class="border-b border-destructive/15">
			<CardTitle class="flex items-center gap-2 text-destructive"><AlertTriangleIcon class="size-5" />Zona Berbahaya</CardTitle>
		</CardHeader>
		<CardContent class="flex flex-col gap-5 sm:flex-row sm:items-center sm:justify-between">
			<div class="space-y-1.5">
				<p class="font-medium">Hapus {data.venue.name}</p>
				<p class="max-w-3xl text-sm text-muted-foreground">Lapangan, jadwal, website, domain, galeri, dan data operasional akan dihapus permanen. Tindakan ini tidak dapat dibatalkan.</p>
			</div>
			<Button variant="destructive" onclick={() => deleteDialog.showModal()}><Trash2Icon />Hapus Usaha</Button>
		</CardContent>
	</Card>
</main>

<dialog bind:this={deleteDialog} class="m-auto w-[calc(100%-2rem)] max-w-lg rounded-xl border bg-background p-0 text-foreground shadow-2xl backdrop:bg-black/50">
	<div class="space-y-5 p-6">
		<div class="flex size-10 items-center justify-center rounded-full bg-destructive/10 text-destructive"><AlertTriangleIcon /></div>
		<div><h2 class="text-lg font-semibold">Hapus {data.venue.name} permanen?</h2><p class="mt-2 text-sm text-muted-foreground">Seluruh data operasional usaha akan dihapus. Integrasi penghapusan belum tersedia pada API saat ini.</p></div>
		<Field.Field>
			<Field.Label for="delete-confirmation">Ketik <span class="font-mono text-destructive">{confirmationText}</span> untuk melanjutkan</Field.Label>
			<Input id="delete-confirmation" bind:value={confirmation} autocomplete="off" placeholder={confirmationText} />
		</Field.Field>
		<div class="flex justify-end gap-2">
			<Button variant="outline" onclick={closeDeleteDialog}>Batal</Button>
			<Button variant="destructive" disabled={confirmation !== confirmationText} onclick={() => { status = 'Endpoint penghapusan usaha belum tersedia.'; closeDeleteDialog(); }}>Hapus Usaha</Button>
		</div>
	</div>
</dialog>
