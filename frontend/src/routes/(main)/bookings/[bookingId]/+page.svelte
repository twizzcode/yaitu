<script lang="ts">
	import { onMount, untrack } from 'svelte';
	import { enhance } from '$app/forms';
	import { invalidateAll } from '$app/navigation';
	import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';
	import QrCodeIcon from '@lucide/svelte/icons/qr-code';
	import { Button } from '#lib/components/ui/button/index.js';
	import { toast } from 'svelte-sonner';

	let { data, form } = $props();

	let now = $state(0);
	let paying = $state(false);
	let payment = $state(untrack(() => data.payment));

	// Sinkronkan payment dari server saat data berubah.
	$effect(() => {
		payment = data.payment;
	});

	$effect(() => {
		if (form?.success && form.payment) {
			payment = form.payment;
		}

		if (form?.message && !form?.success) {
			toast.error(form.message);
			paying = false;
		}
	});

	const rupiah = new Intl.NumberFormat('id-ID', {
		style: 'currency',
		currency: 'IDR',
		maximumFractionDigits: 0
	});

	function formatDate(value: string) {
		return new Intl.DateTimeFormat('id-ID', {
			timeZone: data.booking.timezone,
			weekday: 'long',
			day: 'numeric',
			month: 'long',
			year: 'numeric'
		}).format(new Date(value));
	}

	function formatTime(value: string) {
		return new Intl.DateTimeFormat('id-ID', {
			timeZone: data.booking.timezone,
			hour: '2-digit',
			minute: '2-digit',
			hourCycle: 'h23'
		}).format(new Date(value));
	}

	const statusLabels = {
		pending_payment: 'Menunggu pembayaran',
		confirmed: 'Terkonfirmasi',
		cancelled: 'Dibatalkan',
		expired: 'Kedaluwarsa',
		completed: 'Selesai'
	};

	const remainingMilliseconds = $derived.by(() => {
		if (!data.booking.payment_expires_at || now === 0) {
			return 0;
		}

		return Math.max(
			new Date(data.booking.payment_expires_at).getTime() - now,
			0
		);
	});

	const paymentExpired = $derived(
		data.booking.status === 'pending_payment' &&
			now > 0 &&
			remainingMilliseconds === 0
	);

	const displayedStatus = $derived(
		paymentExpired ? 'expired' : data.booking.status
	);

	const remainingTime = $derived.by(() => {
		if (now === 0) {
			return 'Menghitung...';
		}

		const totalSeconds = Math.floor(remainingMilliseconds / 1000);
		const minutes = Math.floor(totalSeconds / 60);
		const seconds = totalSeconds % 60;

		return `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`;
	});

	function statusClass(status: string) {
		switch (status) {
			case 'confirmed':
				return 'border-emerald-500/30 bg-emerald-500/10 text-emerald-700';
			case 'pending_payment':
				return 'border-amber-500/30 bg-amber-500/10 text-amber-700';
			case 'completed':
				return 'border-blue-500/30 bg-blue-500/10 text-blue-700';
			default:
				return 'border-destructive/30 bg-destructive/10 text-destructive';
		}
	}

	onMount(() => {
		now = Date.now();

		const interval = window.setInterval(() => {
			now = Date.now();
		}, 1000);

		return () => {
			window.clearInterval(interval);
		};
	});

	// Polling status pembayaran selama masih pending_payment.
	onMount(() => {
		const poll = window.setInterval(async () => {
			if (data.booking.status !== 'pending_payment') {
				return;
			}

			try {
				const response = await fetch(
					`/api/payments/status/${data.booking.id}`
				);

				if (!response.ok) {
					return;
				}

				const result = (await response.json()) as {
					payment: typeof payment;
				};

				const next = result.payment;

				if (next && next.status !== payment?.status) {
					payment = next;

					if (next.status === 'settlement' || next.status === 'capture') {
						toast.success('Pembayaran berhasil! Booking dikonfirmasi.');
					}

					await invalidateAll();
				}
			} catch {
				// Abaikan error polling.
			}
		}, 4000);

		return () => {
			window.clearInterval(poll);
		};
	});
</script>

<svelte:head>
	<title>Detail Booking | Lapanganku</title>
	<meta
		name="description"
		content={`Detail booking ${data.booking.court_name} di ${data.booking.venue_name}.`}
	/>
</svelte:head>

