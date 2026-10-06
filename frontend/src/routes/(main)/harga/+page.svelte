<script lang="ts">
	import ArrowUpRightIcon from '@lucide/svelte/icons/arrow-up-right';
	import CheckIcon from '@lucide/svelte/icons/check';
	import FlameIcon from '@lucide/svelte/icons/flame';
	import MinusIcon from '@lucide/svelte/icons/minus';
	import SparklesIcon from '@lucide/svelte/icons/sparkles';
	import XIcon from '@lucide/svelte/icons/x';
	import BadgeMain from '#lib/components/main/badge-main.svelte';

	type Plan = {
		id: string;
		name: string;
		slug: string;
		priceMonthly: number;
		priceYearly: number;
		features: string[];
		maxCourts: number;
		maxStaff: number;
		customDomainEnabled: boolean;
		blogEnabled: boolean;
		galleryEnabled: boolean;
		storageLimitBytes: number;
		advancedAnalyticsEnabled: boolean;
		reportExportEnabled: boolean;
		whatsappEnabled: boolean;
		removeBrandingEnabled: boolean;
		supportPriority: 'standard' | 'fast' | 'priority';
		maxGalleryImages: number;
		popular?: boolean;
	};

	const plans: Plan[] = [
		{
			id: 'plan-basic',
			name: 'Basic',
			slug: 'basic',
			priceMonthly: 189000,
			priceYearly: 1814400,
			features: ['5 Lapangan', 'Booking online', 'Website dasar'],
			maxCourts: 5,
			maxStaff: 0,
			customDomainEnabled: false,
			blogEnabled: false,
			galleryEnabled: false,
			storageLimitBytes: 0,
			advancedAnalyticsEnabled: false,
			reportExportEnabled: false,
			whatsappEnabled: false,
			removeBrandingEnabled: false,
			supportPriority: 'standard',
			maxGalleryImages: 0
		},
		{
			id: 'plan-pro',
			name: 'Pro',
			slug: 'pro',
			priceMonthly: 359000,
			priceYearly: 3446400,
			features: ['15 Lapangan', '2 Staff', 'Analytics', 'Custom domain', 'Blog', 'Galeri', 'Export laporan', 'Notifikasi WhatsApp'],
			maxCourts: 15,
			maxStaff: 2,
			customDomainEnabled: true,
			blogEnabled: true,
			galleryEnabled: true,
			storageLimitBytes: 104857600,
			advancedAnalyticsEnabled: true,
			reportExportEnabled: true,
			whatsappEnabled: true,
			removeBrandingEnabled: false,
			supportPriority: 'fast',
			maxGalleryImages: 100,
			popular: true
		},
		{
			id: 'plan-business',
			name: 'Business',
			slug: 'business',
			priceMonthly: 579000,
			priceYearly: 5558400,
			features: ['25 Lapangan', '5 Staff', 'Analytics lanjutan', 'Custom domain', 'Blog', 'Galeri', 'Export laporan', 'Notifikasi WhatsApp', 'Tanpa branding'],
			maxCourts: 25,
			maxStaff: 5,
			customDomainEnabled: true,
			blogEnabled: true,
			galleryEnabled: true,
			storageLimitBytes: 1073741824,
			advancedAnalyticsEnabled: true,
			reportExportEnabled: true,
			whatsappEnabled: true,
			removeBrandingEnabled: true,
			supportPriority: 'priority',
			maxGalleryImages: 1000
		}
	];

	const descriptions: Record<string, string> = {
		basic: 'Untuk venue yang mulai merapikan operasional dan booking online.',
		pro: 'Untuk venue berkembang yang membutuhkan tim, promosi, dan analitik.',
		business: 'Untuk operasional besar dengan kapasitas dan kontrol maksimal.'
	};

	const supportLabels = {
		standard: 'Standar',
		fast: 'Cepat',
		priority: 'Prioritas'
	};

	let billing = $state<'monthly' | 'yearly'>('monthly');

	function formatPrice(price: number) {
		return new Intl.NumberFormat('id-ID', { maximumFractionDigits: 0 }).format(price);
	}

	function formatStorage(bytes: number) {
		if (!bytes) return '';
		if (bytes >= 1024 ** 3) return `${bytes / 1024 ** 3} GB`;
		return `${bytes / 1024 ** 2} MB`;
	}
</script>

<svelte:head>
	<title>Harga | LapanganKu.id</title>
	<meta name="description" content="Pilih paket LapanganKu.id sesuai kebutuhan operasional venue Anda." />
</svelte:head>

