import { error } from '@sveltejs/kit';
import { getBlogPost, getRelatedPosts } from '#lib/blog.js';
import type { PageServerLoad } from './$types';

export const load = (({ params }) => {
	const post = getBlogPost(params.slug);
	if (!post) error(404, 'Artikel tidak ditemukan');

	return {
		post,
		related: getRelatedPosts(post)
	};
}) satisfies PageServerLoad;
