# AGENTS.md file for smolgit

## Overview
smolgit is a minimalist Git server written in Go with a React/TS frontend. It exposes a REST API over HTTP and an SSH interface.

## Directory layout
- `cmd/` – main entry point.
- `pkg/` – core libraries: ssh, git, config, model, route, etc.
- `frontend/` – React/TS UI, built into `dist/` and embedded via `go:embed`.
- `test/` – Bats integration tests.
- `Makefile` – build & tooling targets.
- `.github/workflows/` – CI.

## Build & Run

| Target | Description | Command |
|--------|-------------|---------|
| `build` | Compile Go binary into `bin/smolgit` | `make build` |
| `build-frontend` | Build React assets into `dist/` | `make build-frontend` |
| `build-docker` | Build Docker image | `make build-docker` |
| `run` | Run compiled binary locally | `make run` |
| `run-docker` | Run smolgit in a container | `make run-docker` |
| `config` | Generate default config to stdout | `make config` |
| `config-docker` | Generate config for Docker | `make config-docker` |
| `clean` | Remove build artifacts | `make clean` |
| `lint` | Run golangci-lint on Go code | `make lint` |
| `test` | Run Go unit tests | `make test` |
| `integration-test` | Run Bats integration tests | `make integration-test` |
| `help` | Show available make targets | `make help` |

### Docker usage
```sh
# Build image
make build-docker

# Generate config
make config-docker

# Run
docker run -it -p 3080:3080 -p 3081:3081 -v $(pwd)/:/etc/smolgit smolgit
```

## Configuration

Configuration is supplied as YAML. The default location is `./config.yaml` or `$HOME/.config/smolgit/config.yaml`. It can be generated with:

```sh
./bin/smolgit config > config.yaml
```

A minimal example:

```yaml
log:
  color: true
  json: false
  level: DEBUG
server:
  disabled: false
  jwt_key: "your-secret-key"
  auth_disabled: false
  addr: ":3080"
  brand: "smolgit"
  auth:
    enabled: false
    accounts:
      - login: user
        password: pass
ssh:
  addr: ":3081"
git:
  path: /tmp/smolgit
  base: "git@my-git-server.lan"
  users:
    - name: "bob"
      permissions: "*"
      keys:
        - ssh-rsa AAAAB3...
```

* `server.jwt_key` is mandatory – without it the server will refuse to start.
* `git.path` can point to a local `git` binary; if omitted the binary defaults to `/usr/local/bin/git`.
* `git.users` controls SSH access; each user has a list of public keys and a permissions string.

## Frontend

The UI is built with Vite + React + TypeScript. Development can be run in isolation:

```sh
cd frontend
pnpm install
pnpm dev   # hot‑reload server at http://localhost:5173
```

Production build (used by the binary) is performed by `make build-frontend`. The resulting `dist/` folder is embedded into the Go binary via `go:embed`.

Linting the frontend code uses ESLint:

```sh
cd frontend
pnpm run lint
```

## Testing

### Go unit tests
```sh
make test
```

### Bats integration tests
```sh
make integration-test
```

The Bats tests start a temporary smolgit instance, exercise the CLI, and verify HTTP responses.

### CI

GitHub Actions runs:
* `make build`
* `make lint`
* `make config`
* `make integration-test`

See `.github/workflows/test.yml` for details.

## Naming & Coding Conventions

| Area | Convention |
|------|------------|
| Go packages | lowercase, no underscores |
| Go identifiers | camelCase for variables, PascalCase for exported types |
| Go structs | Tags use `koanf` style |
| Frontend | PascalCase components, camelCase props |
| Paths | Use relative paths in imports (e.g. `smolgit/pkg/git`) |

## Gotchas

1. **Config path** – if you run `make run` the binary looks for `./config.yaml` or `$HOME/.config/smolgit/config.yaml`. If you use Docker, mount the config into `/etc/smolgit`.
2. **JWT key** – the server refuses to start without `server.jwt_key`. Generate a random string or use `openssl rand -base64 32`.
3. **Git base** – `git.base` should match the SSH prefix used by clients, e.g. `git@my-git-server.lan`. Otherwise `git clone` will fail.
4. **Frontend build** – the binary embeds the built assets. Do **not** copy the `dist` folder manually; always run `make build-frontend` after changes.

## Contribution

Follow the existing style. Run:

```sh
make lint
make test
```

before pushing. Use the PR title format `[smolgit] <Title>`.

## Commit Messages
Follow Conventional Commits:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `test`: Adding missing tests
- `coauth`: Co-authored changes (e.g., "Co-authored-by: Qwen3 <email>")
- Example: `fix(core): format New function to pass gofumpt (coauth: Qwen3)`

## Useful Commands

```sh
# Show Go build info
go list -m -u all

# View lint output
make lint

# View test coverage
go test -cover ./...

# Run frontend dev server
cd frontend && pnpm dev
```
