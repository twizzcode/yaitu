# Deployment Lapanganku (VPS + Docker + Caddy)

Panduan deploy produksi: VPS sendiri, semua service via Docker Compose,
Caddy sebagai reverse proxy dengan SSL otomatis.

## Ringkas: pertama kali deploy

Di server (setelah Docker terpasang & repo sudah di-clone):

```bash
cp .env.production.example .env      # lalu isi nilainya
mkdir -p certs                        # taruh origin.crt & origin.key di sini
docker compose up -d --build          # build + start semua service
docker compose ps                     # pastikan semua "Up"
```

> **Kamu tidak perlu build manual.** `docker compose up -d --build` build
> image langsung di server. Tidak ada langkah `bun run build` / `go build`
> manual.
>
> **Migration jalan otomatis.** Service `migrate` dijalankan sebelum backend
> (backend menunggu `service_completed_successfully`). Tidak perlu langkah
> migration manual.

Update berikutnya cukup:

```bash
git pull
docker compose up -d --build
```

## Arsitektur

```
Platform (lapanganku.id, *.lapanganku.id)      Custom domain (lapangan-a.com)
   Browser → Cloudflare (SSL gratis)             Browser
                  │ Origin Cert, Full(strict)        │ A record
                  ▼                                   ▼
              ┌────────────────────────────────────────────┐
              │  Caddy (:80/:443)                          │
              │  • Origin Cert untuk platform              │
              │  • on-demand TLS custom domain + ask        │
              └───────────────────┬────────────────────────┘
                                  ▼
                      frontend (adapter-node :3000)
                      backend (Go :8080) → postgres
```

- **Subdomain platform** (`<slug>.lapanganku.id`): SSL diurus Cloudflare (proxy).
  Caddy memakai Cloudflare **Origin Certificate** (mode Full strict).
- **Custom domain** (`lapangan-a.com`): Caddy terbitkan sertifikat Let's Encrypt
  otomatis via **on-demand TLS**, dengan endpoint **ask** ke backend agar hanya
  domain terdaftar yang boleh dapat sertifikat.

## Prasyarat

- VPS dengan **IP publik statis**
- Docker + Docker Compose plugin
- Port **80** dan **443** terbuka
- Domain `lapanganku.id` DNS di Cloudflare

## Sizing VPS

Pemakaian memori idle seluruh stack ±100–150 MB (frontend Node ~19 MB,
backend Go ~30 MB, Caddy ~40 MB, Postgres ~24 MB). Jadi **2 vCPU / 4 GB RAM
sudah aman** untuk menjalankan aplikasi pada skala awal.

Catatan penting: **build di server** (`--build`) adalah puncak pemakaian
(Vite/SvelteKit bisa spike ~1–1.5 GB RAM). Agar aman:

1. Tambah **swap 2 GB** (pengaman murah, mencegah OOM saat build):

   ```bash
   sudo fallocate -l 2G /swapfile && sudo chmod 600 /swapfile
   sudo mkswap /swapfile && sudo swapon /swapfile
   echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
   ```

2. Tuning Postgres (default terlalu konservatif bila data tumbuh):

   ```
   shared_buffers = 512MB
   effective_cache_size = 2GB
   work_mem = 16MB
   maintenance_work_mem = 128MB
   ```

3. Hindari build saat traffic ramai. Setelah image stabil, restart tanpa build:
   `docker compose up -d`.

Naikkan ke 4 vCPU / 8 GB bila: traffic tinggi, sering build di server, atau
data/booking sudah besar.

## 1. DNS di Cloudflare

Buat record berikut (proxy **ON / orange** untuk platform):

| Type | Name | Content | Proxy |
|---|---|---|---|
| A | `lapanganku.id` | `<IP VPS>` | Proxied |
| A | `*.lapanganku.id` | `<IP VPS>` | Proxied |

Set **SSL/TLS mode** ke **Full (strict)**.

## 2. Cloudflare Origin Certificate

1. Dashboard Cloudflare → **SSL/TLS → Origin Server → Create Certificate**
2. Hostnames: `lapanganku.id` dan `*.lapanganku.id`
3. Simpan hasilnya ke repo:
   - `certs/origin.crt`
   - `certs/origin.key`

Folder `certs/` sudah masuk `.gitignore` (berisi private key — jangan commit).

## 3. Konfigurasi environment

**Semua environment variable cukup diisi di SATU file: `.env` di root repo.**

