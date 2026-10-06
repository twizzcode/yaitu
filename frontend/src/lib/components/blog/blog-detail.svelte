<script lang="ts">
	import ArrowLeftIcon from '@lucide/svelte/icons/arrow-left';
	import ArrowUpRightIcon from '@lucide/svelte/icons/arrow-up-right';
	import CalendarIcon from '@lucide/svelte/icons/calendar';
	import ClockIcon from '@lucide/svelte/icons/clock';
	import EyeIcon from '@lucide/svelte/icons/eye';
	import MapPinIcon from '@lucide/svelte/icons/map-pin';
	import QuoteIcon from '@lucide/svelte/icons/quote';
	import ShareIcon from '@lucide/svelte/icons/share-2';
	import type { BlogBlock, BlogPost, BlogVenue } from '#lib/blog.js';
	import { Badge } from '#lib/components/ui/badge/index.js';
	import { Button } from '#lib/components/ui/button/index.js';
	import { Card } from '#lib/components/ui/card/index.js';
	import { Separator } from '#lib/components/ui/separator/index.js';

	let {
		post,
		related = [],
		basePath = '',
		venue,
		showBooking = false
	}: {
		post: BlogPost;
		related?: BlogPost[];
		basePath?: string;
		venue?: BlogVenue | null;
		showBooking?: boolean;
	} = $props();

	async function shareArticle() {
		const url = window.location.href;
		if (navigator.share) {
			await navigator.share({ title: post.title, url });
			return;
		}
		await navigator.clipboard.writeText(url);
	}
</script>

