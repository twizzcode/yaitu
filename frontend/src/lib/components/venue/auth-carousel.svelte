<script lang="ts">
	import MapPinIcon from '@lucide/svelte/icons/map-pin';
	import { cn } from '#lib/utils.js';
	import type { Slide } from './types.js';

	let { slides }: { slides: Slide[] } = $props();

	let active = $state(0);

	$effect(() => {
		if (slides.length <= 1) return;

		const timer = setInterval(() => {
			active = (active + 1) % slides.length;
		}, 5000);

		return () => clearInterval(timer);
	});
</script>

<div class="bg-muted relative hidden overflow-hidden lg:block">
	{#each slides as slide, i (slide.src)}
		<div
			class={cn(
				'absolute inset-0 transition-opacity duration-1000 ease-in-out',
				i === active ? 'opacity-100' : 'opacity-0'
			)}
		>
			<img src={slide.src} alt={slide.title} class="h-full w-full object-cover" />
			<div
				class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/40 to-black/20"
			></div>
			<div class="absolute inset-x-0 bottom-0 p-10 text-white">
				<div class="flex items-center gap-2 text-sm text-white/80">
					<MapPinIcon class="size-4" />
					<span>Platform Sewa Lapangan #1</span>
				</div>
				<h2 class="mt-3 text-3xl font-bold tracking-tight">{slide.title}</h2>
				<p class="mt-2 max-w-md text-white/80">{slide.subtitle}</p>
			</div>
		</div>
	{/each}

	<!-- Dots -->
	<div class="absolute bottom-6 left-10 z-10 flex gap-2">
		{#each slides as slide, i (slide.src)}
			<button
				type="button"
				onclick={() => (active = i)}
				aria-label={`Slide ${i + 1}`}
				class={cn(
					'h-1.5 rounded-full bg-white/50 transition-all duration-300',
					i === active ? 'w-8 bg-white' : 'w-1.5 hover:bg-white/80'
				)}
			></button>
		{/each}
	</div>
</div>
