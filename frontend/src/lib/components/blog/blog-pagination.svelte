<script lang="ts">
	import ChevronLeftIcon from '@lucide/svelte/icons/chevron-left';
	import ChevronRightIcon from '@lucide/svelte/icons/chevron-right';
	import { Button } from '#lib/components/ui/button/index.js';

	let {
		page,
		totalPages,
		onPageChange
	}: { page: number; totalPages: number; onPageChange: (page: number) => void } = $props();

	let pages = $derived(Array.from({ length: totalPages }, (_, index) => index + 1));
</script>

<nav aria-label="Paginasi blog" class="flex items-center justify-center gap-1">
	<Button variant="outline" size="icon-lg" disabled={page === 1} onclick={() => onPageChange(page - 1)} aria-label="Halaman sebelumnya">
		<ChevronLeftIcon aria-hidden="true" />
	</Button>
	{#each pages as item (item)}
		<Button variant={item === page ? 'default' : 'outline'} size="icon-lg" onclick={() => onPageChange(item)} aria-label={`Halaman ${item}`} aria-current={item === page ? 'page' : undefined}>
			{item}
		</Button>
	{/each}
	<Button variant="outline" size="icon-lg" disabled={page === totalPages} onclick={() => onPageChange(page + 1)} aria-label="Halaman berikutnya">
		<ChevronRightIcon aria-hidden="true" />
	</Button>
</nav>
