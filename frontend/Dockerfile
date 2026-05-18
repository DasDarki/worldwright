FROM oven/bun:1-alpine AS deps

WORKDIR /app

COPY package.json bun.lockb* package-lock.json* ./
RUN bun install --no-progress

FROM oven/bun:1-alpine AS build

WORKDIR /app
ENV NODE_ENV=production

COPY --from=deps /app/node_modules ./node_modules
COPY . .

RUN bun run build

FROM oven/bun:1-alpine AS runtime

WORKDIR /app

ENV NODE_ENV=production \
    NITRO_PORT=3000 \
    HOST=0.0.0.0

COPY --from=build /app/.output ./.output

EXPOSE 3000

CMD ["bun", "--bun", "run", ".output/server/index.mjs"]
