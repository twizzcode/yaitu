<script lang="ts">
	import CheckIcon from '@lucide/svelte/icons/check';
	import ChevronDownIcon from '@lucide/svelte/icons/chevron-down';
	import SearchIcon from '@lucide/svelte/icons/search';
	import type { BlogPost } from '#lib/blog.js';
	import BadgeMain from '#lib/components/main/badge-main.svelte';
	import { Button } from '#lib/components/ui/button/index.js';
	import * as DropdownMenu from '#lib/components/ui/dropdown-menu/index.js';
	import { Input } from '#lib/components/ui/input/index.js';
	import BlogCard from './blog-card.svelte';
	import BlogPagination from './blog-pagination.svelte';

	const postsPerPage = 6;

	let { posts, basePath = '' }: { posts: BlogPost[]; basePath?: string } = $props();
	let category = $state('Semua');
	let query = $state('');
	let page = $state(1);
	let resultsSection: HTMLElement;

	let availableCategories = $derived([
		'Semua',
		...new Set(posts.map((post) => post.category))
	]);
	let filtered = $derived.by(() => {
		const normalizedQuery = query.trim().toLowerCase();
		return posts.filter((post) => {
			const matchesCategory = category === 'Semua' || post.category === category;
			const matchesQuery = !normalizedQuery || [post.title, post.excerpt, post.author]
				.some((value) => value.toLowerCase().includes(normalizedQuery));
			return matchesCategory && matchesQuery;
		});
	});
	let totalPages = $derived(Math.max(1, Math.ceil(filtered.length / postsPerPage)));
	let currentPage = $derived(Math.min(page, totalPages));
	let pagePosts = $derived(filtered.slice((currentPage - 1) * postsPerPage, currentPage * postsPerPage));

	function selectCategory(nextCategory: string) {
		category = nextCategory;
		page = 1;
	}

	function handleSearch(event: Event) {
		query = (event.currentTarget as HTMLInputElement).value;
		page = 1;
	}

	function changePage(nextPage: number) {
		page = Math.max(1, Math.min(totalPages, nextPage));
		resultsSection?.scrollIntoView({ behavior: 'smooth', block: 'start' });
	}
</script>

<div class="flex flex-col gap-10 pb-24">
	<section class="px-4 pt-12 md:pt-20">
		<div class="mx-auto flex max-w-4xl flex-col items-center text-center">
			<BadgeMain><span>Blog LapanganKu.id</span></BadgeMain>
			<h1 class="mt-7 text-4xl font-bold tracking-tight md:text-5xl">Insight untuk pengelola lapangan.</h1>
			<p class="mt-6 max-w-2xl text-base leading-7 text-muted-foreground md:text-lg">
				Panduan praktis, tips operasional, dan cerita nyata dari pengelola venue yang merapikan bisnisnya bersama LapanganKu.id.
			</p>
		</div>
	</section>

	<section class="mx-auto w-full max-w-7xl px-4 sm:px-8">
		<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
			<DropdownMenu.Root>
				<DropdownMenu.Trigger>
					{#snippet child({ props })}
						<Button {...props} variant="outline" class="w-full justify-between font-normal sm:w-52">
							<span class="inline-flex items-center gap-2">
								<span class="text-muted-foreground">Kategori:</span>
								<span class="font-medium">{category}</span>
							</span>
							<ChevronDownIcon class="size-4 text-muted-foreground" />
						</Button>
					{/snippet}
				</DropdownMenu.Trigger>
				<DropdownMenu.Content align="start" class="w-52">
					{#each availableCategories as item (item)}
						<DropdownMenu.Item onSelect={() => selectCategory(item)} class="justify-between">
							{item}
							{#if category === item}<CheckIcon class="size-4" />{/if}
						</DropdownMenu.Item>
					{/each}
				</DropdownMenu.Content>
			</DropdownMenu.Root>

			<div class="relative w-full sm:w-72">
				<SearchIcon class="absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground" aria-hidden="true" />
				<Input value={query} oninput={handleSearch} placeholder="Cari artikel..." aria-label="Cari artikel" class="h-9 pl-9" />
			</div>
		</div>
	</section>

	<section bind:this={resultsSection} class="mx-auto w-full max-w-7xl scroll-mt-6 px-4 sm:px-8">
		{#if pagePosts.length === 0}
			<div class="flex flex-col items-center justify-center gap-3 rounded-3xl border border-dashed bg-muted/30 px-6 py-20 text-center">
				<SearchIcon class="size-8 text-muted-foreground" />
				<h2 class="text-lg font-semibold">Artikel tidak ditemukan</h2>
				<p class="max-w-sm text-sm text-muted-foreground">Coba ubah kata kunci pencarian atau pilih kategori lain.</p>
			</div>
		{:else}
			<div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
				{#each pagePosts as post (post.id)}
					<BlogCard {post} {basePath} />
				{/each}
			</div>
		{/if}
	</section>

	{#if totalPages > 1}
		<section class="mx-auto w-full max-w-7xl px-4 sm:px-8">
			<BlogPagination page={currentPage} {totalPages} onPageChange={changePage} />
		</section>
	{/if}
</div>
