<script lang="ts">
	import { Badge } from '#lib/components/ui/badge/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '#lib/components/ui/card/index.js';

	type TopField = {
		field: string;
		bookings: number;
		revenue: string;
		utilization: string;
	};

	const sampleData: TopField[] = [
		{ field: 'Lapangan Futsal A', bookings: 96, revenue: 'Rp12,4 jt', utilization: '84%' },
		{ field: 'Lapangan Badminton 1', bookings: 82, revenue: 'Rp8,2 jt', utilization: '76%' },
		{ field: 'Lapangan Basket', bookings: 54, revenue: 'Rp6,8 jt', utilization: '61%' }
	];

	let { data = sampleData }: { data?: TopField[] } = $props();
</script>

<Card>
	<CardHeader>
		<CardTitle>Lapangan Terlaris</CardTitle>
		<p class="text-sm text-muted-foreground">Ranking berdasarkan pendapatan terbayar tahun ini.</p>
	</CardHeader>
	<CardContent>
		<div class="overflow-x-auto rounded-lg border">
			<table class="w-full text-sm">
				<thead class="border-b bg-muted/50 text-left">
					<tr>
						<th class="w-12 px-3 py-2 font-medium">#</th>
						<th class="px-3 py-2 font-medium">Lapangan</th>
						<th class="px-3 py-2 font-medium">Booking</th>
						<th class="px-3 py-2 font-medium">Pendapatan</th>
						<th class="px-3 py-2 text-right font-medium">Utilisasi</th>
					</tr>
				</thead>
				<tbody>
					{#if data.length === 0}
						<tr>
							<td colspan="5" class="px-3 py-10 text-center text-muted-foreground">
								Belum ada data booking terbayar.
							</td>
						</tr>
					{:else}
						{#each data as item, index (item.field)}
							<tr class="border-b last:border-b-0">
								<td class="px-3 py-3">{index + 1}</td>
								<td class="px-3 py-3 font-medium">{item.field}</td>
								<td class="px-3 py-3">{item.bookings.toLocaleString('id-ID')}</td>
								<td class="px-3 py-3">{item.revenue}</td>
								<td class="px-3 py-3 text-right">
									<Badge variant="secondary">{item.utilization}</Badge>
								</td>
							</tr>
						{/each}
					{/if}
				</tbody>
			</table>
		</div>
	</CardContent>
</Card>
