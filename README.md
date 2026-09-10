# Task Manager – Full Stack Application Documentation

A full-stack task management application built with:

- Go + Gin (REST API)
- GORM
- MySQL 8.0 (Dockerized)
- JWT Authentication
- React + Vite frontend

---

# Features

## Authentication
- Signup
- Login
- Logout
- JWT authentication
- Protected routes

## Project Management
- Create projects
- View user projects
- Delete projects
- Project membership system
- Owner/member roles

## Task Management
- Create tasks
- Update tasks
- Delete tasks
- Assign/unassign users
- Task priorities
- Task statuses

## Permissions
- Role-based access control
- Project owner permissions
- Member permissions

## Frontend
- Multi-page React + Vite application (5 separate HTML pages)
- Tailwind CSS dark interface
- Token stored in localStorage, never shown or pasted by hand
- Guarded pages that redirect visitors who are not logged in
- Modals, inline validation and toasts instead of browser dialogs

---

# Tech Stack

## Backend
- Go
- Gin
- GORM
- MySQL
- JWT
- Docker

## Frontend
- React 19
- Vite (multi-page build)
- Tailwind CSS v4

---

# Project Structure

```txt
project-root/
│
├── controllers/
├── middleware/
├── models/
├── initializers/
├── main.go
├── docker-compose.yml
│
├── Dockerfile
│
├── frontend/
│   ├── index.html          # one HTML file per page
│   ├── login.html
│   ├── signup.html
│   ├── projects.html
│   ├── project.html
│   │
│   ├── src/
│   │   ├── entries/        # one React entry point per HTML page
│   │   ├── pages/          # one component per page
│   │   ├── components/     # shared components (ui/ holds the primitives)
│   │   ├── lib/            # api client, session, auth guards, helpers
│   │   └── styles/
│   │
│   ├── public/
│   ├── vite.config.js
│   └── package.json
│
└── README.md
````

---

# API Endpoints

## Authentication

| Method | Endpoint        | Description    |
| ------ | --------------- | -------------- |
| POST   | `/api/signup`   | Create account |
| POST   | `/api/login`    | Login          |
| POST   | `/api/logout`   | Logout         |
| GET    | `/api/validate` | Validate JWT   |

---

## Projects

| Method | Endpoint            | Description         |
| ------ | ------------------- | ------------------- |
| GET    | `/api/projects`     | Get user projects   |
| POST   | `/api/projects`     | Create project      |
| GET    | `/api/projects/:id` | Get project details |
| DELETE | `/api/projects/:id` | Delete project      |

---

## Members

| Method | Endpoint                            | Description   |
| ------ | ----------------------------------- | ------------- |
| POST   | `/api/projects/:id/members`         | Add member    |
| DELETE | `/api/projects/:id/members/:userId` | Remove member |

---

## Tasks

| Method | Endpoint                                        | Description   |
| ------ | ----------------------------------------------- | ------------- |
| GET    | `/api/projects/:id/tasks`                       | Get tasks     |
| POST   | `/api/projects/:id/tasks`                       | Create task   |
| PUT    | `/api/tasks/:id`                                | Update task   |
| DELETE | `/api/tasks/:id`                                | Delete task   |
| POST   | `/api/projects/:projectId/tasks/:taskId/assign` | Assign task   |
| PUT    | `/api/tasks/:taskId/unassign`                   | Unassign task |

---

# Installation

## 1. Clone Repository

```bash
git clone https://github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR.git
cd FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR
```

Contributions are welcome — issues and pull requests are open.

---

# Environment Variables

The project ships a `.env.exempl` template (safe to commit — placeholder values
only, no real secrets). Copy it to `.env` and fill in your real values:

```bash
cp .env.exempl .env
```

```env
DB_USER=db_user
DB_PASSWORD=db_password_change_this
DB_HOST=127.0.0.1
DB_PORT=db_port_change_this
DB_NAME=db_name_change_this

