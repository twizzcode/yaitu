import { PUBLIC_ROOT_DOMAIN } from '$app/env/public';
import type { Reroute } from '@sveltejs/kit/hooks';

const reservedSubdomains = new Set(['www', 'admin', 'api']);

// Cache sederhana untuk resolusi custom domain (hostname -> slug).
// Menghindari panggilan berulang untuk hostname yang sama.
type CacheEntry = { slug: string | null; expiresAt: number };
const customDomainCache = new Map<string, CacheEntry>();
const CACHE_TTL_MS = 60_000;

async function resolveCustomDomain(
	fetchFn: typeof fetch,
	hostname: string
): Promise<string | null> {
	const cached = customDomainCache.get(hostname);

	if (cached && cached.expiresAt > Date.now()) {
		return cached.slug;
	}

	let slug: string | null = null;

	try {
		// Same-origin: diteruskan ke backend oleh route /api/resolve-domain.
		const response = await fetchFn(
			`/api/resolve-domain?hostname=${encodeURIComponent(hostname)}`
		);

		if (response.ok) {
			const data = (await response.json()) as { venue_slug?: string };
			slug = data.venue_slug ?? null;
		}
	} catch {
		slug = null;
	}

	customDomainCache.set(hostname, {
		slug,
		expiresAt: Date.now() + CACHE_TTL_MS
	});

	return slug;
}

export const reroute: Reroute = async ({ url, fetch }) => {
	const hostname = url.hostname.toLowerCase();
	const rootDomain = PUBLIC_ROOT_DOMAIN.toLowerCase();

	// Jangan pernah rewrite endpoint internal (termasuk proxy resolve domain)
	// agar tidak terjadi rekursi.
	if (url.pathname.startsWith('/api/')) {
		return;
	}

	// Hostname lokal / IP langsung: bukan tenant, biarkan apa adanya.
	if (
		hostname === 'localhost' ||
		hostname === '127.0.0.1' ||
		hostname === '[::1]' ||
		hostname === '0.0.0.0'
	) {
		return;
	}

	// Root platform atau www: tidak di-rewrite.
	if (hostname === rootDomain || hostname === `www.${rootDomain}`) {
		return;
	}

	const domainSuffix = `.${rootDomain}`;

	// Subdomain platform: <slug>.<root-domain> -> /venues/<slug>
	if (hostname.endsWith(domainSuffix)) {
		const subdomain = hostname.slice(0, -domainSuffix.length);

		if (
			!subdomain ||
			subdomain.includes('.') ||
			reservedSubdomains.has(subdomain)
		) {
			return;
		}

		return rewriteToVenue(subdomain, url.pathname);
	}

	// Selain itu: kemungkinan custom domain customer. Resolve ke venue.
	const slug = await resolveCustomDomain(fetch, hostname);

	if (!slug) {
		return;
	}

	return rewriteToVenue(slug, url.pathname);
};

function rewriteToVenue(slug: string, pathname: string): string | undefined {
	if (pathname === '/') {
		return `/venues/${slug}`;
	}

	// Path storefront (mis. /courts/[courtId]/book) dipetakan ke route internal.
	if (pathname.startsWith('/courts/')) {
		return `/venues/${slug}${pathname}`;
	}

	// Biarkan path lain (mis. /login, /bookings) apa adanya.
	return;
}
