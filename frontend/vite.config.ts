import tailwindcss from '@tailwindcss/vite';
import adapter from '@sveltejs/adapter-node';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
  server: {
		allowedHosts: ['.lvh.me']
	},
	ssr: {
		// svelte-sonner mengirim file .svelte mentah, jadi harus diproses Vite
		// (bukan di-externalize) agar bisa dikompilasi saat SSR.
		noExternal: ['svelte-sonner']
	},
	plugins: [
		tailwindcss(),
		sveltekit({
			compilerOptions: {
				// Force runes mode for the project, except for libraries. Can be removed in svelte 6.
				runes: ({ filename }) =>
					filename.split(/[/\\]/).includes('node_modules') ? undefined : true
			},

			// adapter-node: build untuk deployment VPS (Node server di balik Caddy).
			adapter: adapter()
		})
	]
});
