*This project has been created as part of the 42 curriculum by awaegaer, macaruan, gbakulin, mafioron and aule-bre.*

# Transcendence: 67 Survivor

## Description

**Transcendence: 67 Survivor** is a real-time web application that hosts a *Vampire Survivor*-style game, wrapped around a full social platform with accounts, friendships and live chat.

Players sign up, log in, and dive into a survival arena where waves of ranged enemies close in from every side. The goal is simple: stay alive as long as possible. The game can be played **solo (1 player)** or in **local co-op (2 players)** on the same machine, and every run feeds a personal progression system (XP, levels, stats, leaderboard).

Around the game sits a complete user-management layer: editable profiles with avatars, a friends system with live online/offline status, private real-time messaging, match history and a friends leaderboard. The whole stack is containerized and shipped with production-grade DevOps tooling (centralized logging, monitoring, alerting, analytics dashboards and automated backups).

### Key features

- **Survival game** (Go + Ebitengine compiled to WebAssembly) — solo and local co-op, enemy waves, weapons, real-time scoring and timer.
- **Account system** — email/password sign-up and login, JWT authentication, hashed & salted passwords.
- **Profiles** — editable username/email, avatar upload with client-side and server-side validation, personal statistics.
- **Friends** — send/accept/reject friend requests, friends list with real-time online status.
- **Live chat** — private real-time messaging over WebSockets, persisted history, anti-spam cooldown.
- **Progression** — XP and level system, per-user statistics (best score, best time, total time), friends leaderboard, match history.
- **Public API** — secured, rate-limited and documented REST API to interact with the database.
- **Internationalization** — full UI available in French, English and Russian with a live language switcher.
- **Cross-browser support** — works on Chrome and additional browsers.
- **Legal pages** — accessible Terms of Service & Privacy Policy.
- **Observability** — centralized logs (ELK), metrics dashboards (Prometheus/Grafana), an advanced analytics dashboard, alerting, and automated backups with disaster-recovery procedures.

---

## Instructions

### Prerequisites

Before running the project, make sure you have the following installed:

- **Docker** (or **Podman**) with the **Docker Compose** plugin
- **GNU Make**
- **OpenSSL** (used to generate the Elasticsearch/Kibana secrets and certificates)
- A POSIX shell (`bash`/`sh`)

> The backend container automatically installs and runs `mkcert` to generate the local HTTPS certificates, so you do **not** need it on your host machine.

### Environment configuration

Credentials and secrets are stored in a local `.env` file that is **ignored by Git**. An example file is provided at the root of the repository: `.env.exemple`.

The Elasticsearch/Kibana passwords are generated automatically by the bootstrap step (`docker/start_elastic.sh`), which writes them into `docker/.env`. For the remaining variables, copy the example file and adjust the values:

```bash
cp .env.exemple docker/.env
```

Variables used by the project:

| Variable | Description |
|---|---|
| `POSTGRES_USER` | PostgreSQL user |
| `POSTGRES_PASSWORD` | PostgreSQL password (min. 16 characters) |
| `POSTGRES_DB` | PostgreSQL database name |
| `JWT_SECRET` | Secret key used to sign JWT tokens |
| `ELASTIC_PASSWORD` | Elasticsearch `elastic` user password (auto-generated) |
| `KIBANA_SYSTEM_PASSWORD` | Kibana system user password (auto-generated) |

### Running the project

The entire stack starts with a **single command** from the repository root:

```bash
make
```

This runs, in order:
1. `make bootstrap` — generates the SSL certificates and bootstraps Elasticsearch (snapshot repository, ILM policy, index template, Kibana password).
2. `make build` — builds the application images.
3. `make up` — starts all services in detached mode.

### Opening the game

> ⚠️ **Do NOT use `localhost` anymore.** Type `hostname -I` in your terminal, copy your IP address, and replace `localhost` in the URL with it.

```bash
hostname -I
```

Then open your browser at:

```
https://<your-ip>:5173/
```

For example:

```
https://10.171.55.106:5173/
```

> Your browser will warn about the self-signed certificate the first time — this is expected for a local development setup. The application is compatible with the latest stable version of Google Chrome.

### Useful Make targets

| Command | Description |
|---|---|
| `make` / `make all` | Bootstrap, build and start everything |
| `make up` | Start the services |
| `make stop` | Stop the services |
| `make ps` | List running services |
| `make build` | Build the images |
| `make clean` | Remove backend/frontend images |
| `make fclean` | Stop and clean images |
| `make death` | Full reset (volumes, certs, backups, prune) |
| `make re` | Full clean rebuild |

