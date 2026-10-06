# Project Handoff: Lapanganku

Dokumen ini menjadi sumber konteks utama untuk developer atau agent berikutnya. Baca sebelum mengubah kode. Kondisi yang dicatat sesuai workspace pada 3 Oktober 2026.

> Ringkasan singkat untuk agent ada di [`AGENTS.md`](../AGENTS.md) (root repo).

## Current Stopping Point

Abstraksi **workspace sudah dihapus**. Model sekarang langsung `User -> Venue`
dengan tabel `venue_members` (role `owner`/`admin`). Satu user dapat mengelola
banyak venue.

Perubahan terakhir:

- Migration ditulis ulang dari awal: migration `3` = `create_venues` (tanpa
  `workspace_id`), migration `4` = `create_venue_members`. Tabel `workspaces`
  dan `workspace_members` sudah tidak ada.
- Backend: `workspace.go` dihapus. Endpoint venue menjadi
  `POST /api/venues` dan `GET /api/venues` (tanpa `workspaceID`).
- Frontend: folder route `(admin)/admin/workspaces` dihapus. Route admin
  sekarang langsung `(admin)/admin/venues/[venueSlug]/...`. Halaman
  `(admin)/admin/venues/new` untuk membuat venue.
- `/admin` redirect ke venue pertama user, atau ke `/admin/venues/new` bila
  belum punya venue.
- Halaman fungsional `courts` dan `operating-hours` dipindah ke struktur
  `[venueSlug]` (sebelumnya hanya ada di struktur `workspaces/[workspaceId]`).

Migrasi dari database versi workspace lama: karena migration ditulis ulang,
DB harus di-reset (`docker compose down -v` lalu `migrate up`). Tidak ada
migrasi data otomatis dari skema workspace lama.

Tugas berikut yang disepakati: **custom domain management** sudah dikerjakan
(API + UI + routing + SSL). Midtrans QRIS sudah diimplementasikan. Sisa
roadmap: **beli domain otomatis** (rencana via registrar eksternal).

Sebelum melanjutkan, perbaikan kecil paling penting adalah stale
`pending_payment` pada availability; detail ada di bagian Known Bugs.

## 1. Tujuan Produk

Lapanganku adalah SaaS multi-tenant untuk pengelolaan venue dan booking lapangan.

Ada dua sisi pengguna:

1. **Owner/admin venue** mengelola venue, lapangan, harga, jam operasional, booking, domain, tim, plan, dan billing.
2. **Customer** membuka storefront venue, memilih lapangan/tanggal/slot, login terpusat, membuat booking, membayar, dan melihat riwayat booking.

Model domain:

```text
User
└── Venue
    ├── Members (owner/admin)
    ├── Domains
    ├── Operating Hours
    └── Courts
        └── Bookings
```

Keputusan produk:

- Satu user dapat memiliki/mengelola banyak venue.
- Akses venue dikontrol lewat `venue_members` (owner/admin).
- Plan (kelak) menempel pada venue.
- Basic/Pro: 10 court per venue; Business: 10 court per venue (angka limit belum final, billing belum dibuat).
- Customer wajib login untuk booking.
- Login terpusat berada di root domain.
- Tenant platform memakai `<venue-slug>.<root-domain>`.
- Custom apex dan custom subdomain direncanakan.

## 2. Arsitektur

```text
Browser
  |
  v
SvelteKit 3 / Svelte 5
  |
  v
Go net/http API :8080
  |
  v
PostgreSQL 18
```

Repo:

```text
lapanganku/
├── backend/       Go API dan migration SQL
├── frontend/      SvelteKit
├── docs/          dokumentasi handoff
├── compose.yaml   PostgreSQL lokal
└── README.md
```

Backend sengaja memakai standard library `net/http`, bukan framework. Database memakai `pgxpool`; tidak ada ORM.

## 3. URL dan Tenant Routing

Development:

```text
http://lvh.me:5173                         root platform
http://lvh.me:5173/admin                   admin
http://arena-futsal-bandung.lvh.me:5173    tenant storefront
http://localhost:8080                      API
```

