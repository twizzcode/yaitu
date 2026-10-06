<script lang="ts">
	import { tick, untrack } from 'svelte';
	import ArrowLeftIcon from '@lucide/svelte/icons/arrow-left';
	import ArrowRightIcon from '@lucide/svelte/icons/arrow-right';
	import CheckIcon from '@lucide/svelte/icons/check';
	import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';
	import UploadIcon from '@lucide/svelte/icons/upload';
	import XIcon from '@lucide/svelte/icons/x';
	import { Button } from '#lib/components/ui/button/index.js';
	import { Input } from '#lib/components/ui/input/index.js';
	import { Textarea } from '#lib/components/ui/textarea/index.js';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '#lib/components/ui/card/index.js';
	import { Field, FieldError, FieldLabel, FieldGroup } from '#lib/components/ui/field/index.js';
	import * as AlertDialog from '#lib/components/ui/alert-dialog/index.js';
	import AdminHeroSection from '#lib/components/admin/admin-hero-section.svelte';
	import SlugField from './slug-field.svelte';
	import AddressFields from './address-fields.svelte';
	import type { AddressSelection, SlugStatus } from './types.js';
	import { cn } from '#lib/utils.js';
	import { toast } from 'svelte-sonner';
	import { convertToWebp, uploadFile } from '#lib/upload.js';
	import {
		slugify,
		validateKodePos,
		validateNik,
		validateUsername,
		validateWhatsapp
	} from '#lib/validation.js';

	type FormValues = {
		namaVenue: string;
		slug: string;
		namaLengkap: string;
		nik: string;
		nomorWhatsApp: string;
	} & AddressSelection;

	let {
		defaultEmail,
		rootDomain,
		hasExistingVenue = false,
		form
	}: {
		defaultEmail: string;
		rootDomain: string;
		hasExistingVenue?: boolean;
		form?: { message?: string } | null;
	} = $props();

	const STEPS = [
		{ id: 1, name: 'Info Venue', description: 'Detail identitas venue kamu' },
		{ id: 2, name: 'Info Pribadi', description: 'Data diri & verifikasi KTP' },
		{ id: 3, name: 'Kontak & Alamat', description: 'Cara pelanggan menemukan kamu' }
	];

	const KTP_MAX_BYTES = 5 * 1024 * 1024;
	const KTP_ALLOWED = ['image/jpeg', 'image/png', 'image/webp', 'image/heic', 'image/heif'];

	// Konfirmasi tambah venue bila user sudah punya venue.
	let confirmed = $state(untrack(() => !hasExistingVenue));

	let currentStep = $state(1);
	let isPending = $state(false);
	let submitError = $state(untrack(() => form?.message ?? ''));

	$effect(() => {
		if (form?.message) {
			submitError = form.message;
			toast.error(form.message);
			// Server action gagal — aktifkan kembali tombol submit.
			isPending = false;
		}
	});

	let formData = $state<FormValues>({
		namaVenue: '',
		slug: '',
		namaLengkap: '',
		nik: '',
		nomorWhatsApp: '',
		alamat: '',
		kota: '',
		kecamatan: '',
		kelurahan: '',
		provinsi: '',
		kodePos: ''
	});

	let ktpFile = $state<File | null>(null);
	let ktpPreview = $state('');
	let ktpProcessing = $state(false);
	let ktpObjectKey = $state('');
	let formElement = $state<HTMLFormElement | null>(null);
	let fieldErrors = $state<Record<string, string>>({});
	let slugStatus = $state<SlugStatus>('idle');

	function updateField<K extends keyof FormValues>(field: K, value: FormValues[K]) {
		formData[field] = value;

		if (fieldErrors[field]) {
			const next = { ...fieldErrors };
			delete next[field];
			fieldErrors = next;
		}
	}

	function validateStep(step: number): Record<string, string> {
		const errors: Record<string, string> = {};

		const required = (field: keyof FormValues, label: string) => {
			if (!String(formData[field] ?? '').trim()) {
				errors[field] = `${label} wajib diisi`;
			}
		};

		if (step === 1) {
			required('namaVenue', 'Nama venue');
			const slugError = validateUsername(formData.slug);
			if (slugError) errors.slug = slugError;
		}

		if (step === 2) {
			required('namaLengkap', 'Nama lengkap');
			const nikError = validateNik(formData.nik);
			if (nikError) errors.nik = nikError;
			const waError = validateWhatsapp(formData.nomorWhatsApp);
			if (waError) errors.nomorWhatsApp = waError;
			if (!ktpFile) {
				errors.ktp = 'Foto KTP wajib diunggah';
			}
		}

		if (step === 3) {
			required('alamat', 'Alamat');
			required('kota', 'Kota');
			required('provinsi', 'Provinsi');
			const kodePosError = validateKodePos(formData.kodePos);
			if (kodePosError) errors.kodePos = kodePosError;
		}

		return errors;
	}

	function validateAll(): Record<string, string> {
		const all: Record<string, string> = {};
		for (let step = 1; step <= STEPS.length; step++) {
			Object.assign(all, validateStep(step));
		}
		return all;
	}

	function firstInvalidStep(errors: Record<string, string>): number {
		return (
			[1, 2, 3].find((step) => Object.keys(validateStep(step)).length > 0) ?? 1
		);
	}

	function nextStep() {
		const errors = validateStep(currentStep);
		fieldErrors = errors;

		if (Object.keys(errors).length > 0) return;
		if (currentStep < STEPS.length) currentStep += 1;
	}

	function prevStep() {
		submitError = '';
		fieldErrors = {};

		if (currentStep > 1) currentStep -= 1;
	}

	function goToStep(step: number) {
		if (step < currentStep) {
			fieldErrors = {};
			currentStep = step;
		}
	}

	async function handleKtpChange(event: Event) {
		const input = event.currentTarget as HTMLInputElement;
		const file = input.files?.[0];

		if (!file) return;

		if (file.size > KTP_MAX_BYTES) {
			const message = 'Ukuran file maksimal 5 MB';
			fieldErrors = { ...fieldErrors, ktp: message };
			toast.error(message);
			input.value = '';
			return;
		}

		if (!KTP_ALLOWED.includes(file.type)) {
			const message = 'Format harus JPG, PNG, WEBP, atau HEIC';
			fieldErrors = {
				...fieldErrors,
				ktp: message
			};
			toast.error(message);
			input.value = '';
			return;
		}

		const next = { ...fieldErrors };
		delete next.ktp;
		fieldErrors = next;

		// Konversi ke WebP di browser. Upload baru dilakukan saat submit
		// agar tidak ada file orphan bila user batal.
		ktpProcessing = true;

		try {
			const webp = await convertToWebp(file);

			if (ktpPreview) URL.revokeObjectURL(ktpPreview);
			ktpFile = webp;
			ktpPreview = URL.createObjectURL(webp);
			// Foto berubah — buang key lama agar diunggah ulang saat submit.
			ktpObjectKey = '';
		} catch (err) {
			const message =
				err instanceof Error ? err.message : 'Gagal memproses foto KTP';
			fieldErrors = { ...fieldErrors, ktp: message };
			toast.error(message);
			input.value = '';
		} finally {
			ktpProcessing = false;
		}
	}

	function clearKtp() {
		if (ktpPreview) URL.revokeObjectURL(ktpPreview);
		ktpFile = null;
		ktpPreview = '';
		ktpObjectKey = '';
		ktpProcessing = false;
	}

	async function handleSubmit(event: SubmitEvent) {
		submitError = '';

		const errors = validateAll();

		if (Object.keys(errors).length > 0) {
			event.preventDefault();
			fieldErrors = errors;
			currentStep = firstInvalidStep(errors);
			toast.error('Periksa kembali data yang belum lengkap');
			return;
		}

		if (!ktpFile) {
			event.preventDefault();
			toast.error('Foto KTP wajib diunggah');
			return;
		}

		// Kalau ktp_key sudah ada, ini submit lanjutan setelah upload selesai:
		// biarkan form terkirim ke server action.
		if (ktpObjectKey) {
			return;
		}

		// Blokir submit bawaan; upload dulu, lalu submit ulang.
		event.preventDefault();
		isPending = true;

		const toastId = toast.loading('Mengunggah foto KTP...');

		try {
			const uploaded = await uploadFile(ktpFile, 'ktp');
			ktpObjectKey = uploaded.objectKey;
			toast.success('Foto KTP berhasil diunggah', { id: toastId });
		} catch (err) {
			const message =
				err instanceof Error ? err.message : 'Gagal mengunggah foto KTP';
			toast.error(message, { id: toastId });
			isPending = false;
			currentStep = 2;
			return;
		}

		await tick();

		// Submit ulang; kali ini ktpObjectKey sudah terisi sehingga lolos
		// pengecekan di atas dan form benar-benar dikirim ke server action.
		formElement?.requestSubmit();
	}

	function autoSlug() {
		// Jangan timpa slug yang sudah diisi user.
		if (formData.slug.trim()) return;

		const generated = slugify(formData.namaVenue);
		if (generated) updateField('slug', generated);
	}

	const progress = $derived(((currentStep - 1) / (STEPS.length - 1)) * 100);

	// Step 1 hanya boleh lanjut bila nama terisi dan slug sudah tersedia.
	const stepOneReady = $derived(
		!!formData.namaVenue.trim() && slugStatus === 'available'
	);
