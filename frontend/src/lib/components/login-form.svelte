<script lang="ts">
	import * as Card from '#lib/components/ui/card/index.js';
	import { Button } from '#lib/components/ui/button/index.js';
	import {
		FieldGroup,
		Field,
		FieldLabel,
		FieldError,
		FieldSeparator,
		FieldDescription
	} from '#lib/components/ui/field/index.js';
	import { Input } from '#lib/components/ui/input/index.js';
	import { cn } from '#lib/utils.js';
	import type { HTMLAttributes } from 'svelte/elements';
	import { page } from '$app/state';

	let {
		class: className,
		form,
		...restProps
	}: HTMLAttributes<HTMLDivElement> & {
		form?: { email?: string; message?: string } | null;
	} = $props();

	const id = $props.id();
	const hasError = $derived(!!form?.message);
	const googleLoginURL = $derived(
		`/auth/google?next=${encodeURIComponent(page.url.searchParams.get('next') ?? '/admin')}`
	);
</script>

<div class={cn('flex flex-col gap-6', className)} {...restProps}>
	<Card.Root class="min-h-[32rem] overflow-hidden p-0">
		<Card.Content class="grid flex-1 p-0 md:grid-cols-2">
			<form method="POST" class="flex flex-col p-6 md:p-8">
				<FieldGroup class="flex-1 justify-between">
					<div class="flex flex-col gap-4">
						<div class="flex flex-col items-center gap-2 text-center">
							<h1 class="text-2xl font-bold">Masuk ke Lapanganku</h1>
							<p class="text-balance text-muted-foreground">
								Kelola venue, lapangan, dan booking kamu.
							</p>
						</div>

						{#if form?.message}
							<FieldError>{form.message}</FieldError>
						{/if}

						<Field data-invalid={hasError}>
							<FieldLabel for="email-{id}">Email</FieldLabel>
							<Input
								id="email-{id}"
								name="email"
								type="email"
								autocomplete="email"
								placeholder="nama@email.com"
								required
								aria-invalid={hasError}
								value={form?.email ?? ''}
							/>
						</Field>

						<Field data-invalid={hasError}>
							<div class="flex items-center">
								<FieldLabel for="password-{id}">Password</FieldLabel>
								<a
									href="/forgot-password"
									class="ms-auto text-sm underline-offset-2 hover:underline"
								>
									Lupa password?
								</a>
							</div>
							<Input
								id="password-{id}"
								name="password"
								type="password"
								autocomplete="current-password"
								placeholder="Minimal 8 karakter"
								required
								aria-invalid={hasError}
							/>
						</Field>

						<Field>
							<Button type="submit">Masuk</Button>
						</Field>

						<FieldSeparator class="*:data-[slot=field-separator-content]:bg-card">
							Atau lanjutkan dengan
						</FieldSeparator>

						<Field>
							<Button
								variant="outline"
								href={googleLoginURL}
							>
								<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
									<path
										d="M12.48 10.92v3.28h7.84c-.24 1.84-.853 3.187-1.787 4.133-1.147 1.147-2.933 2.4-6.053 2.4-4.827 0-8.6-3.893-8.6-8.72s3.773-8.72 8.6-8.72c2.6 0 4.507 1.027 5.907 2.347l2.307-2.307C18.747 1.44 16.133 0 12.48 0 5.867 0 .307 5.387.307 12s5.56 12 12.173 12c3.573 0 6.267-1.173 8.373-3.36 2.16-2.16 2.84-5.213 2.84-7.667 0-.76-.053-1.467-.173-2.053H12.48z"
										fill="currentColor"
									/>
								</svg>
								Masuk dengan Google
							</Button>
						</Field>
					</div>

					<FieldDescription class="text-center">
						Belum punya akun? <a href="/register">Daftar sekarang</a>
					</FieldDescription>
				</FieldGroup>
			</form>

			<div class="relative hidden bg-muted md:block">
				<img
					src="/placeholder.svg"
					alt=""
					class="absolute inset-0 h-full w-full object-cover"
				/>
			</div>
		</Card.Content>
	</Card.Root>

	<FieldDescription class="px-6 text-center">
		Dengan melanjutkan, kamu menyetujui
		<a href="/terms">Ketentuan Layanan</a> dan
		<a href="/privacy">Kebijakan Privasi</a> kami.
	</FieldDescription>
</div>
