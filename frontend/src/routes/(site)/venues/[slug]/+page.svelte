<script lang="ts">
    import TenantHeader from '#lib/components/tenant-header.svelte';

	let { data } = $props();

	const rupiah = new Intl.NumberFormat('id-ID', {
		style: 'currency',
		currency: 'IDR',
		maximumFractionDigits: 0
	});
</script>

<svelte:head>
	<title>{data.venue.name} | Lapanganku</title>
	<meta
		name="description"
		content={`Booking lapangan di ${data.venue.name}.`}
	/>
</svelte:head>

<main class="min-h-screen bg-background text-foreground">
    <TenantHeader user={data.user} />

	<section class="border-b bg-muted">
		<div class="mx-auto max-w-6xl px-6 py-12 md:py-16">
			<p class="text-sm font-medium uppercase tracking-widest text-muted-foreground">
				Booking lapangan
			</p>

			<h1 class="mt-3 text-3xl font-bold tracking-tight md:text-5xl">
				{data.venue.name}
			</h1>

			<p class="mt-4 max-w-2xl text-muted-foreground">
				{data.venue.address}
			</p>

			<div class="mt-5 flex flex-wrap gap-3 text-sm text-muted-foreground">
				<span class="rounded-full border bg-background px-3 py-1">
					{data.venue.timezone}
				</span>

				{#if data.venue.whatsapp}
					<a
						href={`https://wa.me/${data.venue.whatsapp}`}
						target="_blank"
						rel="noreferrer"
						class="rounded-full border bg-background px-3 py-1 transition hover:text-foreground"
					>
						Hubungi WhatsApp
					</a>
				{/if}
			</div>
		</div>
	</section>

	<section class="mx-auto max-w-6xl px-6 py-12">
		<div>
			<h2 class="text-2xl font-bold">
				Pilih lapangan
			</h2>

			<p class="mt-2 text-sm text-muted-foreground">
				Pilih lapangan untuk melihat jadwal yang tersedia.
			</p>
		</div>

		{#if data.venue.courts.length === 0}
			<div class="mt-8 rounded-xl border border-dashed p-10 text-center">
				<h3 class="font-semibold">Belum ada lapangan aktif</h3>

				<p class="mt-2 text-sm text-muted-foreground">
					Silakan hubungi venue untuk informasi lebih lanjut.
				</p>
			</div>
		{:else}
			<div class="mt-8 grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
				{#each data.venue.courts as court}
					<article class="flex flex-col rounded-xl border bg-card p-6 shadow-sm">
						<div class="flex items-start justify-between gap-4">
							<div>
								<h3 class="text-lg font-semibold">
									{court.name}
								</h3>

								<p class="mt-1 capitalize text-sm text-muted-foreground">
									{court.sport}
								</p>
							</div>

							<span class="rounded-full bg-muted px-3 py-1 text-xs font-medium">
								{court.slot_duration_minutes} menit
							</span>
						</div>

						<div class="mt-8">
							<p class="text-xs text-muted-foreground">
								Mulai dari
							</p>

							<p class="mt-1 text-xl font-bold">
								{rupiah.format(court.price_per_slot)}
							</p>

							<p class="text-xs text-muted-foreground">
								per slot
							</p>
						</div>

						<a
							href={`/courts/${court.id}/book`}
							class="mt-6 rounded-md bg-primary px-4 py-2.5 text-center text-sm font-medium text-primary-foreground transition hover:opacity-90"
						>
							Lihat jadwal
						</a>
					</article>
				{/each}
			</div>
		{/if}
	</section>
</main>