</script>

<AlertDialog.Root open={!confirmed} onOpenChange={() => {}}>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Buat venue baru?</AlertDialog.Title>
			<AlertDialog.Description>
				Akun kamu sudah punya venue. Lanjutkan untuk menambah venue baru dengan data dan
				tagihan terpisah.
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={() => (window.location.href = '/admin')}>
				Batal
			</AlertDialog.Cancel>
			<AlertDialog.Action onclick={() => (confirmed = true)}>Tambah Venue</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

{#if confirmed}
	<main class="mx-auto max-w-3xl space-y-5 px-4 py-10 md:py-14">
		<AdminHeroSection
			title="Daftar Sebagai Penyedia Venue"
			description="Lengkapi data di bawah untuk mulai menerima booking."
		/>

		<Card class="w-full overflow-hidden">
			<CardHeader>
				<CardTitle>{STEPS[currentStep - 1].name}</CardTitle>
				<CardDescription>{STEPS[currentStep - 1].description}</CardDescription>
			</CardHeader>
			<CardContent>
				<!-- Stepper -->
				<div class="relative mb-8">
					<div class="bg-muted absolute top-4 left-0 h-0.5 w-full"></div>
					<div
						class="bg-primary absolute top-4 left-0 h-0.5 transition-all duration-500"
						style="width: {progress}%"
					></div>

					<ol class="relative flex justify-between">
						{#each STEPS as step (step.id)}
							{@const isDone = step.id < currentStep}
							{@const isActive = step.id === currentStep}
							<li class="flex flex-col items-center gap-2">
								<button
									type="button"
									onclick={() => goToStep(step.id)}
									disabled={step.id >= currentStep}
									class={cn(
										'flex size-8 items-center justify-center rounded-full border-2 bg-background text-xs font-semibold transition-colors',
										isDone && 'border-primary bg-primary text-primary-foreground',
										isActive && 'border-primary text-primary',
										!isDone && !isActive && 'border-muted text-muted-foreground',
										isDone && 'cursor-pointer'
									)}
								>
									{#if isDone}
										<CheckIcon class="size-4" />
									{:else}
										{step.id}
									{/if}
								</button>
								<span
									class={cn(
										'hidden text-xs font-medium sm:block',
										isActive ? 'text-foreground' : 'text-muted-foreground'
									)}
								>
									{step.name}
								</span>
							</li>
						{/each}
					</ol>
				</div>

				<form bind:this={formElement} id="venue-form" method="POST" onsubmit={handleSubmit}>
					<!-- Nilai lengkap selalu terkirim meski step tidak aktif -->
					<input type="hidden" name="namaVenue" value={formData.namaVenue} />
					<input type="hidden" name="slug" value={formData.slug} />
					<input type="hidden" name="namaLengkap" value={formData.namaLengkap} />
					<input type="hidden" name="nik" value={formData.nik} />
					<input type="hidden" name="nomorWhatsApp" value={formData.nomorWhatsApp} />
					<input type="hidden" name="alamat" value={formData.alamat} />
					<input type="hidden" name="ktp_key" value={ktpObjectKey} />

					<!-- Input file tetap ter-mount untuk memicu upload ke storage.
					     Tanpa `name`, file tidak ikut dikirim ke server action. -->
					<input
						id="ktp"
						type="file"
						accept="image/jpeg,image/png,image/webp,image/heic,image/heif"
						class="hidden"
						onchange={handleKtpChange}
						disabled={isPending || ktpProcessing}
					/>

					<FieldGroup>
						{#if currentStep === 1}
							<Field data-invalid={!!fieldErrors.namaVenue}>
								<FieldLabel for="namaVenue">Nama Venue *</FieldLabel>
								<Input
									id="namaVenue"
									placeholder="Contoh: Arena Futsal Bandung"
									value={formData.namaVenue}
									oninput={(event) => updateField('namaVenue', event.currentTarget.value)}
									onblur={autoSlug}
									aria-invalid={!!fieldErrors.namaVenue}
								/>
								{#if fieldErrors.namaVenue}<FieldError>{fieldErrors.namaVenue}</FieldError>{/if}
							</Field>

							<SlugField
								bind:value={formData.slug}
								suffix={rootDomain}
								error={fieldErrors.slug}
								disabled={isPending}
								onStatusChange={(status) => (slugStatus = status)}
							/>
						{/if}

						{#if currentStep === 2}
							<Field data-invalid={!!fieldErrors.namaLengkap}>
								<FieldLabel for="namaLengkap">Nama Lengkap (sesuai KTP) *</FieldLabel>
								<Input
									id="namaLengkap"
									autocomplete="name"
									placeholder="Nama lengkap sesuai KTP"
									value={formData.namaLengkap}
									oninput={(event) => updateField('namaLengkap', event.currentTarget.value)}
									aria-invalid={!!fieldErrors.namaLengkap}
								/>
								{#if fieldErrors.namaLengkap}<FieldError>{fieldErrors.namaLengkap}</FieldError>{/if}
							</Field>

							<Field data-invalid={!!fieldErrors.nik}>
								<FieldLabel for="nik">NIK *</FieldLabel>
								<Input
									id="nik"
									inputmode="numeric"
									placeholder="16 digit NIK"
									maxlength={16}
									value={formData.nik}
									oninput={(event) =>
										updateField('nik', event.currentTarget.value.replace(/\D/g, '').slice(0, 16))}
									aria-invalid={!!fieldErrors.nik}
								/>
								{#if fieldErrors.nik}<FieldError>{fieldErrors.nik}</FieldError>{/if}
							</Field>

							<Field data-invalid={!!fieldErrors.nomorWhatsApp}>
								<FieldLabel for="nomorWhatsApp">Nomor HP / WhatsApp *</FieldLabel>
								<Input
									id="nomorWhatsApp"
									type="tel"
									placeholder="08xxxxxxxxxx"
									value={formData.nomorWhatsApp}
									oninput={(event) => updateField('nomorWhatsApp', event.currentTarget.value)}
									aria-invalid={!!fieldErrors.nomorWhatsApp}
								/>
								{#if fieldErrors.nomorWhatsApp}
									<FieldError>{fieldErrors.nomorWhatsApp}</FieldError>
								{/if}
							</Field>

							<Field>
								<FieldLabel for="email">Email</FieldLabel>
								<Input id="email" type="email" value={defaultEmail} disabled readonly class="bg-muted/50" />
							</Field>

							<Field data-invalid={!!fieldErrors.ktp}>
								<FieldLabel>Upload KTP *</FieldLabel>

								{#if ktpPreview}
									<div class="border-input flex items-center gap-3 rounded-lg border p-2">
										<img src={ktpPreview} alt="Preview KTP" class="h-20 w-32 rounded-md object-cover" />
										<div class="min-w-0 flex-1">
											<p class="truncate text-sm font-medium">{ktpFile?.name ?? 'Foto KTP'}</p>
											{#if ktpProcessing}
												<p class="text-muted-foreground flex items-center gap-1.5 text-xs">
													<LoaderCircleIcon class="size-3.5 animate-spin" />
													Mengonversi ke WebP...
												</p>
											{:else}
												<p class="flex items-center gap-1.5 text-xs text-emerald-600 dark:text-emerald-400">
													<CheckIcon class="size-3.5" />
													Siap diunggah (WebP)
												</p>
											{/if}
										</div>
										<Button
											type="button"
											variant="ghost"
											size="icon-sm"
											onclick={clearKtp}
											disabled={isPending || ktpProcessing}
										>
											<XIcon />
										</Button>
									</div>
								{:else}
									<label
										for="ktp"
										class={cn(
											'border-input hover:bg-muted/50 flex cursor-pointer flex-col items-center justify-center gap-2 rounded-lg border border-dashed px-4 py-6 text-center transition-colors',
											(isPending || ktpProcessing) && 'pointer-events-none opacity-60'
										)}
									>
										<UploadIcon class="text-muted-foreground size-4" />
										<span class="text-muted-foreground text-xs">
											Klik untuk pilih foto KTP (JPG/PNG, maks 5 MB)
										</span>
									</label>
								{/if}

								{#if fieldErrors.ktp}<FieldError>{fieldErrors.ktp}</FieldError>{/if}
							</Field>
						{/if}

						{#if currentStep === 3}
							<AddressFields
								value={{
									alamat: formData.alamat,
									kota: formData.kota,
									kecamatan: formData.kecamatan,
									kelurahan: formData.kelurahan,
									provinsi: formData.provinsi,
									kodePos: formData.kodePos
								}}
								onChange={(field, value) => updateField(field, value)}
								errors={fieldErrors}
							/>

							<Field data-invalid={!!fieldErrors.alamat}>
								<FieldLabel for="alamat">Alamat Lengkap *</FieldLabel>
								<Textarea
									id="alamat"
									placeholder="Nama jalan, nomor, patokan lokasi venue"
									class="min-h-20"
									value={formData.alamat}
									oninput={(event) => updateField('alamat', event.currentTarget.value)}
									aria-invalid={!!fieldErrors.alamat}
								/>
								{#if fieldErrors.alamat}<FieldError>{fieldErrors.alamat}</FieldError>{/if}
							</Field>
						{/if}

						{#if submitError}
							<p class="text-destructive text-sm" role="alert">{submitError}</p>
						{/if}
					</FieldGroup>
				</form>
			</CardContent>
		</Card>

		<Card class="w-full">
			<CardContent class="flex gap-3">
				{#if currentStep > 1}
					<Button type="button" variant="outline" onclick={prevStep} class="flex-1">
						<ArrowLeftIcon />
						Kembali
					</Button>
				{/if}

				{#if currentStep < STEPS.length}
					<Button
						type="button"
						onclick={nextStep}
						disabled={isPending || (currentStep === 1 && !stepOneReady)}
						class="flex-1"
					>
						{#if currentStep === 1 && slugStatus === 'checking'}
							<LoaderCircleIcon class="animate-spin" />
							Memeriksa...
						{:else}
							Lanjut
							<ArrowRightIcon />
						{/if}
					</Button>
				{:else}
					<Button type="submit" form="venue-form" disabled={isPending} class="flex-1">
						{#if isPending}
							<LoaderCircleIcon class="animate-spin" />
							Memproses...
						{:else}
							Daftar Sekarang
						{/if}
					</Button>
				{/if}
			</CardContent>
		</Card>
	</main>
{/if}
