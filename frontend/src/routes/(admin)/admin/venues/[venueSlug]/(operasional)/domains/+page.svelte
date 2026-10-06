<script lang="ts">
	import { enhance } from '$app/forms';
	import CheckCircleIcon from '@lucide/svelte/icons/circle-check';
	import XCircleIcon from '@lucide/svelte/icons/circle-x';
	import ClockIcon from '@lucide/svelte/icons/clock';
	import SearchIcon from '@lucide/svelte/icons/search';
	import PlusIcon from '@lucide/svelte/icons/plus';
	import CopyIcon from '@lucide/svelte/icons/copy';
	import Trash2Icon from '@lucide/svelte/icons/trash-2';
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import ChevronDownIcon from '@lucide/svelte/icons/chevron-down';
	import GlobeIcon from '@lucide/svelte/icons/globe';
	import { Button } from '#lib/components/ui/button/index.js';
	import { Badge } from '#lib/components/ui/badge/index.js';
	import { Card, CardContent } from '#lib/components/ui/card/index.js';
	import * as Field from '#lib/components/ui/field/index.js';
	import { Input } from '#lib/components/ui/input/index.js';
	import AdminHeroSection from '#lib/components/admin/admin-hero-section.svelte';
	import { toast } from 'svelte-sonner';

	type Domain = {
		id: string;
		hostname: string;
		type: 'platform' | 'custom';
		status: 'pending' | 'active' | 'failed';
		verification_token: string;
		verified_at: string | null;
		created_at: string;
		is_apex: boolean;
	};

	let { data, form } = $props();

	let search = $state('');
	let addOpen = $state(false);
	let expandedId = $state<string | null>(null);
	let verifyingId = $state<string | null>(null);
	let removingId = $state<string | null>(null);

	const platformDomains = $derived(data.domains.filter((d) => d.type === 'platform'));
	const customDomains = $derived(data.domains.filter((d) => d.type === 'custom'));

	const filteredPlatform = $derived(
		platformDomains.filter((d) => d.hostname.includes(search.trim().toLowerCase()))
	);
	const filteredCustom = $derived(
		customDomains.filter((d) => d.hostname.includes(search.trim().toLowerCase()))
	);

	const statusMeta = {
		active: { label: 'Valid Configuration', tone: 'text-emerald-600 dark:text-emerald-400' },
		pending: { label: 'Menunggu verifikasi DNS', tone: 'text-amber-600 dark:text-amber-400' },
		failed: { label: 'Konfigurasi belum benar', tone: 'text-destructive' }
	} as const;

	$effect(() => {
		if (form?.success && form.message) {
			toast.success(form.message);
			addOpen = false;
		}

		if (form?.message && !form?.success) {
			toast.error(form.message);
		}
	});

	function toggle(id: string) {
		expandedId = expandedId === id ? null : id;
	}

	async function copyText(text: string) {
		try {
			await navigator.clipboard.writeText(text);
			toast.success('Disalin ke clipboard');
		} catch {
			toast.error('Gagal menyalin');
		}
	}
</script>

<svelte:head>
	<title>Domain | Lapanganku</title>
	<meta name="description" content="Kelola domain venue." />
</svelte:head>

