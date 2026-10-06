<script lang="ts">
	import CheckIcon from '@lucide/svelte/icons/check';
	import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';
	import { Field, FieldError, FieldLabel } from '#lib/components/ui/field/index.js';
	import RegionCombobox from './region-combobox.svelte';
	import type { AddressSelection, Region } from './types.js';

	let {
		value,
		onChange,
		errors = {}
	}: {
		value: AddressSelection;
		onChange: (field: keyof AddressSelection, value: string) => void;
		errors?: Partial<Record<keyof AddressSelection, string>>;
	} = $props();

	let provinces = $state<Region[]>([]);
	let regencies = $state<Region[]>([]);
	let districts = $state<Region[]>([]);
	let villages = $state<Region[]>([]);

	let selectedProvince = $state<Region | null>(null);
	let selectedRegency = $state<Region | null>(null);
	let selectedDistrict = $state<Region | null>(null);
	let selectedVillage = $state<Region | null>(null);

	let loading = $state({
		provinces: false,
		regencies: false,
		districts: false,
		villages: false,
		kodePos: false
	});
	let loadError = $state(false);

	type KodePosStatus = 'idle' | 'found' | 'not-found' | 'error';
	let kodePosStatus = $state<KodePosStatus>('idle');

	async function fetchRegions(path: string, signal: AbortSignal): Promise<Region[]> {
		const response = await fetch(`/api/regions?path=${encodeURIComponent(path)}`, {
			signal
		});

		if (!response.ok) {
			throw new Error('Gagal memuat data wilayah');
		}

		const payload = (await response.json()) as { data?: Region[] };

		return payload.data ?? [];
	}

	/** Normalisasi nama daerah untuk pencocokan yang toleran. */
	function normalizeName(name: string): string {
		return name
			.toLowerCase()
			.replace(/^(kota|kabupaten|kab\.?|kecamatan|kec\.?|kelurahan|desa)\s+/i, '')
			.replace(/[^a-z0-9]/g, '')
			.trim();
	}

	function findByName(items: Region[], name?: string): Region | null {
		if (!name) return null;

		const target = normalizeName(name);
		if (!target) return null;

		return (
			items.find((item) => normalizeName(item.name) === target) ??
			items.find((item) => normalizeName(item.name).includes(target)) ??
			null
		);
	}

	// Muat provinsi sekali, lalu pulihkan pilihan tersimpan.
	$effect(() => {
		const controller = new AbortController();
		loading.provinces = true;
		loadError = false;

		fetchRegions('/provinces', controller.signal)
			.then((items) => {
				provinces = items;
				const match = findByName(items, value.provinsi);
				if (match) selectedProvince = match;
			})
			.catch((err: unknown) => {
				if ((err as Error)?.name !== 'AbortError') {
					provinces = [];
					loadError = true;
				}
			})
			.finally(() => {
				loading.provinces = false;
			});

		return () => controller.abort();
	});

	$effect(() => {
		const province = selectedProvince;

		if (!province) {
			regencies = [];
			return;
		}

		const controller = new AbortController();
		loading.regencies = true;

		fetchRegions(`/regencies/${province.code}`, controller.signal)
			.then((items) => {
				regencies = items;
				const match = findByName(items, value.kota);
				if (match) selectedRegency = match;
			})
			.catch(() => {
				regencies = [];
			})
			.finally(() => {
				loading.regencies = false;
			});

		return () => controller.abort();
	});

	$effect(() => {
		const regency = selectedRegency;

		if (!regency) {
			districts = [];
			return;
		}

		const controller = new AbortController();
		loading.districts = true;

		fetchRegions(`/districts/${regency.code}`, controller.signal)
			.then((items) => {
				districts = items;
				const match = findByName(items, value.kecamatan);
				if (match) selectedDistrict = match;
			})
			.catch(() => {
				districts = [];
			})
			.finally(() => {
				loading.districts = false;
			});

		return () => controller.abort();
	});

	$effect(() => {
		const district = selectedDistrict;

		if (!district) {
			villages = [];
			return;
		}

		const controller = new AbortController();
		loading.villages = true;

		fetchRegions(`/villages/${district.code}`, controller.signal)
			.then((items) => {
				villages = items;
				const match = findByName(items, value.kelurahan);
				if (match) selectedVillage = match;
			})
			.catch(() => {
				villages = [];
			})
			.finally(() => {
				loading.villages = false;
			});

		return () => controller.abort();
	});

	function selectProvince(item: Region) {
		selectedProvince = item;
		selectedRegency = null;
		selectedDistrict = null;
		selectedVillage = null;
		regencies = [];
		districts = [];
		villages = [];
		kodePosStatus = 'idle';

		onChange('provinsi', item.name);
		onChange('kota', '');
		onChange('kecamatan', '');
		onChange('kelurahan', '');
		onChange('kodePos', '');
	}

	function selectRegency(item: Region) {
		selectedRegency = item;
		selectedDistrict = null;
		selectedVillage = null;
		districts = [];
		villages = [];
		kodePosStatus = 'idle';

		onChange('kota', item.name);
		onChange('kecamatan', '');
		onChange('kelurahan', '');
		onChange('kodePos', '');
	}

	function selectDistrict(item: Region) {
		selectedDistrict = item;
		selectedVillage = null;
		villages = [];
		kodePosStatus = 'idle';

		onChange('kecamatan', item.name);
		onChange('kelurahan', '');
		onChange('kodePos', '');
	}

	function selectVillage(item: Region) {
		selectedVillage = item;
		onChange('kelurahan', item.name);

		// Susun alamat otomatis, jangan timpa alamat yang sudah diketik manual.
		if (!value.alamat.trim()) {
			const parts = [item.name, selectedDistrict?.name, selectedRegency?.name].filter(
				Boolean
			);
			onChange('alamat', parts.join(', '));
		}

		// Ambil kode pos otomatis untuk kelurahan ini.
		lookupKodePos(item.name);
	}

	async function lookupKodePos(kelurahan: string) {
		loading.kodePos = true;
		kodePosStatus = 'idle';
		onChange('kodePos', '');

		try {
			const params = new URLSearchParams({ kelurahan });

			if (selectedDistrict?.name) params.set('kecamatan', selectedDistrict.name);
			if (selectedRegency?.name) params.set('kota', selectedRegency.name);

			const response = await fetch(`/api/kode-pos?${params.toString()}`);

			if (!response.ok) {
				kodePosStatus = 'not-found';
				return;
			}

			const result = (await response.json()) as { code?: string };

			if (result.code) {
				onChange('kodePos', result.code);
				kodePosStatus = 'found';
			} else {
				kodePosStatus = 'not-found';
			}
		} catch {
			kodePosStatus = 'error';
		} finally {
			loading.kodePos = false;
		}
	}
