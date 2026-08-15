---
name: docker-patterns
description: Guides Dockerfile and Docker Compose design for local development and image builds — multi-stage builds, networking, volumes, and container security — primarily for Go services with TypeScript/PHP notes. Use when writing or reviewing a Dockerfile, setting up Docker Compose for local development, troubleshooting container networking or volume issues, or hardening a container image. For rollout strategy and health-probe design once an image is built, see deployment-patterns; for the pipeline that builds and pushes it, see ci-cd-and-automation.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Docker Patterns

This skill covers building the image and running it locally. What to do
with that image in production (rollout strategy, health probes,
rollback) is `deployment-patterns`; how it gets built and pushed in
automation is `ci-cd-and-automation`.

## Multi-stage Dockerfile (Go)

Go's static binary output makes a minimal final image straightforward —
this should be the default shape for any Go service in this
organization:

```dockerfile
# Stage 1: build
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

# Stage 2: minimal runtime — no shell, no package manager, no Go toolchain
FROM gcr.io/distroless/static-debian12:nonroot AS runner
COPY --from=builder /server /server
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/server"]
```

`distroless/static` has no shell — this is a deliberate security choice
(smaller attack surface, nothing for an attacker to `exec` into even with
a code-execution bug), not an oversight. If you genuinely need shell
debugging in the runtime image, use `distroless/static-debian12:debug`
(adds `busybox`) rather than reverting to a full Alpine/Debian base for
every service.

Copying `go.mod`/`go.sum` before the rest of the source lets Docker
cache the `go mod download` layer across builds where only application
code changed — reordering `COPY . .` before dependency install defeats
this caching and slows every rebuild.

## Multi-stage Dockerfile (TypeScript/Node)

```dockerfile
FROM node:22-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci

FROM node:22-alpine AS build
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build && npm prune --production

FROM node:22-alpine AS production
WORKDIR /app
RUN addgroup -g 1001 -S appgroup && adduser -S appuser -u 1001 -G appgroup
USER appuser
COPY --from=build --chown=appuser:appgroup /app/dist ./dist
COPY --from=build --chown=appuser:appgroup /app/node_modules ./node_modules
COPY --from=build --chown=appuser:appgroup /app/package.json ./
ENV NODE_ENV=production
EXPOSE 3000
CMD ["node", "dist/server.js"]
```

## Multi-stage Dockerfile (PHP/Laravel)

```dockerfile
FROM composer:2 AS vendor
WORKDIR /app
COPY composer.json composer.lock ./
RUN composer install --no-dev --no-scripts --no-autoloader --prefer-dist

FROM php:8.3-fpm-alpine AS runner
WORKDIR /var/www
RUN docker-php-ext-install pdo pdo_mysql opcache
COPY --from=vendor /app/vendor ./vendor
COPY . .
RUN composer dump-autoload --optimize --no-dev
RUN chown -R www-data:www-data /var/www
USER www-data
EXPOSE 9000
CMD ["php-fpm"]
```

## Docker Compose for local development

```yaml
services:
  app:
    build:
      context: .
      target: dev
    ports: ["8080:8080"]
    volumes:
      - .:/app
    environment:
      - DATABASE_URL=postgres://postgres:postgres@db:5432/app_dev
    depends_on:
      db:
        condition: service_healthy

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: app_dev
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 3s
      retries: 5

volumes:
  pgdata:
```

`depends_on: condition: service_healthy` (not just `depends_on: [db]`)
waits for Postgres to actually accept connections, not just for the
container process to start — a plain `depends_on` only orders container
*start*, and the app container can boot and try to connect before
Postgres is ready to accept queries.

## Networking

Services on the same Compose network resolve each other by service name
(`db`, `redis`) — no manual IP or `/etc/hosts` management needed.
Segment networks so a database is unreachable from a network segment
that doesn't need it:

```yaml
services:
  api:
    networks: [frontend-net, backend-net]
  db:
    networks: [backend-net]   # unreachable from frontend-net
networks:
  frontend-net:
  backend-net:
```

