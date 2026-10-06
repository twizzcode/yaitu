<script lang="ts">
	import UserMenu from '#lib/components/user-menu.svelte';

	type User = {
		id: string;
		name: string;
		email: string;
	};

	let {
		user
	}: {
		user: User | null;
	} = $props();

</script>

<header class="border-b bg-card">
	<div class="mx-auto flex max-w-6xl items-center justify-between gap-4 px-6 py-4">
		<a href="/" class="text-lg font-bold">
			Lapanganku
		</a>

		<nav
			aria-label="Navigasi pengguna"
			class="flex items-center gap-3 text-sm"
		>
			{#if user}
				<UserMenu {user} showBookings />
			{:else}
				<a
					href="/auth/login?next=/"
					class="rounded-md bg-primary px-4 py-2 font-medium text-primary-foreground"
				>
					Masuk
				</a>
			{/if}
		</nav>
	</div>
</header>
