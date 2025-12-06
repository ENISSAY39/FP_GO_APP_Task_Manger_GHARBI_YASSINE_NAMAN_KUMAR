# Backend API Documentation

This document describes the backend architecture and API endpoints for the Task Manager application built with **Go (Gin)**, **GORM**, and **JWT Authentication**.

---

## 📌 Overview

The backend handles:

* User Authentication (Signup, Login)
* Project Management (Create, List, Delete, Membership Control)
* Task Management (Create, Update, Assign, List)
* Role-Based Permissions

The project follows a modular structure:

```
controllers/
    authController.go
    projectsControllers.go
    tasksControllers.go
    permission_helpers.go
models/
initializers/
middlewares/
routes/
```


<table style="width:100%; border:2px solid #666; border-collapse:collapse;"> <tr style="background:#222; color:#fff; border:2px solid #666;"> <th style="padding:10px; border:2px solid #666;">Method</th> <th style="padding:10px; border:2px solid #666;">Endpoint</th> <th style="padding:10px; border:2px solid #666;">Description</th> <th style="padding:10px; border:2px solid #666;">Auth Required</th> </tr> <!-- AUTH --> <tr style="background:#111; color:#ddd;"> <td style="padding:8px; border:2px solid #666;">POST</td> <td style="padding:8px; border:2px solid #666;">/api/signup</td> <td style="padding:8px; border:2px solid #666;">Create new user</td> <td style="padding:8px; border:2px solid #666;">No</td> </tr> <tr style="background:#181818; color:#ddd;"> <td style="padding:8px; border:2px solid #666;">POST</td> <td style="padding:8px; border:2px solid #666;">/api/login</td> <td style="padding:8px; border:2px solid #666;">Login & receive JWT cookie</td> <td style="padding:8px; border:2px solid #666;">No</td> </tr> <tr style="background:#111; color:#ddd;"> <td style="padding:8px; border:2px solid #666;">GET</td> <td style="padding:8px; border:2px solid #666;">/api/validate</td> <td style="padding:8px; border:2px solid #666;">Validate JWT token</td> <td style="padding:8px; border:2px solid #666;">Yes</td> </tr> <tr style="background:#181818; color:#ddd;"> <td style="padding:8px; border:2px solid #666;">POST</td> <td style="padding:8px; border:2px solid #666;">/api/logout</td> <td style="padding:8px; border:2px solid #666;">Logout user (clear cookie)</td> <td style="padding:8px; border:2px solid #666;">Yes</td> </tr> <!-- PROJECTS --> <tr style="background:#333; color:#fff;"> <td colspan="4" style="text-align:center; padding:10px; border:2px solid #666;"><b>Project Management</b></td> </tr> <tr style="background:#111; color:#ddd;"> <td style="padding:8px; border:2px solid #666;">POST</td> <td style="padding:8px; border:2px solid #666;">/api/projects</td> <td style="padding:8px; border:2px solid #666;">Create new project</td> <td style="padding:8px; border:2px solid #666;">Yes</td> </tr> <tr style="background:#181818; color:#ddd;"> <td style="padding:8px; border:2px solid #666;">GET</td> <td style="padding:8px; border:2px solid #666;">/api/projects</td> <td style="padding:8px; border:2px solid #666;">Get projects of logged-in user</td> <td style="padding:8px; border:2px solid #666;">Yes</td> </tr> <tr style="background:#111; color:#ddd;"> <td style="padding:8px; border:2px solid #666;">GET</td> <td style="padding:8px; border:2px solid #666;">/api/projects/:projectId</td> <td style="padding:8px; border:2px solid #666;">Get project details</td> <td style="padding:8px; border:2px solid #666;">Yes</td> </tr> <tr style="background:#181818; color:#ddd;"> <td style="padding:8px; border:2px solid #666;">DELETE</td> <td style="padding:8px; border:2px solid #666;">/api/projects/:projectId</td> <td style="padding:8px; border:2px solid #666;">Delete project (Owner only)</td> <td style="padding:8px; border:2px solid #666;">Yes</td> </tr> <!-- MEMBERS --> <tr style="background:#333; color:#fff;"> <td colspan="4" style="text-align:center; padding:10px; border:2px solid #666;"><b>Project Members</b></td> </tr> <tr style="background:#111; color:#ddd;"> <td style="padding:8px; border:2px solid #666;">POST</td> <td style="padding:8px; border:2px solid #666;">/api/projects/:projectId/members</td> <td style="padding:8px; border:2px solid #666;">Add project member</td> <td style="padding:8px; border:2px solid #666;">Yes (Owner)</td> </tr> <tr style="background:#181818; color:#ddd;"> <td style="padding:8px; border:2px solid #666;">DELETE</td> <td style="padding:8px; border:2px solid #666;">/api/projects/:projectId/members/:userId</td> <td style="padding:8px; border:2px solid #666;">Remove member</td> <td style="padding:8px; border:2px solid #666;">Yes (Owner)</td> </tr> <!-- TASKS --> <tr style="background:#333; color:#fff;"> <td colspan="4" style="text-align:center; padding:10px; border:2px solid #666;"><b>Task Management</b></td> </tr> <tr style="background:#111; color:#ddd;"> <td style="padding:8px; border:2px solid #666;">POST</td> <td style="padding:8px; border:2px solid #666;">/api/projects/:projectId/tasks</td> <td style="padding:8px; border:2px solid #666;">Create task</td> <td style="padding:8px; border:2px solid #666;">Yes</td> </tr> <tr style="background:#181818; color:#ddd;"> <td style="padding:8px; border:2px solid #666;">GET</td> <td style="padding:8px; border:2px solid #666;">/api/projects/:projectId/tasks</td> <td style="padding:8px; border:2px solid #666;">Get project tasks</td> <td style="padding:8px; border:2px solid #666;">Yes</td> </tr> <tr style="background:#111; color:#ddd;"> <td style="padding:8px; border:2px solid #666;">PUT</td> <td style="padding:8px; border:2px solid #666;">/api/tasks/:taskId</td> <td style="padding:8px; border:2px solid #666;">Update task</td> <td style="padding:8px; border:2px solid #666;">Yes</td> </tr> <tr style="background:#181818; color:#ddd;"> <td style="padding:8px; border:2px solid #666;">DELETE</td> <td style="padding:8px; border:2px solid #666;">/api/tasks/:taskId</td> <td style="padding:8px; border:2px solid #666;">Delete task</td> <td style="padding:8px; border:2px solid #666;">Yes</td> </tr> <tr style="background:#111; color:#ddd;"> <td style="padding:8px; border:2px solid #666;">PUT</td> <td style="padding:8px; border:2px solid #666;">/api/tasks/:taskId/assign</td> <td style="padding:8px; border:2px solid #666;">Assign task</td> <td style="padding:8px; border:2px solid #666;">Yes</td> </tr> <tr style="background:#181818; color:#ddd;"> <td style="padding:8px; border:2px solid #666;">PUT</td> <td style="padding:8px; border:2px solid #666;">/api/tasks/:taskId/unassign</td> <td style="padding:8px; border:2px solid #666;">Unassign task</td> <td style="padding:8px; border:2px solid #666;">Yes</td> </tr> </table>

