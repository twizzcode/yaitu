<script lang="ts">
  import * as Sidebar from "#lib/components/ui/sidebar/index.js";
  let {
    items,
  }: {
    items: {
      name: string;
      url: string;
      isActive: boolean;
      // This should be `Component` after @lucide/svelte updates types
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      icon: any;
    }[];
  } = $props();
</script>
<Sidebar.Group>
  <Sidebar.GroupLabel>Menu Utama</Sidebar.GroupLabel>
  <Sidebar.Menu>
    {#each items as item (item.name)}
      <Sidebar.MenuItem>
          <Sidebar.MenuButton
            isActive={item.isActive}
            tooltipContent={item.name}
          >
          {#snippet child({ props })}
            <a href={item.url} {...props}>
              <item.icon />
              <span>{item.name}</span>
            </a>
          {/snippet}
        </Sidebar.MenuButton>
      </Sidebar.MenuItem>
    {/each}
  </Sidebar.Menu>
</Sidebar.Group>
