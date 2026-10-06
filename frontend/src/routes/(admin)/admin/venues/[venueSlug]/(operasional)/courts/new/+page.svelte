<script lang="ts">
	import { Card, CardContent, CardHeader, CardTitle } from '#lib/components/ui/card/index.js';
	import { Button } from '#lib/components/ui/button/index.js';
	import * as Field from '#lib/components/ui/field/index.js';
	import { Input } from '#lib/components/ui/input/index.js';
	import AdminHeroSection from '#lib/components/admin/admin-hero-section.svelte';

	let { data, form } = $props();
</script>

<svelte:head>
	<title>Tambah Lapangan | Lapanganku</title>
	<meta name="description" content="Tambah lapangan baru." />
</svelte:head>

<main class="mx-auto max-w-3xl space-y-5">
	<AdminHeroSection
		title="Tambah Lapangan"
		description={`Tambahkan lapangan baru untuk ${data.venue.name}.`}
		size="sm"
	/>

	<Card class="w-full overflow-hidden">
		<CardHeader>
			<CardTitle>Informasi Lapangan</CardTitle>
		</CardHeader>
		<CardContent>
			<form method="POST" class="space-y-5">
				{#if form?.message}
					<p
						role="alert"
						class="rounded-xl border border-destructive/30 bg-destructive/10 px-4 py-3 text-sm text-destructive"
					>
						{form.message}
					</p>
				{/if}

				<Field.Field>
					<Field.Label for="name">Nama lapangan *</Field.Label>
					<Input
						id="name"
						name="name"
						required
						minlength={2}
						maxlength={100}
						value={form?.name ?? ''}
						placeholder="Lapangan Futsal A"
					/>
				</Field.Field>

				<Field.Field>
					<Field.Label for="sport">Jenis olahraga *</Field.Label>
					<Input
						id="sport"
						name="sport"
						required
						minlength={2}
						maxlength={50}
						value={form?.sport ?? ''}
						placeholder="futsal"
					/>
				</Field.Field>

				<Field.Field>
					<Field.Label for="price_per_slot">Harga per slot *</Field.Label>
					<Input
						id="price_per_slot"
						name="price_per_slot"
						type="number"
						required
						min="0"
						step="1"
						value={form?.pricePerSlot ?? 0}
						placeholder="150000"
					/>
					<Field.Description>
						Masukkan nominal rupiah tanpa titik atau koma.
					</Field.Description>
				</Field.Field>

				<Field.Field>
					<Field.Label for="slot_duration_minutes">Durasi slot</Field.Label>
					<select
						id="slot_duration_minutes"
						name="slot_duration_minutes"
						required
						class="border-input bg-input/20 dark:bg-input/30 focus-visible:border-ring focus-visible:ring-ring/30 h-9 w-full rounded-md border px-3 py-1 text-sm outline-none transition-colors focus-visible:ring-2"
					>
						<option value="30" selected={form?.slotDurationMinutes === 30}>30 menit</option>
						<option value="60" selected={(form?.slotDurationMinutes ?? 60) === 60}>
							60 menit
						</option>
						<option value="90" selected={form?.slotDurationMinutes === 90}>90 menit</option>
						<option value="120" selected={form?.slotDurationMinutes === 120}>
							120 menit
						</option>
					</select>
				</Field.Field>

				<div class="flex justify-end gap-3 pt-2">
					<Button
						type="button"
						variant="outline"
						href="/admin/venues/{data.venue.slug}/courts"
					>
						Batal
					</Button>
					<Button type="submit">Tambah lapangan</Button>
				</div>
			</form>
		</CardContent>
	</Card>
</main>