### Service ports

| Service | Port |
|---|---|
| Frontend (Vite) | `5173` |
| Backend (Go/Gin) | `8081` |
| PostgreSQL | `5432` |
| Prometheus | `9090` |
| Grafana | `3000` |
| Elasticsearch | `9200` |
| Kibana | `5601` |

---

## Team Information

| Login | Name | Role(s) | Responsibilities |
|---|---|---|---|
| mafioron | Maël | Product Owner + Developer (Backend) | Product vision, feature prioritization, backlog, validation of completed work; backend development |
| aule-bre | Aurèle | Project Manager / Scrum Master + Developer (Game) | Team coordination, meetings & planning, deadlines, risk management; game development |
| macaruan | Maati | Technical Lead / Architect + Developer (Frontend) | Technical architecture, technology-stack decisions, code quality & reviews; frontend development |
| awaegaer | Arthur | Developer (Frontend) | Frontend development |
| gbakulin | Galina | Developer (Database & DevOps) | Database design, containerization and DevOps infrastructure |

All five team members contributed as developers across both the mandatory part and the modules.

---

## Project Management

- **Version control:** the project is hosted on **GitHub**, using a **branch-per-person** workflow. The `main` branch is kept **clean and always functional**; each member develops on their own branch and merges into `main` once the work is reviewed and stable.
- **Communication:** day-to-day communication happens on **Discord**.
- **Key moments:** for important decisions (merges, planning the next set of actions, and project-level decisions), the team organized **in-person meetings at school**.
- **Work breakdown:** the project was split by domain — **frontend** (Maati & Arthur), **backend** (Maël), **game** (Aurèle), and **database & DevOps** (Galina) — with everyone contributing as a developer.

---

## Technical Stack

### Frontend
- **React 19 + TypeScript** — component-based UI, considered a framework in this project's context.
- **Vite** — development server (HMR) and build tool.
- **Tailwind CSS** — utility-first styling solution.
- **React Router** — client-side routing (Single Page Application).
- **i18next / react-i18next** — internationalization (French, English, Russian).

### Backend
- **Go (Golang)** — main backend language.
- **Gin** — backend web framework (routing, middleware, handlers).
- **GORM** — ORM for database access and migrations.
- **golang-jwt** — JWT generation and validation.
- **bcrypt** — password hashing and salting.
- **Gorilla WebSocket** — real-time chat and online presence.
- **zerolog** — structured logging.
- **Prometheus client** — custom application metrics (HTTP requests, unique visitors, active users, visit duration).

### Game
- **Go + Ebitengine**, compiled to **WebAssembly (WASM)** and embedded into the frontend through a `<canvas>` element using `wasm_exec.js`.

### Database
- **PostgreSQL 14**

> **Why PostgreSQL?** Our data is strongly relational — users, friendships, private messages and match records all reference each other — so a relational database was the natural fit. We chose PostgreSQL for its ACID guarantees and reliability under concurrent access (important for a multi-user real-time app), its first-class support in GORM, and the fact that it is free, open-source and battle-tested. It also integrates cleanly with our monitoring stack through `postgres-exporter`.

### DevOps / Infrastructure
- **Docker / Podman + Docker Compose** — containerization, single-command deployment.
- **Elasticsearch + Logstash + Kibana + Filebeat (ELK)** — centralized log management with ILM (retention/rollover) and snapshot/backup policies.
- **Prometheus + Grafana** — metrics collection, dashboards and alerting; with **node-exporter** and **postgres-exporter**.
- **Backup container** — scheduled (cron) backups of PostgreSQL, Prometheus and Elasticsearch data, with a documented disaster-recovery procedure.
- **HTTPS everywhere** — TLS for all external connections (frontend, backend, Grafana, Elasticsearch, Kibana) using locally generated certificates.

### Justification for major technical choices
- **Go + Gin + GORM** for a fast, statically-typed backend with a clean ORM layer.
- **React + Vite + Tailwind** for a modern, responsive and quickly iterable frontend.
- **WebAssembly** lets the game be written in Go (shared language with the backend) while running natively in the browser.
- **PostgreSQL** for a robust relational store (see above).

---

## Database Schema

The database uses four main tables with clearly defined relations.

### Tables

**`users`**
| Field | Type | Notes |
|---|---|---|
| `id` | uint | Primary key |
| `username` | string | Unique |
| `email` | string | Unique |
| `password_hash` | string | Not null (bcrypt) |
| `avatar_url` | string | Default placeholder avatar |

