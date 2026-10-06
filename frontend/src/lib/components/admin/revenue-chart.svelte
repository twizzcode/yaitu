<script lang="ts">
	import { Badge } from '#lib/components/ui/badge/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '#lib/components/ui/card/index.js';

	type RevenueItem = {
		month: string;
		income: number;
		expense: number;
	};

	const sampleData: RevenueItem[] = [
		{ month: 'Jan', income: 9_000_000, expense: 4_000_000 },
		{ month: 'Feb', income: 12_000_000, expense: 5_500_000 },
		{ month: 'Mar', income: 10_500_000, expense: 4_800_000 },
		{ month: 'Apr', income: 15_000_000, expense: 6_000_000 },
		{ month: 'Mei', income: 13_500_000, expense: 5_200_000 },
		{ month: 'Jun', income: 18_000_000, expense: 7_000_000 }
	];

	let { data = sampleData }: { data?: RevenueItem[] } = $props();

	let totalIncome = $derived(data.reduce((sum, item) => sum + item.income, 0));
	let totalExpense = $derived(data.reduce((sum, item) => sum + item.expense, 0));
	let hasData = $derived(totalIncome > 0 || totalExpense > 0);
	let maxValue = $derived(Math.max(...data.flatMap((item) => [item.income, item.expense]), 1));

	function formatRupiahShort(value: number) {
		if (value >= 1_000_000_000) return `Rp${(value / 1_000_000_000).toLocaleString('id-ID')} M`;
		if (value >= 1_000_000) return `Rp${(value / 1_000_000).toLocaleString('id-ID')} jt`;
		return `Rp${value.toLocaleString('id-ID')}`;
	}
</script>

<Card>
	<CardHeader class="flex grid-cols-none flex-row items-start justify-between gap-4">
		<div>
			<CardTitle>Pendapatan Tahun Ini</CardTitle>
			<div class="mt-2 flex items-center gap-2">
				<p class="text-2xl font-semibold tracking-tight">{formatRupiahShort(totalIncome)}</p>
				<span class="text-xs text-muted-foreground">
					{hasData ? 'pendapatan terbayar' : 'belum ada data'}
				</span>
			</div>
		</div>
		<Badge variant="secondary">Pengeluaran {formatRupiahShort(totalExpense)}</Badge>
	</CardHeader>
	<CardContent>
		{#if hasData}
			<div class="flex h-72 flex-col" aria-label="Grafik pendapatan dan pengeluaran per bulan">
				<div class="mb-4 flex justify-end gap-4 text-xs text-muted-foreground">
					<span class="flex items-center gap-1.5"><i class="size-2.5 rounded-sm bg-primary"></i>Pendapatan</span>
					<span class="flex items-center gap-1.5"><i class="size-2.5 rounded-sm bg-muted-foreground/50"></i>Pengeluaran</span>
				</div>
				<div class="grid flex-1 grid-cols-[repeat(var(--items),minmax(0,1fr))] gap-3 border-b" style:--items={data.length}>
					{#each data as item (item.month)}
						<div class="flex min-w-0 items-end justify-center gap-1">
							<div
								class="w-4 rounded-t bg-primary sm:w-6"
								style:height={`${(item.income / maxValue) * 100}%`}
								title={`Pendapatan ${item.month}: ${formatRupiahShort(item.income)}`}
							></div>
							<div
								class="w-4 rounded-t bg-muted-foreground/50 sm:w-6"
								style:height={`${(item.expense / maxValue) * 100}%`}
								title={`Pengeluaran ${item.month}: ${formatRupiahShort(item.expense)}`}
							></div>
						</div>
					{/each}
				</div>
				<div class="grid grid-cols-[repeat(var(--items),minmax(0,1fr))] gap-3 pt-2 text-center text-xs text-muted-foreground" style:--items={data.length}>
					{#each data as item (item.month)}<span>{item.month}</span>{/each}
				</div>
			</div>
		{:else}
			<div class="flex h-72 items-center justify-center rounded-lg border border-dashed">
				<p class="text-sm text-muted-foreground">Belum ada pendapatan tercatat tahun ini.</p>
			</div>
		{/if}
	</CardContent>
</Card>
