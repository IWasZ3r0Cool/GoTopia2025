# Architecture

GoTopia is a server-authoritative multiplayer game with one in-memory game room. The current scope is deliberately small: up to four connected players, building placement, a turn clock, and building income.

## Runtime boundaries

- `internal/game` owns all game rules and mutable state. Its public read API returns deep snapshots, never live maps or pointers.
- `internal/server` owns HTTP, WebSocket clients, message validation, and fan-out. The hub serializes client membership and player joins.
- `cmd/server` only composes the engine, hub, HTTP server, environment configuration, and graceful shutdown.
- `frontend/src/store` owns the browser/server boundary. Components dispatch typed intents and render the last authoritative state; they do not simulate game rules.
- `frontend/src/components` contains semantic, responsive UI components. The main island is interactive; opponent boards are read-only summaries.

## State flow

1. The browser opens `/ws` and sends `JOIN_GAME` with a display name.
2. The hub allocates an opaque player ID. The browser receives `WELCOME`, followed by `GAME_STATE`.
3. The browser sends a typed action such as `BUILD`.
4. The engine validates and applies the action while holding its lock.
5. The engine emits a coalesced update signal. The hub encodes one deep snapshot and broadcasts it.
6. Every client replaces its Redux game state with that snapshot.

Engine update signals are hints, not state. Coalescing them is safe because every broadcast contains the complete latest state.

## Concurrency invariants

- Only the engine mutates game state.
- Every mutation uses the engine write lock.
- Every external read uses `Snapshot`, which deep-copies players, islands, and building maps.
- Only the hub goroutine mutates the client set or allocates player IDs.
- Each WebSocket has one reader and one writer goroutine. Every JSON message is written in its own WebSocket frame.
- Shutdown is context-driven and idempotent.

## Known product boundaries

The server currently has no persistence, authentication, reconnection token, lobby/matchmaking layer, boats, combat, disasters, scoring, or end-game condition. Those are product features, not implied by the current protocol. Add them behind engine methods and protocol messages rather than mutating Redux optimistically.
