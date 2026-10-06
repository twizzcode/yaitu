<script lang="ts">
	import CalendarCheckIcon from '@lucide/svelte/icons/calendar-check';
	import CircleDollarSignIcon from '@lucide/svelte/icons/circle-dollar-sign';
	import ClockIcon from '@lucide/svelte/icons/clock';
	import UsersIcon from '@lucide/svelte/icons/users';
	import type { LucideIcon } from '@lucide/svelte';
	import { Badge } from '#lib/components/ui/badge/index.js';
	import { Card, CardContent } from '#lib/components/ui/card/index.js';

	type Stat = {
		title: string;
		value: string;
		change?: string;
		icon: LucideIcon;
		negative: boolean;
		note?: string;
	};

	const sampleStats: Stat[] = [
		{ title: 'Total booking', value: '248', change: '+12%', icon: CalendarCheckIcon, negative: false },
		{ title: 'Pendapatan', value: 'Rp24,8 jt', change: '+8%', icon: CircleDollarSignIcon, negative: false },
		{ title: 'Jam terisi', value: '186 jam', change: '-3%', icon: ClockIcon, negative: true },
		{ title: 'Pelanggan', value: '142', change: '+16%', icon: UsersIcon, negative: false }
	];

	let { stats = sampleStats }: { stats?: Stat[] } = $props();
</script>

<Card class="py-0">
	<CardContent class="grid gap-0 p-0 sm:grid-cols-2 xl:grid-cols-4">
		{#each stats as stat (stat.title)}
			{@const Icon = stat.icon}
			<div
				class="flex items-start justify-between gap-4 border-b p-4 sm:odd:border-r xl:border-r xl:border-b-0 xl:last:border-r-0"
			>
				<div>
					<p class="font-medium">{stat.title}</p>
					<p class="mt-5 text-2xl font-semibold tracking-tight">{stat.value}</p>
					<div class="mt-1 flex items-center gap-2">
						<p class="text-xs text-muted-foreground">{stat.note ?? '7 hari terakhir'}</p>
						{#if stat.change}
							<Badge variant={stat.negative ? 'destructive' : 'secondary'}>{stat.change}</Badge>
						{/if}
					</div>
				</div>
				<div class="rounded-full border p-3">
					<Icon class="size-5" aria-hidden="true" />
				</div>
			</div>
		{/each}
	</CardContent>
</Card>