Target production:

```text
https://domain.com                         landing/login/admin
https://<venue-slug>.domain.com            platform storefront
https://booking.venue.com                  custom subdomain
https://venue.com                          custom apex
```

`frontend/src/hooks.ts` (`reroute`) menangani subdomain platform **dan** custom
domain:

```text
<slug>.<root-domain> /            -> /venues/{slug}
<slug>.<root-domain> /courts/...  -> /venues/{slug}/courts/...
<custom-domain>      /            -> /venues/{slug}   (via resolver)
```

Custom domain di-resolve lewat `/api/resolve-domain` (proxy server-only ke
`GET /api/public/domains/resolve`) dan di-cache 60 detik. Endpoint `/api/*`
selalu di-skip agar tidak terjadi rekursi.

## 4. Stack

### Backend

- Go module: `lapanganku/backend`
- Go version pada `go.mod`: `1.27.1`
- HTTP: `net/http`
- PostgreSQL pool: `github.com/jackc/pgx/v5/pgxpool`
- Password: Argon2id dari `golang.org/x/crypto/argon2`
- Session/code token: 32 random bytes, base64url; SHA-256 disimpan DB
- Server address saat ini hardcoded `:8080`

### Frontend

- SvelteKit 3
- Svelte 5 runes mode
- TypeScript strict
- Vite 8
- Tailwind CSS 4
- shadcn-svelte / bits-ui
- Bun
- `adapter-node` dipakai untuk deployment VPS (di balik Caddy)

### Infrastruktur

- PostgreSQL 18 melalui Docker Compose
- Cloudflare R2 (S3-compatible) untuk upload KTP via presigned URL
- Deployment: VPS + Docker Compose + Caddy (reverse proxy + auto SSL)
- `adapter-node` untuk frontend; panduan di `docs/DEPLOY.md`
- **Dev vs Prod dipisah:**
  - `compose.dev.yaml` (dev): postgres + migrate, kredensial hardcoded, volume
    `postgres_dev_data` terpisah.
  - `compose.yaml` (prod): full stack (postgres, migrate, backend, frontend, caddy).
  - `.env.local` (dev, gitignored) untuk backend native; `.env` (prod, gitignored)
    untuk Docker Compose.
  - `Makefile` menyediakan `make db`, `make backend`, `make frontend`, `make check`.
- Service `migrate` jalan otomatis sebelum backend (prod) & saat `make db` (dev).
- `.dockerignore` ada di `frontend/` dan `backend/` (mencegah `node_modules`
  ikut ke build context)
- Belum ada CI

## 5. Environment Variables

Ada tiga sumber env, dipisah tegas:

| File | Untuk | Sumber |
|---|---|---|
| `.env.local` (gitignored) | Backend dev native | `.env.local.example` |
| `.env` (gitignored) | Prod Docker Compose | `.env.production.example` |
| default di `frontend/src/env.ts` | Frontend dev | — |

`compose.dev.yaml` **hardcode** kredensial dev (tidak membaca `.env`), sehingga
rahasia produksi tidak pernah tersentuh saat dev.

Backend wajib:

```env
DATABASE_URL=postgres://lapanganku:lapanganku@localhost:5432/lapanganku?sslmode=disable
ROOT_DOMAIN=lvh.me
```

`ROOT_DOMAIN` hanya hostname, tanpa scheme, port, atau wildcard.

Backend opsional (object storage KTP, Cloudflare R2):

```env
R2_ACCOUNT_ID=
R2_ACCESS_KEY_ID=
R2_SECRET_ACCESS_KEY=
R2_BUCKET=
R2_PUBLIC_URL=
```

Bila salah satu dari `R2_ACCOUNT_ID`, `R2_ACCESS_KEY_ID`, `R2_SECRET_ACCESS_KEY`,
atau `R2_BUCKET` kosong, fitur upload dinonaktifkan dan endpoint presign
mengembalikan HTTP 503. Endpoint R2: `https://<account>.r2.cloudflarestorage.com`
dengan region `auto` dan path-style. `R2_PUBLIC_URL` (alias `R2_PUBLIC_BASE_URL`)
dipakai untuk mengembalikan `public_url` objek.

