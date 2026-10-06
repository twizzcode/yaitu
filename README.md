# Lapanganku

SaaS manajemen venue dan booking lapangan. Backend memakai Go + PostgreSQL. Frontend memakai SvelteKit 3, Svelte 5, Tailwind CSS 4, dan shadcn-svelte.

> **Untuk AI agent / developer baru:** baca [`AGENTS.md`](AGENTS.md) dulu.
> Dokumentasi handoff lengkap ada di [`docs/PROJECT-HANDOFF.md`](docs/PROJECT-HANDOFF.md).

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

Pengembangan memakai PostgreSQL via Docker, sedangkan backend & frontend
dijalankan native (hot reload). Prod memakai Docker Compose penuh — lihat
[`docs/DEPLOY.md`](docs/DEPLOY.md).

### Ringkas (Makefile)

```bash
make db        # start PostgreSQL dev + migration (otomatis)
make backend   # terminal 1: jalankan Go API (:8080), baca .env.local
make frontend  # terminal 2: jalankan SvelteKit dev (:5173)
```

`make help` untuk daftar lengkap perintah.

### Manual

#### 1. PostgreSQL + Migration

```bash
docker compose -f compose.dev.yaml up -d
```

Service `migrate` ikut jalan otomatis, jadi tidak perlu install CLI migrate.

#### 2. Backend

Backend membaca variabel dari `.env.local` (salin dari `.env.local.example`):

```bash
cp .env.local.example .env.local   # sekali saja
cd backend
set -a && source ../.env.local && set +a
go run .
```

Fish:

```fish
cd backend
for line in (cat ../.env.local | grep -v '^#' | grep '=')
    set -x (echo $line | cut -d= -f1) (echo $line | cut -d= -f2-)
end
go run .
```

#### 3. Frontend

```bash
cd frontend
bun install
bun run dev --host 0.0.0.0
```

Frontend dev memakai default `lvh.me` / `localhost:8080` (lihat
`frontend/src/env.ts`), jadi tidak butuh file env tambahan.

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

- Semua service via Docker Compose (`caddy`, `frontend`, `backend`, `postgres`, `migrate`).
- Caddy reverse proxy + SSL otomatis (Cloudflare Origin Cert untuk platform,
  on-demand TLS Let's Encrypt untuk custom domain).
- Custom domain: owner arahkan **A record** ke IP VPS, klik Verifikasi.

```bash
cp .env.production.example .env   # isi nilainya
docker compose up -d --build
```

> Dev memakai `compose.dev.yaml` + `.env.local`; prod memakai `compose.yaml` +
> `.env`. Keduanya terpisah agar rahasia produksi tidak terbaca saat dev.

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