{#snippet ArticleBody(blocks: BlogBlock[])}
	<div class="flex flex-col gap-6">
		{#each blocks as block, index (index)}
			{#if block.type === 'heading'}
				<h2 class="mt-4 text-2xl font-semibold tracking-tight">{block.text}</h2>
			{:else if block.type === 'list'}
				<ul class="flex flex-col gap-3">
					{#each block.items as item (item)}
						<li class="flex items-start gap-3 leading-7"><span class="mt-2 size-2 shrink-0 rounded-full bg-primary"></span><span class="text-muted-foreground">{item}</span></li>
					{/each}
				</ul>
			{:else if block.type === 'orderedList'}
				<ol class="flex flex-col gap-3">
					{#each block.items as item, itemIndex (item)}
						<li class="flex items-start gap-3 leading-7"><span class="mt-0.5 flex size-6 shrink-0 items-center justify-center rounded-full bg-primary/10 text-xs font-semibold text-primary">{itemIndex + 1}</span><span class="text-muted-foreground">{item}</span></li>
					{/each}
				</ol>
			{:else if block.type === 'quote'}
				<figure class="my-2 rounded-3xl border bg-muted/40 p-6 md:p-8">
					<QuoteIcon class="size-6 text-primary" />
					<blockquote class="mt-4 text-lg leading-8 font-medium">{block.text}</blockquote>
					{#if block.cite}<figcaption class="mt-4 text-sm text-muted-foreground">- {block.cite}</figcaption>{/if}
				</figure>
			{:else if block.type === 'image'}
				<img src={block.src} alt={block.alt ?? ''} loading="lazy" width="1200" height="675" class="aspect-video w-full rounded-2xl border object-cover" />
			{:else if block.type === 'divider'}
				<hr class="my-2 border-border" />
			{:else}
				<p class="leading-8 text-muted-foreground">{block.text}</p>
			{/if}
		{/each}
	</div>
{/snippet}

{#snippet RelatedSection()}
	{#if related.length > 0}
		<div class="flex flex-col gap-4">
			<div class="flex items-center gap-3"><h2 class="text-lg font-semibold">Artikel Terkait</h2><Separator class="flex-1" /></div>
			<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-1">
				{#each related as item (item.id)}
					<Card class="group/related h-full gap-0 p-0 shadow-xs">
						<a href={`${basePath}/blog/${item.slug}`} class="flex gap-3 p-3 lg:flex-col lg:gap-0 lg:p-0">
							<img src={item.image} alt={item.title} width="340" height="191" loading="lazy" class="aspect-video w-28 shrink-0 rounded-lg object-cover lg:w-full lg:rounded-b-none" />
							<div class="flex min-w-0 flex-1 flex-col gap-1.5 lg:p-4">
								<Badge variant="outline" class="w-fit">{item.category}</Badge>
								<h3 class="line-clamp-2 text-sm leading-snug font-semibold lg:text-base">{item.title}</h3>
								<p class="mt-auto flex items-center gap-1.5 text-xs text-muted-foreground"><ClockIcon class="size-3" /> {item.readTime} · {item.publishedAt}</p>
							</div>
						</a>
					</Card>
				{/each}
			</div>
		</div>
	{/if}
{/snippet}

<div class="pb-24">
	<div class="mx-auto w-full max-w-7xl px-4 pt-10 sm:px-8 md:pt-16">
		<div class="grid gap-12 lg:grid-cols-[minmax(0,1fr)_340px] lg:gap-14">
			<article class="mx-auto w-full min-w-0">
				<Button href={`${basePath}/blog`} variant="ghost" size="sm" class="-ml-2 w-fit text-muted-foreground"><ArrowLeftIcon /> Semua Artikel</Button>
				<h1 class="mt-4 text-3xl leading-tight font-bold tracking-tight md:text-5xl">{post.title}</h1>
				<div class="relative mt-8 aspect-video overflow-hidden rounded-3xl border shadow-sm">
					<img src={post.image} alt={post.title} width="1200" height="675" class="size-full object-cover" />
				</div>

				<div class="mt-8 flex flex-col gap-3 text-sm text-muted-foreground sm:flex-row sm:flex-wrap sm:items-center sm:gap-4">
					<div class="flex min-w-0 items-center gap-2.5">
						<img src={post.authorImg} alt="" width="36" height="36" class="size-9 rounded-full border object-cover" />
						<div><p class="font-medium text-foreground">{post.author}</p>{#if post.authorRole}<p class="text-xs">{post.authorRole}</p>{/if}</div>
					</div>
					<Separator orientation="vertical" class="hidden h-8 sm:block" />
					<div class="flex flex-wrap items-center gap-3">
						<span class="inline-flex items-center gap-1.5"><CalendarIcon class="size-3.5" /> {post.publishedAt}</span>
						<span class="inline-flex items-center gap-1.5"><ClockIcon class="size-3.5" /> {post.readTime} baca</span>
						<span class="inline-flex items-center gap-1.5"><EyeIcon class="size-3.5" /> {post.views ?? 0} dilihat</span>
					</div>
				</div>

				<Separator class="my-8" />
				{@render ArticleBody(post.content)}

				{#if venue}
					<a href={basePath || '/'} class="mt-12 block rounded-3xl border bg-muted/40 p-6 text-center">
						<div class="flex flex-col items-center gap-3">
							{#if venue.logoUrl}<img src={venue.logoUrl} alt={venue.name} width="56" height="56" class="size-14 rounded-full border object-cover" />{:else}<div class="flex size-14 items-center justify-center rounded-full border bg-background font-semibold">{venue.name.slice(0, 2).toUpperCase()}</div>{/if}
							<div><p class="text-xs uppercase tracking-widest text-muted-foreground">Ditulis di</p><p class="font-semibold">{venue.name}</p><p class="mt-1 inline-flex items-center gap-1.5 text-sm text-muted-foreground"><MapPinIcon class="size-3.5" /> {[venue.city, venue.province].filter(Boolean).join(', ')}</p></div>
						</div>
					</a>
				{/if}

				<div class="mt-8 flex flex-wrap items-center justify-between gap-3">
					<div class="flex items-center gap-2"><span class="text-sm text-muted-foreground">Topik:</span><Badge variant="outline">{post.category}</Badge></div>
					<Button variant="outline" size="sm" onclick={shareArticle}><ShareIcon /> Bagikan</Button>
				</div>
				<div class="mt-14 lg:hidden">{@render RelatedSection()}</div>
			</article>

			<aside class="hidden lg:block">
				<div class="sticky top-6 flex flex-col gap-6">
					{@render RelatedSection()}
					{#if showBooking}
						<Card class="gap-0 bg-muted/40 p-6"><h3 class="text-lg font-semibold">Siap booking lapangan?</h3><p class="mt-2 text-sm leading-6 text-muted-foreground">Cek jadwal dan pesan langsung tanpa antre.</p><Button href={`${basePath}#booking`} class="mt-5 w-fit">Booking Sekarang <ArrowUpRightIcon /></Button></Card>
					{/if}
				</div>
			</aside>
		</div>
	</div>
</div>
