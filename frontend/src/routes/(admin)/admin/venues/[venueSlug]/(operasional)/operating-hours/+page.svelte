<script lang="ts">
	import { Card, CardContent, CardHeader, CardTitle } from '#lib/components/ui/card/index.js';
	import { Button } from '#lib/components/ui/button/index.js';
	import * as Field from '#lib/components/ui/field/index.js';
	import AdminHeroSection from '#lib/components/admin/admin-hero-section.svelte';

	let { data, form } = $props();

	const dayNames = ['Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu', 'Minggu'];

	const hours = $derived(form?.hours?.length === 7 ? form.hours : data.hours);
</script>

<svelte:head>
	<title>Jam Operasional | Lapanganku</title>
	<meta name="description" content="Atur jam operasional venue." />
</svelte:head>

<main class="mx-auto max-w-3xl space-y-5">
	<AdminHeroSection
		title="Jam Operasional"
		description={`Atur waktu buka dan tutup ${data.venue.name} untuk setiap hari.`}
		size="sm"
	/>

	<Card class="w-full overflow-hidden">
		<CardHeader>
			<CardTitle>Jadwal Mingguan</CardTitle>
		</CardHeader>
		<CardContent>
			<form method="POST" class="space-y-4">
				{#if form?.message}
					<p
						role="alert"
						class="rounded-xl border border-destructive/30 bg-destructive/10 px-4 py-3 text-sm text-destructive"
					>
						{form.message}
					</p>
				{/if}

				{#each hours as hour (hour.day_of_week)}
					<fieldset
						class="grid items-end gap-4 rounded-xl border p-4 sm:grid-cols-[1fr_auto_auto]"
					>
						<legend class="sr-only">{dayNames[hour.day_of_week - 1]}</legend>

						<div>
							<p class="font-medium">{dayNames[hour.day_of_week - 1]}</p>

							<label class="mt-2 flex items-center gap-2 text-sm text-muted-foreground">
								<input
									name={`closed_${hour.day_of_week}`}
									type="checkbox"
									checked={hour.is_closed}
									class="size-4 accent-primary"
								/>
								Tutup
							</label>
						</div>

						<Field.Field>
							<Field.Label for={`opens_at_${hour.day_of_week}`} class="text-xs">
								Buka
							</Field.Label>
							<input
								id={`opens_at_${hour.day_of_week}`}
								name={`opens_at_${hour.day_of_week}`}
								type="time"
								value={hour.opens_at ?? '08:00'}
								class="border-input bg-input/20 dark:bg-input/30 focus-visible:border-ring focus-visible:ring-ring/30 h-9 rounded-md border px-3 py-1 text-sm outline-none transition-colors focus-visible:ring-2"
							/>
						</Field.Field>

						<Field.Field>
							<Field.Label for={`closes_at_${hour.day_of_week}`} class="text-xs">
								Tutup
							</Field.Label>
							<input
								id={`closes_at_${hour.day_of_week}`}
								name={`closes_at_${hour.day_of_week}`}
								type="time"
								value={hour.closes_at ?? '22:00'}
								class="border-input bg-input/20 dark:bg-input/30 focus-visible:border-ring focus-visible:ring-ring/30 h-9 rounded-md border px-3 py-1 text-sm outline-none transition-colors focus-visible:ring-2"
							/>
						</Field.Field>
					</fieldset>
				{/each}

				<div class="flex justify-end gap-3 pt-2">
					<Button type="button" variant="outline" href="/admin/venues/{data.venue.slug}">
						Batal
					</Button>
					<Button type="submit">Simpan jadwal</Button>
				</div>
			</form>
		</CardContent>
	</Card>
</main>
