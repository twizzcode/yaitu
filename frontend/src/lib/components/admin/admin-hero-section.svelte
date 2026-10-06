<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Card, CardContent } from '#lib/components/ui/card/index.js';

	type HeroTone = 'default' | 'slate';
	type Content = string | Snippet;

	type Props = {
		title: Content;
		description?: Content;
		children?: Snippet;
		class?: string;
		size?: 'default' | 'sm';
		tone?: HeroTone;
	};

	const tones = {
		default: {
			bottom: 'bg-emerald-200/40',
			topLeft: 'bg-yellow-300/40',
			topRight: 'bg-sky-300/25',
			overlay: 'from-primary/5 via-transparent to-sky-200/10'
		},
		slate: {
			bottom: 'bg-slate-300/40',
			topLeft: 'bg-slate-400/30',
			topRight: 'bg-zinc-300/25',
			overlay: 'from-slate-500/5 via-transparent to-zinc-300/10'
		}
	} satisfies Record<HeroTone, Record<string, string>>;

	let {
		title,
		description,
		children,
		class: className,
		size = 'default',
		tone = 'default'
	}: Props = $props();

	let isSmall = $derived(size === 'sm');
	let colors = $derived(tones[tone]);
</script>

<Card class="relative w-full overflow-hidden {className ?? ''}">
	<div
		aria-hidden="true"
		class="pointer-events-none absolute inset-x-0 bottom-0 mx-auto rounded-full blur-3xl {isSmall
			? 'h-24 max-w-xl'
			: 'h-44 max-w-3xl'} {colors.bottom}"
	></div>
	<div
		aria-hidden="true"
		class="pointer-events-none absolute -top-24 -left-24 rounded-full blur-3xl {isSmall
			? 'size-40'
			: 'size-64'} {colors.topLeft}"
	></div>
	<div
		aria-hidden="true"
		class="pointer-events-none absolute -top-16 rounded-full blur-3xl {isSmall
			? '-right-16 size-36'
			: '-right-24 size-56'} {colors.topRight}"
	></div>
	<div
		aria-hidden="true"
		class="pointer-events-none absolute inset-0 bg-gradient-to-br {colors.overlay}"
	></div>

	<CardContent
		class="relative flex flex-col items-center justify-center px-6 text-center {isSmall
			? 'py-6'
			: 'py-10'}"
	>
		<h1
			class="text-balance font-semibold tracking-tight text-foreground {isSmall
				? 'max-w-3xl text-2xl md:text-3xl'
				: 'max-w-4xl text-4xl md:text-5xl'}"
		>
			{#if typeof title === 'function'}
				{@render title()}
			{:else}
				{title}
			{/if}
		</h1>

		{#if description}
			<p class="mt-2 max-w-xl text-sm text-muted-foreground {isSmall ? '' : 'mt-3'}">
				{#if typeof description === 'function'}
					{@render description()}
				{:else}
					{description}
				{/if}
			</p>
		{/if}

		{#if children}
			<div class={isSmall ? 'mt-4' : 'mt-6'}>{@render children()}</div>
		{/if}
	</CardContent>
</Card>
