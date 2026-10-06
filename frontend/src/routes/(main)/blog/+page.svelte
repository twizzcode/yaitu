<script lang="ts">
	import { blogPosts } from '#lib/blog.js';
	import { BlogList } from '#lib/components/blog/index.js';
</script>

<svelte:head>
	<title>Blog | LapanganKu.id</title>
	<meta name="description" content="Panduan praktis dan insight untuk pengelola lapangan olahraga." />
</svelte:head>

<BlogList posts={blogPosts} />
