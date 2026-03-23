# Pictionary Backend

Real-time multiplayer Pictionary backend built with Go, HTTP, and WebSockets.

## Setup

```bash
go mod download
go run main.go
```

## Environment Variables

- `PORT`: HTTP server port (default: `8080`)

## HTTP API

### `POST /api/join`

Join matchmaking and receive player + room identifiers.

Request body:

```json
{
  "name": "Alice"
}
```

Response:

```json
{
  "playerId": "uuid",
  "roomId": "uuid",
  "wsUrl": "ws://localhost:8080/ws?playerId=...&roomId=..."
}
```

### `GET /ws?playerId=...&roomId=...`

Upgrade to WebSocket for real-time gameplay events.

### `GET /health`

Health check endpoint.

Response:

```json
{
  "status": "ok"
}
```

## WebSocket Events

All messages use:

```json
{
  "type": "EVENT_TYPE",
  "payload": {}
}
```

### Client -> Server

- `DRAW`
- `GUESS`
- `SELECT_WORD`
- `CLEAR_CANVAS`

### Server -> Client

- `ROOM_STATE`
- `DRAW_BROADCAST`
- `TURN_START`
- `WORD_OPTIONS`
- `CORRECT_GUESS`
- `ROUND_END`
- `GAME_END`
- `GUESS_MESSAGE`
