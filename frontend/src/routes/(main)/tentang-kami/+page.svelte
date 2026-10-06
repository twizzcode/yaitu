<script lang="ts">
	import ArrowUpRightIcon from '@lucide/svelte/icons/arrow-up-right';
	import CheckCircleIcon from '@lucide/svelte/icons/circle-check';
	import CompassIcon from '@lucide/svelte/icons/compass';
	import HeadphonesIcon from '@lucide/svelte/icons/headphones';
	import MailIcon from '@lucide/svelte/icons/mail';
	import MessageCircleIcon from '@lucide/svelte/icons/message-circle';
	import SendIcon from '@lucide/svelte/icons/send';
	import ShieldCheckIcon from '@lucide/svelte/icons/shield-check';
	import SparklesIcon from '@lucide/svelte/icons/sparkles';
	import TargetIcon from '@lucide/svelte/icons/target';
	import UsersIcon from '@lucide/svelte/icons/users';
	import ZapIcon from '@lucide/svelte/icons/zap';
	import BadgeMain from '#lib/components/main/badge-main.svelte';
	import { Button } from '#lib/components/ui/button/index.js';

	const stats = [
		{ value: '100%', label: 'Berbasis cloud' },
		{ value: '24/7', label: 'Akses dashboard' },
		{ value: '1', label: 'Sistem terintegrasi' },
		{ value: '6+', label: 'Jenis venue olahraga' }
	];

	const visionMission = [
		{
			icon: CompassIcon,
			tag: 'Visi',
			title: 'Jadi standar operasional lapangan olahraga di Indonesia.',
			content: 'Kami ingin setiap venue, dari satu lapangan sampai sport center, bisa dikelola dengan rapi tanpa alat yang rumit.'
		},
		{
			icon: TargetIcon,
			tag: 'Misi',
			title: 'Menyederhanakan setiap alur kerja pengelola lapangan.',
			content: 'Booking, jadwal, pembayaran, pelanggan, dan laporan kami satukan menjadi pengalaman yang cepat dipahami siapa pun.'
		}
	];

	const values = [
		{ icon: CheckCircleIcon, title: 'Sederhana', content: 'Fitur harus mudah dipahami dan menyelesaikan pekerjaan tanpa menambah langkah yang tidak perlu.' },
		{ icon: ShieldCheckIcon, title: 'Dapat dipercaya', content: 'Data operasional harus akurat, aman, dan selalu bisa diandalkan oleh pemilik maupun staff.' },
		{ icon: UsersIcon, title: 'Dekat dengan pengguna', content: 'Masukan pengelola lapangan menjadi dasar keputusan produk yang kami buat.' },
		{ icon: ZapIcon, title: 'Terus membaik', content: 'Kami mengembangkan produk secara bertahap berdasarkan masalah nyata di lapangan.' }
	];

	const team = [
		{ name: 'Twizzcode Team', role: 'Product & Engineering', img: '/avatars/james.png' }
	];

	const contactChannels = [
		{ icon: MessageCircleIcon, label: 'WhatsApp', value: '+62 851-7444-6002', href: 'https://wa.me/6285174446002' },
		{ icon: MailIcon, label: 'Email', value: 'halo@lapanganku.id', href: 'mailto:halo@lapanganku.id' },
		{ icon: HeadphonesIcon, label: 'Dukungan', value: 'Senin - Sabtu', href: 'mailto:halo@lapanganku.id?subject=Dukungan%20LapanganKu.id' },
		{ icon: CompassIcon, label: 'Lokasi', value: 'Indonesia', href: '#kontak' }
	];
</script>

<svelte:head>
	<title>Tentang Kami | LapanganKu.id</title>
	<meta name="description" content="Kenali visi, misi, nilai, dan tim di balik LapanganKu.id." />
</svelte:head>