**Penting — CORS bucket wajib diatur.** Upload memakai metode `PUT`, yang
selalu memicu preflight `OPTIONS`. Tanpa CORS policy di bucket R2, browser
memblokir request dan muncul error `Failed to fetch` (curl tetap bisa, karena
curl tidak menerapkan CORS). Set CORS policy minimal:

```json
[
  {
    "AllowedOrigins": ["http://lvh.me:5173"],
    "AllowedMethods": ["PUT", "GET", "HEAD"],
    "AllowedHeaders": ["*"],
    "ExposeHeaders": ["ETag"],
    "MaxAgeSeconds": 3600
  }
]
```

Di produksi, ganti `AllowedOrigins` dengan origin sebenarnya. CORS bisa diatur
dari dashboard Cloudflare R2 (bucket → Settings → CORS Policy) atau via
`PutBucketCors` (R2 mendukung S3 API ini).

Variabel frontend dideklarasikan di `frontend/src/env.ts` memakai SvelteKit 3
`defineEnvVars`, dengan default untuk pengembangan lokal. Tidak ada
`frontend/.env`; di produksi semua env cukup diisi di `.env` root repo (dibaca
Docker Compose). `PUBLIC_ROOT_DOMAIN`/`PUBLIC_ROOT_URL` bersifat static
(di-inline saat build) sehingga dikirim sebagai build arg dari compose.

**Backend origin:** frontend memakai private env `API_URL` (default
`http://localhost:8080` untuk dev, `http://backend:8080` di Docker Compose).
Semua `+page.server.ts`/`+server.ts` sudah memakai `API_URL`.

## 6. Menjalankan Proyek

**Dev** dan **prod** dipisah:
- Dev: `compose.dev.yaml` (hanya postgres + migrate) + `.env.local`, backend &
  frontend native (hot reload).
- Prod: `compose.yaml` (full stack) + `.env`.

### Dev (ringkas)

```bash
make db        # start PostgreSQL dev + migration otomatis
make backend   # terminal 1: Go API :8080 (baca .env.local)
make frontend  # terminal 2: SvelteKit dev :5173
```

Atau manual:

```bash
cp .env.local.example .env.local
docker compose -f compose.dev.yaml up -d
cd backend && set -a && source ../.env.local && set +a && go run .
cd frontend && bun install && bun run dev --host 0.0.0.0
```

Versi schema terkini: `19`, `dirty=false`. Migration di dev berjalan otomatis
via service `migrate` di `compose.dev.yaml` (tidak perlu CLI).

### Checks

```bash
make check      # backend test + frontend check
```

Atau manual:

```bash
cd backend && go test ./...
cd frontend && bun run check && bun run build
```

Status terakhir: semuanya lulus. Frontend memakai `adapter-node` untuk produksi.

## 7. Database Schema dan Migration

Semua migration ada di `backend/migrations/` dan memakai pasangan `.up.sql` / `.down.sql`.

| Versi | Migration | Isi utama |
|---:|---|---|
| 1 | `create_users` | User, email unik, nullable password hash, email verification timestamp |
| 2 | `create_sessions` | Hashed session token, user, expiry |
| 3 | `create_venues` | Venue, global slug, address, timezone, WhatsApp, active flag |
| 4 | `create_venue_members` | Membership venue dengan role owner/admin |
| 5 | `create_courts` | Court, sport, integer IDR price, slot duration, active flag |
| 6 | `create_operating_hours` | Jadwal Senin-Minggu per venue |
| 7 | `create_bookings` | Booking status, payment expiry, GiST overlap constraint |
| 8 | `create_venue_domains` | Platform/custom hostname, verification status/token hash |
| 9 | `create_auth_codes` | One-time SSO code hash, target host, return path, expiry |
| 10 | `add_session_metadata` | Session hostname dan family ID |
| 11 | `add_auth_code_family` | Family ID pada auth code |
| 12 | `unique_family_hostname_session` | Satu tenant session per family+hostname |
| 13 | `create_user_identities` | Identitas provider eksternal (Google) per user |
| 14 | `create_oauth_states` | State/nonce/PKCE verifier untuk OAuth pusat |
| 15 | `add_venue_ktp_key` | Kolom `ktp_key` (object key KTP di object storage) pada venues |
| 16 | `add_venue_owner_profile` | Kolom profil pemilik & wilayah: `owner_name`, `owner_nik`, `province`, `city`, `district`, `village`, `postal_code` |
| 17 | `add_venue_description_logo` | Kolom `description` dan `logo_key` pada venues |
| 18 | `add_domain_verification_token` | Kolom `verification_token` pada venue_domains (token verifikasi domain) |
| 19 | `create_payments` | Tabel payments: order_id, QRIS, status, idempotency, raw charge/notification |

