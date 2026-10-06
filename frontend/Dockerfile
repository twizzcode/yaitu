# ---------- Build stage ----------
FROM oven/bun:1 AS build

WORKDIR /app

# PUBLIC_* bersifat "static" di SvelteKit: nilainya di-inline saat build.
# Jadi harus tersedia sebagai build arg, bukan hanya runtime env.
ARG PUBLIC_ROOT_DOMAIN=lapanganku.id
ARG PUBLIC_ROOT_URL=https://lapanganku.id
ENV PUBLIC_ROOT_DOMAIN=$PUBLIC_ROOT_DOMAIN
ENV PUBLIC_ROOT_URL=$PUBLIC_ROOT_URL

# Install dependency dulu agar layer cache efisien.
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

# Copy sumber dan build.
COPY . .
RUN bun run build

# ---------- Runtime stage ----------
FROM node:22-alpine AS runtime

WORKDIR /app

ENV NODE_ENV=production
ENV PORT=3000
ENV HOST=0.0.0.0

# adapter-node menghasilkan folder build/ yang sudah self-contained
# (dependency ikut dibundel), jadi tidak perlu npm install lagi.
COPY --from=build /app/build ./build

EXPOSE 3000

CMD ["node", "build"]
