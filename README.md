

# Task Manager – Full Stack Application Documentation 

A full-stack task management application built with:

- Go + Gin (REST API)
- GORM
- MySQL
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

---

# Tech Stack

## Backend
- Go
- Gin
- GORM
- MySQL
- JWT

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
│
├── frontend/
│   ├── src/
│   │   ├── scripts/
│   │   ├── App.jsx
│   │   └── main.jsx
│   ├── public/
│   └── package.json
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

## Backend

```bash
go mod tidy
go run main.go
```

Backend runs on:

```txt
http://localhost:3000
```

---

## Frontend

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

# Environment Variables

Create a `.env` file:

```env
JWT_SECRET=your_secret_key

DB_USER=root
DB_PASSWORD=your_password
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=task_manager
```

---

# Contributors

## Gharbi Yassine

* Authentication
* JWT middleware
* Permissions system
* Projects system
* Frontend React migration
* Frontend modular architecture

## Naman Kumar

* Tasks system
* Task assignment system
* Task models/controllers

---

# Demo Videos (without React front-end)

Part 1:
[https://youtu.be/P7U-sndT01s](https://youtu.be/P7U-sndT01s)

Part 2:
[https://youtu.be/cLzsbtYu3Bc](https://youtu.be/cLzsbtYu3Bc)

---

# Future Improvements

* Full React state management
* Better UI/UX
* Real-time updates
* Notifications
* Docker deployment
* Unit tests
* CI/CD pipeline

