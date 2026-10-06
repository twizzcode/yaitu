<script lang="ts">
	import { Badge } from '#lib/components/ui/badge/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '#lib/components/ui/card/index.js';

	type FieldRevenueItem = {
		name: string;
		value: number;
		change?: string;
		fill?: string;
	};

	type Row = Required<FieldRevenueItem> & { start: number; end: number };

	const colors = [
		'var(--chart-2)',
		'var(--chart-3)',
		'var(--chart-4)',
		'var(--chart-5)',
		'var(--chart-1)'
	];
	const sampleData: FieldRevenueItem[] = [
		{ name: 'Futsal A', value: 12_400_000 },
		{ name: 'Badminton 1', value: 8_200_000 },
		{ name: 'Basket', value: 6_800_000 }
	];

	let { data = sampleData }: { data?: FieldRevenueItem[] } = $props();

	let total = $derived(data.reduce((sum, item) => sum + item.value, 0));
	let rows: Row[] = $derived.by(() => {
		let cursor = 0;
		return data.map((item, index) => {
			const start = cursor;
			cursor += total > 0 ? (item.value / total) * 100 : 0;
			return {
				...item,
				change: item.change ?? '',
				fill: item.fill ?? colors[index % colors.length],
				start,
				end: cursor
			};
		});
	});
	let gradient = $derived(
		rows.map((row) => `${row.fill} ${row.start}% ${row.end}%`).join(', ')
	);

	function formatRupiahShort(value: number) {
		if (value >= 1_000_000_000) return `Rp${(value / 1_000_000_000).toLocaleString('id-ID')} M`;
		if (value >= 1_000_000) return `Rp${(value / 1_000_000).toLocaleString('id-ID')} jt`;
		return `Rp${value.toLocaleString('id-ID')}`;
	}
</script>

<Card>
	<CardHeader>
		<CardTitle>Pendapatan Bulan Ini</CardTitle>
	</CardHeader>
	<CardContent>
		{#if rows.length === 0 || total === 0}
			<div class="flex h-72 items-center justify-center rounded-lg border border-dashed">
				<p class="text-sm text-muted-foreground">Belum ada pendapatan bulan ini.</p>
			</div>
		{:else}
			<div
				class="relative mx-auto aspect-square max-h-56 rounded-full"
				style:background={`conic-gradient(${gradient})`}
				aria-label={`Total pendapatan ${formatRupiahShort(total)}`}
			>
				<div class="absolute inset-[26%] flex flex-col items-center justify-center rounded-full bg-card">
					<span class="text-sm text-muted-foreground">Total</span>
					<strong class="text-xl font-semibold">{formatRupiahShort(total)}</strong>
				</div>
			</div>
			<div class="mt-4 flex flex-col gap-4">
				{#each rows as field (field.name)}
					{@const percent = Math.round((field.value / total) * 100)}
					<div class="flex items-center justify-between gap-4 text-sm">
						<div class="flex items-center gap-2 font-medium">
							<span class="h-4 w-0.5" style:background-color={field.fill}></span>{field.name}
						</div>
						<div class="flex items-center gap-2">
							<span class="font-semibold">{formatRupiahShort(field.value)}</span>
							<Badge variant="secondary">{percent}%</Badge>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</CardContent>
</Card>
