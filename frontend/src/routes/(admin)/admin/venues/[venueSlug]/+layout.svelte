<script lang="ts">
  import { page } from "$app/state";
  import * as Breadcrumb from "#lib/components/ui/breadcrumb/index.js";
  import { ScrollArea } from "#lib/components/ui/scroll-area/index.js";
  import * as Sidebar from "#lib/components/ui/sidebar/index.js";
  import { Separator } from "#lib/components/ui/separator/index.js";
  import AppSidebar from "#lib/components/sidebar/app-sidebar.svelte";

  let { data, children } = $props();

  const pageNames: Record<string, string> = {
    finance: "Keuangan",
    analytics: "Analisis",
    articles: "Berita dan Artikel",
    gallery: "Galeri",
    profile: "Profil",
    courts: "Lapangan",
    domains: "Domain",
    "operating-hours": "Jadwal",
    staff: "Staff",
    customers: "Pelanggan",
  };

  const venueUrl = $derived(`/admin/venues/${data.venue.slug}`);
  const section = $derived(page.url.pathname.slice(venueUrl.length).split("/").filter(Boolean)[0]);
  const currentPage = $derived(section ? (pageNames[section] ?? section) : "Dashboard");
</script>

<Sidebar.Provider class="h-svh overflow-hidden">
  <AppSidebar
    user={data.user}
    venues={data.venues}
    activeVenue={data.venue}
  />

  <Sidebar.Inset class="h-svh min-h-0 overflow-hidden">
    <header
      class="sticky top-0 z-20 flex h-16 shrink-0 items-center gap-2 border-b bg-background transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12"
    >
      <div class="flex items-center gap-2 px-4">
        <Sidebar.Trigger class="-ms-1" />

        <Separator
          orientation="vertical"
          class="me-2 data-vertical:h-4 data-vertical:self-auto"
        />

        <Breadcrumb.Root class="min-w-0">
          <Breadcrumb.List class="flex-nowrap">
            <Breadcrumb.Item class="hidden sm:inline-flex">
              <Breadcrumb.Link href={venueUrl}>{data.venue.name}</Breadcrumb.Link>
            </Breadcrumb.Item>

            <Breadcrumb.Separator class="hidden sm:block" />

            <Breadcrumb.Item class="min-w-0">
              <Breadcrumb.Page class="truncate font-medium">{currentPage}</Breadcrumb.Page>
            </Breadcrumb.Item>
          </Breadcrumb.List>
        </Breadcrumb.Root>
      </div>
    </header>

    <ScrollArea class="min-h-0 flex-1">
      <div class="flex min-h-full flex-col">
        <div class="flex-1 p-4 md:p-6">
          {@render children()}
        </div>

        <footer
          class="flex flex-col justify-between gap-3 border-t px-4 py-5 text-xs font-medium text-muted-foreground sm:flex-row sm:items-center md:px-6"
        >
          <p>&copy; 2026 LapanganKu.id. All rights reserved.</p>
          <nav aria-label="Tautan legal">
            <ul class="flex flex-wrap gap-x-4 gap-y-2">
              <li>
                <a href="/syarat-dan-ketentuan" class="underline transition-colors hover:text-primary">
                  Syarat dan Ketentuan
                </a>
              </li>
              <li>
                <a href="/kebijakan-privasi" class="underline transition-colors hover:text-primary">
                  Kebijakan Privasi
                </a>
              </li>
            </ul>
          </nav>
        </footer>
      </div>
    </ScrollArea>
  </Sidebar.Inset>
</Sidebar.Provider>
