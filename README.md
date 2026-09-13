# Trustyon

**Personal Developer Operating System**

Trustyon is a practical web application for managing projects, tasks, notes, and development activity in one place.

But the goal is bigger than a project manager.

> **GitHub shows what you coded. Trustyon shows how you are growing as a developer.**

Trustyon is being built as a personal development hub where coding activity becomes a visible collection of projects, languages, skills, milestones, and achievements.

## ✨ Vision

Trustyon aims to turn a developer's history into something that feels alive.

Instead of only looking at repositories, Trustyon will make it easy to answer:

- What am I working on?
- What have I built?
- Which languages am I actually using?
- What skills have I developed through real projects?
- How has my development journey changed over time?
- What should I work on next?

The long-term concept is a **Developer Collection**:

```text
Language
   ↓
Skill
   ↓
Project
   ↓
Achievement
   ↓
Developer Journey
```

## 🎯 Product Direction

Trustyon is intentionally different from GitHub.

| GitHub | Trustyon |
| --- | --- |
| Code hosting | Development hub |
| Repositories | Projects & goals |
| Commits | Development activity |
| Contributions | Growth history |
| Repository metadata | Personal collection |
| What I coded | Who I'm becoming |

The objective is not to replace GitHub. Trustyon sits **on top of your development activity** and turns it into a more motivating personal experience.

## 🧩 Current Features

- Project management
- Task management
- Notes
- Dashboard
- GitHub activity history
- GitHub repository synchronization
- GitHub OAuth authentication
- GitHub development statistics
- Language Collection
- Language levels
- Skill unlock UI
- Responsive web interface
- Go health endpoint
- PostgreSQL schema and migrations
- Automated tests and CI
- GitHub Pages demo deployment

## 🛠 Stack

- **Frontend:** React + TypeScript + Vite
- **Backend:** Go REST API
- **Database:** PostgreSQL
- **Infrastructure:** Docker Compose
- **CI/CD:** GitHub Actions
- **Integration:** GitHub OAuth + GitHub API

## 🏗 Architecture

```text
                 ┌─────────────────────┐
                 │   React / TypeScript │
                 │        Vite         │
                 └──────────┬──────────┘
                            │
                            ▼
                    ┌───────────────┐
                    │   REST API    │
                    └───────┬───────┘
                            │
                 ┌──────────┴──────────┐
                 │                     │
                 ▼                     ▼
          ┌─────────────┐      ┌──────────────┐
          │ Go Service  │      │ GitHub API   │
          └──────┬──────┘      └──────────────┘
                 │
                 ▼
          ┌─────────────┐
          │ PostgreSQL  │
          └─────────────┘
```

## 🗺 Roadmap

### Phase 0 — Foundation ✅

Build the base platform and make it runnable.

- [x] Repository structure
- [x] React + TypeScript + Vite frontend
- [x] Go REST API
- [x] PostgreSQL
- [x] Docker Compose
- [x] Health endpoint
- [x] Initial database schema
- [x] CI with GitHub Actions
- [x] GitHub Pages deployment

### Phase 1 — Core Workspace 🚧

Turn Trustyon into a useful daily development workspace.

- [x] Projects CRUD
- [x] Tasks CRUD
- [x] Notes CRUD
- [x] Dashboard
- [x] GitHub activity timeline
- [x] GitHub repository synchronization
- [ ] Better task filtering and sorting
- [ ] Project detail pages
- [ ] Activity search
- [ ] Dashboard customization

### Phase 2 — Developer Collection 🎮

Make development history feel like a collection rather than a database.

- [x] Language Collection
- [x] Language levels
- [x] Skill unlock UI
- [ ] Project achievement cards
- [ ] Achievement badges
- [ ] Language progression history
- [ ] Skill → Project relationships
- [ ] Developer profile / collection page
- [ ] Milestones and streaks
- [ ] Visual development timeline

### Phase 3 — GitHub Intelligence 🔗

Use GitHub as a source of development activity without trying to become another GitHub.

