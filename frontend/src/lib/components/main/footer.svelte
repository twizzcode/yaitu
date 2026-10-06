<script lang="ts">
	import MapPinIcon from '@lucide/svelte/icons/map-pin';

	const sections = [
		{
			title: 'Produk',
			links: [
				{ name: 'Home', href: '/' },
				{ name: 'Harga', href: '/harga' },
				{ name: 'Daftar', href: '/register' }
			]
		},
		{
			title: 'Perusahaan',
			links: [
				{ name: 'Tentang Kami', href: '/tentang-kami' },
				{ name: 'Blog', href: '/blog' }
			]
		},
		{
			title: 'Bantuan',
			links: [
				{ name: 'Blog', href: '/blog' },
				{ name: 'FAQ', href: '/#faq' }
			]
		}
	];

	const legalLinks = [
		{ name: 'Syarat dan Ketentuan', href: '/syarat-dan-ketentuan' },
		{ name: 'Kebijakan Privasi', href: '/kebijakan-privasi' }
	];
</script>

<footer class="bg-muted px-4 py-16 sm:px-6 lg:px-8">
	<div class="mx-auto max-w-3xl text-center">
		<h2 class="text-3xl font-bold tracking-tight drop-shadow-md md:text-4xl lg:text-5xl">
			Siap Mengelola Lapangan<br class="hidden md:block" />
			dengan Lebih Mudah?
		</h2>
		<p class="mt-6 text-base text-muted-foreground drop-shadow-sm md:text-lg">
			Mulai digitalisasi operasional lapangan Anda dan<br class="hidden md:block" />
			berikan pengalaman booking yang lebih baik untuk pelanggan.
		</p>
	</div>

	<div class="mx-auto mt-16 w-full max-w-7xl rounded-xl bg-background p-6 sm:p-10">
		<div class="grid gap-10 sm:grid-cols-2 lg:grid-cols-5">
			<div class="sm:col-span-2">
				<a href="/" class="inline-flex items-center gap-2" aria-label="LapanganKu.id beranda">
					<img src="/favicon.svg" alt="logo" class="size-8" aria-hidden="true" />
					<span class="font-medium">Lapangan<span class="font-bold text-blue-700">Ku</span>.id</span>
				</a>
				<p class="mt-4 max-w-md text-sm font-medium leading-relaxed text-muted-foreground">
					Platform manajemen bisnis lapangan untuk membantu pengelola menyederhanakan booking,
					operasional, pelanggan, dan keuangan dalam satu tempat.
				</p>
			</div>

			{#each sections as section (section.title)}
				<nav aria-label={section.title}>
					<h3 class="text-sm font-semibold">{section.title}</h3>
					<ul class="mt-4 space-y-3 text-sm text-muted-foreground">
						{#each section.links as link (link.href)}
							<li>
								<a href={link.href} class="transition-colors hover:text-foreground">{link.name}</a>
							</li>
						{/each}
					</ul>
				</nav>
			{/each}
		</div>

		<div class="mt-8 flex flex-col justify-between gap-4 border-t border-border pt-8 text-xs font-medium text-muted-foreground md:flex-row md:items-center">
			<p>&copy; 2026 LapanganKu.id. All rights reserved.</p>
			<ul class="flex flex-wrap gap-x-4 gap-y-2">
				{#each legalLinks as link (link.href)}
					<li>
						<a href={link.href} class="underline transition-colors hover:text-primary">{link.name}</a>
					</li>
				{/each}
			</ul>
		</div>
	</div>
</footer>
