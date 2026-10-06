<script lang="ts">
	import Footer from '#lib/components/main/footer.svelte';
	import NavBar from '#lib/components/main/nav-bar.svelte';
	import { ScrollArea } from '#lib/components/ui/scroll-area/index.js';

	let { data, children } = $props();
</script>

<div class="flex h-dvh flex-col overflow-hidden">
	<header class="shrink-0">
		<NavBar user={data.user} />
	</header>

	<ScrollArea class="min-h-0 flex-1">
		{@render children()}
		<Footer />
	</ScrollArea>
</div>
