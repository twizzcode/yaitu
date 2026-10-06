<script lang="ts">
	import CalendarCogIcon from '@lucide/svelte/icons/calendar-cog';
	import PlusIcon from '@lucide/svelte/icons/plus';
	import { Badge } from '#lib/components/ui/badge/index.js';
	import { Button } from '#lib/components/ui/button/index.js';
	import { Card, CardContent } from '#lib/components/ui/card/index.js';
	import AdminHeroSection from '#lib/components/admin/admin-hero-section.svelte';

	let { data } = $props();

	const rupiah = new Intl.NumberFormat('id-ID', {
		style: 'currency',
		currency: 'IDR',
		maximumFractionDigits: 0
	});
</script>

<svelte:head>
	<title>Lapangan {data.venue.name} | Lapanganku</title>
	<meta name="description" content={`Kelola lapangan ${data.venue.name}.`} />
</svelte:head>

<main class="mx-auto max-w-7xl space-y-5">
	<AdminHeroSection
		title="Kelola Lapangan"
		description={`Atur lapangan dan harga untuk ${data.venue.name}.`}
		size="sm"
	/>

	<div class="flex items-center justify-between gap-4">
		<div>
			<h2 class="text-lg font-semibold">Daftar lapangan</h2>
			<p class="text-sm text-muted-foreground">
				{data.courts.length} lapangan terdaftar.
			</p>
		</div>

		<Button href="/admin/venues/{data.venue.slug}/courts/new">
			<PlusIcon />
			Tambah lapangan
		</Button>
	</div>

	{#if data.courts.length === 0}
		<div class="rounded-xl border border-dashed p-10 text-center">
			<h3 class="font-semibold">Belum ada lapangan</h3>
			<p class="mt-2 text-sm text-muted-foreground">
				Tambahkan lapangan pertama untuk venue ini.
			</p>
		</div>
	{:else}
		<div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
			{#each data.courts as court (court.id)}
				<Card class="overflow-hidden">
					<CardContent class="space-y-4">
						<div class="flex items-start justify-between gap-4">
							<div>
								<h3 class="text-lg font-semibold">{court.name}</h3>
								<p class="mt-1 capitalize text-sm text-muted-foreground">
									{court.sport}
								</p>
							</div>

							<Badge variant={court.is_active ? 'secondary' : 'outline'}>
								{court.is_active ? 'aktif' : 'nonaktif'}
							</Badge>
						</div>

						<div>
							<p class="text-lg font-semibold">
								{rupiah.format(court.price_per_slot)}
							</p>
							<p class="text-sm text-muted-foreground">
								per {court.slot_duration_minutes} menit
							</p>
						</div>

						<Button
							variant="outline"
							class="w-full"
							href="/admin/venues/{data.venue.slug}/operating-hours"
						>
							<CalendarCogIcon />
							Atur jam operasional
						</Button>
					</CardContent>
				</Card>
			{/each}
		</div>
	{/if}
</main>
