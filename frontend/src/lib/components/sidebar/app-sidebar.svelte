<script lang="ts">
  import Building2Icon from "@lucide/svelte/icons/building-2";
  import CalendarDaysIcon from "@lucide/svelte/icons/calendar-days";
  import ChartNoAxesCombinedIcon from "@lucide/svelte/icons/chart-no-axes-combined";
  import ContactIcon from "@lucide/svelte/icons/contact";
  import Globe2Icon from "@lucide/svelte/icons/globe-2";
  import ImagesIcon from "@lucide/svelte/icons/images";
  import LayoutDashboardIcon from "@lucide/svelte/icons/layout-dashboard";
  import MapPinnedIcon from "@lucide/svelte/icons/map-pinned";
  import NewspaperIcon from "@lucide/svelte/icons/newspaper";
  import UsersRoundIcon from "@lucide/svelte/icons/users-round";
  import WalletCardsIcon from "@lucide/svelte/icons/wallet-cards";
  import type { ComponentProps } from "svelte";
  import * as Sidebar from "#lib/components/ui/sidebar/index.js";
  import NavMain from "./nav-main.svelte";
  import NavOperasional from "./nav-operasional.svelte";
  import NavUser from "./nav-user.svelte";
  import VenueSwitcher from "./venue-switcher.svelte";
  import { page } from "$app/state";

  type User = {
    name: string;
    email: string;
  };

  type Venue = {
    id: string;
    name: string;
    slug: string;
    address: string;
    timezone: string;
    whatsapp: string;
    is_active: boolean;
  };

  let {
    ref = $bindable(null),
    collapsible = "icon",
    user,
    venues,
    activeVenue,
    ...restProps
  }: ComponentProps<typeof Sidebar.Root> & {
    user: User;
    venues: Venue[];
    activeVenue: Venue;
  } = $props();

  const venueBase = $derived(
    `/admin/venues/${activeVenue.slug}`,
  );

  function isPathActive(url: string, exact = false) {
    return exact
      ? page.url.pathname === url
      : page.url.pathname === url || page.url.pathname.startsWith(`${url}/`);
  }

  const navMain = $derived([
    {
      name: "Dashboard",
      url: venueBase,
      icon: LayoutDashboardIcon,
      isActive: isPathActive(venueBase, true),
    },
    {
      name: "Keuangan",
      url: `${venueBase}/finance`,
      icon: WalletCardsIcon,
      isActive: isPathActive(`${venueBase}/finance`),
    },
    {
      name: "Analisis",
      url: `${venueBase}/analytics`,
      icon: ChartNoAxesCombinedIcon,
      isActive: isPathActive(`${venueBase}/analytics`),
    },
    {
      name: "Berita dan Artikel",
      url: `${venueBase}/articles`,
      icon: NewspaperIcon,
      isActive: isPathActive(`${venueBase}/articles`),
    },
    {
      name: "Galeri",
      url: `${venueBase}/gallery`,
      icon: ImagesIcon,
      isActive: isPathActive(`${venueBase}/gallery`),
    },
    {
      name: "Profil",
      url: `${venueBase}/profile`,
      icon: Building2Icon,
      isActive: isPathActive(`${venueBase}/profile`),
    },
  ]);

  const navOperasional = $derived([
    {
      title: "Lapangan",
      url: `${venueBase}/courts`,
      icon: MapPinnedIcon,
      isActive: isPathActive(`${venueBase}/courts`),
    },
    {
      title: "Domain",
      url: `${venueBase}/domains`,
      icon: Globe2Icon,
      isActive: isPathActive(`${venueBase}/domains`),
    },
    {
      title: "Jadwal",
      url: `${venueBase}/operating-hours`,
      icon: CalendarDaysIcon,
      isActive: isPathActive(`${venueBase}/operating-hours`),
    },
    {
      title: "Staff",
      url: `${venueBase}/staff`,
      icon: ContactIcon,
      isActive: isPathActive(`${venueBase}/staff`),
    },
    {
      title: "Pelanggan",
      url: `${venueBase}/customers`,
      icon: UsersRoundIcon,
      isActive: isPathActive(`${venueBase}/customers`),
    },
  ]);
</script>

<Sidebar.Root bind:ref {collapsible} {...restProps}>
  <Sidebar.Header>
    <VenueSwitcher {venues} {activeVenue} />
  </Sidebar.Header>

  <Sidebar.Content>
    <NavMain items={navMain} />
    <NavOperasional items={navOperasional} />
  </Sidebar.Content>

  <Sidebar.Footer>
    <NavUser
      user={{
        name: user.name,
        email: user.email,
        avatar: "",
      }}
    />
  </Sidebar.Footer>

  <Sidebar.Rail />
</Sidebar.Root>