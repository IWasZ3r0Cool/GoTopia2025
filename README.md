# GoTopia

GoTopia is an open-source, real-time island strategy game inspired by the 1983 Intellivision game *Utopia*. A Go server owns the game state and synchronizes up to four React clients over WebSockets.

The playable core currently supports joining a shared game, receiving a personal island, placing six building types, seeing other connected rulers, collecting farm/factory income, and growing population into housing capacity on 60-second turns.

## Requirements

- Go 1.25.4 or newer
- Node.js 24 or newer
- npm 11 or newer

## Run locally

Install the browser dependencies once:

```sh
make install
```

Start the Go server in one terminal:

```sh
make dev-server
```

Start Vite in another:

```sh
make dev-web
```

Open the URL printed by Vite. Its `/ws` development proxy connects to the Go server at `localhost:8080`.

## Verify and build

```sh
make check
```

`make check` runs Go formatting validation, `go vet`, TypeScript checking, ESLint, Go tests under the race detector, browser unit tests, and both production builds.

To run the production build:

```sh
make build
./bin/gotopia
```

The server serves `frontend/dist` as a single-page application and exposes `GET /healthz` for health checks.

## Configuration

| Variable | Default | Purpose |
| --- | --- | --- |
| `PORT` | `8080` | HTTP and WebSocket listen port |
| `GOTOPIA_STATIC_DIR` | `frontend/dist` | Built frontend directory |
| `GOTOPIA_ALLOWED_ORIGINS` | same-host and loopback origins | Comma-separated additional WebSocket origins |
| `VITE_WS_URL` | same-origin `/ws` | Optional browser WebSocket override at build time |

## Repository map

- `cmd/server`: executable composition and graceful shutdown
- `internal/game`: game rules, state, lifecycle, and unit tests
- `internal/server`: WebSocket/HTTP transport and protocol integration tests
- `frontend/src/store`: typed Redux state and WebSocket boundary
- `frontend/src/components`: accessible game interface and component tests
- `ARCHITECTURE.md`: runtime boundaries, data flow, and concurrency invariants

## Current limits

Games are held in memory and there is one room per server process. Disconnecting removes a player and their island. Persistence, reconnection, matchmaking, boats, combat, disasters, scoring, and victory conditions remain future features.
