# Todo List Client

React + TypeScript + Vite frontend for the Todo List app. It connects to the Go API in `../server` for authentication and todo management.

## Requirements

- Node.js 20 or compatible
- npm
- API server running on `http://localhost:3000`

## Environment

Create `client/.env.local`:

```env
VITE_API_URL=http://localhost:3000
```

## Run

Install dependencies:

```bash
npm install
```

Start the Vite dev server:

```bash
npm run dev
```

The app runs on `http://localhost:5173` by default.

## Scripts

```bash
npm run dev
npm run lint
npm run build
npm run preview
```

## Features

- Register and login with email/password
- Persist auth token in local storage
- Restore the current session with `/api/auth/me`
- List, create, complete, and delete authenticated user todos
