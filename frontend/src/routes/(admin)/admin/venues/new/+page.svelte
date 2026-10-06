<script lang="ts">
	import RegisterVenueForm from '#lib/components/venue/register-venue-form.svelte';

	let { data, form } = $props();
</script>

<svelte:head>
	<title>Tambah Venue | Lapanganku</title>
	<meta
		name="description"
		content="Daftarkan venue dan lapangan kamu di Lapanganku dan mulai terima booking online."
	/>
</svelte:head>

<RegisterVenueForm
	defaultEmail={data.user.email}
	rootDomain={data.rootDomain}
	hasExistingVenue={data.hasExistingVenue}
	{form}
/>
