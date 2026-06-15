FROM node:22-alpine AS dev

WORKDIR /app

# corepack provisions the exact pnpm version pinned in package.json
# ("packageManager"), keeping the container consistent with local installs
ENV COREPACK_ENABLE_DOWNLOAD_PROMPT=0
RUN corepack enable

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile

COPY . .

CMD ["pnpm", "dev", "--host"]
