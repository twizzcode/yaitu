<script lang="ts">
	import CheckIcon from '@lucide/svelte/icons/check';
	import ChevronsUpDownIcon from '@lucide/svelte/icons/chevrons-up-down';
	import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';
	import * as Popover from '#lib/components/ui/popover/index.js';
	import * as Command from '#lib/components/ui/command/index.js';
	import { cn } from '#lib/utils.js';
	import type { Region } from './types.js';

	let {
		value = $bindable<Region | null>(null),
		onSelect,
		items,
		placeholder = 'Pilih',
		searchPlaceholder = 'Cari...',
		emptyText = 'Tidak ada data',
		disabled = false,
		loading = false,
		class: className
	}: {
		value?: Region | null;
		onSelect: (item: Region) => void;
		items: Region[];
		placeholder?: string;
		searchPlaceholder?: string;
		emptyText?: string;
		disabled?: boolean;
		loading?: boolean;
		class?: string;
	} = $props();

	let open = $state(false);
</script>

<Popover.Root bind:open>
	<Popover.Trigger
		{disabled}
		class={cn(
			'border-input bg-input/20 dark:bg-input/30 dark:hover:bg-input/50 focus-visible:border-ring focus-visible:ring-ring/30 flex h-7 w-full items-center justify-between gap-1.5 rounded-md border px-2 py-1.5 text-xs/relaxed whitespace-nowrap outline-none transition-colors focus-visible:ring-2 disabled:cursor-not-allowed disabled:opacity-50',
			className
		)}
	>
		<span class={cn('truncate', !value && 'text-muted-foreground')}>
			{loading ? 'Memuat...' : (value?.name ?? placeholder)}
		</span>
		{#if loading}
			<LoaderCircleIcon class="text-muted-foreground size-3.5 shrink-0 animate-spin" />
		{:else}
			<ChevronsUpDownIcon class="text-muted-foreground size-3.5 shrink-0" />
		{/if}
	</Popover.Trigger>

	<Popover.Content
		class="w-(--bits-popover-anchor-width) min-w-(--bits-popover-anchor-width) p-0"
		align="start"
	>
		<Command.Root>
			<Command.Input placeholder={searchPlaceholder} />

			<Command.List class="max-h-64">
				<Command.Empty>{emptyText}</Command.Empty>

				<Command.Group>
					{#each items as item (item.code)}
						<Command.Item
							value={`${item.name} ${item.code}`}
							onSelect={() => {
								onSelect(item);
								open = false;
							}}
						>
							<span class="truncate">{item.name}</span>
							{#if value?.code === item.code}
								<CheckIcon class="ms-auto" />
							{/if}
						</Command.Item>
					{/each}
				</Command.Group>
			</Command.List>
		</Command.Root>
	</Popover.Content>
</Popover.Root>