**`friendships`**
| Field | Type | Notes |
|---|---|---|
| `id` | uint | Primary key |
| `user_id` | uint | Indexed, FK → `users.id` (requester) |
| `friend_id` | uint | Indexed, FK → `users.id` (target) |
| `status` | varchar(20) | `pending` / `accepted` (default `pending`) |
| `created_at` | timestamp | |
| `updated_at` | timestamp | |

**`direct_messages`**
| Field | Type | Notes |
|---|---|---|
| `id` | uint | Primary key |
| `sender_id` | uint | Indexed, FK → `users.id` |
| `receiver_id` | uint | Indexed, FK → `users.id` |
| `content` | text | Not null |
| `created_at` | timestamp | |

**`matches`**
| Field | Type | Notes |
|---|---|---|
| `id` | uint | Primary key |
| `user_id` | uint | Indexed, FK → `users.id` |
| `duration` | int | Survival time in seconds, not null |
| `score` | int | Not null |
| `created_at` | timestamp | |

### Relations
- A **User** can have many **Friendships** (both as requester `user_id` and as target `friend_id`).
- A **User** can send and receive many **DirectMessages**.
- A **User** can have many **Matches** (one record per game run).

Migrations are handled automatically at startup via GORM's `AutoMigrate` on `User`, `Friendship`, `DirectMessage` and `Match`.

---

## Features List

| Feature | Description | Worked on by |
|---|---|---|
| Sign-up / Login | Email + password authentication, JWT auto-login, session restore | Maël |
| Password security | bcrypt hashing & salting | Maël |
| Profile management | Edit username/email, upload avatar (type/size/dimension validation, front + back) | Maël (back), Maati & Arthur (front) |
| User statistics | Best score, best survival time, total time played | Maël (back), Maati & Arthur (front) |
| XP & level system | XP earned per run, level progression with progress bar | Maël (back), Maati & Arthur (front) |
| Friends system | Send / accept / reject requests, friends list | Maël (back), Maati & Arthur (front) |
| Online presence | Real-time online/offline status via WebSocket | Maël |
| Private chat | Real-time messaging, persisted history, anti-spam cooldown, character limit | Maël (back), Maati & Arthur (front) |
| Match history | Personal list of past runs with date and duration | Maël (back), Maati & Arthur (front) |
| Friends leaderboard | Ranking of the player and their friends by survival record | Maël (back), Maati & Arthur (front) |
| Public API | Secured, rate-limited and documented REST endpoints | Maël |
| Survival game | Vampire-survivor-style game (Go/Ebitengine → WASM), solo & local co-op | Aurèle |
| Game over / retry flow | End-of-run screen, retry without reload, score saved to backend | Aurèle (game), Maati & Arthur (front) |
| Internationalization | Full UI in FR / EN / RU with live switcher | Maati & Arthur |
| Cross-browser support | Verified on Chrome plus additional browsers | Maati & Arthur |
| Terms of Service & Privacy | Accessible legal pages, linked from the footer | Maati & Arthur |
| Multi-user support | Concurrent users, real-time updates across clients | Maël (back), Galina (infra) |
| Centralized logging | ELK stack with Filebeat shipping app logs | Galina |
| Monitoring & alerting | Prometheus metrics, Grafana dashboards, alert rules | Galina |
| Analytics dashboard | Interactive charts and visualizations of platform/game data | Galina |
| Automated backups | Scheduled backups + disaster-recovery procedure | Galina |

---

## Modules

Point system: **Major = 2 pts**, **Minor = 1 pt**. The project must validate **at least 14 points** to pass; every module implemented beyond the required 14 points counts as **bonus**, with a maximum of **5 bonus points**.

### Mandatory modules — 14 points

| Category | Module | Type | Points | Worked on by |
|---|---|---|---|---|
| Web | Use a framework for **both** frontend (React) and backend (Gin) | Major | 2 | Maati & Arthur (front), Maël (back) |
| Web | Implement **real-time features** using WebSockets | Major | 2 | Maël |
| Web | **Allow users to interact** with other users (chat, profile, friends) | Major | 2 | Maël (back), Maati & Arthur (front) |
| Web | Use an **ORM** for the database (GORM) | Minor | 1 | Maël (back), Galina (database) |
| Accessibility & i18n | Support for **multiple languages** (FR, EN, RU = 3 languages) | Minor | 1 | Maati & Arthur |
| User Management | **Standard user management and authentication** | Major | 2 | Maël (back), Maati & Arthur (front) |
| DevOps | **Log management** with ELK (Elasticsearch, Logstash, Kibana) | Major | 2 | Galina |
| DevOps | **Monitoring system** with Prometheus and Grafana | Major | 2 | Galina |