---

## 🔐 Authentication

### **POST /signup** – Create a New User

**Request Body:**

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "secret123"
}
```

**Validations:**

* Name ≥ 2 characters
* Valid email
* Password ≥ 6 characters

**Response:**

```json
{
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

---

### **POST /login** – Login & Get JWT Token

**Request Body:**

```json
{
  "email": "john@example.com",
  "password": "secret123"
}
```

**Response:**

```json
{
  "token": "<jwt-token>",
  "user": { "id": 1, "name": "John Doe", "email": "john@example.com" }
}
```

Token contains:

* `sub`: user ID
* `exp`: expiry (24h)

---

## 📁 Project Management

All project routes require **JWT authentication**.

### **POST /projects** – Create Project

Request Body:

```json
{
  "name": "My Project",
  "description": "Project description"
}
```

The creator becomes:

* Project Owner
* Project Member with role "OWNER"

---

### **GET /projects/my** – Get All Projects of Logged-in User

Returns list of projects where the user is a member.

---

### **GET /projects/:projectId** – Get Project Details

Only accessible if user is a project member.
Loads:

* Members
* Tasks

---

### **DELETE /projects/:projectId** – Delete Project

Only owner can delete.

---

## 👥 Project Members

### **POST /projects/:projectId/members** – Add Member (OWNER Only)

Body:

```json
{
  "user_id": 5,
  "role": "MEMBER" // optional
}
```

### **DELETE /projects/:projectId/members** – Remove Member

Body:

```json
{
  "user_id": 5
}
```

Owner cannot remove themselves.

---

## 📝 Task Management

### **POST /projects/:projectId/tasks** – Create Task

Members only.

Fields:

```json
{
  "title": "Task Name",
  "description": "Optional text",
  "priority": "HIGH", // optional
  "due_date": "2025-01-20T10:00:00Z"
}
```

---

### **GET /projects/:projectId/tasks** – List Tasks

Includes task assignees.

---

### **PATCH /tasks/:taskId** – Update Task

Allowed by:

* Task creator OR
* Project owner

Partial update fields:

```json
{
  "title": "New title",
  "status": "DONE",
  "priority": "LOW"
}
```

---

## 🔑 Roles & Permission Logic

Defined in `permission_helpers.go`:

* `IsProjectOwner(projectID, userID)`
* `IsProjectMember(projectID, userID)`
* `AddProjectMember()`
* `RemoveProjectMember()`

Roles:

* **OWNER** – full permissions
* **MEMBER** – limited to tasks + viewing

---

## 🧱 Models (Summary)

* User
* Project
* Task
* ProjectMember (pivot table)

---

## 🚀 Running the Server

### Environment Variables

```
JWT_SECRET=your_secret_key
DB_URL=mysql://user:password@tcp(127.0.0.1:3306)/dbname
```

### Start Server

```
go run main.go
```

---

## 📬 Contact

For improvements or bugs, please reach out.
