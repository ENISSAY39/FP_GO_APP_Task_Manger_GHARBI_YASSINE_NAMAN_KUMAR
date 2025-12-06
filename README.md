📌 Overview

The backend handles:

User Authentication (Signup, Login, Logout, Token Validation)

Project Management (Create, List, View, Delete, Member Control)

Task Management (Create, Update, Assign, Unassign, List)

Role-Based Permissions

Automatic Frontend Serving (Vite/React)

📁 Project Structure
controllers/
    authController.go
    projectsControllers.go
    tasksControllers.go
    permission_helpers.go

models/
    User.go
    Project.go
    ProjectMember.go
    Task.go
    TaskAssignee.go

initializers/
    loadEnv.go
    database.go
    sync.go

middleware/
    auth.go

frontend/
    dist/ or build/ (auto-served if exists)

main.go

⚙️ Environment Setup
1. Clone the repository
git clone <your-repository-url>
cd <project-folder>

2. Install dependencies
go mod tidy

3. Create .env file
DB_URL="root:password@tcp(localhost:3306)/task_manager?parseTime=true"
JWT_SECRET="your-secret-key"
PORT=3000

4. Start the server
go run main.go


Server runs at:
http://localhost:3000

🌐 Frontend Handling

The backend automatically serves your frontend from these folders:

frontend/dist
frontend/build
public


If found:

Static assets served from /

Fallback to index.html for SPA routing

If not found:

No frontend build found (API-only mode)

🔑 Authentication API
Signup
POST /api/signup

Login
POST /api/login

Validate Token
GET /api/validate

Logout
POST /api/logout

📁 Projects API
Create Project
POST /api/projects

List User Projects
GET /api/projects

Project Details
GET /api/projects/:projectId

Delete Project (Owner only)
DELETE /api/projects/:projectId

👥 Project Members API
Add Member
POST /api/projects/:projectId/members

Remove Member
DELETE /api/projects/:projectId/members/:userId

📝 Task API
Create Task
POST /api/projects/:projectId/tasks

List Tasks
GET /api/projects/:projectId/tasks

Update Task
PUT /api/tasks/:taskId

Delete Task
DELETE /api/tasks/:taskId

Assign Task
PUT /api/tasks/:taskId/assign

Unassign Task
PUT /api/tasks/:taskId/unassign

🎯 Permission Rules
Action	Owner	Member
Create task	✔	✔
Update task	✔	Creator only
Delete task	✔	Creator only
Add/remove members	✔	✘
Delete project	✔	✘
🛢 Database Models
User

ID, Name, Email, PasswordHash

Project

ID, Name, Description, OwnerID

ProjectMember

ProjectID, UserID, Role

Task

Title, Description, Status, Priority, ProjectID, CreatorID

TaskAssignee

TaskID, UserID

👨‍💻 Developers

Yassine Gharbi & Naman Kumar
