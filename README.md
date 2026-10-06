# Lapanganku

SaaS manajemen venue dan booking lapangan. Backend memakai Go + PostgreSQL. Frontend memakai SvelteKit 3, Svelte 5, Tailwind CSS 4, dan shadcn-svelte.

Dokumentasi handoff lengkap untuk developer atau agent berikutnya ada di:

- [`docs/PROJECT-HANDOFF.md`](docs/PROJECT-HANDOFF.md)

## Model Produk

```text
User
└── Venue (bisnis/lokasi, subdomain/custom domain, member owner/admin)
    └── Court (lapangan, harga, durasi slot)
        └── Booking
```

Satu user dapat memiliki/mengelola banyak venue. Akses dikontrol lewat
`venue_members` (role `owner`/`admin`). Abstraksi workspace sudah dihapus.

Rencana limit plan per venue (belum diimplementasikan):

| Plan | Maks. lapangan per venue |
|---|---:|
| Basic | 3 |
| Pro | 10 |
| Business | 10 |

Limit plan dan billing belum diimplementasikan.

## URL Lokal

```text
http://lvh.me:5173                         landing/root platform
http://lvh.me:5173/admin                   dashboard admin
http://<venue-slug>.lvh.me:5173            storefront venue
http://localhost:8080                      Go API
```

`lvh.me` dan subdomainnya mengarah ke `127.0.0.1`, sehingga wildcard subdomain bisa dites tanpa mengubah `/etc/hosts`.

## Menjalankan Lokal

### 1. PostgreSQL

```bash
docker compose up -d
docker compose ps
```

### 2. Migration

Proyek memakai CLI [`golang-migrate`](https://github.com/golang-migrate/migrate). Dari folder `backend`:

```bash
/home/twizzcode/go/bin/migrate \
  -path migrations \
  -database "postgres://lapanganku:lapanganku@localhost:5432/lapanganku?sslmode=disable" \
  up
```

Cek versi:

```bash
/home/twizzcode/go/bin/migrate \
  -path migrations \
  -database "postgres://lapanganku:lapanganku@localhost:5432/lapanganku?sslmode=disable" \
  version
```

Versi terkini: `19`.

### 3. Backend

Bash/Zsh:

```bash
cd backend
export DATABASE_URL="postgres://lapanganku:lapanganku@localhost:5432/lapanganku?sslmode=disable"
export ROOT_DOMAIN="lvh.me"
go run .
```

Fish:

```fish
cd backend
set -x DATABASE_URL "postgres://lapanganku:lapanganku@localhost:5432/lapanganku?sslmode=disable"
set -x ROOT_DOMAIN "lvh.me"
go run .
```

### 4. Frontend

```bash
cd frontend
bun install
bun run dev --host 0.0.0.0
```

## Verifikasi

```bash
cd backend
go test ./...
```

```bash
cd frontend
bun run check
bun run build
```

Status terakhir:

- Backend test lulus.
- Frontend check: 0 error, 0 warning.
- Frontend production build lulus.
- `adapter-node` dipakai untuk deployment produksi (VPS + Caddy).

## Deployment Produksi

Panduan lengkap ada di [`docs/DEPLOY.md`](docs/DEPLOY.md). Ringkas:

- Semua service via Docker Compose (`caddy`, `frontend`, `backend`, `postgres`).
- Caddy reverse proxy + SSL otomatis (Cloudflare Origin Cert untuk platform,
  on-demand TLS Let's Encrypt untuk custom domain).
- Custom domain: owner arahkan **A record** ke IP VPS, klik Verifikasi.

```bash
cp .env.production.example .env   # isi nilainya
docker compose up -d --build
```

## Status Fitur

Sudah tersedia:

- Register/login email-password dengan Argon2id.
- Session DB-backed dan logout.
- Google OAuth (OIDC + PKCE) untuk login pusat.
- Venue (banyak venue per user) dengan member owner/admin.
- Court, harga, durasi slot, dan jam operasional.
- Storefront venue melalui subdomain platform.
- Availability berbasis timezone dan pencegahan bentrok booking di PostgreSQL.
- Booking customer, detail booking, dan riwayat booking.
- Central SSO root domain ke tenant subdomain memakai one-time code.
- Tenant session host-only dan local logout.
- Schema domain platform/custom dan resolver domain aktif.
- Custom domain: tambah, verifikasi (A record), hapus, routing, dan SSL otomatis
  (Caddy on-demand TLS + ask endpoint).
- Pembayaran booking via Midtrans QRIS (Core API) dengan webhook tervalidasi.

Belum tersedia:

- Beli domain otomatis (registrar eksternal).
- Plan, subscription, limit plan, dan billing.
- Invite/manage admin venue.
- Admin booking management.
- Landing page produk sebenarnya; `/` masih starter page.

Jangan commit `cookies.txt` atau credential/session lokal.