**Mandatory subtotal: 14 points.** ✅

### Bonus modules — beyond 14 points (capped at 5 bonus points)

| Category | Module | Type | Points | Worked on by |
|---|---|---|---|---|
| Gaming & UX | **Complete web-based game** where users can play solo/duo | Major | 2 | Aurèle |
| Web | **Public API** with secured API key, rate limiting, documentation and 5+ endpoints | Major | 2 | Maël |
| Data & Analytics | **Advanced analytics dashboard** with data visualization | Major | 2 | Galina |
| User Management | **Game statistics and match history** (requires the game module) | Minor | 1 | Maël (back), Aurèle (game) |
| Accessibility & i18n | Support for **additional browsers** | Minor | 1 | Maati & Arthur |
| DevOps | **Health check & status page** with automated backups and disaster recovery | Minor | 1 | Galina |

**Bonus implemented: 9 points** — counted up to the subject's maximum of **5 bonus points**.

**Total claimed: 23 points** (14 mandatory + 9 bonus, of which a maximum of 5 bonus points are awarded).

### How the modules were implemented & why

- **Framework (front + back):** React/TypeScript on the frontend and Go/Gin on the backend. Chosen because they give a structured architecture, a rich ecosystem and fast iteration — exactly what a team project of this scope needs.
- **Real-time (WebSockets):** a Gorilla WebSocket hub broadcasts online presence and routes private chat messages. We needed instant updates (who is online, incoming messages) that polling could not deliver smoothly.
- **User interaction:** chat, profile pages and a full friends system. This is the social backbone of the platform and a prerequisite for the chat and leaderboard features.
- **Public API:** a secured, rate-limited and documented REST API exposing more than five endpoints across GET/POST/PUT/DELETE. It was added so the database can be queried programmatically in a safe, controlled way (API key + rate limiting to prevent abuse).
- **ORM (GORM):** removes hand-written SQL, handles migrations automatically and keeps the data layer maintainable as the schema evolves.
- **Multiple languages:** the whole UI is translatable (i18next) into French, English and Russian, with a live switcher — chosen to make the platform accessible to a wider audience.
- **Additional browsers:** beyond Chrome, the app was tested and fixed on additional browsers to guarantee a consistent experience regardless of the user's setup.
- **Standard user management & authentication:** secure sign-up/login, editable profiles, avatar upload and online status — the foundation every other feature relies on.
- **Game statistics & match history:** each run is stored and aggregated into stats and history, giving players a sense of progression and a reason to come back.
- **Web-based game (solo/duo):** the core experience. We deliberately built a **cooperative PvE survival game rather than a competitive PvP one**, because as a group this direction spoke to us more and was more motivating to build together. Marked as a bonus module.
- **ELK log management:** centralizes application logs for searching and analysis (Kibana), with retention and snapshot policies — essential for debugging a distributed, containerized app.
- **Prometheus & Grafana monitoring:** collects metrics and exposes dashboards with alert rules, so we can see the health of the backend, database and host at a glance.
- **Health check & backups:** container health checks, automated scheduled backups and a documented disaster-recovery procedure to protect data and detect outages.
- **Advanced analytics dashboard:** interactive charts and visualizations of platform and game data (real-time updates, customizable ranges) built on top of our metrics/logging infrastructure, turning raw data into actionable insight.

---

## Individual Contributions

- **Maël (mafioron) — Backend (Go / Gin).** Authentication and JWT, password security, friends system, the WebSocket hub for chat and online presence, match/statistics/level endpoints, the public REST API, profile updates and avatar handling, Prometheus metrics middleware.
  - *Challenge:* synchronizing online-presence state across reconnecting WebSocket clients without leaking stale connections. Solved by tracking clients per user ID in the hub and closing/replacing previous connections on re-registration.

- **Aurèle (aule-bre) — Game (Go / Ebitengine → WASM).** The survival game itself: world/entities, players and enemies, weapons and bullets, collisions, enemy waves and spawning, movement, scoring/timer, game-over and reset flow, and the WASM ↔ JavaScript bridge.
  - *Challenge:* preventing the game from launching twice when React re-mounted the component, which spawned duplicate canvases. Solved with a cancellation guard around the asynchronous WASM instantiation.

