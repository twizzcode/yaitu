# AGENTS.md — Panduan untuk AI Agent / Developer

> **Baca file ini dulu.** Lalu baca [`docs/PROJECT-HANDOFF.md`](docs/PROJECT-HANDOFF.md)
> untuk konteks lengkap. Dokumen ini ringkas; handoff adalah sumber detail.

## Apa proyek ini

**Lapanganku** — SaaS multi-tenant manajemen venue & booking lapangan olahraga.
Satu user bisa punya banyak venue; tiap venue punya lapangan (court) dan booking.
Model: `User → Venue (member owner/admin) → Court → Booking`.

## Stack

| Bagian | Teknologi |
|---|---|
| Backend | Go (standard library `net/http`, tanpa framework), `pgxpool`, tanpa ORM |
| Frontend | SvelteKit 3, Svelte 5 (runes), Tailwind CSS 4, shadcn-svelte, Bun |
| DB | PostgreSQL 18 |
| Storage | Cloudflare R2 (presigned URL) untuk KTP/logo |
| Bayar | Midtrans QRIS (Core API) |
| Deploy | Docker Compose + Caddy (reverse proxy + auto SSL) |
| Adapter | `@sveltejs/adapter-node` |

## Struktur

```text
backend/     Go API + migrations/ (SQL, versi terakhir: 19)
frontend/    SvelteKit
docs/        PROJECT-HANDOFF.md (detail), DEPLOY.md (produksi)
compose.yaml        Prod (postgres, migrate, backend, frontend, caddy)
compose.dev.yaml    Dev (postgres, migrate)
Makefile            Perintah dev (make help)
.env                Prod (gitignored)      .env.local  Dev (gitignored)
```

## Menjalankan

**Dev** (backend & frontend native, hot reload):

```bash
make db        # start PostgreSQL dev + migration otomatis
make backend   # terminal 1: Go API :8080 (baca .env.local)
make frontend  # terminal 2: SvelteKit :5173
```

URL: `http://lvh.me:5173` (platform), `<slug>.lvh.me:5173` (storefront), `http://localhost:8080` (API).

**Prod** (di VPS): `docker compose up -d --build` (pakai `.env`). Lihat `docs/DEPLOY.md`.

**Verifikasi wajib sebelum selesai:** `make check` (backend `go test` + frontend `svelte-check`).

## Aturan penting (JANGAN dilanggar)

- Model **`User → Venue`**; **tidak ada workspace** (sudah dihapus). Jangan
  reintroduksi.
- Akses venue **selalu** lewat `venue_members` (owner/admin). Jangan percaya
  `venueID` dari klien tanpa cek keanggotaan (`requireVenueAccess` / `hasVenueAccess`).
- Jangan simpan token session / auth code mentah di DB (selalu hash SHA-256).
- Jangan taruh session token di URL.
- Jangan set cookie `Domain=.domain.com` (session host-only by design).
- Jangan percaya harga dari frontend; hitung di backend.
- Pertahankan **DB overlap constraint** booking (race safety).
- Custom domain hanya aktif setelah verifikasi DNS (A record).
- Endpoint ask TLS (`/api/internal/tls/allow/{token}`) wajib dijaga
  `SAAS_ASK_TOKEN`.
- `frontend/src/hooks.ts` tidak boleh me-rewrite `/api/*` (anti rekursi).
- Jangan commit rahasia: `.env`, `.env.local`, `certs/`, `backend/main.txt`,
  `cookies.txt` (semua sudah di `.gitignore`).

## Konvensi backend

- Env dibaca via `os.Getenv`. Fitur eksternal **opsional**: bila env kosong,
  fitur dinonaktifkan (log pesan), **bukan** `log.Fatal`. Contoh: Google OAuth,
  R2, Midtrans.
- Env wajib: `DATABASE_URL`, `ROOT_DOMAIN`.
- Handler: ambil user via `getAuthenticatedUser`, cek akses venue, validasi
  input, lalu query. Error → `http.Error` dengan pesan Indonesia.
- Semua query pakai `$1, $2, ...` (parameterized).

## Konvensi frontend

- Env frontend dideklarasikan di `frontend/src/env.ts` (`defineEnvVars`).
  - `PUBLIC_ROOT_DOMAIN`/`PUBLIC_ROOT_URL`: static (di-inline saat build →
    dikirim sebagai build arg dari compose).
  - `API_URL`, `SERVER_PUBLIC_IP`: private (server-only).
- `+page.server.ts`/`+server.ts` memanggil backend lewat `API_URL`
  (dev: `http://localhost:8080`, Docker: `http://backend:8080`). Jangan
  hardcode `localhost:8080`.
- Form action kirim `Cookie: session=...` ke backend.
- Notifikasi pakai `svelte-sonner` (`toast`).

## Endpoint utama (ringkas)

Lihat `docs/PROJECT-HANDOFF.md` §9 untuk lengkap.

```text
Auth     POST /api/auth/{register,login,logout}, GET /api/auth/me
Venue    POST|GET /api/venues, GET|PUT /api/venues/{slug}
Court    POST|GET /api/venues/{venueID}/courts
Hours    PUT|GET /api/venues/{venueID}/operating-hours
Domain   GET|POST /api/venues/{slug}/domains,
         DELETE /api/venues/{slug}/domains/{id},
         POST /api/venues/{slug}/domains/{id}/verify
Upload   POST /api/uploads/presign
Booking  POST|GET /api/bookings, GET /api/bookings/{id}
Bayar    POST /api/payments, GET /api/bookings/{id}/payment,
         POST /api/webhooks/midtrans
Public   GET /api/public/venues/{slug}, GET /api/public/domains/resolve
SSO      POST /api/sso/{authorize,exchange}
TLS ask  GET /api/internal/tls/allow/{token}?domain=...
```

## Database

- Migration di `backend/migrations/` (pasangan `.up.sql`/`.down.sql`).
- Versi terkini: **19**.
- Jalan otomatis via service `migrate` (dev: `make db`; prod: saat `up`).
- `down` (stop) **tidak** menghapus data; `down -v` (make `db-reset`) menghapusnya.

## Alur kerja yang diharapkan

1. Baca `docs/PROJECT-HANDOFF.md` (§14 Known Bugs, §17 Catatan Agent).
2. Jalankan `make check` dulu untuk baseline.
3. Buat perubahan kecil & terfokus. Ikuti konvensi di atas.
4. Verifikasi: `make check` (+ uji manual bila perlu).
5. Update dokumentasi bila mengubah arsitektur/endpoint/env.

## Known issues & roadmap

Lihat `docs/PROJECT-HANDOFF.md` §13 (belum dibuat), §14 (known bugs),
§15 (next work). Ringkas: beli domain otomatis, plan/subscription, admin
booking management, CI, cookie `Secure` di produksi.
