import { defineEnvVars } from '@sveltejs/kit/env';

// Semua environment variable punya default untuk pengembangan lokal.
// Di produksi, nilainya datang dari `.env` di root repo (lihat
// `.env.production.example` dan docs/DEPLOY.md).
export const variables = defineEnvVars({
	PUBLIC_ROOT_DOMAIN: {
		public: true,
		static: true,
		schema: (value) => value || 'lvh.me'
	},
	PUBLIC_ROOT_URL: {
		public: true,
		static: true,
		schema: (value) => value || 'http://lvh.me:5173'
	},
	// Origin backend API (server-only). Di produksi diarahkan ke service backend.
	API_URL: {
		public: false,
		static: false,
		schema: (value) => value || 'http://localhost:8080'
	},
	// IP publik server (server-only), ditampilkan sebagai instruksi A record.
	SERVER_PUBLIC_IP: {
		public: false,
		static: false,
		schema: (value) => value ?? ''
	}
});
