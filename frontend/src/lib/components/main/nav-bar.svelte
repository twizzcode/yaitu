<script lang="ts">
	import MenuIcon from '@lucide/svelte/icons/menu';
	import MapPinIcon from '@lucide/svelte/icons/map-pin';
	import { Button, buttonVariants } from '#lib/components/ui/button/index.js';
	import * as Sheet from '#lib/components/ui/sheet/index.js';
	import UserMenu from '#lib/components/user-menu.svelte';

	type User = {
		id: string;
		name: string;
		email: string;
	};

	let { user }: { user: User | null } = $props();

	const links = [
		{ href: '/', label: 'Beranda' },
		{ href: '/harga', label: 'Harga' },
		{ href: '/blog', label: 'Blog' },
		{ href: '/demo', label: 'Demo' },
		{ href: '/tentang-kami', label: 'Tentang Kami' }
	];
</script>

<nav class="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/80">
	<div class="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
		<a href="/" class="flex items-center gap-2 font-semibold" aria-label="LapanganKu beranda">
    		<img src="/favicon.svg" alt="logo" class="size-8" aria-hidden="true" />
			<span class="font-medium">Lapangan<span class="font-bold text-blue-700">Ku</span>.id</span>
		</a>

		<div class="hidden items-center gap-1 md:flex">
			{#each links as link}
				<Button href={link.href} size="lg" class="py-4 px-4 text-sm" variant="ghost">{link.label}</Button>
			{/each}
		</div>

		<div class="hidden items-center gap-2 md:flex">
			{#if user}
				<UserMenu {user} />
			{:else}
				<Button href="/login" variant="outline">Masuk</Button>
				<Button href="/register">Daftar</Button>
			{/if}
		</div>

		<Sheet.Root>
			<Sheet.Trigger
				class={`${buttonVariants({ variant: 'outline', size: 'icon-lg' })} md:hidden`}
				aria-label="Buka menu navigasi"
			>
				<MenuIcon aria-hidden="true" />
			</Sheet.Trigger>
			<Sheet.Content class="w-72" side="right">
				<Sheet.Header>
					<Sheet.Title class="flex items-center gap-2">
    					<img src="/favicon.svg" alt="logo" class="size-4" aria-hidden="true" />
                  		<span class="font-xs">Lapangan<span class="font-bold text-blue-700">Ku</span>.id</span>
					</Sheet.Title>
					<Sheet.Description>Menu navigasi utama</Sheet.Description>
				</Sheet.Header>

				<div class="flex flex-col gap-2 px-4">
					{#each links as link}
						<Button href={link.href} variant="ghost" class="justify-start">{link.label}</Button>
					{/each}
				</div>

				<Sheet.Footer class="flex items-center justify-end gap-2">
					{#if user}
						<span class="min-w-0 flex-1 truncate text-sm text-muted-foreground">{user.name}</span>
						<UserMenu {user} />
					{:else}
						<Button href="/login" variant="outline" class="flex-1">Masuk</Button>
						<Button href="/register" class="flex-1">Daftar</Button>
					{/if}
				</Sheet.Footer>
			</Sheet.Content>
		</Sheet.Root>
	</div>
</nav>