- **Maati (macaruan) — Frontend (React / TypeScript / Tailwind) & Technical Lead.** Application architecture and routing, pages and components, the user/WebSocket context, internationalization wiring, the Vite proxy/HTTPS setup, and overall code quality and reviews.
  - *Challenge:* getting the Vite dev server, the Go backend and the WebSocket connection to coexist over HTTPS without CORS issues. Solved with a Vite proxy that forwards `/api`, `/uploads` and `/ws` to the backend.

- **Arthur (awaegaer) — Frontend (React / TypeScript / Tailwind).** Pages (Profile, Chat, Match History, Terms of Service), reusable components (auth modal, language switcher), responsive styling, and the front-side form validation.
  - *Challenge:* implementing avatar upload with reliable client-side validation (type, size and dimensions) before sending the file to the backend. Solved by validating the image asynchronously via the `Image` API before submission.

- **Galina (gbakulin) — Database & DevOps / Containerization.** PostgreSQL setup and schema, Docker/Podman and `docker-compose` configuration, the full ELK stack, Prometheus/Grafana monitoring, alerting and the analytics dashboard, automated backups and disaster-recovery procedure, and SSL/TLS certificate generation.
  - *Challenge:* wiring TLS across every service (Elasticsearch, Kibana, Logstash, backend) with a shared CA while keeping the single-command startup working. Solved with a dedicated certificate-bootstrap step and a scripted Elasticsearch initialization.

Any other useful information (usage documentation, credits) is welcome.

---

## Resources

Documentation, references and tutorials that supported the creation of this project:

**Frontend**
- React — https://react.dev/
- TypeScript — https://www.typescriptlang.org/docs/
- Vite — https://vitejs.dev/
- Tailwind CSS — https://tailwindcss.com/docs
- React Router — https://reactrouter.com/
- i18next / react-i18next — https://www.i18next.com/ and https://react.i18next.com/

**Backend**
- Go — https://go.dev/doc/
- Gin — https://gin-gonic.com/docs/
- GORM — https://gorm.io/docs/
- golang-jwt — https://github.com/golang-jwt/jwt
- bcrypt (golang.org/x/crypto) — https://pkg.go.dev/golang.org/x/crypto/bcrypt
- Gorilla WebSocket — https://github.com/gorilla/websocket
- zerolog — https://github.com/rs/zerolog

**Game**
- Ebitengine — https://ebitengine.org/
- Go WebAssembly — https://github.com/golang/go/wiki/WebAssembly

**Database**
- PostgreSQL — https://www.postgresql.org/docs/

**DevOps / Infrastructure**
- Docker — https://docs.docker.com/
- Docker Compose — https://docs.docker.com/compose/
- Podman — https://docs.podman.io/
- Prometheus — https://prometheus.io/docs/
- Grafana — https://grafana.com/docs/
- Elastic Stack (ELK) — https://www.elastic.co/guide/index.html
- mkcert — https://github.com/FiloSottile/mkcert

### Use of AI

AI tools were used as assistants throughout the project for **debugging**, **code review / proofreading**, and **understanding concepts related to the subject**. All AI-assisted output was reviewed, tested and validated by the team before being integrated.

---

## Compliance with Mandatory Requirements

| Requirement | Status |
|---|---|
| Web app with **frontend + backend + database** | React (front), Go/Gin (back), PostgreSQL (db) |
| **Single-command** deployment via containers | `make` (Docker / Podman + Compose) |
| Runs on latest **Google Chrome**, no console errors | Verified |
| **Privacy Policy & Terms of Service** accessible | Linked from the footer; the legal page covers both ToS and data privacy |
| **Multi-user** support | Concurrent sessions, real-time updates across clients via WebSocket |
| Responsive, accessible **frontend** | Tailwind, tested on desktop and mobile/tablet |
| **CSS framework / styling solution** | Tailwind CSS |
| Credentials in `.env` (Git-ignored) + example file | `.env.exemple` provided; `.env` ignored by Git |
| Clear **database schema** with relations | See *Database Schema* above |
| Secure **auth** (email/password, hashed + salted) | bcrypt hashing & salting + JWT |
| **Form/input validation** on frontend **and** backend | Client-side (React) and server-side (Go) validation |
| **HTTPS** for all backend connections | TLS on every external connection |

---

## Notes & Known Limitations

- The game is **local** co-op (two players on the same keyboard: Player 1 on arrow keys, Player 2 on `WASD`); remote two-player gameplay over the network is not implemented.
- HTTPS uses **locally generated** certificates, so browsers will show a security warning on first use in development.
- Remember to open the app using your machine's IP (`hostname -I`) rather than `localhost`.