<section class="px-4 py-16 md:py-24">
	<div class="mx-auto flex max-w-3xl flex-col items-center text-center">
		<BadgeMain><span>Paket Harga</span></BadgeMain>
		<h1 class="mt-5 text-4xl font-bold tracking-tight md:text-5xl">
			Harga fleksibel untuk setiap lapangan
		</h1>
		<p class="mt-5 text-base leading-7 text-muted-foreground md:text-lg">
			Pilih paket sesuai kebutuhan operasional. Mulai dari booking online, pengelolaan jadwal,
			sampai laporan bisnis dalam satu dashboard.
		</p>

		<div class="mt-8 inline-flex rounded-full border bg-muted p-1">
			<button
				type="button"
				onclick={() => (billing = 'monthly')}
				class="rounded-full px-5 py-2 text-sm font-medium"
				class:bg-background={billing === 'monthly'}
				class:text-foreground={billing === 'monthly'}
				class:shadow-sm={billing === 'monthly'}
				class:text-muted-foreground={billing !== 'monthly'}
				aria-pressed={billing === 'monthly'}
			>
				Bulanan
			</button>
			<button
				type="button"
				onclick={() => (billing = 'yearly')}
				class="rounded-full px-5 py-2 text-sm font-medium"
				class:bg-background={billing === 'yearly'}
				class:text-foreground={billing === 'yearly'}
				class:shadow-sm={billing === 'yearly'}
				class:text-muted-foreground={billing !== 'yearly'}
				aria-pressed={billing === 'yearly'}
			>
				Tahunan
				<span class="ml-2 rounded-full bg-primary px-2 py-0.5 text-[10px] text-primary-foreground">
					Hemat 20%
				</span>
			</button>
		</div>
	</div>

	<div class="mx-auto mt-14 grid max-w-6xl gap-6 lg:grid-cols-3">
		{#each plans as plan (plan.id)}
			{@const price = billing === 'monthly' ? plan.priceMonthly : plan.priceYearly}
			{@const monthlyEquivalent = Math.round(plan.priceYearly / 12)}
			{@const yearlySaving = plan.priceMonthly * 12 - plan.priceYearly}
			<article
				class={`overflow-hidden rounded-3xl border bg-card ${plan.popular ? 'border-primary/50 shadow-xl' : ''}`}
			>
				<div class={`p-6 ${plan.popular ? 'bg-muted/50' : ''}`}>
					<div class="flex items-center justify-between gap-4">
						<h2 class="text-lg font-semibold">{plan.name}</h2>
						{#if plan.popular}
							<span class="inline-flex items-center gap-1 rounded-full bg-primary px-3 py-1 text-xs font-medium text-primary-foreground">
								<FlameIcon class="size-3.5" aria-hidden="true" /> Terpopuler
							</span>
						{/if}
					</div>

					<p class="mt-3 min-h-12 text-sm leading-6 text-muted-foreground">
						{descriptions[plan.slug]}
					</p>
					<div class="mt-5 flex items-end gap-2">
						<span class="text-4xl font-bold tracking-tight">Rp{formatPrice(price)}</span>
						<span class="pb-1 text-sm text-muted-foreground">
							{billing === 'monthly' ? '/bulan' : '/tahun'}
						</span>
					</div>
					<p class="mt-2 min-h-5 text-xs text-muted-foreground">
						{billing === 'yearly'
							? `Setara Rp${formatPrice(monthlyEquivalent)}/bulan · hemat Rp${formatPrice(yearlySaving)}`
							: 'Ditagih setiap bulan'}
					</p>

					<a href="/register" class="mt-6 flex h-12 w-full items-center justify-between rounded-full bg-primary py-1 pr-1 pl-5 text-sm font-medium text-primary-foreground">
						Berlangganan sekarang
						<span class="flex size-10 items-center justify-center rounded-full bg-primary-foreground text-primary">
							<ArrowUpRightIcon class="size-5" aria-hidden="true" />
						</span>
					</a>
				</div>

				<div class="border-t p-6">
					<p class="font-semibold">Termasuk dalam paket:</p>
					<ul class="mt-4 space-y-3 text-sm text-muted-foreground">
						{#each plan.features as feature (feature)}
							<li class="flex items-start gap-3">
								<span class="mt-0.5 flex size-5 shrink-0 items-center justify-center rounded-full bg-primary text-primary-foreground">
									<CheckIcon class="size-3.5" aria-hidden="true" />
								</span>
								<span>{feature}</span>
							</li>
						{/each}
					</ul>
				</div>
			</article>
		{/each}
	</div>

	<div class="mx-auto mt-24 max-w-6xl">
		<div class="max-w-2xl">
			<span class="inline-flex items-center gap-2 text-xs font-semibold uppercase tracking-[0.16em] text-primary">
				<SparklesIcon class="size-3.5" aria-hidden="true" /> Bandingkan paket
			</span>
			<h2 class="mt-5 text-3xl font-bold tracking-tight md:text-5xl">
				Satu venue, paket yang tumbuh bersama bisnis
			</h2>
			<p class="mt-4 text-base leading-7 text-muted-foreground md:text-lg">
				Bandingkan kapasitas tim, media, promosi, dan dukungan. Upgrade saat operasional berkembang.
			</p>
		</div>

		<div class="mt-10 overflow-x-auto rounded-3xl border bg-card">
			<table class="w-full min-w-[760px] text-left text-sm">
				<thead>
					<tr class="border-b bg-muted/70">
						<th class="w-[34%] px-6 py-6 font-semibold">Kapasitas & fitur</th>
						{#each plans as plan (plan.id)}
							<th class="px-5 py-6">
								<span class="block text-base font-semibold">{plan.name}</span>
								<span class="mt-1 block text-xs font-normal text-muted-foreground">
									Mulai Rp{formatPrice(plan.priceMonthly)}/bulan
								</span>
							</th>
						{/each}
					</tr>
				</thead>
				<tbody class="divide-y divide-border">
					<tr>
						<th class="px-6 py-4 font-medium">Maksimal lapangan</th>
						{#each plans as plan}<td class="border-l px-5 py-4 text-muted-foreground">{plan.maxCourts} lapangan</td>{/each}
					</tr>
					<tr class="bg-muted/20">
						<th class="px-6 py-4 font-medium">Staff</th>
						{#each plans as plan}<td class="border-l px-5 py-4 text-muted-foreground">{plan.maxStaff || 'Pemilik saja'}</td>{/each}
					</tr>
					<tr>
						<th class="px-6 py-4 font-medium">Custom domain</th>
						{#each plans as plan}<td class="border-l px-5 py-4">{#if plan.customDomainEnabled}<CheckIcon class="size-4 text-emerald-600" aria-label="Ya" />{:else}<XIcon class="size-4 text-muted-foreground" aria-label="Tidak" />{/if}</td>{/each}
					</tr>
					<tr class="bg-muted/20">
						<th class="px-6 py-4 font-medium">Blog dan galeri</th>
						{#each plans as plan}<td class="border-l px-5 py-4">{#if plan.blogEnabled && plan.galleryEnabled}<CheckIcon class="size-4 text-emerald-600" aria-label="Ya" />{:else}<XIcon class="size-4 text-muted-foreground" aria-label="Tidak" />{/if}</td>{/each}
					</tr>
					<tr>
						<th class="px-6 py-4 font-medium">Galeri dan penyimpanan</th>
						{#each plans as plan}<td class="border-l px-5 py-4 text-muted-foreground">{#if plan.maxGalleryImages}{plan.maxGalleryImages} foto · {formatStorage(plan.storageLimitBytes)}{:else}<MinusIcon class="size-4" aria-label="Tidak tersedia" />{/if}</td>{/each}
					</tr>
					<tr class="bg-muted/20">
						<th class="px-6 py-4 font-medium">Analitik dan export laporan</th>
						{#each plans as plan}<td class="border-l px-5 py-4">{#if plan.advancedAnalyticsEnabled && plan.reportExportEnabled}<CheckIcon class="size-4 text-emerald-600" aria-label="Ya" />{:else}<XIcon class="size-4 text-muted-foreground" aria-label="Tidak" />{/if}</td>{/each}
					</tr>
					<tr>
						<th class="px-6 py-4 font-medium">Notifikasi WhatsApp</th>
						{#each plans as plan}<td class="border-l px-5 py-4">{#if plan.whatsappEnabled}<CheckIcon class="size-4 text-emerald-600" aria-label="Ya" />{:else}<XIcon class="size-4 text-muted-foreground" aria-label="Tidak" />{/if}</td>{/each}
					</tr>
					<tr class="bg-muted/20">
						<th class="px-6 py-4 font-medium">Tanpa branding</th>
						{#each plans as plan}<td class="border-l px-5 py-4">{#if plan.removeBrandingEnabled}<CheckIcon class="size-4 text-emerald-600" aria-label="Ya" />{:else}<XIcon class="size-4 text-muted-foreground" aria-label="Tidak" />{/if}</td>{/each}
					</tr>
					<tr>
						<th class="px-6 py-4 font-medium">Dukungan</th>
						{#each plans as plan}<td class="border-l px-5 py-4 text-muted-foreground">{supportLabels[plan.supportPriority]}</td>{/each}
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</section>