</script>

{#if loadError}
	<p
		class="border-destructive/30 bg-destructive/5 text-destructive rounded-lg border px-3 py-2 text-xs"
	>
		Gagal memuat data wilayah. Periksa koneksi lalu muat ulang halaman.
	</p>
{/if}

<Field data-invalid={!!errors.provinsi}>
	<FieldLabel>Provinsi *</FieldLabel>
	<RegionCombobox
		bind:value={selectedProvince}
		onSelect={selectProvince}
		items={provinces}
		placeholder="Pilih provinsi"
		searchPlaceholder="Cari provinsi..."
		loading={loading.provinces}
	/>
	{#if errors.provinsi}<FieldError>{errors.provinsi}</FieldError>{/if}
</Field>

<div class="grid gap-4 sm:grid-cols-2">
	<Field data-invalid={!!errors.kota}>
		<FieldLabel>Kota/Kabupaten *</FieldLabel>
		<RegionCombobox
			bind:value={selectedRegency}
			onSelect={selectRegency}
			items={regencies}
			placeholder="Pilih kota/kab"
			searchPlaceholder="Cari kota/kab..."
			disabled={!selectedProvince}
			loading={loading.regencies}
		/>
		{#if errors.kota}<FieldError>{errors.kota}</FieldError>{/if}
	</Field>

	<Field data-invalid={!!errors.kecamatan}>
		<FieldLabel>Kecamatan</FieldLabel>
		<RegionCombobox
			bind:value={selectedDistrict}
			onSelect={selectDistrict}
			items={districts}
			placeholder="Pilih kecamatan"
			searchPlaceholder="Cari kecamatan..."
			disabled={!selectedRegency}
			loading={loading.districts}
		/>
		{#if errors.kecamatan}<FieldError>{errors.kecamatan}</FieldError>{/if}
	</Field>
</div>

<div class="grid gap-4 sm:grid-cols-2">
	<Field data-invalid={!!errors.kelurahan}>
		<FieldLabel>Kelurahan</FieldLabel>
		<RegionCombobox
			bind:value={selectedVillage}
			onSelect={selectVillage}
			items={villages}
			placeholder="Pilih kelurahan"
			searchPlaceholder="Cari kelurahan..."
			disabled={!selectedDistrict}
			loading={loading.villages}
		/>
		{#if errors.kelurahan}<FieldError>{errors.kelurahan}</FieldError>{/if}
	</Field>

	<Field data-invalid={!!errors.kodePos}>
		<FieldLabel for="kodePos">Kode Pos *</FieldLabel>

		<div class="relative">
			<input
				id="kodePos"
				name="kodePos"
				inputmode="numeric"
				value={value.kodePos}
				oninput={(event) =>
					onChange(
						'kodePos',
						event.currentTarget.value.replace(/\D/g, '').slice(0, 5)
					)}
				placeholder={loading.kodePos ? 'Mencari kode pos...' : 'Terisi otomatis'}
				maxlength={5}
				disabled={loading.kodePos}
				class="border-input bg-input/20 dark:bg-input/30 focus-visible:border-ring focus-visible:ring-ring/30 aria-invalid:border-destructive aria-invalid:ring-destructive/20 h-7 w-full rounded-md border px-2 py-0.5 text-sm transition-colors outline-none focus-visible:ring-2 aria-invalid:ring-2 disabled:opacity-70 md:text-xs/relaxed"
				aria-invalid={!!errors.kodePos}
			/>

			{#if loading.kodePos}
				<LoaderCircleIcon
					class="text-muted-foreground pointer-events-none absolute top-1/2 right-2 size-3.5 -translate-y-1/2 animate-spin"
				/>
			{:else if kodePosStatus === 'found'}
				<CheckIcon
					class="pointer-events-none absolute top-1/2 right-2 size-3.5 -translate-y-1/2 text-emerald-600 dark:text-emerald-400"
				/>
			{/if}
		</div>

		{#if errors.kodePos}
			<FieldError>{errors.kodePos}</FieldError>
		{:else if kodePosStatus === 'not-found'}
			<p class="text-muted-foreground text-xs/relaxed">
				Kode pos tidak ditemukan otomatis — isi manual.
			</p>
		{:else if kodePosStatus === 'error'}
			<p class="text-muted-foreground text-xs/relaxed">
				Gagal mengambil kode pos — isi manual.
			</p>
		{:else if kodePosStatus === 'found'}
			<p class="text-xs/relaxed text-emerald-600 dark:text-emerald-400">
				Terisi otomatis dari kelurahan.
			</p>
		{:else}
			<p class="text-muted-foreground text-xs/relaxed">
				Terisi otomatis setelah memilih kelurahan.
			</p>
		{/if}
	</Field>
</div>

<!-- Nilai wilayah tersembunyi agar ikut terkirim lewat form POST -->
<input type="hidden" name="provinsi" value={value.provinsi} />
<input type="hidden" name="kota" value={value.kota} />
<input type="hidden" name="kecamatan" value={value.kecamatan} />
<input type="hidden" name="kelurahan" value={value.kelurahan} />