Bind exposed ports to localhost during development
(`"127.0.0.1:5432:5432"`), and omit `ports:` entirely for anything that
should only be reachable within the Docker network in a
production-adjacent compose file.

## Volumes

- **Named volume** (`pgdata:`) — Docker-managed, persists across
  container restarts; use for anything stateful (databases, caches).
- **Bind mount** (`.:/app`) — maps a host directory in; use for source
  code during development to enable hot reload.
- **Anonymous volume** (`/app/node_modules`) — protects a
  container-generated directory from being shadowed by a bind mount
  covering its parent.

## Container security

```dockerfile
FROM golang:1.23.4-alpine3.20   # 1. pin an exact version, never :latest
# ... build ...
FROM gcr.io/distroless/static-debian12:nonroot  # 2. no shell, minimal surface
USER nonroot:nonroot            # 3. never run as root
```

```yaml
services:
  app:
    security_opt: ["no-new-privileges:true"]
    read_only: true              # root filesystem is immutable
    tmpfs: ["/tmp"]               # writable scratch space only where declared
    cap_drop: ["ALL"]
```

**Secrets**: never bake into an image layer (even a `RUN rm secret.txt`
leaves it recoverable in an earlier layer). Inject via environment
variables from a non-committed `.env`, or Docker/orchestrator secrets —
see `ci-cd-and-automation` for how CI secrets specifically should be
scoped separately from production secrets.

## .dockerignore

```
.git
.env
.env.*
node_modules
dist
coverage
*.log
Dockerfile*
docker-compose*.yml
tests/
```

Missing this means `COPY . .` ships `.git` history, `.env` files, and
test fixtures into the image — inflating size and, worse, potentially
baking secrets from a local `.env` into a shipped layer.

## Debugging commands

```bash
docker compose logs -f app                 # follow logs
docker compose exec app sh                 # shell into a running container
docker compose exec app nslookup db        # DNS resolution check
docker network inspect <project>_default   # inspect network membership
docker compose down -v                     # stop + remove volumes (DESTRUCTIVE — confirm before running)
```

## Gotchas

- `depends_on` without a `condition` only waits for the dependency
  container to *start*, not for the service inside it to be *ready* —
  this is the most common cause of "works on rebuild, fails on first
  `up`" flakiness in Compose setups.
- A distroless/scratch final image has no shell — `docker compose exec
  app sh` will fail with "executable file not found." This is
  intentional hardening, not a bug; use the `:debug` variant temporarily,
  or add ephemeral debug tooling via `docker run --rm -it
  --pid=container:<id> nicolaka/netshoot` instead of permanently adding a
  shell to the production image.
- `COPY . .` before installing dependencies busts the dependency-layer
  cache on every source change, turning a 2-second incremental rebuild
  into a full reinstall every time.
- Anonymous volumes (`/app/node_modules`) silently mask whatever the
  image built at that path with whatever was there from a previous
  container run — `docker compose down -v` (not just `down`) is required
  to actually clear it out when dependencies change.

## Real-world grounding

Google's `distroless` base images (open-sourced 2017) formalized the
minimal-attack-surface argument used above: a production container
should contain only the application and its runtime dependencies, not a
package manager, shell, or coreutils an attacker could use post-exploit
— this is the direct justification for defaulting Go services to
`distroless/static` rather than `alpine`, which still carries a full
shell and package manager most services never need at runtime.

## Verification

- [ ] Multi-stage build separates build tools from the runtime image
- [ ] Runtime image pins an exact version tag, never `:latest`
- [ ] Container runs as a non-root user
- [ ] `depends_on` uses `condition: service_healthy` where startup order matters, not a bare dependency list
- [ ] `.dockerignore` excludes `.git`, `.env`, and test fixtures
- [ ] No secret is baked into an image layer, even one later removed in a subsequent layer
- [ ] Database/internal ports are not exposed beyond localhost or the Docker network unless required
