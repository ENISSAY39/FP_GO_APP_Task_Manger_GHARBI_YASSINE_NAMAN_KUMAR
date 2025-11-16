# Task Manager - Final Project - 
# By Gharbi Yassine & Naman Kumar 


## Description
Simple Task Manager API built with Go, Gin, Gorm, MySQL.
Features: Auth (JWT cookie), Projects CRUD, Tasks CRUD with assignment and statuses.

## Tech stack
- Go + Gin + GORM
- MySQL
- bcrypt for password hashing
- JWT for auth

## Models
- User, Project, Task (see models folder)

## How to run (dev)
1. Copy `.env.example` to `.env` and set DB_DSN, SECRET_KEY, PORT
2. `go run main.go` (or `go run .`)
3. The API runs on `:8080`

## API Endpoints
- POST /signup
- POST /login
- GET /projects
- POST /projects
- GET /projects/:id
- PUT /projects/:id
- DELETE /projects/:id
- POST /projects/:projectId/tasks
- GET /projects/:projectId/tasks
- GET /tasks/:id
- PUT /tasks/:id
- DELETE /tasks/:id
- PUT /tasks/:id/assign
- PUT /tasks/:id/status

## Division of work
- Person A: Models, DB, Projects endpoints, migrations
- Person B: Auth, middleware, Tasks endpoints, README & tests

## Deliverables
- Public GitHub repo whith youtube demo 

