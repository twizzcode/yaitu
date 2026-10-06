<script lang="ts">
	import LayoutDashboardIcon from '@lucide/svelte/icons/layout-dashboard';
	import CalendarDaysIcon from '@lucide/svelte/icons/calendar-days';
	import LogOutIcon from '@lucide/svelte/icons/log-out';
	import * as Avatar from '#lib/components/ui/avatar/index.js';
	import * as DropdownMenu from '#lib/components/ui/dropdown-menu/index.js';
	import { Button } from '#lib/components/ui/button/index.js';

	type User = {
		name: string;
		email: string;
	};

	let {
		user,
		showBookings = false
	}: {
		user: User;
		showBookings?: boolean;
	} = $props();

	const initials = $derived(
		user.name
			.split(/\s+/)
			.filter(Boolean)
			.slice(0, 2)
			.map((word) => word[0]?.toUpperCase())
			.join('')
	);
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger>
		{#snippet child({ props })}
			<Button {...props} variant="ghost" size="icon-lg" aria-label={`Buka menu ${user.name}`}>
				<Avatar.Root class="size-8">
					<Avatar.Fallback class="bg-primary text-primary-foreground">
						{initials}
					</Avatar.Fallback>
				</Avatar.Root>
			</Button>
		{/snippet}
	</DropdownMenu.Trigger>

	<DropdownMenu.Content class="w-64" align="end">
		<DropdownMenu.Label class="font-normal">
			<div class="grid gap-0.5">
				<span class="truncate text-sm font-medium">{user.name}</span>
				<span class="truncate text-xs text-muted-foreground">{user.email}</span>
			</div>
		</DropdownMenu.Label>
		<DropdownMenu.Separator />

		<DropdownMenu.Item>
			{#snippet child({ props })}
				<a {...props} href="/admin" class="flex w-full items-center gap-2">
					<LayoutDashboardIcon />
					Dashboard
				</a>
			{/snippet}
		</DropdownMenu.Item>

		{#if showBookings}
			<DropdownMenu.Item>
				{#snippet child({ props })}
					<a {...props} href="/my-bookings" class="flex w-full items-center gap-2">
						<CalendarDaysIcon />
						Booking saya
					</a>
				{/snippet}
			</DropdownMenu.Item>
		{/if}

		<DropdownMenu.Separator />
		<form method="POST" action="/auth/logout">
			<DropdownMenu.Item variant="destructive">
				{#snippet child({ props })}
					<button {...props} type="submit" class="flex w-full items-center gap-2">
						<LogOutIcon />
						Keluar
					</button>
				{/snippet}
			</DropdownMenu.Item>
		</form>
	</DropdownMenu.Content>
</DropdownMenu.Root>
