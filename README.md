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
- React + Vite architecture
- Modular frontend structure
- API-based communication
- Dynamic project/task rendering

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
- React
- Vite
- Modular JavaScript architecture

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
├── frontend/
│   ├── src/
│   │   ├── scripts/
│   │   │   ├── api.js
│   │   │   ├── auth.js
│   │   │   ├── projects.js
│   │   │   ├── tasks.js
│   │   │   ├── members.js
│   │   │   └── utils.js
│   │   │
│   │   ├── styles/
│   │   ├── App.jsx
│   │   ├── main.jsx
│   │   └── main.js
│   │
│   ├── public/
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
git clone <repository_url>
cd <repository_name>
```

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

---

# Database Migration

GORM automatically migrates the database schema at startup:

```go
initializers.DB.AutoMigrate(...)
```

No manual SQL setup required.

---

# Frontend Architecture

The frontend was migrated from a single large JavaScript file to a modular React
architecture.

## Main Modules

| File          | Responsibility          |
| ------------- | ----------------------- |
| `auth.js`     | Login / Signup / Logout |
| `projects.js` | Project management      |
| `tasks.js`    | Task management         |
| `members.js`  | Members management      |
| `api.js`      | API fetch helper        |
| `utils.js`    | Shared utilities        |

---

# Contributors

## Gharbi Yassine

* Authentication system
* JWT middleware
* Permissions system
* Projects system
* React frontend migration
* Frontend modular architecture
* Dockerized MySQL migration
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

* Full React state management
* Better UI/UX
* Real-time updates
* Notifications
* CI/CD pipeline
* Unit tests
* PostgreSQL support
* Frontend routing
* Dark mode