<div class="flex flex-col gap-20 pb-24 md:gap-28">
	<section class="mx-auto w-full max-w-7xl px-4 pt-12 sm:px-8 md:pt-20">
		<div class="mx-auto flex max-w-4xl flex-col items-center text-center">
			<BadgeMain><SparklesIcon class="size-3.5 text-primary" /><span>Tentang LapanganKu.id</span></BadgeMain>
			<h1 class="mt-7 text-4xl font-bold tracking-tight md:text-5xl lg:text-6xl">
				Kami bantu pengelola lapangan<br class="hidden md:block" /> bekerja lebih tenang.
			</h1>
			<p class="mt-6 max-w-2xl text-base leading-7 text-muted-foreground md:text-lg">
				LapanganKu.id lahir dari pengalaman nyata mengelola lapangan olahraga, dari jadwal yang bentrok,
				rekap pembayaran manual, sampai pelanggan yang menunggu balasan admin. Kami ubah semuanya menjadi satu alur sederhana.
			</p>
			<div class="mt-8 flex flex-wrap items-center justify-center gap-3">
				<Button href="/harga" size="lg" class="px-6 font-semibold shadow-lg">Lihat Paket <ArrowUpRightIcon /></Button>
				<Button href="#kontak" size="lg" variant="outline" class="px-6 font-semibold shadow-lg">Hubungi Kami</Button>
			</div>
		</div>
	</section>

	<section class="mx-auto w-full max-w-7xl px-4 sm:px-8">
		<div class="grid grid-cols-2 gap-px overflow-hidden rounded-3xl border bg-border p-px shadow-sm md:grid-cols-4">
			{#each stats as stat (stat.label)}
				<div class="flex flex-col items-center justify-center gap-1 bg-background px-4 py-8 text-center">
					<span class="text-3xl font-bold tracking-tight md:text-4xl">{stat.value}</span>
					<span class="text-xs font-medium text-muted-foreground md:text-sm">{stat.label}</span>
				</div>
			{/each}
		</div>
	</section>

	<section class="mx-auto w-full max-w-7xl px-4 sm:px-8">
		<div class="flex flex-col items-center gap-4 text-center">
			<BadgeMain><span>Visi & Misi</span></BadgeMain>
			<h2 class="max-w-2xl text-3xl font-semibold tracking-tight md:text-4xl">Arah yang kami tuju dan cara kami mencapainya.</h2>
			<p class="max-w-xl leading-7 text-muted-foreground">Komitmen kami untuk membuat operasional lapangan olahraga lebih sederhana, rapi, dan mudah dikelola.</p>
		</div>

		<div class="mt-12 overflow-hidden rounded-3xl border bg-card shadow-sm">
			<div class="grid divide-y md:grid-cols-2 md:divide-x md:divide-y-0">
				{#each visionMission as card (card.tag)}
					<div class="flex min-h-80 flex-col items-center justify-center p-8 text-center md:p-10">
						<div class="flex size-12 items-center justify-center rounded-xl border bg-background text-primary shadow-xs"><card.icon class="size-6" strokeWidth={1.8} /></div>
						<span class="mt-6 text-xs font-semibold uppercase tracking-widest text-primary">{card.tag}</span>
						<h3 class="mt-3 max-w-lg text-2xl font-semibold tracking-tight md:text-3xl">{card.title}</h3>
						<p class="mt-4 max-w-lg leading-7 text-muted-foreground">{card.content}</p>
					</div>
				{/each}
			</div>
		</div>
	</section>

	<section class="mx-auto w-full max-w-7xl px-4 sm:px-8">
		<div class="flex flex-col items-center gap-4 text-center">
			<BadgeMain><span>Nilai yang Kami Pegang</span></BadgeMain>
			<h2 class="max-w-2xl text-3xl font-semibold tracking-tight md:text-4xl">Prinsip yang membimbing setiap keputusan kami.</h2>
			<p class="max-w-xl leading-7 text-muted-foreground">Empat hal sederhana yang selalu kami jaga saat membangun dan mengembangkan LapanganKu.id.</p>
		</div>
		<div class="mx-auto mt-12 grid max-w-5xl gap-x-12 gap-y-8 md:grid-cols-2">
			{#each values as value (value.title)}
				<div class="flex gap-4">
					<div class="flex size-11 shrink-0 items-center justify-center rounded-xl border text-primary"><value.icon class="size-5" strokeWidth={1.8} /></div>
					<div><h3 class="font-semibold">{value.title}</h3><p class="mt-2 text-sm leading-6 text-muted-foreground">{value.content}</p></div>
				</div>
			{/each}
		</div>
	</section>

	<section class="mx-auto w-full max-w-7xl px-4 sm:px-8">
		<div class="flex flex-col items-center gap-4 text-center">
			<BadgeMain><span>Tim di Balik LapanganKu.id</span></BadgeMain>
			<h2 class="max-w-2xl text-3xl font-semibold tracking-tight md:text-4xl">Orang-orang yang membangun dan merawat produknya.</h2>
			<p class="max-w-xl leading-7 text-muted-foreground">Tim kecil yang menikmati berbicara langsung dengan pengelola lapangan.</p>
		</div>
		<div class="mx-auto mt-12 grid max-w-sm grid-cols-1">
			{#each team as member (member.name)}
				<div class="flex flex-col items-center rounded-2xl border bg-card p-6 text-center shadow-xs">
					<div class="rounded-full bg-gradient-to-br from-primary/20 to-primary/5 p-1"><img src={member.img} alt={member.name} width="120" height="120" class="size-20 rounded-full border object-cover" /></div>
					<h3 class="mt-4 text-base font-semibold">{member.name}</h3><p class="mt-1 text-xs text-muted-foreground">{member.role}</p>
				</div>
			{/each}
		</div>
	</section>

	<section id="kontak" class="mx-auto w-full max-w-7xl scroll-mt-6 px-4 sm:px-8">
		<div class="overflow-hidden rounded-[2rem] border bg-muted/40 p-8 shadow-sm md:p-12">
			<div class="grid gap-10 lg:grid-cols-2 lg:gap-16">
				<div class="flex h-full flex-col items-start gap-4">
					<BadgeMain><span>Kontak</span></BadgeMain>
					<h2 class="max-w-2xl text-3xl font-semibold tracking-tight md:text-4xl">Ada pertanyaan? Ngobrol saja dulu.</h2>
					<p class="max-w-md leading-7 text-muted-foreground">Tim kami siap membantu setup, migrasi data, atau menjelaskan fitur yang cocok untuk venue Anda.</p>
					<div class="mt-auto flex w-full flex-wrap gap-3 pt-6">
						<Button href="https://wa.me/6285174446002" size="lg" class="px-6 font-semibold shadow-lg"><SendIcon /> Chat via WhatsApp</Button>
						<Button href="mailto:halo@lapanganku.id" size="lg" variant="outline" class="px-6 font-semibold shadow-lg">Kirim Email</Button>
					</div>
				</div>
				<div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
					{#each contactChannels as channel (channel.label)}
						<a href={channel.href} class="flex h-full flex-col gap-4 rounded-2xl border bg-background p-5 shadow-xs transition hover:border-primary/40 hover:shadow-md">
							<div class="flex size-10 items-center justify-center rounded-xl border bg-muted/40 text-primary"><channel.icon class="size-5" strokeWidth={1.7} /></div>
							<div><p class="text-xs font-medium text-muted-foreground">{channel.label}</p><p class="mt-1 text-sm font-semibold">{channel.value}</p></div>
						</a>
					{/each}
				</div>
			</div>
		</div>
	</section>
</div>
