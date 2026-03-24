# Pictionary (Realtime Multiplayer MVP)

## Project Structure

- `pictionary-backend/` - Go WebSocket + REST backend
- `pictionary-frontend/` - React + TypeScript + Vite frontend

## Local Development

### Backend

1. Open terminal in `pictionary-backend/`
2. Copy `.env.example` to `.env` (optional for local defaults)
3. Run:

```bash
go run main.go
```

Backend runs on `http://localhost:8080` by default.

### Frontend

1. Open terminal in `pictionary-frontend/`
2. Copy `.env.example` to `.env` (optional)
3. Install and run:

```bash
npm install
npm run dev
```

Frontend runs on `http://localhost:5173`.

## Environment Variables

### Backend (`pictionary-backend/.env`)

- `PORT` - HTTP server port (default `8080`)
- `ALLOWED_ORIGINS` - Comma-separated CORS origins (example: `http://localhost:5173`)

### Frontend (`pictionary-frontend/.env`)

- `VITE_API_URL` - Backend HTTP origin for API proxy/deploy
- `VITE_WS_URL` - Backend WebSocket origin

## Deployment

### Backend on Railway

1. Create a Railway project from `pictionary-backend/`
2. Railway uses `Dockerfile` for container build
3. Set env vars:
   - `PORT=8080`
   - `ALLOWED_ORIGINS=https://your-frontend.vercel.app`
4. Deploy and copy backend URL (for example `https://your-backend.railway.app`)

### Frontend on Vercel

1. Import `pictionary-frontend/` as a Vercel project
2. Keep `vercel.json` rewrite to forward `/api/*` to Railway backend
3. Set env vars:
   - `VITE_API_URL=https://your-backend.railway.app`
   - `VITE_WS_URL=wss://your-backend.railway.app`
4. Deploy