Relasi:

```text
users
├── sessions
├── auth_codes
├── user_identities
├── venue_members ── venues
│                    ├── venue_domains
│                    ├── operating_hours
│                    └── courts ── bookings
└── bookings
```

`oauth_states` berdiri sendiri tanpa FK ke `users`.

Booking overlap dicegah pada DB dengan exclusion constraint terhadap rentang `[starts_at, ends_at)` untuk status `pending_payment` dan `confirmed`.

## 8. Backend Files

```text
main.go              startup, config, route registration, basic auth handlers
password.go          Argon2id hash/verify
password_test.go     satu unit test hashing password
session.go           session token dan authenticated session/user lookup
venue.go             create/list/get venue, venue_members, platform-domain transaction
court.go             create/list court dan venue access helper
operating_hours.go   get/upsert schedule
availability.go      public fixed-slot availability
booking.go           create/list/detail customer booking
public.go            public venue and active courts
domain.go            hostname normalization, platform-domain sync, resolver
domain_manage.go     list/add/delete/verify custom domain (A record + ask TLS)
sso.go               one-time-code authorize/exchange
google_auth.go       Google OIDC + PKCE login pusat
midtrans.go          pembayaran QRIS (Core API charge + webhook + status)
storage.go           Cloudflare R2 presigned URL untuk upload KTP
```

## 9. API Endpoints

### Health

```text
GET /health
```

### Auth

```text
POST /api/auth/register
POST /api/auth/login
GET  /api/auth/me
POST /api/auth/logout
```

Password login membuat session tujuh hari. Session cookie `HttpOnly`, `SameSite=Lax`; backend saat ini masih hardcode `Secure=false`.

### SSO

```text
POST /api/sso/authorize
POST /api/sso/exchange
```

- Authorize hanya menerima central session dengan `session.hostname == ROOT_DOMAIN`.
- Code berlaku 60 detik dan hanya dapat digunakan sekali.
- Exchange memvalidasi target domain masih aktif.
- Tenant session memakai `family_id` sama dengan central session.
- Satu row tenant session per `(family_id, hostname)`.

### Venue

```text
POST /api/venues
GET  /api/venues
GET  /api/venues/{venueSlug}
PUT  /api/venues/{venueSlug}
```

- `POST /api/venues` membuat row venue, membership owner (`venue_members`),
  dan platform domain secara transaction. Tidak ada lagi `workspaceID`.
  Menerima data inti (nama, slug, alamat, timezone, WhatsApp, `ktp_key`) dan
  data pemilik/wilayah (`owner_name`, `owner_nik`, `province`, `city`,
  `district`, `village`, `postal_code`).
- `GET /api/venues` mengembalikan semua venue tempat user terdaftar sebagai
  member, lengkap dengan `role`, `ktp_key`, `ktp_url`, dan data pemilik/wilayah.
- `GET /api/venues/{venueSlug}` mengembalikan venue milik member aktif.
- `PUT /api/venues/{venueSlug}` memperbarui profil venue (nama, alamat,
  WhatsApp, `ktp_key`, `logo_key`, `description`, data pemilik/wilayah).
  `ktp_key`/`logo_key` kosong berarti tidak mengubah objek yang sudah ada.
  `slug` tidak dapat diubah lewat endpoint ini.

