<script lang="ts">
	import MousePointerClickIcon from '@lucide/svelte/icons/mouse-pointer-click';
	import UsersIcon from '@lucide/svelte/icons/users';
	import { Badge } from '#lib/components/ui/badge/index.js';
	import { Card, CardContent } from '#lib/components/ui/card/index.js';

	type TrafficStat = {
		title: string;
		value: string;
		change: string;
		negative: boolean;
	};

	const sampleStats: TrafficStat[] = [
		{ title: 'Tayangan halaman', value: '12.840', change: '+12,5%', negative: false },
		{ title: 'Pengunjung', value: '4.280', change: '+8,2%', negative: false }
	];

	let { stats = sampleStats, live = false }: { stats?: TrafficStat[]; live?: boolean } = $props();
</script>

<Card class="py-0">
	<CardContent class="grid gap-0 p-0 sm:grid-cols-2">
		{#each stats as stat, index (stat.title)}
			{@const Icon = index === 0 ? MousePointerClickIcon : UsersIcon}
			<div
				class="flex items-start justify-between gap-4 border-b p-4 sm:border-b-0 sm:[&:first-child]:border-r"
			>
				<div>
					<p class="font-medium">{stat.title}</p>
					<p class="mt-5 text-2xl font-semibold tracking-tight">{stat.value}</p>
					<div class="mt-1 flex items-center gap-2">
						<Badge variant={live ? 'secondary' : 'outline'}>{live ? 'Umami' : 'kosong'}</Badge>
						<Badge variant={stat.negative ? 'destructive' : 'secondary'}>{stat.change}</Badge>
					</div>
				</div>
				<div class="rounded-full border p-3">
					<Icon class="size-5" aria-hidden="true" />
				</div>
			</div>
		{/each}
	</CardContent>
</Card>
