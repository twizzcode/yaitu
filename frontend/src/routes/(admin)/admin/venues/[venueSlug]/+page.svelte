<script lang="ts">
    import AdminHeroSection from "#lib/components/admin/admin-hero-section.svelte";
    import FieldRevenueChart from "#lib/components/admin/field-revenue-chart.svelte";
    import RevenueChart from "#lib/components/admin/revenue-chart.svelte";
    import StatsCards from "#lib/components/admin/stats-cards.svelte";
    import TopSportsTable from "#lib/components/admin/top-sports-table.svelte";
    import TrafficCards from "#lib/components/admin/traffic-cards.svelte";

	let { data } = $props();
</script>

<svelte:head>
	<title>Dashboard {data.venue.name} | Lapanganku</title>
	<meta
		name="description"
		content={`Dashboard pengelolaan ${data.venue.name}.`}
	/>
</svelte:head>

<main class="mx-auto max-w-7xl">
    <div class="flex flex-col gap-5">
        
        <AdminHeroSection
            title="Ringkasan Bisnis Anda"
            description="Pantau performa lapangan, pendapatan, dan aktivitas booking dalam satu tampilan."
        />
    
        <TrafficCards />

        <StatsCards />

        <div class="grid gap-5 xl:grid-cols-[2fr_1fr]">
          <RevenueChart />
          <FieldRevenueChart />
        </div>

        <TopSportsTable />
    </div>
	
</main>