JWT_SECRET=development_secret_key_change_this
```

`.env` is gitignored and must never be committed — it's used both when running
the backend directly (`go run main.go`) and by the Dockerized `app` service
below.

> When running via Docker Compose, `DB_HOST` is automatically overridden to
> `mysql` (the service name) for the `app` container, so you don't need a
> separate `.env` for Docker.

---

# 🚀 Run everything with Docker (recommended)

`docker-compose.yml` starts the **whole stack** in one command: MySQL, the Go
backend, and the React frontend (built and served by the backend). No more
starting the database, the backend and the frontend by hand in three terminals.

```bash
docker compose up --build -d
```

This builds and starts two containers:

| Service | What it runs                                                                          | Port   |
| ------- | -------------------------------------------------------------------------------------- | ------ |
| `mysql` | MySQL 8.0 database                                                                     | `3306` |
| `app`   | Go backend (API) + the frontend built with `npm run build`, served as static files     | `3000` |

Once it's up, open:

```txt
http://localhost:3000
```

Useful commands:

```bash
docker compose ps          # check container status
docker compose logs -f app # follow backend logs
docker compose down        # stop everything (the DB volume is kept)
```

> The `app` image is built from the multi-stage `Dockerfile` at the project
> root: stage 1 builds the frontend (`frontend/dist`), stage 2 builds the Go
> binary, and the final image just runs that binary, which serves the API under
> `/api` and the built frontend for everything else (see `main.go`). Your local
> `.env` file is mounted into the container read-only — it is never copied into
> the image.

---

# Manual setup (without Docker)

Useful during active development, for hot reload on the frontend and `go run` on
the backend.

## 1. Start MySQL only

```bash
docker compose up -d mysql
```

## 2. Backend

Install Go dependencies:

```bash
go mod tidy
```

Run backend:

```bash
go run main.go
```

Backend runs on:

```txt
http://localhost:3000
```

## 3. Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend runs on:

```txt
http://localhost:5173
```

The dev server proxies `/api` to the Go backend on port 3000, so log in and
everything else works exactly as it does in the Docker setup.

## 4. Seed test data (optional)

A fresh database has no accounts, so every test starts with a signup form. The
seed fills it instead:

```bash
go run ./cmd/seed
```

It creates five accounts on `@example.test`, three projects and twelve tasks,
then prints the login table. Every seeded account shares the password
`secret123`.

The dataset is built so that each term in `CONTEXT.md` has a visible instance:
all three statuses and priorities, an Overdue task, a Due soon one, an Orphan (a
task carrying a Reminder that nobody is assigned to), a task with two assignees,
a task that is `DONE` past its due date — which must *not* read as Overdue — and
a project with no tasks at all, for the empty state. Due dates are offsets from
the moment you run it, so the dataset never decays into "everything is overdue".

Re-running resets: seeded rows are deleted and rebuilt. Only `@example.test`
accounts and the projects they own are ever touched — your own account and
projects are never in range.

```bash
go run ./cmd/seed -join you@example.com   # also add an existing account as MEMBER
go run ./cmd/seed -force                  # seed a DB_HOST that is not local
```

The seed refuses to run against a non-local `DB_HOST` unless `-force` is passed:
it deletes rows for real, and should be hard to point anywhere but a development
machine.

---

# Database Migration

GORM automatically migrates the database schema at startup:

```go
initializers.DB.AutoMigrate(...)
```

No manual SQL setup required.

---

# Frontend Architecture

The frontend is a **multi-page application**: each screen is its own HTML file
with its own React entry point, and moving between screens is a normal browser
navigation. There is no client-side router.

## Pages

| Page            | File            | Access                      | Content                                    |
| --------------- | --------------- | --------------------------- | ------------------------------------------ |
| Landing         | `index.html`    | public (redirects if logged in) | Presentation and links to log in / sign up |
| Log in          | `login.html`    | public                      | Login form                                 |
| Sign up         | `signup.html`   | public                      | Account creation, then automatic login     |
| Projects        | `projects.html` | protected                   | Project list, creation, deletion           |
| Project detail  | `project.html?id=<id>` | protected            | Members and tasks of one project           |

## Folders

| Folder            | Responsibility                                                        |
| ----------------- | --------------------------------------------------------------------- |
| `src/entries/`    | One entry per HTML page: runs the access guard, then mounts React      |
| `src/pages/`      | One component per page, holding that page's state and API calls        |
| `src/components/` | Shared components — `ui/` holds the primitives (Button, Modal, Field…) |
| `src/lib/`        | `api.js` (HTTP client), `session.js` (localStorage), `auth.js` (login / signup / guards), `format.js`, `constants.js` |

## Authentication

The JWT is stored in `localStorage` and attached automatically to every request
by `src/lib/api.js` — it is never displayed or pasted by hand. Protected pages
run their guard **before** React mounts, so a logged-out visitor is redirected
to `/login.html?next=…` without ever seeing the page. When the API answers `401`
on a request that carried a token, the session is cleared and the user is sent
back to the login page.

## API calls

All calls use a relative path (`/api/...`), so the same build works in both
setups: in production the Go server serves the pages and the API from the same
origin, and in development the Vite dev server proxies `/api` to
`http://localhost:3000` (see `vite.config.js`). **There is no API base URL to
configure anywhere in the interface.**

---

# Contributors

## Gharbi Yassine

* Authentication system
* JWT middleware
* Permissions system
* Projects system
* React frontend migration
* Multi-page frontend architecture (Tailwind CSS)
* Dockerized MySQL migration
* Full Docker setup (backend + frontend)
* README/documentation

## Naman Kumar

* Tasks system
* Task assignment system
* Task models/controllers

---

# Demo Videos (without React front-end)

Part 1: [https://youtu.be/P7U-sndT01s](https://youtu.be/P7U-sndT01s)

Part 2: [https://youtu.be/cLzsbtYu3Bc](https://youtu.be/cLzsbtYu3Bc)

---

# Future Improvements

* Real-time updates
* Notifications
* CI/CD pipeline
* Unit tests
* PostgreSQL support
* Dark mode