<main class="min-h-screen bg-muted/40 text-foreground">
	<header class="border-b bg-card">
		<div class="mx-auto flex max-w-5xl items-center justify-between px-6 py-4">
			<a href="/" class="text-lg font-bold">
				Lapanganku
			</a>

			<a
				href="/my-bookings"
				class="text-sm font-medium text-muted-foreground transition hover:text-foreground"
			>
				Booking saya
			</a>
		</div>
	</header>

	<div class="mx-auto max-w-5xl px-6 py-10">
		<div class="flex flex-col gap-5 sm:flex-row sm:items-start sm:justify-between">
			<div>
				<p class="text-sm font-medium uppercase tracking-widest text-muted-foreground">
					Detail booking
				</p>

				<h1 class="mt-2 text-3xl font-bold tracking-tight">
					{data.booking.court_name}
				</h1>

				<p class="mt-2 text-muted-foreground">
					{data.booking.venue_name}
				</p>
			</div>

			<span
				class={[
					'w-fit rounded-full border px-4 py-2 text-sm font-medium',
					statusClass(displayedStatus)
				]}
			>
				{statusLabels[displayedStatus]}
			</span>
		</div>

		<div class="mt-8 grid gap-6 lg:grid-cols-[1fr_22rem]">
			<section class="rounded-xl border bg-card p-6">
				<h2 class="text-lg font-semibold">
					Rincian jadwal
				</h2>

				<dl class="mt-6 divide-y">
					<div class="grid gap-1 py-4 sm:grid-cols-[10rem_1fr]">
						<dt class="text-sm text-muted-foreground">
							Venue
						</dt>

						<dd class="font-medium">
							{data.booking.venue_name}
						</dd>
					</div>

					<div class="grid gap-1 py-4 sm:grid-cols-[10rem_1fr]">
						<dt class="text-sm text-muted-foreground">
							Lapangan
						</dt>

						<dd class="font-medium">
							{data.booking.court_name}
						</dd>
					</div>

					<div class="grid gap-1 py-4 sm:grid-cols-[10rem_1fr]">
						<dt class="text-sm text-muted-foreground">
							Tanggal
						</dt>

						<dd class="font-medium">
    						{formatDate(data.booking.starts_at)}
						</dd>
					</div>

					<div class="grid gap-1 py-4 sm:grid-cols-[10rem_1fr]">
						<dt class="text-sm text-muted-foreground">
							Jam
						</dt>

						<dd class="font-medium">
    						{formatTime(data.booking.starts_at)}
							–
							{formatTime(data.booking.ends_at)}
						</dd>
					</div>

					<div class="grid gap-1 py-4 sm:grid-cols-[10rem_1fr]">
						<dt class="text-sm text-muted-foreground">
							Zona waktu
						</dt>

						<dd class="font-medium">
							{data.booking.timezone}
						</dd>
					</div>

					<div class="grid gap-1 py-4 sm:grid-cols-[10rem_1fr]">
						<dt class="text-sm text-muted-foreground">
							ID booking
						</dt>

						<dd class="break-all font-mono text-sm">
							{data.booking.id}
						</dd>
					</div>
				</dl>

				<a
					href={`/venues/${data.booking.venue_slug}`}
					class="mt-6 inline-flex rounded-md border px-4 py-2 text-sm font-medium transition hover:bg-muted"
				>
					Kembali ke venue
				</a>
			</section>

			<aside class="h-fit rounded-xl border bg-card p-6">
				<h2 class="text-lg font-semibold">
					Pembayaran
				</h2>

				<div class="mt-5 border-b pb-5">
					<p class="text-sm text-muted-foreground">
						Total pembayaran
					</p>

					<p class="mt-1 text-3xl font-bold">
						{rupiah.format(data.booking.total_amount)}
					</p>
				</div>

				{#if displayedStatus === 'pending_payment'}
					<div class="mt-5 rounded-lg border border-amber-500/30 bg-amber-500/10 p-4">
						<p class="text-sm text-amber-800">
							Selesaikan pembayaran sebelum waktu habis.
						</p>

						<p class="mt-2 font-mono text-2xl font-bold text-amber-900">
							{remainingTime}
						</p>
					</div>

					{#if payment && payment.qr_url}
						<div class="mt-5 rounded-lg border p-4">
							<p class="flex items-center gap-2 text-sm font-medium">
								<QrCodeIcon class="size-4" />
								Scan QRIS untuk membayar
							</p>

							<img
								src={payment.qr_url}
								alt="QRIS pembayaran"
								class="mx-auto mt-4 h-56 w-56 rounded-lg border bg-white p-2 object-contain"
							/>

							<p class="mt-3 text-center text-xs text-muted-foreground">
								Scan dengan aplikasi apa pun yang mendukung QRIS (Gojek, OVO, Dana,
								ShopeePay, m-banking).
							</p>

							<p class="mt-2 text-center text-xs text-muted-foreground">
								Menunggu pembayaran... halaman ini otomatis diperbarui.
							</p>
						</div>
					{:else}
						<form
							method="POST"
							action="?/pay"
							use:enhance={() => {
								paying = true;
								return async ({ update }) => {
									await update();
									paying = false;
								};
							}}
						>
							<Button type="submit" disabled={paying} class="mt-5 w-full py-6">
								{#if paying}
									<LoaderCircleIcon class="animate-spin" />
									Membuat QRIS...
								{:else}
									<QrCodeIcon />
									Bayar dengan QRIS
								{/if}
							</Button>
						</form>

						<p class="mt-3 text-center text-xs text-muted-foreground">
							Bayar dengan QRIS (Gojek, OVO, Dana, ShopeePay, m-banking).
						</p>
					{/if}
				{:else if displayedStatus === 'confirmed'}
					<div class="mt-5 rounded-lg border border-emerald-500/30 bg-emerald-500/10 p-4">
						<p class="font-medium text-emerald-800">
							Pembayaran berhasil
						</p>

						<p class="mt-1 text-sm text-emerald-700">
							Booking kamu sudah dikonfirmasi.
						</p>
					</div>
				{:else if displayedStatus === 'expired'}
					<div class="mt-5 rounded-lg border border-destructive/30 bg-destructive/10 p-4">
						<p class="font-medium text-destructive">
							Waktu pembayaran habis
						</p>

						<p class="mt-1 text-sm text-destructive">
							Pilih slot kembali untuk membuat booking baru.
						</p>
					</div>

					<a
						href={`/venues/${data.booking.venue_slug}`}
						class="mt-5 block rounded-md bg-primary px-4 py-3 text-center text-sm font-medium text-primary-foreground"
					>
						Pilih jadwal lain
					</a>
				{:else}
					<div class="mt-5 rounded-lg border bg-muted p-4">
						<p class="font-medium">
							{statusLabels[displayedStatus]}
						</p>
					</div>
				{/if}
			</aside>
		</div>
	</div>
</main>