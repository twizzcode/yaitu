<script lang="ts">
	import { BlogDetail } from '#lib/components/blog/index.js';

	let { data } = $props();
</script>

<svelte:head>
	<title>{data.post.title} | LapanganKu.id</title>
	<meta name="description" content={data.post.excerpt} />
</svelte:head>

<BlogDetail post={data.post} related={data.related} />
