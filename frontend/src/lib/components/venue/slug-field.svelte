<script lang="ts">
	import CheckIcon from '@lucide/svelte/icons/check';
	import CircleXIcon from '@lucide/svelte/icons/circle-x';
	import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';
	import XIcon from '@lucide/svelte/icons/x';
	import { Input } from '#lib/components/ui/input/index.js';
	import { Field, FieldError, FieldLabel } from '#lib/components/ui/field/index.js';
	import { validateUsername } from '#lib/validation.js';
	import type { SlugStatus } from './types.js';

	let {
		value = $bindable(''),
		label = 'Alamat publik',
		suffix,
		error,
		disabled = false,
		checkEnabled = true,
		debounceMs = 400,
		onStatusChange
	}: {
		value?: string;
		label?: string;
		suffix: string;
		error?: string;
		disabled?: boolean;
		checkEnabled?: boolean;
		debounceMs?: number;
		onStatusChange?: (status: SlugStatus) => void;
	} = $props();

	let status = $state<SlugStatus>('idle');
	let statusMessage = $state('');
	let inputRef = $state<HTMLInputElement | null>(null);

	// Nilai internal (non-reaktif) agar tidak memicu ulang effect pengecekan.
	let lastChecked = '';

	$effect(() => {
		onStatusChange?.(status);
	});

	// Cek ketersediaan otomatis dengan debounce saat user mengetik.
	$effect(() => {
		const current = value;
		const isDisabled = disabled;
		const enabled = checkEnabled;

		if (isDisabled || !enabled) {
			return;
		}

		const trimmed = current.trim().toLowerCase();
		const validationError = validateUsername(trimmed);

		if (!trimmed) {
			status = 'idle';
			statusMessage = '';
			lastChecked = '';
			return;
		}

		if (validationError) {
			status = 'invalid';
			statusMessage = validationError;
			lastChecked = '';
			return;
		}

		if (trimmed === lastChecked) {
			return;
		}

		status = 'checking';
		statusMessage = '';

		const timer = setTimeout(() => {
			checkAvailability(trimmed);
		}, debounceMs);

		return () => clearTimeout(timer);
	});

	async function checkAvailability(slug: string) {
		try {
			const response = await fetch(
				`/api/check-slug?slug=${encodeURIComponent(slug)}`
			);

			// Abaikan respons usang bila user sudah mengetik lagi.
			if (value.trim().toLowerCase() !== slug) return;

			const result = (await response.json()) as {
				status: SlugStatus;
				message?: string;
			};

			status = result.status;
			statusMessage = result.message ?? '';
			lastChecked = slug;
		} catch {
			if (value.trim().toLowerCase() !== slug) return;

			status = 'error';
			statusMessage = 'Gagal memeriksa alamat publik';
			lastChecked = '';
		}
	}

	function handleInput(event: Event) {
		const target = event.currentTarget as HTMLInputElement;
		const next = target.value.toLowerCase().replace(/[^a-z0-9-]/g, '');

		// Sinkronkan DOM bila karakter tidak valid (mis. spasi) dihapus.
		if (target.value !== next) {
			target.value = next;
		}

		value = next;
	}

	function clear() {
		value = '';
		status = 'idle';
		statusMessage = '';
		lastChecked = '';
		inputRef?.focus();
	}

	const invalid = $derived(!!error || status === 'taken' || status === 'invalid');
	const message = $derived(error || statusMessage);
</script>

<Field data-invalid={invalid}>
	<FieldLabel for="slug">{label} *</FieldLabel>

	<div
		class="border-input bg-input/20 dark:bg-input/30 focus-within:border-ring focus-within:ring-ring/30 flex items-center overflow-hidden rounded-md border transition-colors focus-within:ring-2"
	>
		<Input
			bind:ref={inputRef}
			id="slug"
			type="text"
			autocomplete="off"
			spellcheck={false}
			inputmode="url"
			placeholder="lapangan-anda"
			value={value}
			oninput={handleInput}
			disabled={disabled}
			aria-invalid={invalid}
			class="h-8 rounded-none border-0 bg-transparent shadow-none focus-visible:ring-0 aria-invalid:ring-0"
		/>

		<span class="text-muted-foreground shrink-0 border-l px-2 text-xs/relaxed">
			.{suffix}
		</span>

		{#if value}
			<button
				type="button"
				onclick={clear}
				disabled={disabled}
				aria-label="Bersihkan alamat publik"
				class="text-muted-foreground hover:text-foreground shrink-0 border-l px-2 transition-colors disabled:opacity-50"
			>
				<XIcon class="size-3.5" />
			</button>
		{/if}
	</div>

	<div class="min-w-0">
		{#if status === 'checking'}
			<p class="text-muted-foreground flex items-center gap-1.5 text-xs/relaxed">
				<LoaderCircleIcon class="size-3.5 animate-spin" />
				Memeriksa ketersediaan...
			</p>
		{:else if status === 'available'}
			<p class="flex items-center gap-1.5 text-xs/relaxed text-emerald-600 dark:text-emerald-400">
				<CheckIcon class="size-3.5" />
				{statusMessage}
			</p>
		{:else if invalid && message}
			<FieldError>{message}</FieldError>
		{:else if message}
			<p class="text-muted-foreground flex items-center gap-1.5 text-xs/relaxed">
				<CircleXIcon class="size-3.5" />
				{message}
			</p>
		{:else}
			<p class="text-muted-foreground text-xs/relaxed">
				Alamat publik venue kamu, tanpa spasi.
			</p>
		{/if}
	</div>
</Field>