```bash
cp .env.production.example .env
# edit .env: isi password DB, IP VPS, kredensial Google/R2, ACME_EMAIL,
# SAAS_ASK_TOKEN, dan kredensial Midtrans
```

Compose membaca `.env` ini dan mendistribusikannya ke semua service. Jadi kamu
**tidak perlu** mengisi env terpisah di `frontend/.env` atau `backend/`.

Catatan penting:

- `ROOT_DOMAIN` dan `PUBLIC_ROOT_URL` di-**inline saat build** frontend
  (SvelteKit static env). Compose sudah mengirimkannya sebagai **build arg**,
  jadi tetap cukup dari `.env` — tapi bila kamu mengubahnya, wajib
  `docker compose build frontend` ulang (bukan sekadar `restart`).
- `API_URL` default `http://backend:8080` (nama service Docker). Biasanya tidak
  perlu diubah.

Buat token ask:

```bash
openssl rand -hex 32   # salin hasilnya ke SAAS_ASK_TOKEN
```

### Midtrans

1. Dashboard Midtrans → **Settings → Access Keys** → salin Server & Client Key.
2. Isi `MIDTRANS_SERVER_KEY`, `MIDTRANS_CLIENT_KEY` di `.env`.
3. `MIDTRANS_IS_PRODUCTION=false` untuk sandbox, `true` untuk produksi.
4. Dashboard Midtrans → **Settings → Configuration → Payment Notification URL**:
   `https://<ROOT_DOMAIN>/api/webhooks/midtrans`
5. Untuk produksi, **Core API perlu diaktivasi** (ajukan permintaan ke Midtrans).

## 4. Jalankan (build + start)

> **Penting:** di server kamu **TIDAK perlu build manual**. Perintah
> `docker compose up -d --build` sudah otomatis build image **di server**
> (Docker yang build, bukan kamu). Jadi cukup satu perintah.

```bash
docker compose up -d --build
docker compose ps
```

Cek log:

```bash
docker compose logs -f caddy
docker compose logs -f backend
```

## 5. Migration database

Migration **jalan otomatis** saat `docker compose up` (service `migrate`
dijalankan sebelum backend). Untuk menjalankannya manual:

```bash
docker compose run --rm migrate
```

Cek versi:

```bash
docker compose run --rm migrate version
```

## 6. Verifikasi

- `https://lapanganku.id` → landing
- `https://<slug>.lapanganku.id` → storefront venue
- Tambah custom domain di `/admin/venues/<slug>/domains`, ikuti instruksi A record,
  klik Verifikasi → status `active`, lalu akses `https://lapangan-a.com`.

## Alur Custom Domain

1. Owner input domain di halaman **Domain**.
2. Sistem deteksi **apex** (`lapangan-a.com`) atau **subdomain**
   (`booking.lapangan-a.com`).
3. Owner buat record **A** → IP VPS (ditampilkan di UI).
4. Owner klik **Verifikasi** → backend cek A record.
5. Request pertama masuk → Caddy panggil endpoint **ask** → backend izinkan →
   sertifikat Let's Encrypt terbit otomatis.

### Kenapa A record (bukan CNAME)?

Caddy yang menerbitkan SSL, jadi customer cukup arahkan **A record ke IP VPS**.
Ini juga menyelesaikan masalah **apex domain** yang tidak bisa memakai CNAME.

### Endpoint ask

`GET /api/internal/tls/allow/{SAAS_ASK_TOKEN}?domain=<hostname>`

- Dipanggil Caddy sebelum menerbitkan sertifikat.
- Mengizinkan hanya hostname yang ada di `venue_domains` (status `pending`/`active`).
- Mencegah abuse penerbitan sertifikat.

## Update aplikasi

```bash
git pull
docker compose up -d --build
docker compose run --rm migrate   # bila ada migration baru
```

Catatan: bila mengubah `ROOT_DOMAIN`/`PUBLIC_ROOT_URL` di `.env`, wajib build
ulang frontend (sudah tercakup di `--build`).

## Troubleshooting

- **Sertifikat custom domain gagal**: cek `docker compose logs caddy`. Pastikan
  A record sudah mengarah ke IP VPS dan `SAAS_ASK_TOKEN` sama di Caddy & backend.
- **Platform tidak bisa diakses**: pastikan record Cloudflare proxied dan SSL
  mode Full (strict) serta Origin Certificate terpasang di `certs/`.
- **Ask endpoint 403**: domain belum terdaftar di `venue_domains`, atau token salah.