### Domain

```text
GET    /api/venues/{venueSlug}/domains
POST   /api/venues/{venueSlug}/domains
DELETE /api/venues/{venueSlug}/domains/{domainID}
POST   /api/venues/{venueSlug}/domains/{domainID}/verify
GET    /api/internal/tls/allow/{token}?domain=...
```

- Wajib session + akses venue (owner/admin).
- `GET` mengembalikan semua domain venue (platform dulu, lalu custom), plus
  `is_apex` (hostname 2 label = apex).
- `POST` menambah domain custom dengan status `pending`. Domain platform (di
  bawah `ROOT_DOMAIN`) ditolak.
- `verify` mengecek **A record** hostname via `net.Resolver.LookupHost`,
  dicocokkan dengan `SERVER_PUBLIC_IP`. Sukses → `active`; gagal → `failed` +
  HTTP 422. (Mekanisme DNS TXT lama sudah diganti.)
- `DELETE` hanya menghapus domain custom.
- `GET /api/internal/tls/allow/{token}` adalah **endpoint ask** untuk Caddy
  on-demand TLS. Mengizinkan hanya hostname yang terdaftar (`pending`/`active`).
  Dilindungi `SAAS_ASK_TOKEN`.
- SSL: platform pakai Cloudflare (Origin Certificate, Full strict). Custom
  domain pakai Caddy on-demand TLS (Let's Encrypt).
- UI: `(admin)/admin/venues/[venueSlug]/(operasional)/domains`, bergaya daftar
  domain (mirip Vercel): search, tambah, baris expandable, salin record DNS
  (A record → IP VPS), verifikasi, hapus.
- Routing custom domain: `frontend/src/hooks.ts` (`reroute`) memanggil
  `/api/resolve-domain` (proxy server-only) untuk memetakan hostname → venue.

### Upload

```text
POST /api/uploads/presign
```

- Wajib session.
- Body: `{ "content_type": "image/png", "kind": "ktp" }`.
- Hanya menerima `image/jpeg`, `image/png`, `image/webp`, `image/heic`, `image/heif`.
- Mengembalikan `{ upload_url, object_key, public_url, expires_in }`; URL berlaku 15 menit.
- Object key berformat `ktp/<user_id>/<random>.<ext>`.
- `ktp_key` pada create venue divalidasi harus berprefix `ktp/<user_id>/` milik
  user yang login.
- Frontend memakai proxy `POST /api/uploads/presign` (SvelteKit) karena cookie
  session bersifat HttpOnly.

**Alur upload KTP di frontend** (`register-venue-form.svelte` + `lib/upload.ts`):

1. Saat user memilih file, gambar **dikonversi ke WebP di browser** memakai
   canvas (`convertToWebp`), bukan langsung diunggah.
2. Upload ke R2 baru dilakukan **saat klik "Daftar Sekarang"**. Jadi bila user
   batal atau menutup tab, tidak ada file orphan di storage.
3. Setelah upload sukses, `object_key` dikirim sebagai `ktp_key` pada
   `POST /api/venues`.
4. `PUT` ke R2 memicu preflight CORS; bucket wajib punya CORS policy (lihat
   bagian Environment Variables).

Catatan: tidak ada auto-cleanup file orphan di storage. Karena upload terjadi
saat submit, risikonya kecil, tetapi kegagalan setelah upload dan sebelum venue
terbuat (mis. error DB) masih menyisakan objek. Kalau perlu, tambahkan lifecycle
rule di R2 untuk menghapus objek lama di prefix `ktp/`.

### Court

```text
POST /api/venues/{venueID}/courts
GET  /api/venues/{venueID}/courts
```

### Operating Hours

```text
PUT /api/venues/{venueID}/operating-hours
GET /api/venues/{venueID}/operating-hours
```

Hari memakai ISO: `1=Senin` sampai `7=Minggu`.

### Public

```text
GET /api/public/venues/{slug}
GET /api/public/domains/resolve?hostname=...
GET /api/courts/{courtID}/availability?date=YYYY-MM-DD
```

Availability:

- Publik.
- Maksimal 30 hari ke depan.
- Timezone-aware.
- Slot mengikuti `slot_duration_minutes` dari jam buka.
- Booking `pending_payment`/`confirmed` dianggap memblokir.

### Booking Customer

```text
POST /api/bookings
GET  /api/bookings
GET  /api/bookings/{bookingID}
```

Create booking:

- Wajib session.
- Harga dihitung backend.
- Satu slot tetap per booking.
- Status awal `pending_payment`.
- Hold pembayaran 15 menit.
- Konflik overlap menjadi HTTP 409.

### Pembayaran (Midtrans QRIS)

```text
POST /api/payments
GET  /api/bookings/{bookingID}/payment
POST /api/webhooks/midtrans
```

- `POST /api/payments` membuat transaksi QRIS via Midtrans Core API
  (`POST /v2/charge`, `payment_type=qris`). Wajib session + pemilik booking.
  Idempotent: bila sudah ada payment `pending` untuk booking, dikembalikan
  yang ada tanpa membuat transaksi baru. Booking harus `pending_payment` dan
  belum kedaluwarsa.
- `GET /api/bookings/{bookingID}/payment` mengembalikan payment terbaru
  (dipakai polling UI).
- `POST /api/webhooks/midtrans` menerima notifikasi Midtrans. **Signature
  diverifikasi** dengan `sha512(order_id + status_code + gross_amount +
  server_key)`. Status dipetakan: `settlement`/`capture` → booking
  `confirmed`; `expire`/`cancel`/`deny` → booking `expired`. Order tak dikenal
  dibalas 200 (idempotent, tanpa crash).
- UI: halaman `/bookings/[bookingId]` menampilkan QR, tombol "Bayar dengan
  QRIS", dan polling tiap 4 detik sampai status berubah.
- Env: `MIDTRANS_SERVER_KEY`, `MIDTRANS_CLIENT_KEY`, `MIDTRANS_IS_PRODUCTION`.
  Bila server key kosong, fitur pembayaran dinonaktifkan (HTTP 503).
- Notification URL di dashboard Midtrans: `https://<root>/api/webhooks/midtrans`.

## 10. Frontend Routes

### Root dan Auth

```text
/                       masih starter SvelteKit, landing page belum dibuat
/register               register email/password
/login                  central login
/sso/authorize          bridge central SSO
/auth/login             mulai login tenant
/auth/callback          exchange code dan set tenant cookie
/auth/logout            local tenant logout
```

### Admin

```text
/admin                                        redirect ke venue pertama / /admin/venues/new
/admin/venues/new                             buat venue
/admin/venues/[venueSlug]                     dashboard venue
/admin/venues/[venueSlug]/courts              daftar lapangan
/admin/venues/[venueSlug]/courts/new          tambah lapangan
/admin/venues/[venueSlug]/operating-hours     jam operasional
/admin/venues/[venueSlug]/finance             placeholder
/admin/venues/[venueSlug]/analytics           placeholder
/admin/venues/[venueSlug]/articles            placeholder
/admin/venues/[venueSlug]/gallery             placeholder
/admin/venues/[venueSlug]/profile             form profil (belum tersambung backend)
/admin/venues/[venueSlug]/domains             placeholder
/admin/venues/[venueSlug]/staff               placeholder
/admin/venues/[venueSlug]/customers           placeholder
```

Struktur route `(admin)/admin/workspaces/...` sudah dihapus. Semua route admin
kini berada di bawah `venues/[venueSlug]`.

**Known incomplete:** halaman `finance`, `analytics`, `articles`, `gallery`,
`domains`, `staff`, `customers` masih placeholder. Halaman `profile` hanya
memvalidasi input di klien; endpoint update profil belum ada. Halaman
`venues/[venueSlug]/+page.svelte` (dashboard) memakai komponen chart dengan
data contoh hardcoded.

### Storefront Customer

Internal route:

```text
/venues/[slug]
/venues/[slug]/courts/[courtId]/book
```

Browser platform tenant:

```text
https://<slug>.domain.com/
https://<slug>.domain.com/courts/[courtId]/book
```

Customer account:

```text
/my-bookings
/bookings/[bookingId]
```

## 11. Session dan Central SSO

### Cookie Policy

Keputusan: semua cookie session host-only. Jangan set `Domain=.domain.com`.

```text
domain.com                         central/admin session
venue.domain.com                   tenant session
booking.custom-domain.com          tenant session
custom-apex.com                    tenant session
```

Custom domain tidak dapat berbagi cookie platform karena browser security model. Central SSO menyelesaikan UX login sekali dengan code exchange.

### Flow

```text
tenant /auth/login
→ buat state cookie host-only
→ root /sso/authorize
→ central login jika perlu
→ backend POST /api/sso/authorize
→ redirect tenant /auth/callback?code&state
→ state check
→ backend POST /api/sso/exchange
→ set tenant session cookie host-only
→ redirect return_path
```

Keamanan:

- State URL harus sama dengan state cookie.
- State cookie `HttpOnly`, `SameSite=Lax`, 5 menit.
- Auth code disimpan sebagai SHA-256, berlaku 60 detik, dihapus atomik saat exchange.
- Return path harus relative dan tidak diawali `//`.
- Target host harus domain aktif dan venue aktif.
- Tenant token tidak pernah masuk URL.

Logout tenant hanya mencabut tenant session; central session tetap aktif. `logout-all` belum dibuat.

## 12. Fitur yang Sudah Jalan

- Email/password register dan login.
- Argon2id password hashing.
- DB session dan current-user lookup.
- Central-to-tenant SSO pada platform subdomain.
- Tenant user indicator/header dan local logout.
- Google OAuth (OIDC + PKCE) untuk login pusat.
- Venue create/list/get dengan member owner/admin (banyak venue per user).
- Venue backend create/list dan platform domain generation.
- Court create/list.
- Operating hours get/update.
- Public venue storefront.
- Availability date/slot.
- Customer booking create/list/detail.
- Booking countdown UI.
- Domain resolver untuk domain aktif.
- Platform-domain startup backfill/sync.
- Custom domain management (list/add/verify/hapus) + routing + SSL otomatis.

## 13. Belum Dibuat

Prioritas produk:

1. Beli domain otomatis (rencana: registrar eksternal seperti Citra Host).
2. Plan/subscription/limit enforcement.
3. Admin booking management.
4. Team invite/member management (tambah admin ke venue).
5. Landing page SaaS sebenarnya.
6. CI (build/test otomatis).

Fitur lain belum dibuat:

- Forgot password.
- Email verification flow.
- Terms/privacy pages.
- Cancel/refund.
- Edit/deactivate venue/court.
- Global logout/session-family revoke.
- Registrar/domain purchase.

## 14. Known Bugs dan Risiko

1. **Stale pending booking — SUDAH DIPERBAIKI.** `availability.go` kini hanya
   menganggap `pending_payment` aktif bila `payment_expires_at > NOW()`.

2. **Cookie production belum konsisten aman.** Backend login dan frontend central login hardcode `Secure=false`. Gunakan environment-aware secure cookie sebelum deploy HTTPS.

3. **Backend API URL frontend — SUDAH DIPERBAIKI.** Semua `+page.server.ts` /
   `+server.ts` kini memakai private env `API_URL` (default dev
   `http://localhost:8080`, produksi `http://backend:8080` lewat compose).

4. **Custom-domain hook sudah dibuat.** `hooks.ts` meresolve custom domain
   lewat `/api/resolve-domain` dan merutekan ke storefront venue. Sisa: uji
   end-to-end di staging dengan domain HTTPS sungguhan.

5. **Midtrans QRIS sudah diimplementasikan.** Sisa: uji dengan kredensial
   sandbox asli, lalu aktifkan Core API di dashboard Midtrans untuk produksi
   (perlu request aktivasi ke Midtrans).

6. **Data profil venue sudah dipersist.** Form daftar venue di
   `/admin/venues/new` dan halaman profil `/admin/venues/[venueSlug]/profile`
   menyimpan `owner_name`, `owner_nik`, `ktp_key`, `logo_key`, `description`,
   dan detail wilayah (`province`, `city`, `district`, `village`,
   `postal_code`) lewat `POST`/`PUT /api/venues`. Yang belum: `slug` tidak
   dapat diubah setelah dibuat, dan venue yang dibuat sebelum migration 16/17
   punya kolom profil kosong (isi lewat halaman profil).

7. **Tidak ada auto-cleanup objek orphan di R2.** Upload KTP terjadi saat
   submit sehingga risikonya kecil, tetapi kegagalan setelah upload dan sebelum
   venue terbuat masih menyisakan objek. Pertimbangkan lifecycle rule R2 untuk
   prefix `ktp/`.

### Security/robustness

- Belum ada rate limit login/register/SSO.
- Belum ada CSRF token/origin checks untuk semua mutation.
- HTTP server belum punya read/write/idle timeout dan graceful shutdown.
- JSON decoding/body-size policy belum konsisten.
- Beberapa malformed UUID dapat menjadi 500.
- WhatsApp belum dinormalisasi/divalidasi.
- Multi-member per venue belum ada API invite; baru owner saat create venue.
- OAuth-only user kelak perlu penanganan nullable `password_hash` pada login.
- TLS custom domain belum dirancang untuk provider deployment.

### UI

- `/` masih starter page.
- Dashboard metric masih hardcoded nol.
- Admin menampilkan hostname `.lapanganku.com` hardcoded pada beberapa halaman, bukan env root domain.
- Jam operasional UI menulis `WIB` walau timezone venue bisa lain.
- Tombol Google menuju endpoint yang belum ada.

## 15. Next Recommended Work

Urutan aman:

1. Uji Midtrans QRIS dengan kredensial sandbox asli (charge + webhook end-to-end).
2. Uji end-to-end custom domain di staging HTTPS (DNS + Caddy + routing).
3. Beli domain otomatis (rencana: registrar eksternal seperti Citra Host, bukan
   Cloudflare Registrar).
4. Plan/subscription dan enforcement limit.

## 16. Testing Gap

Saat ini hanya ada `password_test.go`. Tambahkan test kecil untuk logic berisiko:

- `normalizeHostname`
- `normalizeReturnPath`
- one-time code dipakai satu kali
- target host mismatch
- expired SSO code
- slot boundary/timezone
- booking overlap concurrent
- stale payment hold
- signature webhook Midtrans
- owner/admin authorization
- custom-domain resolver

Frontend belum punya unit/E2E test. Minimal E2E nanti:

```text
central login
→ tenant SSO
→ choose court/date/slot
→ create booking
→ booking detail
→ my bookings
→ tenant logout
```

## 17. Catatan Agent Berikutnya

- Model sekarang `User -> Venue`; tidak ada workspace. Plan/billing (kelak)
  menempel di venue.
- Akses venue selalu lewat `venue_members` (owner/admin). Jangan percaya
  `venueID` dari klien tanpa cek keanggotaan.
- Jangan menyimpan token session atau auth code mentah di DB.
- Jangan menaruh session token di URL.
- Jangan set cookie `Domain=.domain.com`; custom domain tidak akan tercakup dan blast radius membesar.
- Jangan mempercayai harga dari frontend.
- Pertahankan DB overlap constraint untuk race safety.
- Jangan membuat custom domain aktif sebelum verifikasi DNS (A record).
- Endpoint ask TLS (`/api/internal/tls/allow`) wajib dijaga `SAAS_ASK_TOKEN`;
  jangan izinkan hostname yang tidak terdaftar (mencegah abuse sertifikat).
- Jangan menerima arbitrary `return_to`/external redirect.
- Jangan mengubah route admin menjadi tenant subdomain; admin tetap root-domain path.
- `frontend/src/hooks.ts` (`reroute`) tidak boleh me-rewrite `/api/*` agar
  tidak rekursi.
- Sebelum coding, jalankan checks dan baca bagian known bugs di atas.
