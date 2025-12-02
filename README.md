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
3. The API runs on `The API runs on http://localhost:3000`

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
- Yassine : Project Model&Controller ,Projects Model&Controlle0r Task DB, endpoints, migrations
- Naman : Users Model&Controller , Auth, middleware, Tasks endpoints, README & tests

## Deliverables
- Public GitHub repo whith youtube demo <>



Dans une tache tous les utilisateur ont le meme role 
la tache se definis quand je crée la tache 


    Dans le login une étape de vérif pour savoir si on redirige vers interface : 
    - utilisateur (savoir les projet et les taches dans les quels il est affécté )
    - Chef de projet (une liste des chef de projet sera donc faite dans le back end dans le require Auth ou je sais pas comment faire ducoup )




