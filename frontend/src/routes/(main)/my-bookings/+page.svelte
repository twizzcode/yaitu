<script lang="ts">
	let { data } = $props();

	const rupiah = new Intl.NumberFormat('id-ID', {
		style: 'currency',
		currency: 'IDR',
		maximumFractionDigits: 0
	});

	const statusLabels = {
		pending_payment: 'Menunggu pembayaran',
		confirmed: 'Terkonfirmasi',
		cancelled: 'Dibatalkan',
		expired: 'Kedaluwarsa',
		completed: 'Selesai'
	};

	function formatDate(value: string, timezone: string) {
		return new Intl.DateTimeFormat('id-ID', {
			timeZone: timezone,
			weekday: 'long',
			day: 'numeric',
			month: 'long',
			year: 'numeric'
		}).format(new Date(value));
	}

	function formatTime(value: string, timezone: string) {
		return new Intl.DateTimeFormat('id-ID', {
			timeZone: timezone,
			hour: '2-digit',
			minute: '2-digit',
			hourCycle: 'h23'
		}).format(new Date(value));
	}

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
</script>

<svelte:head>
	<title>Booking Saya | Lapanganku</title>
	<meta name="description" content="Daftar booking lapangan kamu." />
</svelte:head>

<main class="min-h-screen bg-muted/40 text-foreground">
	<header class="border-b bg-card">
		<div class="mx-auto flex max-w-5xl items-center justify-between px-6 py-4">
			<a href="/" class="text-lg font-bold">
				Lapanganku
			</a>

			<a
				href="/admin"
				class="text-sm font-medium text-muted-foreground transition hover:text-foreground"
			>
				Dashboard
			</a>
		</div>
	</header>

	<div class="mx-auto max-w-5xl px-6 py-10">
		<div>
			<p class="text-sm font-medium uppercase tracking-widest text-muted-foreground">
				Riwayat transaksi
			</p>

			<h1 class="mt-2 text-3xl font-bold tracking-tight">
				Booking saya
			</h1>

			<p class="mt-2 text-muted-foreground">
				Lihat status dan rincian semua booking lapanganmu.
			</p>
		</div>

		{#if data.bookings.length === 0}
			<section
				class="mt-8 rounded-xl border border-dashed bg-card p-10 text-center"
			>
				<h2 class="text-lg font-semibold">
					Belum ada booking
				</h2>

				<p class="mt-2 text-sm text-muted-foreground">
					Pilih venue dan lapangan untuk membuat booking pertama.
				</p>

				<a
					href="/"
					class="mt-6 inline-flex rounded-md bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground"
				>
					Cari lapangan
				</a>
			</section>
		{:else}
			<section class="mt-8 space-y-4">
				{#each data.bookings as booking}
					<article class="rounded-xl border bg-card p-6 shadow-sm">
						<div class="flex flex-col gap-5 sm:flex-row sm:items-start sm:justify-between">
							<div>
								<p class="text-sm font-medium text-muted-foreground">
									{booking.venue_name}
								</p>

								<h2 class="mt-1 text-xl font-semibold">
									{booking.court_name}
								</h2>

								<p class="mt-3 text-sm">
									{formatDate(
										booking.starts_at,
										booking.timezone
									)}
								</p>

								<p class="mt-1 text-sm text-muted-foreground">
									{formatTime(
										booking.starts_at,
										booking.timezone
									)}
									–
									{formatTime(
										booking.ends_at,
										booking.timezone
									)}
								</p>
							</div>

							<div class="flex flex-col items-start gap-3 sm:items-end">
								<span
									class={[
										'rounded-full border px-3 py-1 text-xs font-medium',
										statusClass(booking.status)
									]}
								>
									{statusLabels[booking.status]}
								</span>

								<p class="text-lg font-bold">
									{rupiah.format(booking.total_amount)}
								</p>
							</div>
						</div>

						<div class="mt-6 flex flex-wrap gap-3 border-t pt-5">
							<a
								href={`/bookings/${booking.id}`}
								class="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground"
							>
								Lihat detail
							</a>

							<a
								href={`/venues/${booking.venue_slug}`}
								class="rounded-md border px-4 py-2 text-sm font-medium transition hover:bg-muted"
							>
								Lihat venue
							</a>
						</div>
					</article>
				{/each}
			</section>
		{/if}
	</div>
</main>