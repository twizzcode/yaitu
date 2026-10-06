<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import ChevronsUpDownIcon from "@lucide/svelte/icons/chevrons-up-down";
  import MapPinIcon from "@lucide/svelte/icons/map-pin";
  import PlusIcon from "@lucide/svelte/icons/plus";
  import * as DropdownMenu from "#lib/components/ui/dropdown-menu/index.js";
  import * as Sidebar from "#lib/components/ui/sidebar/index.js";
  import { useSidebar } from "#lib/components/ui/sidebar/index.js";

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
    venues,
    activeVenue,
  }: {
    venues: Venue[];
    activeVenue: Venue;
  } = $props();

  const sidebar = useSidebar();

  const preservablePaths = new Set([
    "",
    "/finance",
    "/analytics",
    "/articles",
    "/gallery",
    "/profile",
    "/courts",
    "/courts/new",
    "/domains",
    "/operating-hours",
    "/staff",
    "/customers",
  ]);

  function venueUrl(venueSlug: string) {
    const currentBase = `/admin/venues/${activeVenue.slug}`;
    const targetBase = `/admin/venues/${venueSlug}`;
    const suffix = page.url.pathname.startsWith(currentBase)
      ? page.url.pathname.slice(currentBase.length)
      : "";

    return `${targetBase}${preservablePaths.has(suffix) ? suffix : ""}`;
  }
</script>

<Sidebar.Menu>
  <Sidebar.MenuItem>
    <DropdownMenu.Root>
      <DropdownMenu.Trigger>
        {#snippet child({ props })}
          <Sidebar.MenuButton
            {...props}
            size="lg"
            class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
          >
            <div
              class="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground"
            >
              <MapPinIcon class="size-4" />
            </div>

            <div class="grid flex-1 text-start text-sm leading-tight">
              <span class="truncate font-medium">{activeVenue.name}</span>
              <span class="truncate text-xs">{activeVenue.slug}</span>
            </div>

            <ChevronsUpDownIcon class="ms-auto" />
          </Sidebar.MenuButton>
        {/snippet}
      </DropdownMenu.Trigger>

      <DropdownMenu.Content
        class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
        align="start"
        side={sidebar.isMobile ? "bottom" : "right"}
        sideOffset={4}
      >
        <DropdownMenu.Label class="text-xs text-muted-foreground">
          Venue
        </DropdownMenu.Label>

        {#each venues as venue (venue.id)}
          <DropdownMenu.Item
            onSelect={() => goto(venueUrl(venue.slug))}
            class="gap-2 p-2"
          >
            <div class="flex size-6 items-center justify-center rounded-md border">
              <MapPinIcon class="size-3.5 shrink-0" />
            </div>

            <div class="grid flex-1">
              <span class="truncate">{venue.name}</span>
              <span class="truncate text-xs text-muted-foreground">
                {venue.slug}
              </span>
            </div>

            {#if venue.id === activeVenue.id}
              <span class="text-xs text-muted-foreground">Aktif</span>
            {/if}
          </DropdownMenu.Item>
        {/each}

        <DropdownMenu.Separator />

        <DropdownMenu.Item
          onSelect={() => goto(`/admin/venues/new`)}
          class="gap-2 p-2"
        >
          <div
            class="flex size-6 items-center justify-center rounded-md border bg-transparent"
          >
            <PlusIcon class="size-4" />
          </div>

          <div class="font-medium text-muted-foreground">Tambah venue</div>
        </DropdownMenu.Item>
      </DropdownMenu.Content>
    </DropdownMenu.Root>
  </Sidebar.MenuItem>
</Sidebar.Menu>