<main class="mx-auto max-w-7xl space-y-5">
	<AdminHeroSection
		title="Domain"
		description="Hubungkan domain kustom ke venue kamu dan pantau status verifikasinya."
		size="sm"
	/>

	<!-- Toolbar -->
	<div class="flex flex-col gap-3 sm:flex-row sm:items-center">
		<div class="relative flex-1">
			<SearchIcon
				class="text-muted-foreground pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2"
			/>
			<Input class="h-9 pl-9" placeholder="Cari domain" bind:value={search} />
		</div>

		<div class="flex gap-2">
			<Button variant="outline" class="h-9" onclick={() => (addOpen = !addOpen)}>
				<PlusIcon />
				Tambah Domain
			</Button>
			<Button
				class="h-9"
				href={`https://${data.venue.slug}.${data.rootDomain}`}
				target="_blank"
			>
				<GlobeIcon />
				Buka Venue
			</Button>
		</div>
	</div>

	<!-- Form tambah domain -->
	{#if addOpen}
		<Card>
			<CardContent>
				<form
					method="POST"
					action="?/add"
					use:enhance
					class="flex flex-col gap-3 sm:flex-row sm:items-end"
				>
					<Field.Field class="flex-1">
						<Field.Label for="hostname">Domain</Field.Label>
						<Input
							id="hostname"
							name="hostname"
							placeholder="lapangan-a.com atau booking.lapangan-a.com"
							required
							autocomplete="off"
							spellcheck={false}
						/>
						<Field.Description>
							Masukkan domain atau subdomain yang kamu miliki.
						</Field.Description>
					</Field.Field>
					<Button type="submit">Tambah</Button>
				</form>
			</CardContent>
		</Card>
	{/if}

	{#if data.domains.length === 0}
		<div class="rounded-xl border border-dashed p-10 text-center">
			<h3 class="font-semibold">Belum ada domain</h3>
			<p class="mt-2 text-sm text-muted-foreground">
				Tambahkan domain kustom untuk venue ini.
			</p>
		</div>
	{:else}
		<div class="space-y-3">
			<!-- Domain platform (read-only) -->
			{#each filteredPlatform as domain (domain.id)}
				<Card class="overflow-hidden py-0">
					<CardContent class="flex items-center justify-between gap-4 p-4">
						<div class="flex min-w-0 items-center gap-3">
							<CheckCircleIcon class="size-5 shrink-0 text-emerald-600 dark:text-emerald-400" />
							<div class="min-w-0">
								<p class="truncate font-medium">{domain.hostname}</p>
								<p class="text-xs text-emerald-600 dark:text-emerald-400">
									Valid Configuration
								</p>
							</div>
						</div>
						<Badge variant="secondary" class="shrink-0">Platform</Badge>
					</CardContent>
				</Card>
			{/each}

			<!-- Domain kustom -->
			{#each filteredCustom as domain (domain.id)}
				{@const meta = statusMeta[domain.status]}
				<Card class="overflow-hidden py-0">
					<CardContent class="p-0">
						<!-- Baris ringkas -->
						<button
							type="button"
							onclick={() => toggle(domain.id)}
							class="flex w-full items-center justify-between gap-4 p-4 text-left transition-colors hover:bg-muted/40"
						>
							<div class="flex min-w-0 items-center gap-3">
								{#if domain.status === 'active'}
									<CheckCircleIcon class="size-5 shrink-0 text-emerald-600 dark:text-emerald-400" />
								{:else if domain.status === 'pending'}
									<ClockIcon class="size-5 shrink-0 text-amber-600 dark:text-amber-400" />
								{:else}
									<XCircleIcon class="size-5 shrink-0 text-destructive" />
								{/if}

								<div class="min-w-0">
									<p class="truncate font-medium">{domain.hostname}</p>
									<p class="text-xs {meta.tone}">{meta.label}</p>
								</div>
							</div>

							<div class="flex shrink-0 items-center gap-2">
								{#if domain.status === 'active'}
									<Badge variant="secondary">Custom</Badge>
								{:else}
									<Badge variant="outline">Pending</Badge>
								{/if}
								<ChevronDownIcon
									class="text-muted-foreground size-4 transition-transform {expandedId ===
									domain.id
										? 'rotate-180'
										: ''}"
								/>
							</div>
						</button>

						<!-- Detail -->
						{#if expandedId === domain.id}
							<div class="border-t p-4 md:p-6">
								<div class="space-y-5">
									<Field.Field>
										<Field.Label for={`domain-${domain.id}`}>Domain</Field.Label>
										<Input
											id={`domain-${domain.id}`}
											value={domain.hostname}
											readonly
											disabled
										/>
									</Field.Field>

									<div>
										<p class="text-sm font-medium">Konfigurasi DNS</p>
										<p class="mt-1 text-sm text-muted-foreground">
											{#if domain.is_apex}
												Buat record <strong>A</strong> berikut di DNS domain kamu, lalu klik
												Verifikasi.
											{:else}
												Buat record <strong>CNAME</strong> atau <strong>A</strong> berikut di
												DNS domain kamu, lalu klik Verifikasi.
											{/if}
											Propagasi DNS bisa memakan waktu beberapa menit.
										</p>

										<div class="mt-3 space-y-3">
											<Field.Field>
												<Field.Label for={`dns-name-${domain.id}`}>Nama (Host)</Field.Label>
												<div class="flex gap-2">
													<Input
														id={`dns-name-${domain.id}`}
														value={domain.is_apex ? '@' : domain.hostname.split('.')[0]}
														readonly
														class="font-mono text-xs"
													/>
													<Button
														type="button"
														variant="outline"
														size="icon"
														onclick={() =>
															copyText(domain.is_apex ? '@' : domain.hostname.split('.')[0])}
														aria-label="Salin nama record"
													>
														<CopyIcon />
													</Button>
												</div>
											</Field.Field>

											<Field.Field>
												<Field.Label for={`dns-value-${domain.id}`}>
													{domain.is_apex ? 'Type A — Value (IP)' : 'Type A — Value (IP)'}
												</Field.Label>
												<div class="flex gap-2">
													<Input
														id={`dns-value-${domain.id}`}
														value={data.serverIp}
														readonly
														class="font-mono text-xs"
													/>
													<Button
														type="button"
														variant="outline"
														size="icon"
														onclick={() => copyText(data.serverIp)}
														aria-label="Salin nilai record"
													>
														<CopyIcon />
													</Button>
												</div>
												<Field.Description>
													{domain.is_apex
														? 'Apex domain wajib memakai record A ke IP di atas.'
														: 'Bisa juga memakai CNAME ke ' + domain.hostname + ' (bila apex sudah diarahkan).'}
												</Field.Description>
											</Field.Field>
										</div>
									</div>

									<div class="flex flex-wrap items-center justify-between gap-3 border-t pt-4">
										<form
											method="POST"
											action="?/remove"
											use:enhance={() => {
												removingId = domain.id;
												return async ({ update }) => {
													await update();
													removingId = null;
												};
											}}
										>
											<input type="hidden" name="domain_id" value={domain.id} />
											<Button
												type="submit"
												variant="destructive"
												disabled={removingId === domain.id}
											>
												<Trash2Icon />
												{removingId === domain.id ? 'Menghapus...' : 'Hapus'}
											</Button>
										</form>

										<form
											method="POST"
											action="?/verify"
											use:enhance={() => {
												verifyingId = domain.id;
												return async ({ update }) => {
													await update();
													verifyingId = null;
												};
											}}
										>
											<input type="hidden" name="domain_id" value={domain.id} />
											<Button type="submit" disabled={verifyingId === domain.id}>
												<RefreshCwIcon
													class={verifyingId === domain.id ? 'animate-spin' : ''}
												/>
												{verifyingId === domain.id ? 'Memverifikasi...' : 'Verifikasi'}
											</Button>
										</form>
									</div>
								</div>
							</div>
						{/if}
					</CardContent>
				</Card>
			{/each}

			{#if filteredPlatform.length === 0 && filteredCustom.length === 0}
				<div class="rounded-xl border border-dashed p-8 text-center text-sm text-muted-foreground">
					Tidak ada domain yang cocok dengan pencarian.
				</div>
			{/if}
		</div>
	{/if}
</main>
