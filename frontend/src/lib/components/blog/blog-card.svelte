<script lang="ts">
	import ArrowUpRightIcon from '@lucide/svelte/icons/arrow-up-right';
	import ClockIcon from '@lucide/svelte/icons/clock';
	import UserIcon from '@lucide/svelte/icons/user';
	import type { BlogPost } from '#lib/blog.js';
	import { Badge } from '#lib/components/ui/badge/index.js';
	import { Card } from '#lib/components/ui/card/index.js';

	let { post, basePath = '' }: { post: BlogPost; basePath?: string } = $props();
</script>

<Card class="group/card h-full gap-0 p-0 shadow-xs transition duration-300 hover:-translate-y-1 hover:shadow-xl">
	<a href={`${basePath}/blog/${post.slug}`} class="flex h-full flex-col">
		<div class="relative aspect-video overflow-hidden">
			<img
				src={post.image}
				alt={post.title}
				width="800"
				height="450"
				loading="lazy"
				class="size-full object-cover transition duration-500 group-hover/card:scale-105"
			/>
			<Badge class="absolute top-3 left-3" variant="secondary">{post.category}</Badge>
		</div>

		<div class="flex flex-1 flex-col gap-3 p-5">
			<h2 class="text-lg leading-snug font-semibold tracking-tight transition group-hover/card:text-primary">
				{post.title}
			</h2>
			<p class="line-clamp-2 text-sm leading-6 text-muted-foreground">{post.excerpt}</p>
			<div class="mt-auto flex items-center gap-2.5 pt-2">
				<img src={post.authorImg} alt="" width="32" height="32" loading="lazy" class="size-8 rounded-full border object-cover" />
				<div class="min-w-0">
					<p class="flex items-center gap-1.5 truncate text-sm font-medium">
						<UserIcon class="size-3.5 text-muted-foreground" aria-hidden="true" /> {post.author}
					</p>
					<p class="flex items-center gap-1.5 text-xs text-muted-foreground">
						<ClockIcon class="size-3" aria-hidden="true" /> {post.readTime} · {post.publishedAt}
					</p>
				</div>
				<span class="ml-auto flex size-8 shrink-0 items-center justify-center rounded-full border text-primary transition group-hover/card:bg-primary group-hover/card:text-primary-foreground">
					<ArrowUpRightIcon class="size-4" aria-hidden="true" />
				</span>
			</div>
		</div>
	</a>
</Card>
