<script lang="ts">
	import RegisterForm from '#lib/components/register-form.svelte';

	let { form } = $props();
</script>

<svelte:head>
	<title>Daftar | Lapanganku</title>
	<meta name="description" content="Buat akun Lapanganku." />
</svelte:head>

<div class="flex min-h-svh flex-col items-center justify-center bg-muted p-6 md:p-10">
	<div class="w-full max-w-sm md:max-w-3xl">
		<RegisterForm {form} />
	</div>
</div>
