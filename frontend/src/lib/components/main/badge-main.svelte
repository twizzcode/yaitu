<script lang="ts">
    let { children } = $props();
</script>

<div class="relative inline-flex w-fit overflow-hidden rounded-full py-1 px-4 border text-xs sm:text-sm justify-center items-center gap-1 hover:scale-101 duration-25">
    {@render children()}
</div>