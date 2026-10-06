<script lang="ts">
    import TenantHeader from '#lib/components/tenant-header.svelte';

	let { data, form } = $props();

	let selectedSlot = $state<string | null>(null);

	const rupiah = new Intl.NumberFormat('id-ID', {
		style: 'currency',
		currency: 'IDR',
		maximumFractionDigits: 0
	});

	function formatTime(value: string) {
		return new Intl.DateTimeFormat('id-ID', {
			timeZone: data.venue.timezone,
			hour: '2-digit',
			minute: '2-digit',
			hourCycle: 'h23'
		}).format(new Date(value));
	}
</script>

<svelte:head>
	<title>Booking {data.court.name} | Lapanganku</title>
	<meta
		name="description"
		content={`Pilih jadwal ${data.court.name} di ${data.venue.name}.`}
	/>
</svelte:head>

<main class="min-h-screen bg-background text-foreground">
    <TenantHeader user={data.user} />

	<section class="border-b bg-muted">
		<div class="mx-auto max-w-6xl px-6 py-10">
			<a
				href="/"
				class="text-sm font-medium text-muted-foreground transition hover:text-foreground"
			>
				← Kembali ke venue
			</a>

			<div class="mt-6 flex flex-col gap-6 md:flex-row md:items-end md:justify-between">
				<div>
					<p class="text-sm font-medium uppercase tracking-widest text-muted-foreground">
						{data.venue.name}
					</p>

					<h1 class="mt-2 text-3xl font-bold tracking-tight md:text-4xl">
						{data.court.name}
					</h1>

					<p class="mt-2 capitalize text-muted-foreground">
						{data.court.sport}
					</p>
				</div>

				<div class="rounded-xl border bg-background px-5 py-4">
					<p class="text-xs text-muted-foreground">
						Harga per slot
					</p>

					<p class="mt-1 text-xl font-bold">
						{rupiah.format(data.court.price_per_slot)}
					</p>

					<p class="text-xs text-muted-foreground">
						{data.court.slot_duration_minutes} menit
					</p>
				</div>
			</div>
		</div>
	</section>

	<section class="mx-auto grid max-w-6xl gap-8 px-6 py-10 lg:grid-cols-[1fr_20rem]">
		<div>
			<div>
				<h2 class="text-xl font-bold">
					Pilih tanggal
				</h2>

				<p class="mt-1 text-sm text-muted-foreground">
					Jadwal ditampilkan dalam zona waktu {data.venue.timezone}.
				</p>
			</div>

			<form method="GET" class="mt-5 flex max-w-sm gap-3">
				<input
					name="date"
					type="date"
					required
					value={data.date}
					class="h-10 flex-1 rounded-md border bg-background px-3 text-sm outline-none focus:border-ring focus:ring-2 focus:ring-ring/30"
				/>

				<button
					type="submit"
					class="rounded-md bg-primary px-5 text-sm font-medium text-primary-foreground transition hover:opacity-90"
				>
					Lihat jadwal
				</button>
			</form>

			<div class="mt-10">
				<h2 class="text-xl font-bold">
					Pilih jam
				</h2>

				<p class="mt-1 text-sm text-muted-foreground">
					Slot abu-abu sudah terisi atau telah lewat.
				</p>

				{#if data.availability.slots.length === 0}
					<div class="mt-5 rounded-xl border border-dashed p-10 text-center">
						<h3 class="font-semibold">
							Tidak ada jadwal
						</h3>

						<p class="mt-2 text-sm text-muted-foreground">
							Venue tutup atau jam operasional belum diatur.
						</p>
					</div>
				{:else}
					<div class="mt-5 grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4">
						{#each data.availability.slots as slot}
							<button
								type="button"
								disabled={!slot.available}
								aria-pressed={selectedSlot === slot.starts_at}
								onclick={() => {
									selectedSlot = slot.starts_at;
								}}
								class={[
									'rounded-lg border px-4 py-3 text-sm font-medium transition',
									selectedSlot === slot.starts_at
										? 'border-primary bg-primary text-primary-foreground'
										: slot.available
											? 'bg-card hover:border-primary'
											: 'cursor-not-allowed bg-muted text-muted-foreground opacity-60'
								]}
							>
								{formatTime(slot.starts_at)}
								<span class="mx-1">–</span>
								{formatTime(slot.ends_at)}
							</button>
						{/each}
					</div>
				{/if}
			</div>
		</div>

		<aside class="h-fit rounded-xl border bg-card p-6 lg:sticky lg:top-6">
			<h2 class="font-semibold">
				Ringkasan booking
			</h2>

			<dl class="mt-5 space-y-4 text-sm">
				<div>
					<dt class="text-muted-foreground">Venue</dt>
					<dd class="mt-1 font-medium">{data.venue.name}</dd>
				</div>

				<div>
					<dt class="text-muted-foreground">Lapangan</dt>
					<dd class="mt-1 font-medium">{data.court.name}</dd>
				</div>

				<div>
					<dt class="text-muted-foreground">Tanggal</dt>
					<dd class="mt-1 font-medium">{data.date}</dd>
				</div>

				<div>
					<dt class="text-muted-foreground">Jam</dt>
					<dd class="mt-1 font-medium">
						{#if selectedSlot}
							{formatTime(selectedSlot)}
						{:else}
							Belum dipilih
						{/if}
					</dd>
				</div>

				<div class="border-t pt-4">
					<dt class="text-muted-foreground">Total</dt>
					<dd class="mt-1 text-xl font-bold">
						{rupiah.format(data.court.price_per_slot)}
					</dd>
				</div>
			</dl>

			{#if form?.message}
				<p
					role="alert"
					class="mt-5 rounded-md border border-destructive/30 bg-destructive/10 px-3 py-2 text-sm text-destructive"
				>
					{form.message}
				</p>
			{/if}

			<form method="POST" class="mt-6">
				<input
					type="hidden"
					name="starts_at"
					value={selectedSlot ?? ''}
				/>

				<button
					type="submit"
					disabled={!selectedSlot}
					class="w-full rounded-md bg-primary px-4 py-3 text-sm font-medium text-primary-foreground transition hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-50"
				>
					Lanjutkan booking
				</button>
			</form>

			<p class="mt-4 text-center text-xs text-muted-foreground">
				Kamu akan diminta login jika belum masuk.
			</p>
		</aside>
	</section>
</main>