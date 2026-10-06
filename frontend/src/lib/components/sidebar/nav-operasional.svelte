<script lang="ts">
  import * as Sidebar from "#lib/components/ui/sidebar/index.js";
  let {
    items,
  }: {
    items: {
      title: string;
      url: string;
      isActive: boolean;
      // this should be `Component` after @lucide/svelte updates types
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      icon: any;
    }[];
  } = $props();
</script>
<Sidebar.Group>
  <Sidebar.GroupLabel>Operasional</Sidebar.GroupLabel>
  <Sidebar.Menu>
    {#each items as item (item.title)}
      <Sidebar.MenuItem>
          <Sidebar.MenuButton
            isActive={item.isActive}
            tooltipContent={item.title}
          >
          {#snippet child({ props })}
            <a href={item.url} {...props}>
              <item.icon />
              <span>{item.title}</span>
            </a>
          {/snippet}
        </Sidebar.MenuButton>
      </Sidebar.MenuItem>
    {/each}
  </Sidebar.Menu>
</Sidebar.Group>