- [x] GitHub OAuth
- [x] Repository sync
- [x] Activity sync
- [x] Development statistics
- [ ] Better commit classification
- [ ] Pull request / issue insights
- [ ] Language usage analytics
- [ ] Project activity scoring
- [ ] Automatic skill suggestions
- [ ] Contribution heatmap inside Trustyon

### Phase 4 — Personal Development Engine 🧠

Help answer the question: **"What should I build next?"**

- [ ] Personal development goals
- [ ] Learning paths
- [ ] Skill gap detection
- [ ] Suggested next projects
- [ ] Project difficulty / experience levels
- [ ] Weekly development review
- [ ] Progress summaries
- [ ] Developer journal

### Phase 5 — Portfolio Mode 🌐

Turn the collection into something worth showing other people.

- [ ] Public developer profile
- [ ] Shareable project pages
- [ ] Public achievement collection
- [ ] Public development timeline
- [ ] Portfolio export
- [ ] Open Graph / social preview cards
- [ ] Custom profile URL

### Phase 6 — Production 🚀

Make Trustyon reliable enough to use continuously.

- [ ] Production authentication/session storage
- [ ] Fine-grained GitHub permissions
- [ ] API error handling improvements
- [ ] Comprehensive API tests
- [ ] Frontend component tests
- [ ] Database migration system
- [ ] Observability / logging improvements
- [ ] Production deployment
- [ ] Backup and recovery strategy
- [ ] Security review

## 🧪 Development Philosophy

Trustyon is itself a development experiment.

The project is intentionally built across multiple technologies so that the application can become a record of the developer learning those technologies.

For example:

```text
TypeScript → React UI → Trustyon Frontend
Go        → REST API → Trustyon Backend
Postgres  → Data     → Trustyon Storage
Docker    → Infra    → Trustyon Deployment
GitHub    → Activity → Developer Collection
```

Future technologies should be introduced only when they solve a real problem. The goal is **purposeful polyglot development**, not collecting languages for the sake of collecting languages.

## 🔐 GitHub OAuth Setup

Create an OAuth App in GitHub and configure:

- Authorization callback URL: `http://localhost:8080/api/auth/github/callback`
- `GITHUB_CLIENT_ID`
- `GITHUB_CLIENT_SECRET`
- Optional `GITHUB_REDIRECT_URI`
- Optional `TRUSTYON_WEB_URL` (defaults to `http://localhost:5173`)
- Optional `TRUSTYON_COOKIE_SECURE=true` when running behind HTTPS
- `GITHUB_USERNAME` is used for unauthenticated public activity

The OAuth flow currently requests `read:user repo` so authenticated activity can access repositories. The development implementation keeps the OAuth access token in an HttpOnly cookie. **Production deployment should move to server-side encrypted session storage and narrower GitHub permissions.**

## ▶️ Run Locally

```bash
docker compose up --build
```

Then open `http://localhost:5173` for the web app and connect GitHub from the dashboard.

For the frontend only:

```bash
cd apps/web
npm install
npm run dev
```

## 🌐 Live Demo

The current frontend demo is deployed with GitHub Pages:

**https://daikiito-dk.github.io/trustyon/**

The Pages deployment is a static demo. GitHub OAuth and the Go/PostgreSQL backend require a separately deployed API environment.

## 📁 Repository Structure

```text
trustyon/
├─ apps/
│  └─ web/              # React + Vite + TypeScript
├─ services/
│  └─ api/              # Go REST API
│     ├─ cmd/server/
│     ├─ internal/
│     │  ├─ handler/
│     │  ├─ service/
│     │  ├─ repository/
│     │  └─ model/
│     └─ migrations/
├─ docs/
├─ .github/workflows/
├─ docker-compose.yml
├─ Makefile
└─ README.md
```

## 📌 Status

🚧 **Early development — v0.1**

Trustyon is actively evolving from a CRUD-based development dashboard into a **personal developer growth and collection platform**.

The roadmap is intentionally iterative. Features may move between phases as the product is used and real development needs become clearer.

---

**Trustyon** — *Build projects. Collect skills. See your journey.*
