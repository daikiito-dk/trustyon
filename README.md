# Trustyon

Personal Developer Operating System.

Trustyon is a practical web application for managing projects, tasks, notes, and development activity in one place.

## Stack

- Frontend: React + TypeScript + Vite
- Backend: Go REST API
- Database: PostgreSQL
- Infrastructure: Docker Compose
- CI/CD: GitHub Actions
- Integration: GitHub OAuth + GitHub API

## MVP

- Project management
- Task management
- Notes
- Dashboard
- GitHub activity history
- GitHub OAuth authentication
- Health endpoint
- PostgreSQL migrations
- Automated tests and CI

## Architecture

```text
React / TypeScript
        │
        ▼
     REST API
        │
        ├──── GitHub OAuth / API
        │
        ▼
    Go service
        │
        ▼
   PostgreSQL
```

## GitHub OAuth setup

Create an OAuth App in GitHub and configure:

- Authorization callback URL: `http://localhost:8080/api/auth/github/callback`
- `GITHUB_CLIENT_ID`
- `GITHUB_CLIENT_SECRET`
- Optional `GITHUB_REDIRECT_URI`
- Optional `TRUSTYON_WEB_URL` (defaults to `http://localhost:5173`)
- Optional `TRUSTYON_COOKIE_SECURE=true` when running behind HTTPS
- `GITHUB_USERNAME` is used for unauthenticated public activity

The OAuth flow requests `read:user repo` so the authenticated activity endpoint can access the user's repositories. The current development implementation keeps the OAuth access token in an HttpOnly cookie; production deployment should move to server-side encrypted session storage before exposing the app publicly.

## Run

```bash
docker compose up --build
```

Then open `http://localhost:5173` for the web app and connect GitHub from the dashboard.

## Status

🚧 Early development — v0.1
